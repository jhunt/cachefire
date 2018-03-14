package main

import (
	"encoding/json"
	"os"
	"strings"
	"fmt"
	"net/http"
)

func reply(w http.ResponseWriter, code int, thing interface{}) {
	b, err := json.Marshal(thing)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to marshal json response: %s\n", err)
		w.WriteHeader(500)
		fmt.Fprintf(w, `{"error":"an internal json error has occurred"}`)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(b)
	w.Write([]byte{'\n'})
}

func oops(w http.ResponseWriter, code int, err string) {
	reply(w, code, struct {
		Error string `json:"error"`
	}{
		Error: err,
	})
}

type API struct {
	Username  string
	Password string
}

func (api API) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	user, pass, authn := req.BasicAuth()
	if !authn {
		oops(w, 401, "unauthorized.")
		return
	}
	if user != api.Username || pass != api.Password {
		oops(w, 403, "forbidden.")
		return
	}

	if req.URL.Path == "/v1/metrics" {
		if req.Method == "GET" {
			Lock.Lock()
			defer Lock.Unlock()

			reply(w, 200, Metrics)
			return
		}

		oops(w, 405, "method not allowed.")
		return
	}

	if strings.HasPrefix(req.URL.Path, "/v1/metrics/") {
		id := strings.TrimPrefix(req.URL.Path, "/v1/metrics/")

		if req.Method == "GET" {
			Lock.Lock()
			defer Lock.Unlock()

			m, ok := Metrics[id]
			if !ok {
				oops(w, 404, "metric not found.")
				return
			}
			reply(w, 200, m)
			return
		}

		oops(w, 405, "method not allowed.")
		return
	}

	oops(w, 404, "not found.")
}
