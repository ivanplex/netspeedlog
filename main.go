package main

import (
	"log/slog"
	"netspeedlog/netspeedlog"
	"os"
	"time"
)

func main() {
	textHander := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(textHander)

	nsl := netspeedlog.New(logger)

	for {
		go nsl.SpeedTest()
		time.Sleep(5 * time.Minute)
	}
}
