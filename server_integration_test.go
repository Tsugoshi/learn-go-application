package poker_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	poker "github.com/tsugoshi/learn-go-application"
)

func TestRecordingWinsAndRetrievingThem(t *testing.T) {

	database, cleanDatabase := createTempFile(t, `[]`)
	defer cleanDatabase()

	store, err := poker.NewFileSystemPlayerStore(database)
	assertNoError(t, err)

	server := poker.MustMakePlayerServer(t, store, &poker.GameSpy{})
	player := "Pepper"

	server.ServeHTTP(httptest.NewRecorder(), poker.NewPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), poker.NewPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), poker.NewPostWinRequest(player))

	t.Run("get score", func(t *testing.T) {

		response := httptest.NewRecorder()
		server.ServeHTTP(response, poker.NewGetScoreRequest(player))
		poker.AssertStatus(t, response.Code, http.StatusOK)
		poker.AssertResponseBody(t, response.Body.String(), "3")
	})

	t.Run("get league", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, poker.NewGetLeagueRequest())

		poker.AssertStatus(t, response.Code, http.StatusOK)
		got := poker.GetLeagueFromResponse(t, response.Body)

		want := []poker.Player{
			{player, 3},
		}

		poker.AssertLeague(t, got, want)
	})
}
