package main

import (
	fmt "github.com/jhunt/go-ansi"
	"net/http"
	"os"

	"github.com/jhunt/go-cli"
	env "github.com/jhunt/go-envirotron"
	"github.com/jhunt/go-firehose"
	"github.com/jhunt/go-log"
)

func main() {
	var opts struct {
		Config   string `cli:"-c, --config"`
		Username string `cli:"-u, --username"  env:"CACHE_FIRE_USERNAME"`
		Password string `cli:"-p, --password"  env:"CACHE_FIRE_PASSWORD"`
		Port     string `cli:"--port"          env:"PORT"`
		LogLevel string `cli:"-l, --log-level" env:"LOG_LEVEL"`
		Debug    bool   `cli:"-D, --debug"     env:"DEBUG"`
	}
	opts.Username = "cachefire"
	opts.Port = "3000"
	opts.LogLevel = "warning"
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

	if opts.Debug {
		opts.LogLevel = "debug"
	}
	log.SetupLogging(log.LogConfig{
		Type: "console",
		Level: opts.LogLevel,
	})

	go firehose.Go(&Nozzle{}, opts.Config)
	http.Handle("/", API{
		Username: opts.Username,
		Password: opts.Password,
	})
	http.ListenAndServe(":"+opts.Port, nil)
}
