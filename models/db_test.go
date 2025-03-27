package models

import (
	"errors"
	"github.com/google/uuid"
	"sync"
	"testing"
)

func Equal[T comparable](t *testing.T, actual, expected T) {
	t.Helper()

	if actual != expected {
		t.Errorf("got: %v; expected: %v", actual, expected)
	}
}

var validTicketOptionId, _ = uuid.Parse("969f4317-09f4-4b15-b8be-a87d40fb56fb")
var invalidTicketOptionId, _ = uuid.Parse("11111111-1111-1111-1111-111111111111")
var validUserId, _ = uuid.Parse("d6abe829-c28c-44ec-bee6-3183f2c53fef")

func TestGetTicketOption(t *testing.T) {
	tests := []struct {
		name           string
		ticketOptionId uuid.UUID
		expected       uuid.UUID
		errorType      error
	}{
		{
			name:           "Valid ID",
			ticketOptionId: validTicketOptionId,
			expected:       validTicketOptionId,
			errorType:      nil,
		},
		{
			name:           "Invalid ID",
			ticketOptionId: invalidTicketOptionId,
			expected:       uuid.Nil,
			errorType:      ErrNoRecord,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := newTestDB(t)

			ticketOption, err := db.GetTicketOption(tt.ticketOptionId)

			Equal(t, ticketOption.ID, tt.expected)

			if !errors.Is(err, tt.errorType) {
				t.Errorf("got: %v; expected: %v", err.Error(), tt.errorType)
			}
		})
	}
}

func TestCreatePurchase(t *testing.T) {
	db := newTestDB(t)

	var wg sync.WaitGroup
	successCount := 0
	var mu sync.Mutex

	// Make 5 concurrent purchases of 3 tickets each (total 15 tickets)
	// Since only 10 tickets are available (see testdata/setup.sql),
	// some purchases should fail
	for i := range 5 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_, err := db.CreatePurchase(validTicketOptionId, 3, validUserId)
			if err == nil {
				mu.Lock()
				successCount++
				mu.Unlock()
			}
		}(i)
	}

	wg.Wait()

	// Verify that only 3 purchases succeeded (3 purchases of 3 tickets each = 9 tickets)
	if successCount != 3 {
		t.Errorf("Expected 3 successful purchases, got %d", successCount)
	}

	// Verify that the ticket_option has an allocation of 1
	ticketOption, err := db.GetTicketOption(validTicketOptionId)
	remainingAllocation := ticketOption.Allocation
	if err != nil {
		t.Fatalf("Failed to get ticket allocation: %v", err)
	}

	if remainingAllocation != 1 {
		t.Errorf("Expected 1 ticket remaining, got %d", remainingAllocation)
	}

	// Check that purchases have been created
	// TODO: Add function to get purchases
	var totalPurchased int
	err = db.DB.QueryRow("SELECT SUM(quantity) FROM purchases WHERE ticket_option_id = $1", validTicketOptionId).Scan(&totalPurchased)
	if err != nil {
		t.Fatalf("Failed to query total purchased: %v", err)
	}

	if totalPurchased != 9 {
		t.Errorf("Ticket accounting error: expected 9 purchased, got %d", totalPurchased)
	}

	if (totalPurchased + remainingAllocation) != 10 {
		t.Errorf("Ticket accounting error: %d purchased + %d remaining != 10 initial", totalPurchased, remainingAllocation)
	}

}
