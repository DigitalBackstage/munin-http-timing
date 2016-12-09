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

// DoPing calls the pinger and prints the response for munin
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

		fmt.Fprint(buf, requestInfoToString(info, config.GetGraphName()))
		totals = append(totals, info.TotalString())
	}

	fmt.Fprintf(buf, "multigraph %s\n", config.GetGraphName())
	for _, value := range totals {
		fmt.Fprintf(buf, value)
	}
	fmt.Fprint(buf, "\n")

	return buf.String(), nil
}

// requestInfoToString returns the timings following the Munin multigraph protocol
// It prints the fields in a specific order, it must match the one in
// graphOrder in config.go
func requestInfoToString(t *pinger.RequestInfo, graphName string) string {
	t.Lock()
	defer t.Unlock()

	buf := &bytes.Buffer{}
	fmt.Fprintf(buf, "multigraph %s.%s\n", graphName, t.Name)

	if t.IsOk() {
		fmt.Fprintf(buf, "resolving.value %v\n", pinger.ToMillisecond(t.Resolving))
		fmt.Fprintf(buf, "connecting.value %v\n", pinger.ToMillisecond(t.Connecting))
		fmt.Fprintf(buf, "sending.value %v\n", pinger.ToMillisecond(t.Sending))
		fmt.Fprintf(buf, "waiting.value %v\n", pinger.ToMillisecond(t.Waiting))
		fmt.Fprintf(buf, "receiving.value %v\n", pinger.ToMillisecond(t.Receiving))
	} else {
		fmt.Fprint(buf, "resolving.value U\n")
		fmt.Fprint(buf, "connecting.value U\n")
		fmt.Fprint(buf, "sending.value U\n")
		fmt.Fprint(buf, "waiting.value U\n")
		fmt.Fprint(buf, "receiving.value U\n")
	}

	fmt.Fprint(buf, "\n")

	return buf.String()
}
