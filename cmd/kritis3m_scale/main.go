package main

import (
	"os"
	"time"

	"github.com/Laboratory-for-Safe-and-Secure-Systems/kritis3m_scale/cmd/cli"
	"github.com/jagottsicher/termcolor"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	var colors bool
	switch l := termcolor.SupportLevel(os.Stderr); l {
	case termcolor.Level16M:
		colors = true
	case termcolor.Level256:
		colors = true
	case termcolor.LevelBasic:
		colors = true
	case termcolor.LevelNone:
		colors = false
	default:
		// no color, return text as is.
		colors = false
	}

	// Adhere to no-color.org manifesto of allowing users to
	// turn off color in cli/services

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
		NoColor:    !colors,
	})

	cli.Execute()
}
