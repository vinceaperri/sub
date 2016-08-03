package main

import (
	"fmt"
	"sync"
)

var bold_red = "\x1B[1;31m"
var bold_purple = "\x1B[1;35m"
var bold_blue = "\x1B[1;34m"
var reset_color = "\x1B[0m"

type logger struct {
	mux *sync.Mutex
}

func new_logger() *logger {
	return &logger{&sync.Mutex{}}
}

func (l *logger) ok(format string, a ...interface{}) (int, error) {
	l.mux.Lock()
	defer l.mux.Unlock()
	return fmt.Printf(bold_blue + format + reset_color, a...)
}

func (l *logger) good(format string, a ...interface{}) (int, error) {
	l.mux.Lock()
	defer l.mux.Unlock()
	return fmt.Printf(bold_purple + format + reset_color, a...)
}

func (l *logger) bad(format string, a ...interface{}) (int, error) {
	l.mux.Lock()
	defer l.mux.Unlock()
	return fmt.Printf(bold_red + format + reset_color, a...)
}
