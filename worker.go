package main

import (
	"context"
	"log"

	"github.com/frozenpine/ngerest"
	"gitlab.quantdo.cn/yuanyang/goload/utils"
)

func worker() {
	var auth context.Context

	if !dryRun {
		if apiKey == "" || apiSecret == "" {
			idMap := utils.NewIdentityMap()

			login := make(map[string]string)

			if err := idMap.CheckIdentity(identity, login); err != nil {
				log.Fatalln(err)
			}

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
	}

	for _, order := range orderList {
		if dryRun {
			log.Println("[DRY-RUN]", (*order).String())
		} else {
			result := makeOrder(auth, order.Symbol, order.Price, order.GetQuantity())

			orderResults = append(orderResults, result)
		}
	}
}
