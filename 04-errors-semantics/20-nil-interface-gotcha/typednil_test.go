package main

import (
	"errors"
	"testing"
)

func TestBrokenIsNotNil(t *testing.T) {
	err := DoThing()
	if err == nil {
		t.Fatal("expected err != nil due to typed-nil trap")
	}
	t.Logf("trap confirmed: %T %#v", err, err)
}

func TestFixedIsNil(t *testing.T) {
	err := DoThingFixed()
	if err != nil {
		t.Fatalf("expected nil, got %T %#v", err, err)
	}
}

func TestErrorsAsThroughWrapping(t *testing.T) {
	err := DoThingWrapped()
	var me *MyError
	if !errors.As(err, &me) {
		t.Fatal("errors.As should extract *MyError")
	}
	if me == nil {
		t.Fatal("extracted *MyError should not be nil")
	}
	if me.Op != "read" {
		t.Fatalf("unexpected Op: %q", me.Op)
	}
}
