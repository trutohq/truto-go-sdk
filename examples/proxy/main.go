// Calls an integration's own API through the Proxy API.
//
//	TRUTO_API_TOKEN=... INTEGRATED_ACCOUNT_ID=... go run ./examples/proxy
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	truto "github.com/trutohq/truto-go-sdk"
)

func main() {
	ctx := context.Background()
	client := truto.NewClient(os.Getenv("TRUTO_API_TOKEN"))
	deals := truto.ProxyResource{Resource: "deals", IntegratedAccountID: os.Getenv("INTEGRATED_ACCOUNT_ID")}

	first, err := client.Proxy.List(ctx, deals).First()
	if err == truto.Done {
		fmt.Println("no deals")
		return
	}
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("first deal:", (*first)["id"])

	updated, err := client.Proxy.Update(ctx, deals, fmt.Sprint((*first)["id"]), map[string]any{"stage": "won"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("updated:", *updated)
}
