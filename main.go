package main

import (
	"github.com/frozenpine/nge4go/swagger"
	"context"
	"flag"
	"fmt"
	"strconv"
	"strings"

	"github.com/antihax/optional"
	"github.com/frozenpine/nge4go"
	"github.com/myzhan/boomer"
)

const (
	defaultHost = "http://trade"
	defaultURI  = "/api/v1"
)

var (
	client            *nge4go.swagger.APIClient
	evnetBus          boomer.Events
	rootCtx, stopFunc = context.WithCancel(context.Background())

	host, baseURI string

	eventHatchComplete = make(chan bool)
)

func initClient() {
	client = nge4go.swagger.NewAPIClient(nge4go.swagger.NewConfiguration())
}

func initArgs() {
	flag.StringVar(&host, "host", defaultHost, "Host to take pressure.")
	flag.StringVar(&baseURI, "base", defaultURI, "Default api base URI.")
}

func worker() {
	auth := context.WithValue(rootCtx, ngeSw.ContextAPIKey, ngeSw.APIKey{
		Key:    "o9EijdKODwl4a26JI0B2",
		Secret: "Lfp4qu5P5Ot63NIba9Wm132a3sCaAet3N7KJr0DrtJ54r6ZnHo6FrV89sG68q4mOK4dia52Epu5H1uUxJEc8KraKZ16B79EJB4V"})

	<-eventHatchComplete

	method := "REST"
	name := "OrderNew"

	for {
		ordOpts := nge4go.swagger.OrderNewOpts{
			Side:     optional.NewString("Buy"),
			OrderQty: optional.NewFloat32(1),
			Price:    optional.NewFloat64(3536)}

		start := boomer.Now()
		ord, rsp, err := client.OrderApi.OrderNew(auth, "XBTUSD", &ordOpts)
		elapsed := start - boomer.Now()

		if rsp.StatusCode > 300 || err != nil {
			boomer.RecordFailure(method, name, elapsed, err.Error())
		}

		bodyLength, err := strconv.Atoi(rsp.Header.Get("Content-Length"))
		if err != nil {
			bodyLength = 0
		}
		boomer.RecordSuccess(method, name, elapsed, int64(bodyLength))
	}
}

func main() {
	if !flag.Parsed() {
		flag.Parse()
	}
	client.ChangeBasePath(
		fmt.Sprintf("%s/%s", strings.TrimRight(host, "/"), strings.TrimLeft(baseURI, "/")))

	boomer.Events.Subscribe("boomer:hatch_complete", func() { eventHatchComplete <- true })

	task := &boomer.Task{
		Name: "Order_new",
		Fn:   worker}

	boomer.Run(task)
}

func init() {
	initClient()
}
