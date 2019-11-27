package middleware

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"phonebook/pkg/httperror"
	"phonebook/pkg/logger"
	"phonebook/pkg/validator"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
)

var loggers = make(map[string]*logger.Logger)
var lock sync.Mutex

//LogAndInstrumentation wrap endpoint function with logger.Log and logger.Instrumentation
//This middleware will measure request_count and request_latency_microseconds
func LogAndInstrumentation(kitLogger log.Logger, namespace, subsystem, action, domain string) endpoint.Middleware {
	var logObj logger.Logger

	key := fmt.Sprintf("%s_%s", namespace, subsystem)

	if val, ok := loggers[key]; ok {
		logObj = *val
	} else {
		lock.Lock()
		logObj = logger.New(nil, nil,
			kitLogger,
		)

		loggers[key] = &logObj
		lock.Unlock()
	}

	return func(f endpoint.Endpoint) endpoint.Endpoint {
		keyvals := make([]interface{}, 0)
		keyvals = append(keyvals,
			"function", action,
			"domain", domain,
		)
		return logObj.Instrumentation(logObj.Log(f, keyvals...), keyvals...)
	}
}

func Nop() endpoint.Middleware {
	return func(f endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			return f(ctx, request)
		}
	}
}

//Validator wrap endpoint function to execute validator.v9
func Validator() endpoint.Middleware {
	return func(f endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			err = validator.DefaultValidator()(request)
			if err != nil {
				return nil, &httperror.ErrorWithStatusCode{
					Err:        err.Error(),
					StatusCode: http.StatusUnprocessableEntity,
				}
			}
			return f(ctx, request)
		}
	}
}
