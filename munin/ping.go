package munin

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/DigitalBackstage/munin-http-timing/config"
	"github.com/DigitalBackstage/munin-http-timing/pinger"
)

var stderr = log.New(os.Stderr, "", 0)

// DoPing does the actual stats gathering (HTTP requests) and prints it for munin
func DoPing(config config.Config) (string, error) {
	rand.Seed(time.Now().Unix())

	if len(config.URIs) <= 0 {
		return "", errors.New("No URIs provided.")
	}

	totals := []string{}
	queue := make(chan *pinger.RequestInfo, len(config.URIs))
	pinger.DoParallelPings(config, queue)

	buf := &bytes.Buffer{}

	for i := 0; i < len(config.URIs); i++ {
		info := <-queue

		if info.Error != nil {
			stderr.Print(info.Error)
		}

		fmt.Fprint(buf, info)
		totals = append(totals, info.TotalString())
	}

	fmt.Fprint(buf, "multigraph timing\n")
	for _, value := range totals {
		fmt.Fprintf(buf, value)
	}
	fmt.Fprint(buf, "\n")

	return buf.String(), nil
}
