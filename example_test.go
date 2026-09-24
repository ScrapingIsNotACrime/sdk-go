package sinac_test

import (
	"context"
	"errors"
	"fmt"
	"log"

	sinac "github.com/ScrapingIsNotACrime/sdk-go"
)

func Example() {
	client, err := sinac.NewClient(sinac.WithAPIKey("sinac_..."))
	if err != nil {
		log.Fatal(err)
	}
	profile, err := client.Instagram.Profile(context.Background(), "nasa")
	if errors.Is(err, sinac.ErrNotFound) {
		fmt.Println("no such profile")
		return
	}
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(profile.Username, profile.Followers)
}

func ExamplePage_All() {
	client, err := sinac.NewClient() // reads SCRAPINGISNOTACRIME_API_KEY
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	page, err := client.GitHub.Followers(ctx, "torvalds", &sinac.GitHubListParams{Limit: 100})
	if err != nil {
		log.Fatal(err)
	}
	// Each page fetched is one billed request, so stop at a bound.
	n := 0
	for user, err := range page.All(ctx) {
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(user.Username)
		if n++; n >= 250 {
			break
		}
	}
}
