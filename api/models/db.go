package models

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type Database struct {
	DB *sql.DB
}

func DB(host, port, user, password, dbname string) (*Database, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", dsn)

	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	d := &Database{DB: db}
	return d, nil
}

func (d *Database) GetTicketOption(id uuid.UUID) (TicketOption, error) {
	stmt := "SELECT id, name, description, allocation FROM ticket_options WHERE id = $1"
	row := d.DB.QueryRow(stmt, id)

	var ticketOption TicketOption
	err := row.Scan(&ticketOption.ID, &ticketOption.Name, &ticketOption.Description, &ticketOption.Allocation)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return TicketOption{}, ErrNoRecord
		} else {
			return TicketOption{}, err
		}
	}

	return ticketOption, nil
}

func (d *Database) CreateTicketOption(name string, description string, allocation int) (TicketOption, error) {
	stmt := "INSERT INTO ticket_options (name, description, allocation) VALUES ($1, $2, $3) RETURNING id"
	result := d.DB.QueryRow(stmt, name, description, allocation)

	var ticketOptionID uuid.UUID
	err := result.Scan(&ticketOptionID)

	if err != nil {
		log.Print(err.Error())
		return TicketOption{}, err
	}

	ticketOption := TicketOption{ID: ticketOptionID, Name: name, Description: description, Allocation: allocation}
	return ticketOption, nil
}

func (d *Database) CreatePurchase(ticketOptionId uuid.UUID, requestedQuantity int, userId uuid.UUID) (Purchase, error) {
	tx, err := d.DB.Begin()
	if err != nil {
		log.Print(err.Error())
		return Purchase{}, err
	}
	defer tx.Rollback()

	var currentAllocation int
	err = tx.QueryRow("SELECT allocation FROM ticket_options WHERE id = $1 FOR UPDATE", ticketOptionId).Scan(&currentAllocation)

	if requestedQuantity > currentAllocation {
		log.Print("Insufficient quantity in allocation")
		return Purchase{}, ErrInvalidQuantity
	}

	newAllocation := currentAllocation - requestedQuantity
	_, err = tx.Exec("UPDATE ticket_options SET allocation = $1, updated_at = NOW() WHERE id = $2", newAllocation, ticketOptionId)
	if err != nil {
		log.Print(err.Error())
		return Purchase{}, err
	}

	var purchaseID uuid.UUID
	err = tx.QueryRow("INSERT INTO purchases (ticket_option_id, user_id, quantity) VALUES ($1, $2, $3) RETURNING id",
		ticketOptionId, userId, requestedQuantity).Scan(&purchaseID)

	if err != nil {
		log.Print(err.Error())
		return Purchase{}, err
	}

	err = tx.Commit()
	if err != nil {
		log.Print(err.Error())
		return Purchase{}, err
	}

	purchase := Purchase{purchaseID, ticketOptionId, userId, requestedQuantity}
	return purchase, nil

}
