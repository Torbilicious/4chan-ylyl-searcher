package main

import (
	"github.com/moshee/go-4chan-api/api"
	"fmt"
	"time"
	"strings"
)

var (
	threads []*api.Thread
)


func ExampleVariables() {
	// All requests will be made with HTTPS
	api.SSL = true

	// will be pulled up to 10 seconds when first used
	api.UpdateCooldown = 5 * time.Second
}


func UpdateThreads(pages ...int) {

	for _, page := range pages {
		newThreads, err := api.GetIndex("gif", page)
		if err != nil {
			panic(err)
		}

		for _, thread := range newThreads {

			firstPost := thread.Posts[0]

			//fmt.Println(firstPost.Subject)

			if strings.Contains(strings.ToLower(firstPost.Subject), "ylyl") {

				threads = append(threads, thread)
			}
		}
		fmt.Printf("Searched %v threads.\n", len(newThreads))
	}


}

func PrintThreads() {

	if len(threads) > 0 {
		//fmt.Print(threads)

		fmt.Printf("Found %v threads totally!", len(threads))
	}

}

func main() {

	UpdateThreads(1, 2, 3, 4, 5)

	ExampleVariables()

	PrintThreads()
}