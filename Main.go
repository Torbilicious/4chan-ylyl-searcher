package main

import (
	"github.com/moshee/go-4chan-api/api"
	"fmt"
	"time"
	"strings"
)

var (
	threads     []*api.Thread
	searchText  string
	searchBoard string
)


func InitVariables() {
	// All requests will be made with HTTPS
	api.SSL = true

	// will be pulled up to 10 seconds when first used
	api.UpdateCooldown = 5 * time.Second

	searchText = "ylyl"
	searchBoard = "gif"
}


func UpdateThreads(pages ...int) {

	for _, page := range pages {
		newThreads, err := api.GetIndex(searchBoard, page)
		if err != nil {
			panic(err)
		}

		for _, thread := range newThreads {

			firstPost := thread.Posts[0]

			//fmt.Println(firstPost.Subject)

			if strings.Contains(strings.ToLower(firstPost.Subject), searchText) {

				threads = append(threads, thread)
			}
		}
		fmt.Printf("Searched %v threads on page %v.\n", len(newThreads), page)
	}
}

func PrintThreads() {

	if len(threads) > 0 {

		fmt.Printf("Found %v threads totally!\n", len(threads))

		fmt.Println("")

		for _, thread := range threads {

			fmt.Printf("URL: https://boards.4chan.org/gif/thread/%v \nName: %v\n\n", thread.Id(), thread.Posts[0].Subject)
		}
	} else {

		fmt.Println("No threads were found.")
	}
}

func main() {

	UpdateThreads(1, 2, 3, 4, 5)

	InitVariables()

	PrintThreads()
}
