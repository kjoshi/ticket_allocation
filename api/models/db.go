package models

import (
	"database/sql"
	"errors"
	"fmt"

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

	ticketOption := TicketOption{Name: name, Description: description, Allocation: allocation}
	if err != nil {
		ticketOption.ID = uuid.Nil
		return ticketOption, err
	}

	ticketOption.ID = ticketOptionID
	return ticketOption, nil
}
