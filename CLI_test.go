package poker_test

import (
	"strings"
	"testing"

	poker "github.com/tsugoshi/learn-go-application"
)

func TestCLI(t *testing.T) {

	t.Run("Record chris win from user input", func(t *testing.T) {
		player := "Chris"
		in := strings.NewReader(player + " wins\n")

		playerStore := &poker.StubPlayerStore{}
		cli := poker.NewCLI(playerStore, in)
		cli.PlayPoker()

		poker.AssertPlayerWin(t, playerStore, player)
	})

	t.Run("Record Cleo win from user input", func(t *testing.T) {
		player := "Cleo"
		in := strings.NewReader(player + " wins\n")

		playerStore := &poker.StubPlayerStore{}
		cli := poker.NewCLI(playerStore, in)
		cli.PlayPoker()

		poker.AssertPlayerWin(t, playerStore, player)
	})

}
