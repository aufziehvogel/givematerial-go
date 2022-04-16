package givematlib

import (
	"bytes"
	"testing"
)

func TestReadStatus(t *testing.T) {
	var buffer bytes.Buffer
	buffer.WriteString("hello\nthis\nis\nsimple")

	status, err := readStatus(&buffer)
	if err != nil {
		t.Error("Failed to read status data")
	}

	want := []string{"hello", "this", "is", "simple"}
	if len(status.known) != len(want) {
		t.Fatalf(
			"Status storage returned wrong number of items, got %d, want %d",
			len(status.known),
			len(want),
		)
	}
	if status.known[0] != "hello" {
		t.Errorf(
			"Status storage was incorrect, got %v, want %v",
			status.known,
			want,
		)
	}
}
