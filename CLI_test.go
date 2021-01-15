package poker_test

import (
	"bytes"
	"io"
	"strings"
	"testing"

	poker "github.com/tsugoshi/learn-go-application"
)

type GameSpy struct {
	startedWith    int
	startCalled    bool
	finishedWith   string
	finishedCalled bool
}

func (g *GameSpy) Start(numberOfPlayers int) {
	g.startedWith = numberOfPlayers
	g.startCalled = true
}

func (g *GameSpy) Finish(winner string) {
	g.finishedWith = winner
	g.finishedCalled = true
}

func TestCLI(t *testing.T) {
	t.Run("it prompts the user to enter number of players", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		in := strings.NewReader("7\n")

		game := &GameSpy{}
		cli := poker.NewCLI(game, in, stdout)
		cli.PlayPoker()

		assertMessagesSentToUser(t, stdout, poker.PlayerGreeting)
		assertGameStartedWith(t, game, 7)
	})

	t.Run("start game with 8 players and record 'Cleo' as winner", func(t *testing.T) {
		game := &GameSpy{}
		dummyStdOut := &bytes.Buffer{}
		in := userSends("8", "Cleo wins")
		cli := poker.NewCLI(game, in, dummyStdOut)

		cli.PlayPoker()

		assertGameStartedWith(t, game, 8)
		assertFinishCalledWith(t, game, "Cleo")
	})

	t.Run("it prints an error when a non numeric value is entered and does not start the game", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		in := strings.NewReader("Pies\n")
		game := &GameSpy{}

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

func assertGameStartedWith(t *testing.T, game *GameSpy, numberOfPlayersWanted int) {
	t.Helper()
	if game.startedWith != numberOfPlayersWanted {
		t.Errorf("wanted Start called with %d but got %d", numberOfPlayersWanted, game.startedWith)
	}
}

func assertGameNotFinished(t *testing.T, game *GameSpy) {
	t.Helper()
	if game.finishedCalled {
		t.Errorf("game should not have finished")
	}
}

func assertGameNotStarted(t *testing.T, game *GameSpy) {
	t.Helper()
	if game.startCalled {
		t.Errorf("game should not have started")
	}
}

func assertFinishCalledWith(t *testing.T, game *GameSpy, winner string) {
	t.Helper()
	if game.finishedWith != winner {
		t.Errorf("expected finish called with %q but got %q", winner, game.finishedWith)
	}
}
