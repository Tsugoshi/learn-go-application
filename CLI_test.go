package poker_test

import (
	"bytes"
	"io"
	"strings"
	"testing"

	poker "github.com/tsugoshi/learn-go-application"
)

func TestCLI(t *testing.T) {
	t.Run("it prompts the user to enter number of players", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		in := strings.NewReader("7\n")

		game := &poker.GameSpy{}
		cli := poker.NewCLI(game, in, stdout)
		cli.PlayPoker()

		assertMessagesSentToUser(t, stdout, poker.PlayerGreeting)
		poker.AssertGameStartedWith(t, game, 7)
	})

	t.Run("start game with 8 players and record 'Cleo' as winner", func(t *testing.T) {
		game := &poker.GameSpy{}
		dummyStdOut := &bytes.Buffer{}
		in := userSends("8", "Cleo wins")
		cli := poker.NewCLI(game, in, dummyStdOut)

		cli.PlayPoker()

		poker.AssertGameStartedWith(t, game, 8)
		poker.AssertFinishCalledWith(t, game, "Cleo")
	})

	t.Run("it prints an error when a non numeric value is entered and does not start the game", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		in := strings.NewReader("Pies\n")
		game := &poker.GameSpy{}

		cli := poker.NewCLI(game, in, stdout)
		cli.PlayPoker()

		assertGameNotStarted(t, game)
		assertMessagesSentToUser(t, stdout, poker.PlayerGreeting, poker.BadStartInput)
	})
}

func assertMessagesSentToUser(t testing.TB, stdout *bytes.Buffer, messages ...string) {
	t.Helper()
	want := strings.Join(messages, "")
	got := stdout.String()
	if got != want {
		t.Errorf("got %q sent to stdout but expected %+v", got, messages)
	}
}

func userSends(messages ...string) io.Reader {
	return strings.NewReader(strings.Join(messages, "\n"))
}

func assertGameNotFinished(t *testing.T, game *poker.GameSpy) {
	t.Helper()
	if game.FinishedCalled {
		t.Errorf("game should not have finished")
	}
}

func assertGameNotStarted(t *testing.T, game *poker.GameSpy) {
	t.Helper()
	if game.StartCalled {
		t.Errorf("game should not have started")
	}
}
