package main

import (
	"testing"

	"bytes"

	"time"

	. "github.com/bborbe/assert"
)

func TestDoFail(t *testing.T) {
	writer := bytes.NewBufferString("")
	err := do(writer, func(directory string, prefix string, keepAmount int) error {
		return nil
	}, "", "", -1, time.Minute, false, "/tmp/lock")
	if err = AssertThat(err, NotNilValue()); err != nil {
		t.Fatal(err)
	}
}

func TestDoSuccess(t *testing.T) {
	writer := bytes.NewBufferString("")
	err := do(writer, func(directory string, prefix string, keepAmount int) error {
		return nil
	}, "/tmp", "backup", 5, time.Minute, true, "/tmp/lock")
	if err = AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
}
