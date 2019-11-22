package phonebook

import (
	// internal golang package
	"context"
	"errors"

	// internal package
	"phonebook/internal/phonebook/model"
	"phonebook/internal/phonebook/repository"

	// thirdparty package
	"github.com/go-kit/kit/log"
	"github.com/google/uuid"
)

type Service struct {
	Actor  string
	Logger log.Logger
	repo   repository.Interface
}

func NewService(repo repository.Interface, actor string, logger log.Logger) *Service {
	return &Service{
		Actor:  actor,
		Logger: logger,
		repo:   repo,
	}
}

// CreatePhoneAddress ...
func (svc *Service) CreatePhoneAddress(ctx context.Context, data *model.PhoneBook) error {
	id, err := uuid.NewRandom()
	if nil != err {
		return err
	}

	data.ID = id

	params := &model.GetPhoneList{
		PhoneNumber: data.PhoneNumber,
	}

	result, err := svc.repo.ListPhoneBook(ctx, params)
	if nil != err {
		return err
	}

	if result != nil {
		return errors.New("phone_already_registered")
	}

	err = svc.repo.AddingPerson(ctx, data)
	if nil != err {
		return err
	}

	return nil
}

// FetchData fetching data of phone book
func (svc *Service) FetchData(ctx context.Context, params *model.GetPhoneList) ([]*model.PhoneBook, error) {
	result, err := svc.repo.ListPhoneBook(ctx, params)
	if nil != err {
		return nil, err
	}

	return result, nil
}

// RemoveData , remove profile on phone book
func (svc *Service) RemoveData(ctx context.Context, data *model.PhoneBook) error {
	res, err := svc.repo.FetchByID(ctx, data.ID)
	if nil != err {
		return err
	}

	if res == nil {
		return errors.New("profile_not_exist")
	}

	err = svc.repo.RemoveData(ctx, &model.PhoneBook{
		ID: data.ID,
	})

	if nil != err {
		return err
	}

	return nil
}

// UpdateData update data profile
func (svc *Service) UpdateData(ctx context.Context, data *model.PhoneBook) error {
	res, err := svc.repo.FetchByID(ctx, data.ID)
	if nil != err {
		return err
	}

	if res == nil {
		return errors.New("profile_not_exist")
	}

	err = svc.repo.UpdatePerson(ctx, data)
	if nil != err {
		return err
	}

	return nil
}
