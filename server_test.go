package poker_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	poker "github.com/tsugoshi/learn-go-application"
)

func TestGETPlayers(t *testing.T) {

	store := poker.StubPlayerStore{
		map[string]int{
			"Pepper": 20,
			"Floyd":  10,
		},
		nil,
		nil,
	}
	server := poker.MustMakePlayerServer(t, &store, dummyGame)

	t.Run("returns Pepper's score", func(t *testing.T) {
		request := poker.NewGetScoreRequest("Pepper")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := "20"

		poker.AssertStatus(t, response.Code, http.StatusOK)
		poker.AssertResponseBody(t, got, want)
	})

	t.Run("returns Floyd's score", func(t *testing.T) {
		request := poker.NewGetScoreRequest("Floyd")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := "10"

		poker.AssertStatus(t, response.Code, http.StatusOK)
		poker.AssertResponseBody(t, got, want)
	})

	t.Run("returns 404 on missing players", func(t *testing.T) {
		request := poker.NewGetScoreRequest("Apollo")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Code
		want := http.StatusNotFound

		poker.AssertStatus(t, got, want)
	})
}

func TestStoreWins(t *testing.T) {
	store := poker.StubPlayerStore{
		map[string]int{},
		nil,
		nil,
	}
	server := poker.MustMakePlayerServer(t, &store, dummyGame)

	t.Run("it records wins on POST", func(t *testing.T) {
		player := "Pepper"
		request := poker.NewPostWinRequest(player)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		poker.AssertStatus(t, response.Code, http.StatusAccepted)

		poker.AssertPlayerWin(t, &store, player)
	})
}

var dummyGame = &poker.GameSpy{}

func TestLeague(t *testing.T) {

	t.Run("it returns the league table as JSON", func(t *testing.T) {

		wantedLeague := []poker.Player{
			{"Chris", 1},
			{"Cornel", 2},
			{"DiCaprio", 30},
		}

		store := poker.StubPlayerStore{nil, nil, wantedLeague}
		server := poker.MustMakePlayerServer(t, &store, dummyGame)

		request := poker.NewGetLeagueRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := poker.GetLeagueFromResponse(t, response.Body)

		poker.AssertContentType(t, response, poker.JsonContentType)
		poker.AssertStatus(t, response.Code, http.StatusOK)
		poker.AssertLeague(t, got, wantedLeague)
	})
}

func TestGame(t *testing.T) {
	t.Run("GET /game returs 200", func(t *testing.T) {
		server := poker.MustMakePlayerServer(t, &poker.StubPlayerStore{}, dummyGame)
		request := poker.NewGetGameRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		poker.AssertStatus(t, response.Code, http.StatusOK)
	})

	t.Run("start game with 3 players and declare Ruth as winner", func(t *testing.T) {
		wantedBlindAlerter := "Blind is 100"
		winner := "Ruth"
		tenMS := 10 * time.Millisecond
		game := &poker.GameSpy{BlindAlerter: []byte(wantedBlindAlerter)}

		server := httptest.NewServer(poker.MustMakePlayerServer(t, &poker.StubPlayerStore{}, game))
		ws := poker.MustDialWS(t, "ws"+strings.TrimPrefix(server.URL, "http")+"/ws")
		defer server.Close()
		defer ws.Close()

		poker.WriteWSMessage(t, ws, "3")
		poker.WriteWSMessage(t, ws, winner)

		poker.AssertGameStartedWith(t, game, 3)
		poker.AssertFinishCalledWith(t, game, winner)

		within(t, tenMS, func() { assertWebSocketGotMessage(t, ws, wantedBlindAlerter) })

	})

}

func assertWebSocketGotMessage(t *testing.T, ws *websocket.Conn, wantedMessage string) {
	_, message, _ := ws.ReadMessage()
	if string(message) != wantedMessage {
		t.Errorf("got %s, wanted %s", string(message), wantedMessage)
	}

}

func within(t testing.TB, d time.Duration, assert func()) {
	t.Helper()

	done := make(chan struct{}, 1)
	go func() {
		assert()
		done <- struct{}{}
	}()

	select {
	case <-time.After(d):
		t.Error("Timed out")
	case <-done:
	}
}
