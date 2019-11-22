package model

import (
	"time"

	"github.com/google/uuid"
)

// PhoneBook ...
type PhoneBook struct {
	ID             uuid.UUID  `db:"id" json:"id" validate:"required"`
	Fullname       *string    `db:"fullname" json:"fullname"`
	PhoneNumber    *string    `db:"phone_number" json:"phone_number"`
	Address        *string    `db:"address" json:"address"`
	CreatedDateUTC *time.Time `db:"created_date_utc" json:"created_date_utc"`
	CreatedBy      *string    `db:"created_by" json:"created_by"`
	UpdatedDateUTC *time.Time `db:"updated_date_by" json:"updated_date_utc"`
	UpdatedBy      *string    `db:"updated_by" json:"updated_by"`
	DeletedDateUTC *time.Time `db:"deleted_date_utc" json:"deleted_date_utc"`
	DeletedBy      *string    `db:"deleted_by" json:"deleted_by"`
}

// GetPhoneList ...
type GetPhoneList struct {
	ID          uuid.UUID `json:"id" httpquery:"id"`
	Fullname    *string   `json:"fullname" httpquery:"fullname"`
	PhoneNumber *string   `json:"phone_number" httpquery:"phone_number"`
	Address     *string   `json:"address" httpquery:"address"`
	OffsetID    uuid.UUID `json:"offset_id" httpquery:"offset_id"`
	Limit       *int      `json:"limit" httpquery:"limit"`
}
