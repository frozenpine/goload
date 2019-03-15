package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/frozenpine/ngerest"
	"github.com/myzhan/boomer"
)

const (
	defaultHost = "http://trade"
	defaultURI  = "/api/v1"
)

var (
	client            *ngerest.APIClient
	eventBus          = boomer.Events
	rootCtx, stopFunc = context.WithCancel(context.Background())

	host, baseURI string

	eventHatchComplete = make(chan bool)

	method = "REST"
	name   = "OrderNew"
)

func initClient() {
	client = ngerest.NewAPIClient(ngerest.NewConfiguration())
}

func initArgs() {
	flag.StringVar(&host, "host", defaultHost, "Host to take pressure.")
	flag.StringVar(&baseURI, "base", defaultURI, "Default api base URI.")
}

func makeOrder(auth context.Context) *ngerest.Order {
	// ordOpts := ngerest.OrderNewOpts{
	// 	Side:     optional.NewString("Buy"),
	// 	OrderQty: optional.NewFloat32(float32(1)),
	// 	Price:    optional.NewFloat64(float64(3536))}

	ordOpts := ngerest.OrderNewOpts{}

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

	return &ord
}

func worker() {
	if !flag.Parsed() {
		flag.Parse()
	}
	client.ChangeBasePath(host)

	auth := context.WithValue(rootCtx, ngerest.ContextAPIKey, ngerest.APIKey{
		Key:    "o9EijdKODwl4a26JI0B2",
		Secret: "Lfp4qu5P5Ot63NIba9Wm132a3sCaAet3N7KJr0DrtJ54r6ZnHo6FrV89sG68q4mOK4dia52Epu5H1uUxJEc8KraKZ16B79EJB4V"})

	if !<-eventHatchComplete {
		log.Fatalln("Hatch failed.")
	}

	for {
		select {
		case <-rootCtx.Done():
			log.Println("Exit")
		default:
			makeOrder(auth)
		}
	}
}

func main() {
	if !flag.Parsed() {
		flag.Parse()
	}
	client.ChangeBasePath(
		fmt.Sprintf("%s/%s", strings.TrimRight(host, "/"), strings.TrimLeft(baseURI, "/")))

	eventBus.Subscribe("boomer:hatch_complete", func() { eventHatchComplete <- true })
	eventBus.Subscribe("boomer:hatch_stop", func() { stopFunc() })

	task := &boomer.Task{
		Name: "Order_new",
		Fn:   worker}

	boomer.Run(task)
}

func init() {
	initClient()

	initArgs()
}
