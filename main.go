package main

import (
	"flag"

	"github.com/sirupsen/logrus"
	"github.com/tjgq/sane"

	"github.com/babolivier/scanner/config"
	"github.com/babolivier/scanner/http"
	"github.com/babolivier/scanner/scanner"
	"github.com/babolivier/scanner/webdav"
)

var (
	cfgPath = flag.String("c", "config.yaml", "Path to the configuration file")
)

func main() {
	// Parse the command-line arguments.
	flag.Parse()

	// Configure the logger to log the full date.
	logrus.SetFormatter(
		&logrus.TextFormatter{
			TimestampFormat: "2006-02-01 15:04:05.999",
			FullTimestamp:   true,
		},
	)

	// Parse the configuration file.
	cfg, err := config.NewConfig(*cfgPath)
	if err != nil {
		panic(err)
	}

	// Instantiate the WebDAV client.
	webDAVClient := webdav.NewClient(cfg.WebDAV)

	// Instantiate the scanner.
	s, err := scanner.NewScanner(cfg.Scanner, webDAVClient)
	if err != nil {
		panic(err)
	}

	// Close the SANE connection and release all resources in use by SANE when exiting.
	defer sane.Exit()

	// Start the HTTP server.
	if err = http.ListenAndServe(cfg.HTTP, s, webDAVClient); err != nil {
		panic(err)
	}
}
