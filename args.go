package main

import (
	"flag"
	"log"
	"strings"

	"gitlab.quantdo.cn/yuanyang/goload/utils"
)

func initArgs() {
	flag.BoolVar(&noBoomer, "noboomer", false, "Runnning in cli mode with out boomer.")

	flag.StringVar(&host, "host", defaultHost, "Host to take pressure.")
	flag.StringVar(&baseURI, "base", defaultURI, "Default api base URI.")

	flag.StringVar(&symbol, "symbol", defaultSymbol, "Order symbol.")
	flag.Int64Var(&quantity, "quantity", defaultQuantity, "Order quantity.")
	flag.Var(&side, "side", "Order quantity.")
	flag.Float64Var(&price, "price", defaultPrice, "Order price.")
	flag.IntVar(&precision, "precis", defaultPrecision, "Precesion for random price.")
	flag.Float64Var(&basePrice, "base-price", defaultPrice, "Base price for random price.")
	flag.Int64Var(&maxQuantity, "max-quantity", defaultMaxQuantity, "Max quantity for random quantity.")

	flag.StringVar(&apiKey, "key", "", "API-Key for order.")
	flag.StringVar(&apiSecret, "secret", "", "API-Secret for order.")
	flag.StringVar(&identity, "identity", defaultIdentity, "Identity for login.")
	flag.StringVar(&password, "pass", defaultPassword, "Password for login.")

	flag.BoolVar(&randPrice, "rand-price", false, "Generate random price[BASE_PRICE.00, BASE_PRICE.99].")
	flag.BoolVar(&randQuantity, "rand-quant", false, "Generate random quantity[1, MAX_QUANTITY].")
	flag.BoolVar(&randSide, "rand-side", false, "Generate order in random side.")
	flag.BoolVar(&bothSide, "both-side", false, "Generate order in both side.")

	flag.IntVar(&count, "count", 1, "Order count per worker.")

	flag.BoolVar(&dryRun, "dry-run", false, "Dry-run under test.")
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

func init() {
	initArgs()
}
