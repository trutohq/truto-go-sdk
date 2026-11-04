// Reads an integration's own API through the Proxy API: the first page of a
// resource and one record by ID. It only reads; nothing is changed.
//
//	TRUTO_API_TOKEN=... INTEGRATED_ACCOUNT_ID=... go run ./examples/proxy contacts
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	truto "github.com/trutohq/truto-go-sdk"
)

func main() {
	resource := "contacts"
	if len(os.Args) > 1 {
		resource = os.Args[1]
	}
	ctx := context.Background()
	client := truto.NewClient(os.Getenv("TRUTO_API_TOKEN"))
	r := truto.ProxyResource{Resource: resource, IntegratedAccountID: os.Getenv("INTEGRATED_ACCOUNT_ID")}

	page, err := client.Proxy.List(ctx, r, truto.WithQuery(truto.Params{"limit": 10})).NextPage()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d %s on the first page (more: %v)\n", len(page.Items), resource, page.NextCursor != "")
	if len(page.Items) == 0 {
		return
	}

	id := fmt.Sprint(page.Items[0]["id"])
	record, err := client.Proxy.Get(ctx, r, id)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("first record:", *record)
}
