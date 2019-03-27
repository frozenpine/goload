package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/frozenpine/ngerest"
	"gitlab.quantdo.cn/yuanyang/goload/utils"
)

func worker() {
	var auth context.Context

	if apiKey == "" || apiSecret == "" {
		idMap := utils.NewIdentityMap()

		login := make(map[string]string)

		if err := idMap.CheckIdentity(identity, login); err != nil {
			log.Fatalln(err)
		}

		// pubKey, _, err := client.UserApi.GetPublicKey(rootCtx)
		pubKey, _, err := client.KeyExchange.GetPublicKey(rootCtx)
		if err != nil {
			log.Fatalln(err)
		}

		login["password"] = pubKey.Encrypt(password)

		auth, _, err = client.User.UserLogin(rootCtx, login)

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
		sides = []utils.OrderSide{utils.Buy, utils.Sell}
	} else {
		if randSide {
			sides = []utils.OrderSide{utils.RandomSide()}
		} else {
			sides = []utils.OrderSide{utils.OrderSide(side)}
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

			ord := makeOrder(auth, symbol, price, quantity*sideValue.Value())

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
