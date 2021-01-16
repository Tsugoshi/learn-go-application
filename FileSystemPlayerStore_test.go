package poker_test

import (
	"io/ioutil"
	"os"
	"testing"

	poker "github.com/tsugoshi/learn-go-application"
)

func TestFileSystemStore(t *testing.T) {

	t.Run("league sorted", func(t *testing.T) {

		database, cleanDatabase := createTempFile(t,
			`[{"Name": "Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 33}]`)

		defer cleanDatabase()

		store, err := poker.NewFileSystemPlayerStore(database)
		assertNoError(t, err)

		want := []poker.Player{
			{"Chris", 33},
			{"Cleo", 10},
		}

		got := store.GetLeague()
		poker.AssertLeague(t, got, want)

		got = store.GetLeague()
		poker.AssertLeague(t, got, want)
	})

	t.Run("/get player score", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t,
			`[{"Name": "Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 33}]`)

		defer cleanDatabase()

		store, err := poker.NewFileSystemPlayerStore(database)
		assertNoError(t, err)

		got := store.GetPlayerScore("Chris")
		want := 33

		assertScoreEqual(t, got, want)
	})

	t.Run("store wins for existing players", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t,
			`[{"Name": "Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 33}]`)

		defer cleanDatabase()

		store, err := poker.NewFileSystemPlayerStore(database)
		assertNoError(t, err)

		store.RecordWin("Chris")

		got := store.GetPlayerScore("Chris")
		want := 34

		assertScoreEqual(t, got, want)
	})

	t.Run("store wins for new player", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t,
			`[{"Name": "Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 33}]`)

		defer cleanDatabase()

		store, err := poker.NewFileSystemPlayerStore(database)
		assertNoError(t, err)

		store.RecordWin("Samantha")

		got := store.GetPlayerScore("Samantha")
		want := 1
		assertScoreEqual(t, got, want)

	})

	t.Run("works with empty file", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, "")

		defer cleanDatabase()

		_, err := poker.NewFileSystemPlayerStore(database)
		assertNoError(t, err)
	})
}

func assertScoreEqual(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("Got %d, wanted %d", got, want)
	}
}

func createTempFile(t *testing.T, initialData string) (*os.File, func()) {
	t.Helper()
	tmpFile, err := ioutil.TempFile("", "db")
	if err != nil {
		t.Fatalf("error while creating temp db file. %v", err)
	}

	tmpFile.Write([]byte(initialData))

	removeFile := func() {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
	}
	return tmpFile, removeFile
}

func assertNoError(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("didn't expect error, but got one: %v", err)
	}
}
