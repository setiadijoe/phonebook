package http

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/iancoleman/strcase"

	"phonebook/pkg/common"
	"phonebook/pkg/httperror"
	"phonebook/pkg/validator"
)

//DecodeOptions executed before decode process
type DecodeOptions func(ctx context.Context, model interface{}, request *http.Request) error

//DecodeParam decode model with DecodeOptions
type DecodeParam struct {
	Model   interface{}
	Options []DecodeOptions
}

//Decode generate a decode function to decode request body (json) to model
func Decode(model interface{}) func(context.Context, *http.Request) (request interface{}, err error) {
	return func(ctx context.Context, r *http.Request) (interface{}, error) {
		if model == nil {
			return nil, nil
		}

		var _model interface{}
		var err error

		param, ok := model.(DecodeParam)
		if ok {
			_model, _ = common.DeepCopy(param.Model)
			for _, option := range param.Options {
				err = option(ctx, _model, r)
				if err != nil {
					httperr := &httperror.ErrorWithStatusCode{
						Err:        err.Error(),
						StatusCode: http.StatusUnprocessableEntity,
					}
					return nil, httperr
				}
			}
		} else {
			_model, _ = common.DeepCopy(model)
		}

		if r.ContentLength != 0 {
			contentType := r.Header["Content-Type"]
			if common.StringInSlice("application/json", contentType) {
				_model, err = ParseJSON(ctx, r, _model)
				if err != nil {
					httperr := &httperror.ErrorWithStatusCode{
						Err:        err.Error(),
						StatusCode: http.StatusUnprocessableEntity,
					}
					return nil, httperr
				}
			}
		}

		err = getURLParamUsingTag(ctx, _model, r)
		if err != nil {
			httperr := &httperror.ErrorWithStatusCode{
				Err:        err.Error(),
				StatusCode: http.StatusUnprocessableEntity,
			}
			return nil, httperr
		}

		err = GetQueryUsingTag(ctx, _model, r)
		if err != nil {
			httperr := &httperror.ErrorWithStatusCode{
				Err:        err.Error(),
				StatusCode: http.StatusUnprocessableEntity,
			}
			return nil, httperr
		}

		err = GetHeaderUsingTag(ctx, _model, r)
		if err != nil {
			httperr := &httperror.ErrorWithStatusCode{
				Err:        err.Error(),
				StatusCode: http.StatusUnprocessableEntity,
			}
			return nil, httperr
		}

		err = validator.DefaultValidator()(_model)
		if err != nil {
			httperr := &httperror.ErrorWithStatusCode{
				Err:        err.Error(),
				StatusCode: http.StatusUnprocessableEntity,
			}
			return nil, httperr
		}

		return _model, nil
	}
}

//DecodeJSON generate a decode function to decode request body (json) to model
func DecodeJSON(model interface{}) func(context.Context, *http.Request) (request interface{}, err error) {
	return func(ctx context.Context, r *http.Request) (interface{}, error) {
		if model == nil {
			return nil, nil
		}

		var _model interface{}
		var err error

		param, ok := model.(DecodeParam)
		if ok {
			_model, _ = common.DeepCopy(param.Model)
			for _, option := range param.Options {
				err = option(ctx, _model, r)
				if err != nil {
					httperr := &httperror.ErrorWithStatusCode{
						Err:        err.Error(),
						StatusCode: http.StatusUnprocessableEntity,
					}
					return nil, httperr
				}
			}
		} else {
			_model, _ = common.DeepCopy(model)
		}

		if r.ContentLength != 0 {
			_model, err = ParseJSON(ctx, r, _model)
			if err != nil {
				httperr := &httperror.ErrorWithStatusCode{
					Err:        err.Error(),
					StatusCode: http.StatusUnprocessableEntity,
				}
				return nil, httperr
			}
		}

		err = getURLParamUsingTag(ctx, _model, r)
		if err != nil {
			httperr := &httperror.ErrorWithStatusCode{
				Err:        err.Error(),
				StatusCode: http.StatusUnprocessableEntity,
			}
			return nil, httperr
		}

		err = GetQueryUsingTag(ctx, _model, r)
		if err != nil {
			httperr := &httperror.ErrorWithStatusCode{
				Err:        err.Error(),
				StatusCode: http.StatusUnprocessableEntity,
			}
			return nil, httperr
		}

		err = GetHeaderUsingTag(ctx, _model, r)
		if err != nil {
			httperr := &httperror.ErrorWithStatusCode{
				Err:        err.Error(),
				StatusCode: http.StatusUnprocessableEntity,
			}
			return nil, httperr
		}

		err = validator.DefaultValidator()(_model)
		if err != nil {
			httperr := &httperror.ErrorWithStatusCode{
				Err:        err.Error(),
				StatusCode: http.StatusUnprocessableEntity,
			}
			return nil, httperr
		}

		return _model, nil
	}
}

