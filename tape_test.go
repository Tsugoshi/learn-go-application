package poker_test

import (
	"io/ioutil"
	"testing"

	poker "github.com/tsugoshi/learn-go-application"
)

func TestTape_Write(t *testing.T) {
	t.Helper()
	file, clean := createTempFile(t, "12345")
	defer clean()

	tape := &poker.Tape{file}
	tape.Write([]byte("abc"))

	file.Seek(0, 0)

	newFileContents, _ := ioutil.ReadAll(file)
	got := string(newFileContents)
	want := "abc"
	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}
