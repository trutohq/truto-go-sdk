// Creates a link token for a tenant and lists the tenant's connected accounts.
//
//	TRUTO_API_TOKEN=... go run ./examples/link_token acme
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	truto "github.com/trutohq/truto-go-sdk"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: link_token <tenant-id>")
	}
	tenantID := os.Args[1]
	ctx := context.Background()
	client := truto.NewClient(os.Getenv("TRUTO_API_TOKEN"))

	lt, err := client.LinkTokens.Create(ctx, truto.LinkTokenRequest{TenantID: truto.Ptr(tenantID)})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("link token:", *lt.LinkToken)

	accounts, err := client.IntegratedAccounts.List(ctx,
		truto.WithQuery(truto.Params{"tenant_id": tenantID})).All()
	if err != nil {
		log.Fatal(err)
	}
	for _, a := range accounts {
		name := ""
		if a.Integration != nil && a.Integration.Name != nil {
			name = *a.Integration.Name
		}
		status := ""
		if a.Status != nil {
			status = string(*a.Status)
		}
		fmt.Printf("%s %s %s\n", *a.ID, name, status)
	}
}
