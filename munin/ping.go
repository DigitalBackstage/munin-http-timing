package munin

import (
	"errors"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/DigitalBackstage/munin-http-timing/config"
	"github.com/DigitalBackstage/munin-http-timing/pinger"
)

var stderr = log.New(os.Stderr, "", 0)

// DoPing calls the pinger and returns the response formatted for munin
func DoPing(config config.Config) (string, error) {
	rand.Seed(time.Now().Unix())

	if len(config.URIs) <= 0 {
		return "", errors.New("No URIs provided.")
	}

	requests := make([]*pinger.RequestInfo, 0, len(config.URIs))
	queue := make(chan *pinger.RequestInfo, len(config.URIs))
	pinger.DoParallelPings(config, queue)

	for i := 0; i < len(config.URIs); i++ {
		info := <-queue
		if info.Error != nil {
			stderr.Print(info.Error)
		}

		requests = append(requests, info)
	}

	return formatMultigraph(requests, config.GetGraphName()), nil
}
