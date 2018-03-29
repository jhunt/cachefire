package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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
	Username string
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

			switch req.URL.Query().Get("fmt") {
			case "telegraf":
				l := make([]Metric, 0)

				for job := range Metrics {
					for idx := range Metrics[job] {
						for name := range Metrics[job][idx] {
							m := Metrics[job][idx][name]

							l = append(l, Metric{
								Type:     m.Type,
								Job:      job,
								Index:    idx,
								Name:     name,
								LastSeen: m.LastSeen,

								Value: m.Value,
								Unit:  m.Unit,

								Tally: m.Tally,
							})
						}
					}
				}
				reply(w, 200, l)

			case "", "default":
				reply(w, 200, Metrics)

			default:
				oops(w, 400, "invalid format.")
			}
			return
		}

		oops(w, 405, "method not allowed.")
		return
	}

	oops(w, 404, "not found.")
}
