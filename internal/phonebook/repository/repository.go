package repository

import (
	// internal golang package
	"context"

	// internal package
	"phonebook/internal/phonebook/model"

	"github.com/google/uuid"
)

type Interface interface {
	ListPhoneBook(ctx context.Context, params *model.GetPhoneList) ([]*model.PhoneBook, error)
	AddingPerson(ctx context.Context, data *model.PhoneBook) error
	UpdatePerson(ctx context.Context, data *model.PhoneBook) error
	RemoveData(ctx context.Context, data *model.PhoneBook) error
	FetchByID(ctx context.Context, id uuid.UUID) (*model.PhoneBook, error)
}
