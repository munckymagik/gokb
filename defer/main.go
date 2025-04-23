package main

import (
	"errors"
	"fmt"
)

func op() (int, error) {
	return 0, errors.New("second err")
}

func finishDeferred(err *error) {
	log("deferred err ptr", *err, err)
}

func finish(err error) {
	log("deferred err normal", err, &err)
}

func log(label string, err error, perr *error) {
	fmt.Printf("%-20s memory address: %p, msg: %s\n", label, perr, err)
}

func main() {
	var err error = errors.New("first err")
	log("first err", err, &err)
	defer finishDeferred(&err)
	defer finish(err)

	result, err := op()
	_ = result

	log("second err", err, &err)
}
