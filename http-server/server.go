package httpserver

import (
	"fmt"
	"net/http"
	"strings"
)

// func PlayerServer(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprint(w, "20")
// }

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

	w.WriteHeader(http.StatusNotFound)

	fmt.Fprint(w, p.store.GetPlayerScore(player))
}
