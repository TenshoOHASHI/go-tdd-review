package httpserver_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// 　必要な呼び出し用の関数を用意
type PlayerStore interface {
	GetPlayerScore(name string) int
}

// サーバはインターフェースに依存
type PlayerServer struct {
	store PlayerStore
}

// 実際にサーバから関数を呼び出す
func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	player := strings.TrimPrefix(r.URL.Path, "/players/")

	score := p.store.GetPlayerScore(player)

	if score == 0 {
		w.WriteHeader(http.StatusNotFound)
	}

	fmt.Fprint(w, p.store.GetPlayerScore(player))

}

type Handler interface {
	ServeHttp(w http.ResponseWriter, r http.Request)
}

func ListenAndServe(addr string, handler Handler) error {
	return nil
}

// func PlayerServer(w http.ResponseWriter, r *http.Request) {

// 	// /players/{name}
// 	player := strings.TrimPrefix(r.URL.Path, "/players/")
// 	fmt.Fprint(w, "20")

// 	// 実際のURLを確認して処理
// 	fmt.Fprint(w, GetPlayerScore(player))

// }

func GetPlayerScore(name string) int {
	if name == "Pepper" {
		return 20
	}
	if name == "Floyd" {
		return 10
	}

	return 0
}

func newGetScoreRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/player/%s", name), nil)
	return req
}

func assertResponseBody(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response body is wrong, got %q want %q", got, want)
	}
}

// func TestGETPlayers(t *testing.T) {
// 	// ここでハンドラーを初期化
// 	server := &PlayerServer{}
// 	t.Run("returns Pepper's score", func(t *testing.T) {
// 		request := newGetScoreRequest("Pepper")
// 		response := httptest.NewRecorder()

// 		server.ServeHTTP(response, request)

// 		assertResponseBody(t, response.Body.String(), "20")
// 	})

// 	t.Run("returns Floyd's score", func(t *testing.T) {
// 		request := newGetScoreRequest("Floyd")
// 		response := httptest.NewRecorder()

// 		server.ServeHTTP(response, request)

// 		assertResponseBody(t, response.Body.String(), "10")
// 	})

// 	t.Run("returns 404 on missing players", func(t *testing.T) {

// 		// ここでデータを受け取って、バリデーションやリクエスト構造体に変換する
// 		request := newGetScoreRequest("Apollo")
// 		response := httptest.NewRecorder()
// 		// Browser → Request → Handler → Fprint → ResponseWriter(Body/buf) → Browser
// 		// サーバに送信、Responseの中にあるBodyフィールドのbuf型にデータ保存する
// 		// BodyはWriteメソッドを満たすため、Fprint内部で、渡ってきたResponseRecorderの構造体にあるBodyフィールドにw.Writeメソッドを通して、データをbufに保存してる
// 		server.ServeHTTP(response, request)

// 		assertStatus(t, response.Code, http.StatusNotFound)
// 	})
// }

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}

// server_test.go
type StubPlayerStore struct {
	scores map[string]int
}

func (s *StubPlayerStore) GetPlayerScore(name string) int {
	score := s.scores[name]
	return score
}

// server_test.go
func TestGETPlayers(t *testing.T) {
	store := StubPlayerStore{
		map[string]int{
			"Pepper": 20,
			"Floyd":  10,
		},
	}
	server := &PlayerServer{&store}

	t.Run("returns Pepper's score", func(t *testing.T) {
		request := newGetScoreRequest("Pepper")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResponseBody(t, response.Body.String(), "20")
	})

	t.Run("returns Floyd's score", func(t *testing.T) {
		request := newGetScoreRequest("Floyd")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResponseBody(t, response.Body.String(), "10")
	})
}

// server_test.go
func TestStoreWins(t *testing.T) {
	store := StubPlayerStore{
		map[string]int{},
	}
	server := &PlayerServer{&store}

	t.Run("it returns accepted on POST", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/players/Pepper", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusAccepted)
	})
}
