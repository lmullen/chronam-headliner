package headliner

import (
	_ "embed"
	"fmt"
	"net/http"
)

//go:embed index.html
var indexHTML string

func (a *App) RootHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, indexHTML)
	}
}
