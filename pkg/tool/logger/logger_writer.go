package logger

import "io"

type Writer struct {
	logger Logger
}

func NewWriter(logger Logger) io.Writer {
	return Writer{
		logger: logger,
	}
}

func (w Writer) Write(p []byte) (n int, err error) {
	w.logger.Info(string(p))
	return 0, nil
}
