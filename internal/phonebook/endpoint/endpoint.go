package endpoint

import (
	// internal golang package
	"context"

	// internal package
	"phonebook/internal/global"
	"phonebook/internal/phonebook"
	"phonebook/internal/phonebook/model"
	"phonebook/pkg/queryable"

	// thirdparty package
	"github.com/go-kit/kit/endpoint"
)

// FetchData ...
func FetchData(svc *phonebook.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		err = queryable.RunInTransaction(ctx, global.DB(), func(ctx context.Context) error {
			reqData := request.(*model.GetPhoneList)
			response, err = svc.FetchData(ctx, reqData)
			return err
		})
		return response, err
	}
}

// Add ...
func Add(svc *phonebook.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		err = queryable.RunInTransaction(ctx, global.DB(), func(ctx context.Context) error {
			reqData := request.(*model.PhoneBook)
			return svc.CreatePhoneAddress(ctx, reqData)
		})

		return nil, err
	}
}

// Update ...
func Update(svc *phonebook.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		err = queryable.RunInTransaction(ctx, global.DB(), func(ctx context.Context) error {
			reqData := request.(*model.PhoneBook)
			return svc.UpdateData(ctx, reqData)
		})
		return nil, err
	}
}

// Remove ...
func Remove(svc *phonebook.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		err = queryable.RunInTransaction(ctx, global.DB(), func(ctx context.Context) error {
			reqData := request.(*model.PhoneBook)
			return svc.RemoveData(ctx, reqData)
		})
		return nil, err
	}
}
