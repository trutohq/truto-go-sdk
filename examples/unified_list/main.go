// Lists CRM contacts for one integrated account through the Unified API,
// page by page, and decodes them into a local type.
//
//	TRUTO_API_TOKEN=... INTEGRATED_ACCOUNT_ID=... go run ./examples/unified_list
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	truto "github.com/trutohq/truto-go-sdk"
)

type Contact struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func main() {
	ctx := context.Background()
	client := truto.NewClient(os.Getenv("TRUTO_API_TOKEN"))

	contacts := client.Unified.List(ctx, truto.UnifiedResource{
		Model:               "crm",
		Resource:            "contacts",
		IntegratedAccountID: os.Getenv("INTEGRATED_ACCOUNT_ID"),
	}, truto.WithQuery(truto.Params{"limit": 50}))

	for record, err := range contacts.Seq() {
		if err != nil {
			log.Fatal(err)
		}
		c, err := truto.Decode[Contact](record)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(c.ID, c.FirstName, c.LastName)
	}
}