//GetURLParam built-in DecodeOptions for decode using url params
func GetURLParam(params []string) DecodeOptions {
	return func(ctx context.Context, model interface{}, r *http.Request) error {
		var err error
		typ := reflect.TypeOf(model).Elem()
		for i := 0; i < typ.NumField(); i++ {
			name := typ.Field(i).Name
			name = strcase.ToSnake(name)
			if common.StringInSlice(name, params) {
				err = getURLParam(ctx, model, r, name, i)
				if err != nil {
					return err
				}
			}
		}
		return nil
	}
}

func getURLParamUsingTag(ctx context.Context, model interface{}, r *http.Request) error {
	var err error
	typ := reflect.TypeOf(model).Elem()
	for i := 0; i < typ.NumField(); i++ {
		tag := typ.Field(i).Tag.Get("httpurl")
		if tag != "" {
			err = getURLParam(ctx, model, r, tag, i)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func getURLParam(ctx context.Context, model interface{}, r *http.Request, param string, valIdx int) error {
	value := chi.URLParam(r, param)
	if value == "" {
		return nil
	}

	return fillFieldValue(model, value, valIdx)
}

//GetQueryUsingTag ...
func GetQueryUsingTag(ctx context.Context, model interface{}, r *http.Request) error {
	var err error
	typ := reflect.TypeOf(model).Elem()
	for i := 0; i < typ.NumField(); i++ {
		tag := typ.Field(i).Tag.Get("httpquery")
		if tag != "" {
			err = getQuery(ctx, model, r, tag, i)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func getQuery(ctx context.Context, model interface{}, r *http.Request, query string, valIdx int) error {
	value := r.URL.Query().Get(query)
	if value == "" {
		return nil
	}

	return fillFieldValue(model, value, valIdx)
}

//GetHeaderUsingTag ...
func GetHeaderUsingTag(ctx context.Context, model interface{}, r *http.Request) error {
	var err error
	typ := reflect.TypeOf(model).Elem()
	for i := 0; i < typ.NumField(); i++ {
		tag := typ.Field(i).Tag.Get("httpheader")
		if tag != "" {
			err = getHeader(ctx, model, r, tag, i)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func getHeader(ctx context.Context, model interface{}, r *http.Request, header string, valIdx int) error {
	value := r.Header.Get(header)
	if value == "" {
		return nil
	}

	return fillFieldValue(model, value, valIdx)
}

func fillFieldValue(model interface{}, value string, valIdx int) error {
	val := reflect.ValueOf(model).Elem()

	switch valtype := val.Field(valIdx).Type().String(); valtype {
	case "string":
		val.Field(valIdx).SetString(value)
	case "int":
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		val.Field(valIdx).SetInt(v)
	case "int64":
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		val.Field(valIdx).SetInt(v)
	case "bool":
		v, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		val.Field(valIdx).SetBool(v)
	case "uuid.UUID":
		v, err := uuid.Parse(value)
		if err != nil {
			return err
		}
		val.Field(valIdx).Set(reflect.ValueOf(v))
	}

	return nil
}

type err interface {
	error() error
}

//Encode generate a encode function to encode response to json
func Encode() func(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, response interface{}) error {
		if e, ok := response.(err); ok && e.error() != nil {
			return e.error()
		}

		// set status code to 204 when the response is nil
		if response == nil {
			w.WriteHeader(http.StatusNoContent)
			json.NewEncoder(w).Encode("")
			return nil
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(response)
		return nil
	}
}

//ParseJSON parse request body (json) to model
func ParseJSON(ctx context.Context, request *http.Request, model interface{}) (interface{}, error) {
	err := json.NewDecoder(request.Body).Decode(model)

	if err != nil {
		return nil, err
	}

	return model, nil
}
