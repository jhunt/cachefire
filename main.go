package main

import (
	fmt "github.com/jhunt/go-ansi"
	"net/http"
	"os"

	"github.com/jhunt/go-cli"
	env "github.com/jhunt/go-envirotron"
	"github.com/jhunt/go-firehose"
)

func main() {
	var opts struct {
		Config   string `cli:"-c, --config"`
		Username string `cli:"-u, --username" env:"CACHE_FIRE_USERNAME"`
		Password string `cli:"-p, --password" env:"CACHE_FIRE_PASSWORD"`
		Port     string `cli:"--port"         env:"PORT"`
	}
	opts.Username = "cachefire"
	opts.Port = "3000"
	env.Override(&opts)
	_, _, err := cli.Parse(&opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "@R{!!! %s}\n", err)
		os.Exit(1)
	}

	if opts.Password == "" {
		fmt.Fprintf(os.Stderr, "@R{!!! missing required --password parameter}\n")
		os.Exit(2)
	}

	go firehose.Go(&Nozzle{}, opts.Config)
	http.Handle("/", API{
		Username: opts.Username,
		Password: opts.Password,
	})
	http.ListenAndServe(":"+opts.Port, nil)
}
