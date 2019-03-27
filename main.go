package main

import (
	"context"
	"flag"
	"log"
	"strconv"

	"gitlab.quantdo.cn/yuanyang/goload/utils"

	"github.com/antihax/optional"
	"github.com/frozenpine/ngerest"
	"github.com/myzhan/boomer"
)

var (
	noBoomer bool

	client *ngerest.APIClient

	rootCtx, stopFunc = context.WithCancel(context.Background())

	host, baseURI string
	symbol        string
	quantity      int64
	price         float64
	side          = defaultSide
	sides         []utils.OrderSide
	precision     int
	basePrice     float64
	maxQuantity   int64
	orderList     []*Order

	apiKey, apiSecret  string
	identity, password string

	randPrice, randQuantity, randSide, bothSide bool

	count int

	method = "New"
	name   = "Order"
)

// func randomOrders() ([]*Order, error) {

// }

func makeOrder(auth context.Context, ordSym string, ordPrice float64, ordVol int64) *ngerest.Order {
	var side string
	if ordVol > 0 {
		side = utils.Buy.String()
	} else {
		side = utils.Sell.String()
		ordVol = -ordVol
	}

	if noBoomer {
		log.Printf("%s [%s]: %d@%."+strconv.Itoa(precision)+"f\n", side, ordSym, ordVol, ordPrice)
	}

	ordOpts := ngerest.OrderNewOpts{
		Side:     optional.NewString(side),
		OrderQty: optional.NewFloat32(float32(ordVol)),
		Price:    optional.NewFloat64(ordPrice)}

	start := boomer.Now()
	ord, rsp, err := client.Order.OrderNew(auth, ordSym, &ordOpts)
	elapsed := boomer.Now() - start

	if err != nil {
		boomer.RecordFailure(method, name, elapsed, err.Error())
		return nil
	}

	bodyLength, err := strconv.Atoi(rsp.Header.Get("Content-Length"))
	if err != nil {
		bodyLength = 0
	}
	boomer.RecordSuccess(method, name, elapsed, int64(bodyLength))

	return &ord
}

func main() {
	if !flag.Parsed() {
		flag.Parse()
	}
	validateArgs()

	if noBoomer {
		worker()
	} else {
		task := &boomer.Task{
			Name: name,
			Fn:   worker,
		}
		boomer.Run(task)
	}
}

func init() {
	client = ngerest.NewAPIClient(ngerest.NewConfiguration())
}
