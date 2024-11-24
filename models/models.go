package models

import (
	"time"

	"github.com/google/uuid"
)

type BaseModel struct {
	ID        uuid.UUID  `db:"id" json:"id"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}

type User struct {
	BaseModel
	Fullname string `db:"fullname" json:"fullname"`
	Email    string `db:"email" json:"email"`
	Password string `db:"password" json:"-"`
	Phone    string `db:"phone" json:"phone"`
}

type Location struct {
	BaseModel
	Name       string `db:"name" json:"name"`
	Address    string `db:"address" json:"address"`
	City       string `db:"city" json:"city"`
	State      string `db:"state" json:"state"`
	Country    string `db:"country" json:"country"`
	PostalCode string `db:"postal_code" json:"postal_code"`
}

type Schedule struct {
	BaseModel
	Title       string    `db:"title" json:"title"`
	Description string    `db:"description" json:"description"`
	StartDate   time.Time `db:"start_date" json:"start_date"`
	EndDate     time.Time `db:"end_date" json:"end_date"`
}

type Event struct {
	BaseModel
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	LocationID  uuid.UUID `db:"location_id" json:"location_id"`
	ScheduleID  uuid.UUID `db:"schedule_id" json:"schedule_id"`
	Location    *Location `db:"-" json:"location,omitempty"`
	Schedule    *Schedule `db:"-" json:"schedule,omitempty"`
}

type TicketType struct {
	BaseModel
	EventID        uuid.UUID `db:"event_id" json:"event_id"`
	Name           string    `db:"name" json:"name"`
	Description    string    `db:"description" json:"description"`
	Price          float64   `db:"price" json:"price"`
	Quota          int       `db:"quota" json:"quota"`
	RemainingQuota int       `db:"remaining_quota" json:"remaining_quota"`
	Event          *Event    `db:"-" json:"event,omitempty"`
}

type Transaction struct {
	BaseModel
	UserID        uuid.UUID           `db:"user_id" json:"user_id"`
	EventID       uuid.UUID           `db:"event_id" json:"event_id"`
	TotalAmount   float64             `db:"total_amount" json:"total_amount"`
	Status        string              `db:"status" json:"status"`
	PaymentMethod string              `db:"payment_method" json:"payment_method"`
	PaymentStatus string              `db:"payment_status" json:"payment_status"`
	User          *User               `db:"-" json:"user,omitempty"`
	Event         *Event              `db:"-" json:"event,omitempty"`
	Details       []TransactionDetail `db:"-" json:"details,omitempty"`
}

type TransactionDetail struct {
	BaseModel
	TransactionID  uuid.UUID    `db:"transaction_id" json:"transaction_id"`
	TicketTypeID   uuid.UUID    `db:"ticket_type_id" json:"ticket_type_id"`
	Quantity       int          `db:"quantity" json:"quantity"`
	PricePerTicket float64      `db:"price_per_ticket" json:"price_per_ticket"`
	Subtotal       float64      `db:"subtotal" json:"subtotal"`
	Transaction    *Transaction `db:"-" json:"transaction,omitempty"`
	TicketType     *TicketType  `db:"-" json:"ticket_type,omitempty"`
}
