package logger

import (
	"context"
	"fmt"
	"log"
	"strings"
)

type (
	Logger interface {
		WithContext(ctx context.Context) Logger
		WithLocation() Logger
		WithPreffix(prefix string) Logger
		Error(v interface{})
		Fatal(v interface{})
		Info(v interface{})
	}
	logger struct {
		ctx          context.Context
		preffix      string
		projectDir   string
		withLocation bool
	}
)

func New(projectDir string) Logger {
	return logger{
		projectDir: projectDir,
	}
}

func (l logger) WithContext(ctx context.Context) Logger {
	l.ctx = ctx
	return l
}

func (l logger) WithPreffix(preffix string) Logger {
	l.preffix = fmt.Sprintf("[%s]", preffix)
	return l
}

func (l logger) WithLocation() Logger {
	l.withLocation = true
	return l
}

func (l logger) Error(v interface{}) {
	log.Printf("[error]%s: %v", l.getAdditionalData(), v)
}

func (l logger) Fatal(v interface{}) {
	log.Fatalf("[fatal]%s: %v", l.getAdditionalData(), v)
}

func (l logger) Info(v interface{}) {
	log.Printf("[info]%s: %v", l.getAdditionalData(), v)
}

func (l logger) getAdditionalData() string {
	data := make([]string, 0)
	if l.preffix != "" {
		data = append(data, l.preffix)
	}
	if l.withLocation {
		fn, file, line := getLocation()
		file = strings.Replace(file, l.projectDir, "", 1)
		data = append(data, fmt.Sprintf("[%s:%d][%s]", file, line, fn))
	}
	return strings.Join(data, "")
}
