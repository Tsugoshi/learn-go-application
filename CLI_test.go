package poker_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	poker "github.com/tsugoshi/learn-go-application"
)

type scheduledAlert struct {
	at     time.Duration
	amount int
}

type SpyBlindAlerter struct {
	alerts []scheduledAlert
}

func (s scheduledAlert) String() string {
	return fmt.Sprintf("%d amount chips at %v", s.amount, s.at)
}

func (s *SpyBlindAlerter) ScheduleAlertAt(duration time.Duration, amount int) {
	s.alerts = append(s.alerts, scheduledAlert{duration, amount})
}

func TestCLI(t *testing.T) {

	dummyBlindAlerter := &SpyBlindAlerter{}
	t.Run("Record chris win from user input", func(t *testing.T) {
		player := "Chris"
		in := strings.NewReader(player + " wins\n")

		playerStore := &poker.StubPlayerStore{}
		cli := poker.NewCLI(playerStore, in, dummyBlindAlerter)
		cli.PlayPoker()

		poker.AssertPlayerWin(t, playerStore, player)
	})

	t.Run("Record Cleo win from user input", func(t *testing.T) {
		player := "Cleo"
		in := strings.NewReader(player + " wins\n")

		playerStore := &poker.StubPlayerStore{}
		cli := poker.NewCLI(playerStore, in, dummyBlindAlerter)
		cli.PlayPoker()

		poker.AssertPlayerWin(t, playerStore, player)
	})

	t.Run("It shedules printing of blind value", func(t *testing.T) {
		player := "Chris"
		in := strings.NewReader(player + " wins\n")
		playerStore := &poker.StubPlayerStore{}
		blindAlerter := &SpyBlindAlerter{}

		cli := poker.NewCLI(playerStore, in, blindAlerter)
		cli.PlayPoker()

		cases := []scheduledAlert{
			{0 * time.Second, 100},
			{10 * time.Minute, 200},
			{20 * time.Minute, 300},
			{30 * time.Minute, 400},
			{40 * time.Minute, 500},
			{50 * time.Minute, 600},
			{60 * time.Minute, 800},
			{70 * time.Minute, 1000},
			{80 * time.Minute, 2000},
			{90 * time.Minute, 4000},
			{100 * time.Minute, 8000},
		}

		for i, want := range cases {
			t.Run(fmt.Sprint(want), func(t *testing.T) {

				if len(blindAlerter.alerts) <= i {
					t.Fatalf("alert %d was not scheduled %v", i, blindAlerter.alerts)
				}

				got := blindAlerter.alerts[i]

				assertScheduledAlert(t, got, want)
			})
		}

	})
}

func assertScheduledAlert(t *testing.T, got, want scheduledAlert) {

	if got.amount != want.amount {
		t.Errorf("Expected %d amount , got %d", want.amount, got.amount)
	}

	if got.at != want.at {
		t.Errorf("expected scheduled at %v, got %v", want.at, got.at)
	}
}
