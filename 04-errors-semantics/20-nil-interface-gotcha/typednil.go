package main

import "fmt"

type MyError struct {
	Op string
}

func (e *MyError) Error() string {
	return e.Op
}


func DoThing() error {
	var err *MyError
	return err
}


func DoThingFixed() error {
	var err *MyError
	if err == nil {
		return nil
	}
	return err
}


func DoThingWrapped() error {
	return fmt.Errorf("context: %w", &MyError{Op: "read"})
}
