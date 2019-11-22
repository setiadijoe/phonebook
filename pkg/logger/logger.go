package logger

import (
	"context"
	"encoding/json"
	"regexp"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"

	"bytes"
	"io"
	"runtime"
	"strconv"

	"github.com/go-kit/kit/log/level"
)

//Request ...
type Request struct {
	Ctx    context.Context
	Method string
	Action string
	Origin string
	Params interface{}
}

//Callback ...
type Callback func(request Request) (interface{}, error)

//Logger ...
type Logger struct {
	callUpdate     chan interface{}
	callError      chan error
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	logger         log.Logger
}

type KitxLogger struct {
	caller log.Valuer
	stderr log.Logger
	stdout log.Logger
	level  level.Value
	filter level.Option
}

func (logger KitxLogger) SetCaller(caller log.Valuer) KitxLogger {
	logger.caller = caller
	return logger
}

func NewKitxLogger(errorWriter io.Writer, outWriter io.Writer) KitxLogger {
	logger := KitxLogger{
		stderr: log.NewLogfmtLogger(errorWriter),
		stdout: log.NewLogfmtLogger(outWriter),
		level:  level.InfoValue(),
		caller: CallerRegex("canfazz"),
	}
	logger.stderr = log.With(logger.stderr, "ts", log.DefaultTimestampUTC)
	logger.stdout = log.With(logger.stdout, "ts", log.DefaultTimestampUTC)
	return logger
}

func Caller(depth int, path int) log.Valuer {
	return func() interface{} {
		_, file, line, _ := runtime.Caller(depth)
		idx := bytes.LastIndexByte([]byte(file), byte('/'))
		for i := 1; i < path; i++ {
			if idx == -1 {
				break
			}
			idx = bytes.LastIndexByte([]byte(file[:idx]), byte('/'))
		}
		result := file[idx+1:] + ":" + strconv.Itoa(line)
		if idx != -1 {
			result = ".../" + result
		}
		return result
	}
}

func CallerRegex(regex string) log.Valuer {
	re := regexp.MustCompile(regex)
	return func() interface{} {
		skip := 1

		_, file, line, ok := runtime.Caller(skip)
		for ok {
			if re.MatchString(file) {
				locs := re.FindStringIndex(file)
				file = file[locs[0]:]
				break
			}

			skip = skip + 1
			_, file, line, ok = runtime.Caller(skip)
		}

		return file + ":" + strconv.Itoa(line)
	}
}

// New create gokit layer Logger
func New(counter metrics.Counter, latency metrics.Histogram, logger log.Logger) Logger {
	return Logger{
		callUpdate:     make(chan interface{}),
		callError:      make(chan error),
		requestCount:   counter,
		requestLatency: latency,
		logger:         logger,
	}
}

//Instrumentation ...
func (m Logger) Instrumentation(
	f func(ctx context.Context, request interface{}) (interface{}, error),
	keyvals ...interface{},

) func(ctx context.Context, request interface{}) (interface{}, error) {
	return func(ctx context.Context, request interface{}) (resp interface{}, err error) {
		defer func(begin time.Time) {
			labelValues := make([]string, len(keyvals))
			for i := 0; i < len(keyvals); i++ {
				labelValues[i] = keyvals[i].(string)
			}

			if err != nil {
				labelValues = append(labelValues, "status", "failed")
			} else {
				labelValues = append(labelValues, "status", "success")
			}

			m.requestCount.With(labelValues...).Add(1)
			m.requestLatency.With(labelValues...).Observe(time.Since(begin).Seconds())
		}(time.Now())
		return f(ctx, request)
	}
}

//Log ...
func (m Logger) Log(
	f func(ctx context.Context, request interface{}) (interface{}, error),
	keyvals ...interface{},
) func(ctx context.Context, request interface{}) (interface{}, error) {
	return func(ctx context.Context, request interface{}) (resp interface{}, err error) {
		defer func(begin time.Time) {
			kv := make([]interface{}, len(keyvals))
			for i := 0; i < len(keyvals); i++ {
				kv[i] = keyvals[i]
			}

			jsonString, _ := json.Marshal(request)
			kv = append(kv,
				"params", string(jsonString),
				"took", time.Since(begin).String(),
			)

			if nil != err {
				kv = append(kv, "err", err.Error())
			}
			_ = m.logger.Log(kv...)
		}(time.Now())
		return f(ctx, request)
	}
}

func (logger KitxLogger) Log(keyvals ...interface{}) error {
	if len(keyvals) == 0 {
		return nil
	}
	if len(keyvals)%2 == 1 {
		keyvals = append(keyvals, log.ErrMissingValue)
	}

	keyvals = append(keyvals, "caller", logger.caller())

	var currentLogger log.Logger
	err := logger.getValByKey("err", keyvals...)
	if err != nil {
		currentLogger = level.NewInjector(logger.stderr, logger.level)
		keyvals = append(keyvals, "level", "error")
	} else {
		currentLogger = level.NewInjector(logger.stdout, logger.level)
	}

	return currentLogger.Log(keyvals...)
}

func (logger KitxLogger) getValByKey(searchKey string, keyvals ...interface{}) interface{} {
	for i := 0; i < len(keyvals); i += 2 {
		key, ok := keyvals[i].(string)
		if !ok {
			continue
		}

		if key != searchKey {
			continue
		}

		return keyvals[i+1]
	}

	return nil
}
