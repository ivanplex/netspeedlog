package netspeedlog

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/showwin/speedtest-go/speedtest"
)

type Nsl struct {
	logger          *slog.Logger
	speedTestClient *speedtest.Speedtest
	// current lowest latency server
	curLowLatServer *speedtest.Server
}

func New(logger *slog.Logger) *Nsl {
	nsl := &Nsl{
		logger:          logger,
		speedTestClient: speedtest.New(),
	}
	return nsl
}

func (n *Nsl) SpeedTest() {
	if isInternetDown() {
		n.logger.Info("internet is down")
		return
	}

	n.refreshLowestLatencyServer()
	t := n.curLowLatServer

	t.PingTest(nil)
	t.DownloadTest()
	t.UploadTest()
	n.logger.Info(
		"speed test result",
		"latency", fmt.Sprintf("%d", t.Latency.Milliseconds()),
		"download", fmt.Sprintf("%.2f", t.DLSpeed),
		"upload", fmt.Sprintf("%.2f", t.ULSpeed))
}

func (n *Nsl) refreshLowestLatencyServer() error {
	serverList, err := n.speedTestClient.FetchServers()
	if err != nil {
		n.logger.Error("unable to fetch speedtest server list",
			"error", err)
	}

	targets, err := serverList.FindServer([]int{})
	if err != nil {
		n.logger.Error("unable to select a speed test server with the lowest latency",
			"error", err)
	}

	n.curLowLatServer = targets[0]
	n.logger.Info("updated lowest latency server",
		"server", n.curLowLatServer.String())
	return nil
}

func isInternetDown() (ok bool) {
	_, err := http.Get("http://clients3.google.com/generate_204")
	return err != nil
}
