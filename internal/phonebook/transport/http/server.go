package http

import (
	"net/http"

	"phonebook/internal/phonebook"
	"phonebook/internal/phonebook/endpoint"
	"phonebook/internal/phonebook/model"
	"phonebook/pkg/server"

	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
)

func Create(svc *phonebook.Service, logger log.Logger, opts []kithttp.ServerOption) http.Handler {
	end := endpoint.Add(svc)
	var serverLogger *server.Logger
	if nil != logger {
		serverLogger = &server.Logger{
			Logger:    logger,
			Namespace: "phonebook",
			Subsystem: "add_new_profile",
			Action:    "POST",
		}
	}
	return server.NewHTTPServer(end, server.HTTPOption{
		DecodeModel: &model.PhoneBook{},
		Logger:      serverLogger,
	}, opts...)
}

func Get(svc *phonebook.Service, logger log.Logger, opts []kithttp.ServerOption) http.Handler {
	end := endpoint.Add(svc)
	var serverLogger *server.Logger
	if nil != logger {
		serverLogger = &server.Logger{
			Logger:    logger,
			Namespace: "phonebook",
			Subsystem: "get_list_phone",
			Action:    "GET",
		}
	}
	return server.NewHTTPServer(end, server.HTTPOption{
		DecodeModel: &model.GetPhoneList{},
		Logger:      serverLogger,
	}, opts...)
}

func Update(svc *phonebook.Service, logger log.Logger, opts []kithttp.ServerOption) http.Handler {
	end := endpoint.Add(svc)
	var serverLogger *server.Logger
	if nil != logger {
		serverLogger = &server.Logger{
			Logger:    logger,
			Namespace: "phonebook",
			Subsystem: "update_profile",
			Action:    "PUT",
		}
	}
	return server.NewHTTPServer(end, server.HTTPOption{
		DecodeModel: &model.PhoneBook{},
		Logger:      serverLogger,
	}, opts...)
}

func Remove(svc *phonebook.Service, logger log.Logger, opts []kithttp.ServerOption) http.Handler {
	end := endpoint.Add(svc)
	var serverLogger *server.Logger
	if nil != logger {
		serverLogger = &server.Logger{
			Logger:    logger,
			Namespace: "phonebook",
			Subsystem: "remove_profile",
			Action:    "PUT",
		}
	}
	return server.NewHTTPServer(end, server.HTTPOption{
		DecodeModel: &model.PhoneBook{},
		Logger:      serverLogger,
	}, opts...)
}
