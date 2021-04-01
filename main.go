package main

import (
	"os"
	"strings"

	"github.com/shiena/ansicolor"

	"github.com/SierraSoftworks/minback-cleanup/commands"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var version = "v1.0.0"

func main() {
	app := cli.NewApp()
	app.Author = "Benjamin Pannell"
	app.Email = "admin@sierrasoftworks.com"
	app.Copyright = "Sierra Softworks © 2017"
	app.Version = version

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "log-level",
			Usage: "DEBUG|INFO|WARN|ERROR",
			Value: "INFO",
		},
	}

	app.Commands = []cli.Command{
		commands.Cleanup,
	}

	app.Before = func(c *cli.Context) error {
		logLevel := c.GlobalString("log-level")
		switch strings.ToUpper(logLevel) {
		case "DEBUG":
			log.SetLevel(log.DebugLevel)
		case "INFO":
			log.SetLevel(log.InfoLevel)
		case "WARN":
			log.SetLevel(log.WarnLevel)
		case "ERROR":
			log.SetLevel(log.ErrorLevel)
		default:
			log.SetLevel(log.InfoLevel)
		}

		return nil
	}

	log.SetOutput(ansicolor.NewAnsiColorWriter(os.Stdout))

	err := app.Run(os.Args)
	if err != nil {
		os.Exit(1)
	}
}
