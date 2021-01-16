package poker

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

type StubPlayerStore struct {
	Scores   map[string]int
	WinCalls []string
	League   League
}

func (s *StubPlayerStore) GetPlayerScore(name string) int {
	score := s.Scores[name]
	return score
}

func (s *StubPlayerStore) RecordWin(name string) {
	s.WinCalls = append(s.WinCalls, name)
}

func (s *StubPlayerStore) GetLeague() League {
	return s.League
}

type GameSpy struct {
	StartedWith    int
	StartCalled    bool
	BlindAlerter   []byte
	FinishedWith   string
	FinishedCalled bool
}

func (g *GameSpy) Start(numberOfPlayers int, out io.Writer) {
	g.StartedWith = numberOfPlayers
	g.StartCalled = true
	out.Write(g.BlindAlerter)
}

func (g *GameSpy) Finish(winner string) {
	g.FinishedWith = winner
	g.FinishedCalled = true
}

type ScheduledAlert struct {
	At     time.Duration
	Amount int
}

type SpyBlindAlerter struct {
	alerts []ScheduledAlert
}

func (s ScheduledAlert) String() string {
	return fmt.Sprintf("%d amount chips at %v", s.Amount, s.At)
}

func (s *SpyBlindAlerter) ScheduleAlertAt(duration time.Duration, amount int, to io.Writer) {
	s.alerts = append(s.alerts, ScheduledAlert{duration, amount})
}

func AssertLeague(t *testing.T, got, want []Player) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, wanted %v", got, want)
	}
}

func AssertStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}

const JsonContentType = "application/json"

func AssertContentType(t *testing.T, response *httptest.ResponseRecorder, want string) {
	t.Helper()
	if response.Result().Header.Get("content-type") != want {
		t.Errorf("response did not have content-type of %s, got %v", want, response.Result().Header)
	}
}

func NewGetScoreRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/players/%s", name), nil)
	return req
}

func NewPostWinRequest(name string) *http.Request {
	request, _ := http.NewRequest(http.MethodPost, "/players/Pepper", nil)
	return request
}
func NewGetLeagueRequest() *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/league", nil)
	return request
}

func NewGetGameRequest() *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/game", nil)
	return request
}

func GetLeagueFromResponse(t *testing.T, body io.Reader) (league []Player) {
	t.Helper()

	league, _ = NewLeague(body)
	return
}

func AssertResponseBody(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response body is wrong, got %q want %q", got, want)
	}
}

func AssertPlayerWin(t *testing.T, store *StubPlayerStore, winner string) {
	t.Helper()

	if len(store.WinCalls) < 1 {
		t.Fatal("expected a win call but didn't get any")
	}

	got := store.WinCalls[0]
	if got != winner {
		t.Errorf("didn't record correct winner, got %q, want %q", got, winner)
	}
}

func CheckSchedulingCases(cases []ScheduledAlert, t *testing.T, blindAlerter *SpyBlindAlerter) {
	for i, want := range cases {
		t.Run(fmt.Sprint(want), func(t *testing.T) {

			if len(blindAlerter.alerts) <= i {
				t.Fatalf("alert %d was not scheduled %v", i, blindAlerter.alerts)
			}

			got := blindAlerter.alerts[i]

			AssertScheduledAlert(t, got, want)
		})
	}
}

func AssertScheduledAlert(t *testing.T, got, want ScheduledAlert) {

	if got.Amount != want.Amount {
		t.Errorf("Expected %d amount , got %d", want.Amount, got.Amount)
	}

	if got.At != want.At {
		t.Errorf("expected scheduled at %v, got %v", want.At, got.At)
	}
}

func MustMakePlayerServer(t *testing.T, store PlayerStore, game Game) *PlayerServer {
	server, err := NewPlayerServer(store, game)
	if err != nil {
		t.Fatal("problem creating playerServer", err)
	}
	return server
}

func MustDialWS(t *testing.T, url string) *websocket.Conn {
	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatalf("error dialing ws on url %s, %v", url, err)
	}
	return ws
}

func WriteWSMessage(t *testing.T, conn *websocket.Conn, message string) {
	t.Helper()
	err := conn.WriteMessage(websocket.TextMessage, []byte(message))

	if err != nil {
		t.Fatalf("could not send message over ws connection, %v", err)
	}
}

func AssertFinishCalledWith(t testing.TB, game *GameSpy, winner string) {
	t.Helper()
	passed := retryUntil(500*time.Millisecond, func() bool {
		return game.FinishedWith == winner
	})

	if !passed {
		t.Errorf("expected finish called with %q, but got %q", winner, game.FinishedWith)
	}
}

func AssertGameStartedWith(t *testing.T, game *GameSpy, numberOfPlayersWanted int) {
	t.Helper()

	passed := retryUntil(500*time.Millisecond, func() bool {
		return game.StartedWith == numberOfPlayersWanted
	})

	if !passed {
		t.Errorf("wanted Start called with %d but got %d", numberOfPlayersWanted, game.StartedWith)
	}
}

func retryUntil(d time.Duration, f func() bool) bool {
	deadline := time.Now().Add(d)
	for time.Now().Before(deadline) {
		if f() {
			return true
		}
	}
	return false
}
