package munin

import (
	"bytes"
	"fmt"
	"sort"
	"time"

	"github.com/DigitalBackstage/munin-http-timing/pinger"
)

// requestInfoToString returns the timings following the Munin multigraph protocol
// It prints the fields in a specific order, it must match the one in
// graphOrder in config.go
func formatRequestInfo(t *pinger.RequestInfo, graphName string) string {
	t.Lock()
	defer t.Unlock()

	buf := &bytes.Buffer{}
	fmt.Fprintf(buf, "multigraph %s.%s\n", graphName, t.Name)

	if t.IsOk() {
		fmt.Fprintf(buf, "resolving.value %v\n", toMillisecond(t.Resolving))
		fmt.Fprintf(buf, "connecting.value %v\n", toMillisecond(t.Connecting))
		fmt.Fprintf(buf, "sending.value %v\n", toMillisecond(t.Sending))
		fmt.Fprintf(buf, "waiting.value %v\n", toMillisecond(t.Waiting))
		fmt.Fprintf(buf, "receiving.value %v\n", toMillisecond(t.Receiving))
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

// TotalString returns the <name>_total.value line for this RequestInfo
func formatRequestInfoTotal(t *pinger.RequestInfo) string {
	t.Lock()
	defer t.Unlock()

	value := "U"
	if t.IsOk() {
		value = fmt.Sprintf("%v", toMillisecond(t.Total))
	}

	return fmt.Sprintf("%s_total.value %v\n", t.Name, value)
}

func toMillisecond(d time.Duration) int64 {
	return int64(d / time.Millisecond)
}

func formatMultigraph(requests []*pinger.RequestInfo, graphName string) string {
	sort.Sort(requestByName(requests))

	buf := &bytes.Buffer{}
	for i := range requests {
		fmt.Fprint(buf, formatRequestInfo(requests[i], graphName))
	}

	fmt.Fprintf(buf, "multigraph %s\n", graphName)
	for _, value := range requests {
		fmt.Fprintf(buf, formatRequestInfoTotal(value))
	}
	fmt.Fprint(buf, "\n")

	return buf.String()
}

type requestByName []*pinger.RequestInfo

func (a requestByName) Len() int           { return len(a) }
func (a requestByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a requestByName) Less(i, j int) bool { return a[i].Name < a[j].Name }
