package log

import (
	"context"
	"fmt"
	"go-clean/src/lib/appcontext"
	"go-clean/src/lib/errors"
	"os"
	"sync"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var once = sync.Once{}

type Interface interface {
	Error(ctx context.Context, obj interface{})
	Info(ctx context.Context, obj interface{})
	Fatal(ctx context.Context, obj interface{})
	Debug(ctx context.Context, obj interface{})
}

type Config struct {
	Level string
}

type logger struct {
	log zerolog.Logger
}

func Init(cfg Config) Interface {
	var zerologging zerolog.Logger
	once.Do(func() {
		level, err := zerolog.ParseLevel(cfg.Level)
		if err != nil {
			log.Fatal().Msg(fmt.Sprintf("failed to parse level : %v", err))
		}

		zerologging = zerolog.New(os.Stdout).
			With().
			Timestamp().
			CallerWithSkipFrameCount(3).
			Logger().
			Level(level)
	})

	return &logger{log: zerologging}
}

func (l *logger) Info(ctx context.Context, obj interface{}) {
	l.log.Info().Fields(getContextFields(ctx)).Msg(fmt.Sprint(getCaller(obj)))
}

func (l *logger) Debug(ctx context.Context, obj interface{}) {
	l.log.Debug().Fields(getContextFields(ctx)).Msg(fmt.Sprint(getCaller(obj)))
}

func (l *logger) Error(ctx context.Context, obj interface{}) {
	l.log.Error().Fields(getContextFields(ctx)).Msg(fmt.Sprint(getCaller(obj)))
}

func (l *logger) Fatal(ctx context.Context, obj interface{}) {
	l.log.Fatal().Fields(getContextFields(ctx)).Msg(fmt.Sprint(getCaller(obj)))
}

func getCaller(obj interface{}) interface{} {
	switch tr := obj.(type) {
	case error:
		filename, line, msg, err := errors.GetCaller(tr)
		if err == nil {
			obj = fmt.Sprintf("%s:%#v --- %s", filename, line, msg)
		}
	case string:
		obj = tr
	default:
		obj = fmt.Sprintf("%#v", tr)
	}

	return obj
}

func getContextFields(ctx context.Context) map[string]interface{} {
	cf := map[string]interface{}{}

	cf["request_id"] = appcontext.GetRequestID(ctx)

	return cf
}
