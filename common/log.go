package common

import (
	"fmt"
	"io"
	slog "log"
	"time"
)

var Log = log{}

type log struct{}

func (this *log) Stdout(format string, v ...interface{}) {
	fmt.Printf(format+"\n", v...)
}

func (this *log) Debug(format string, v ...interface{}) {
	var t = make([]interface{}, 0)
	t = append(t, this.now())
	t = append(t, v...)
	fmt.Printf("%s DEBUG "+format+"\n", t...)
}

func (this *log) Info(format string, v ...interface{}) {
	var t = make([]interface{}, 0)
	t = append(t, this.now())
	t = append(t, v...)
	fmt.Printf("%s INFO "+format+"\n", t...)
}

func (this *log) Warning(format string, v ...interface{}) {
	var t = make([]interface{}, 0)
	t = append(t, this.now())
	t = append(t, v...)
	fmt.Printf("%s WARNING "+format+"\n", t...)
}

func (this *log) Error(format string, v ...interface{}) {
	var t = make([]interface{}, 0)
	t = append(t, this.now())
	t = append(t, v...)
	fmt.Printf("%s ERROR "+format+"\n", t...)
}

func (l *log) now() string {
	return time.Now().Format("2006/01/02 15:04:05")
}

type dummyWriter struct {
	io.Writer
}

func (this dummyWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}
func NewNilLogger() *slog.Logger {
	return slog.New(dummyWriter{}, "", slog.LstdFlags)
}
