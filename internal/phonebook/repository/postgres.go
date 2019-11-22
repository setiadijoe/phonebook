package repository

import (
	// internal golang package
	"bytes"
	"context"
	"time"

	// internal package
	"phonebook/internal/global"
	"phonebook/internal/phonebook/model"

	// thirdparty package
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository struct{}

func NewPostgres() *Repository {
	return &Repository{}
}

// ListPhoneBook , get list of phone book
func (r *Repository) ListPhoneBook(ctx context.Context, getparams *model.GetPhoneList) ([]*model.PhoneBook, error) {
	result := make([]*model.PhoneBook, 0)

	var rows *sqlx.Rows
	var params = make([]interface{}, 0)
	var err error
	var buffer, selectBuffer bytes.Buffer
	var first bool
	var count int
	var createdDateUTC time.Time

	buffer.WriteString(`SELECT id, fullname, phone_number,
	address, created_date_utc, created_by, updated_date_utc, updated_by
	 FROM phone_book WHERE`)

	if uuid.Nil != getparams.OffsetID {
		selectBuffer.WriteString(`SELECT created_date_utc FROM phone_book WHERE id = :id`)
		rows, err = global.DB().NamedQuery(selectBuffer.String(), getparams.OffsetID)
		if nil != err {
			return nil, err
		}

		for rows.Next() {
			err = rows.Scan(&createdDateUTC)
			if nil != err {
				return nil, err
			}
		}
	}

	if uuid.Nil != getparams.ID {
		first = false
		count++
		buffer.WriteString(`id = :id`)
		params = append(params, getparams.ID)
	} else if nil != getparams.Fullname {
		if !first {
			buffer.WriteString(`AND`)
		}
		buffer.WriteString(`fullname ilike :fullname`)
		params = append(params, getparams.Fullname)
		first = false
	} else if nil != getparams.PhoneNumber {
		if !first {
			buffer.WriteString(`AND`)
		}
		buffer.WriteString(`phone_number like :phone_number`)
		params = append(params, getparams.Address)
	} else if uuid.Nil != getparams.OffsetID {
		if !first {
			buffer.WriteString(`AND`)
		}
		buffer.WriteString(`created_date_utc > :created_date_utc  `)
		params = append(params, createdDateUTC)
	}

	buffer.WriteString(` deleted_date_utc IS NULL `)

	if nil != getparams.Limit {
		buffer.WriteString(`LIMIT :limit`)
		params = append(params, getparams.Limit)
	}

	buffer.WriteString(`ORDER BY created_date_utc`)

	if ctx == nil {
		rows, err = global.DB().NamedQuery(buffer.String(), params)
	} else {
		rows, err = global.DB().NamedQueryContext(ctx, buffer.String(), params)
	}

	if nil != err {
		return nil, err
	}

	for rows.Next() {
		phone := &model.PhoneBook{}
		err = rows.Scan(&phone.ID, &phone.Fullname, &phone.PhoneNumber, &phone.Address, &phone.CreatedDateUTC, &phone.CreatedBy, &phone.UpdatedDateUTC, &phone.UpdatedBy)
		if nil != err {
			return nil, err
		}
		result = append(result, phone)
	}

	return result, nil
}

// AddingPerson , adding new person to phone book
func (r *Repository) AddingPerson(ctx context.Context, data *model.PhoneBook) error {
	var err error
	var buffer bytes.Buffer
	var args = make([]interface{}, 0)

	buffer.WriteString(` INSERT INTO phone_book ( id, fullname, phone_number, address, created_by, updated_by) 
	VALUES (:id, :fullname, :phone_number, :address, :created_by, :updated_by)`)
	args = append(args, data.ID, data.Fullname, data.PhoneNumber, data.Address, data.CreatedBy, data.UpdatedBy)
	if nil != ctx {
		_, err = global.DB().NamedExecContext(ctx, buffer.String(), args)
	} else {
		_, err = global.DB().NamedExec(buffer.String(), args)
	}

	if err != nil {
		return err
	}

	return nil
}

// UpdatePerson , update profile of person
func (r *Repository) UpdatePerson(ctx context.Context, data *model.PhoneBook) error {
	var err error
	var buffer bytes.Buffer
	var args = make([]interface{}, 0)
	var first bool

	buffer.WriteString(` UPDATE phone_book SET `)
	if nil != data.Fullname {
		buffer.WriteString(` fullname = :fullname `)
		args = append(args, data.Fullname)
		first = false
	} else if nil != data.PhoneNumber {
		if !first {
			buffer.WriteString(`,`)
		}
		buffer.WriteString(` phone_number = :phone_number `)
		args = append(args, data.PhoneNumber)
		first = false
	} else if nil != data.Address {
		if !first {
			buffer.WriteString(`, `)
		}
		buffer.WriteString(` address = :address`)
		args = append(args, data.Address)
		first = false
	} else if nil != data.UpdatedBy {
		buffer.WriteString(` updated_by = :updated_by`)
		args = append(args, data.UpdatedBy)
	}

	buffer.WriteString(`WHERE id = :id`)
	args = append(args, data.ID)

	if nil != ctx {
		_, err = global.DB().NamedExecContext(ctx, buffer.String(), args)
	} else {
		_, err = global.DB().NamedExec(buffer.String(), args)
	}

	if nil != err {
		return err
	}

	return nil
}

// RemoveData , remove profile but set deleted date utc
func (r *Repository) RemoveData(ctx context.Context, data *model.PhoneBook) error {
	var err error
	var buffer bytes.Buffer
	var args = make([]interface{}, 0)

	buffer.WriteString(`UPDATE phone_book SET deleted_date_utc = :deleted_date_utc , deleted_by = :deleted_by
	WHERE id = :id`)

	args = append(args, data.DeletedDateUTC, data.DeletedBy, data.ID)

	if nil != ctx {
		_, err = global.DB().NamedExecContext(ctx, buffer.String(), args)
	} else {
		_, err = global.DB().NamedExec(buffer.String(), args)
	}

	if nil != err {
		return err
	}

	return nil

}

// FetchByID , get one profile from phone book
func (r *Repository) FetchByID(ctx context.Context, id uuid.UUID) (*model.PhoneBook, error) {
	var err error
	var buffer bytes.Buffer
	var result = &model.PhoneBook{}

	buffer.WriteString(`SELECT * FROM phone_book WHERE id = $1 AND deleted_date_utc IS NULL`)
	if nil != ctx {
		err = global.DB().GetContext(ctx, *result, buffer.String(), id)
	} else {
		err = global.DB().Get(&result, buffer.String(), id)
	}

	if nil != err {
		return nil, err
	}

	return result, nil

}
