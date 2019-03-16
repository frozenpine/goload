package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"strconv"
	"strings"

	"gitlab.quantdo.cn/yuanyang/goload/utils"

	"github.com/antihax/optional"
	"github.com/frozenpine/ngerest"
	"github.com/myzhan/boomer"
)

const (
	defaultHost      = "http://trade"
	defaultURI       = "/api/v1"
	defaultSymbol    = "XBTUSD"
	defaultQuantity  = int64(1)
	defaultPrice     = float64(3536)
	defaultPrecision = 2
	defaultSide      = "Buy"

	defaultIdentity = "yuanyang@quantdo.com.cn"
	defaultPassword = "quantdo123456"
)

var (
	noBoomer bool

	client *ngerest.APIClient

	rootCtx, stopFunc = context.WithCancel(context.Background())

	host, baseURI string
	symbol        string
	quantity      int64
	price         float64
	side          string
	sides         []int64
	precision     int
	basePrice     float64
	maxQuantity   = int64(100)

	apiKey, apiSecret  string
	identity, password string

	randPrice, randQuantity, randSide, bothSide bool

	count int

	method = "New"
	name   = "Order"
)

func initClient() {
	client = ngerest.NewAPIClient(ngerest.NewConfiguration())
}

func initArgs() {
	flag.BoolVar(&noBoomer, "noboomer", false, "Runnning in cli mode with out boomer.")

	flag.StringVar(&host, "host", defaultHost, "Host to take pressure.")
	flag.StringVar(&baseURI, "base", defaultURI, "Default api base URI.")

	flag.StringVar(&symbol, "symbol", defaultSymbol, "Order symbol.")
	flag.Int64Var(&quantity, "quantity", defaultQuantity, "Order quantity.")
	flag.StringVar(&side, "side", defaultSide, "Order side.")
	flag.Float64Var(&price, "price", defaultPrice, "Order price.")
	flag.IntVar(&precision, "precis", defaultPrecision, "Precesion for random price.")
	flag.Float64Var(&basePrice, "base-price", defaultPrice, "Base price for random price.")
	flag.Int64Var(&maxQuantity, "max-quantity", maxQuantity, "Max quantity for random quantity.")

	flag.StringVar(&apiKey, "key", "", "API-Key for order.")
	flag.StringVar(&apiSecret, "secret", "", "API-Secret for order.")
	flag.StringVar(&identity, "identity", defaultIdentity, "Identity for login.")
	flag.StringVar(&password, "pass", defaultPassword, "Password for login.")

	flag.BoolVar(&randPrice, "rand-price", false, "Generate random price[BASE_PRICE.00, BASE_PRICE.99].")
	flag.BoolVar(&randQuantity, "rand-quant", false, "Generate random quantity[1, MAX_QUANTITY].")
	flag.BoolVar(&randSide, "rand-side", false, "Generate order in random side.")
	flag.BoolVar(&bothSide, "both-side", false, "Generate order in both side.")

	flag.IntVar(&count, "count", 1, "Order count per worker.")
}

func validateArgs() {
	basePath := strings.Join([]string{strings.Trim(host, "/"), strings.Trim(baseURI, "/")}, "/")
	log.Println("Change host to:", basePath)
	client.ChangeBasePath(basePath)

	if err := utils.CheckSymbol(symbol); err != nil {
		log.Fatalln(err)
	}
	if err := utils.CheckPrice(price); err != nil {
		log.Fatalln(err)
	}
	if err := utils.CheckQuantity(quantity); err != nil {
		log.Fatalln(err)
	}

	if err := utils.MatchSide(&side, quantity); err != nil {
		log.Fatalln(err)
	}
}

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
	ord, rsp, err := client.OrderApi.OrderNew(auth, ordSym, &ordOpts)
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

func worker() {
	var auth context.Context

	if apiKey == "" || apiSecret == "" {
		idMap := utils.NewIdentityMap()

		login := make(map[string]string)

		if err := idMap.CheckIdentity(identity, login); err != nil {
			log.Fatalln(err)
		}

		pubKey, _, err := client.UserApi.GetPublicKey(rootCtx)
		if err != nil {
			log.Fatalln(err)
		}

		login["password"] = pubKey.Encrypt(password)

		auth, _, err = client.UserApi.UserLogin(rootCtx, login)

		if err != nil {
			log.Fatalln(err)
		}
	} else {
		auth = context.WithValue(
			rootCtx, ngerest.ContextAPIKey, ngerest.APIKey{
				Key:    apiKey,
				Secret: apiSecret,
			})
	}

	if bothSide {
		sides = []int64{1, -1}
	} else {
		if randSide {
			sides = []int64{utils.RandomSide().Value()}
		} else {
			sides = []int64{utils.OrderSide(side).Value()}
		}
	}

	for count > 0 {
		for idx, sideValue := range sides {
			if !(bothSide && idx > 0) {
				if randPrice {
					utils.RandomPrice(&price, precision, basePrice)
				}

				if randQuantity {
					utils.RandomQuantity(&quantity, maxQuantity)
				}
			}

			ord := makeOrder(auth, symbol, price, quantity*sideValue)

			if noBoomer {
				if ord != nil {
					result, _ := json.Marshal(ord)
					log.Println(string(result))
				} else {
					log.Println("making order failed.")
				}
				count--
			}
		}
	}
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
	initClient()
	initArgs()
}
