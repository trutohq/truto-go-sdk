// Shows timeouts, retry hooks and error handling.
//
//	TRUTO_API_TOKEN=... go run ./examples/errors_and_retries
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	truto "github.com/trutohq/truto-go-sdk"
)

func main() {
	ctx := context.Background()
	client := truto.NewClient(os.Getenv("TRUTO_API_TOKEN"),
		truto.WithTimeout(10*time.Second),
		truto.WithMaxRetries(5),
		truto.WithOnRetry(func(i truto.RetryInfo) {
			log.Printf("retry %d after status %d, waiting %s", i.RetryCount, i.StatusCode, i.Wait)
		}),
		truto.WithOnRateLimited(func(i truto.RateLimitInfo) {
			log.Printf("rate limited, waiting %s", i.RetryAfter)
		}),
	)

	_, err := client.IntegratedAccounts.Get(ctx, "does-not-exist", truto.WithNoRetry())
	var apiErr *truto.Error
	switch {
	case truto.IsNotFound(err):
		fmt.Println("not found")
	case errors.Is(err, truto.ErrTimeout):
		fmt.Println("timed out")
	case errors.As(err, &apiErr):
		fmt.Println("API error", apiErr.StatusCode, apiErr.Message())
	case err != nil:
		log.Fatal(err)
	}
}
