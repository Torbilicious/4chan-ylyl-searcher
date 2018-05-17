//
// ██╗   ██╗██╗  ██╗   ██╗██╗                   ███████╗███████╗ █████╗ ██████╗  ██████╗██╗  ██╗███████╗██████╗
// ╚██╗ ██╔╝██║  ╚██╗ ██╔╝██║                   ██╔════╝██╔════╝██╔══██╗██╔══██╗██╔════╝██║  ██║██╔════╝██╔══██╗
//  ╚████╔╝ ██║   ╚████╔╝ ██║         █████╗    ███████╗█████╗  ███████║██████╔╝██║     ███████║█████╗  ██████╔╝
//   ╚██╔╝  ██║    ╚██╔╝  ██║         ╚════╝    ╚════██║██╔══╝  ██╔══██║██╔══██╗██║     ██╔══██║██╔══╝  ██╔══██╗
//    ██║   ███████╗██║   ███████╗              ███████║███████╗██║  ██║██║  ██║╚██████╗██║  ██║███████╗██║  ██║
//    ╚═╝   ╚══════╝╚═╝   ╚══════╝              ╚══════╝╚══════╝╚═╝  ╚═╝╚═╝  ╚═╝ ╚═════╝╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝
//
// @author Torben Binder
// @author Max Bebök

package main

import (
	"github.com/moshee/go-4chan-api/api"
	"gopkg.in/cheggaaa/pb.v1"
	"fmt"
	"time"
	"strings"
	"strconv"
	"os"
	"io/ioutil"
	"encoding/json"
)

const EXIT_SUCCESS = 0
const EXIT_ERROR   = 1

const BASE_URL_THREAD = "https://boards.4chan.org"
const BASE_URL_MEDIA  = "https://i.4cdn.org"

// use https://mholt.github.io/json-to-go/ to generate the structure
// names must start with an upper-case letter for some reasons...
type Config struct {
	SearchText []string `json:"searchText"`
	Blacklist []string `json:"blacklist"`
	SearchBoard string `json:"searchBoard"`
	FilePath string `json:"filePath"`
}

var (
	threads []*api.Thread
	config 	Config
	dnList Download_List
)

/**
 * reads a JSON config file into the global "config"
 * @type {string} path path to the JSON config
 * @return bool true on success, false on any error
 */
func ReadConfig(path string) bool {

	file, err := ioutil.ReadFile(path)
    if err != nil {
        fmt.Printf("Error parsing config: %v\n", err)
        return false
    }

    if err := json.Unmarshal(file, &config); err != nil {
        fmt.Printf("Error mapping config: %v\n", err)
    }

	return true
}

func InitVariables() {

	api.SSL = true // All requests will be made with HTTPS
	api.UpdateCooldown = 5 * time.Second // will be pulled up to 10 seconds when first used
}


func SearchThreads(endPage int) bool {

	bar := pb.StartNew(endPage)
	foundThread := false

	for page := 1; page <= endPage; page++ {

		newThreads, err := api.GetIndex(config.SearchBoard, page)
		if err != nil {
			panic(err)
		}

		THREAD_LOOP:
		for _, thread := range newThreads {

			firstPost := thread.Posts[0]
			wholeText := strings.ToLower(firstPost.Subject + firstPost.Comment)

			for _, txt := range config.Blacklist {
				if (strings.Contains(wholeText, txt)) {
					continue THREAD_LOOP
				}
			}

			for _, txt := range config.SearchText {
				if (strings.Contains(wholeText, txt)) {
					threads = append(threads, thread)
					foundThread = true
					break
				}
			}
		}
		bar.Increment()
	}

	bar.FinishPrint("Done!")
	return foundThread
}

/**
 * create a list of all files to download
 * and fetch statistics for the number and size of files
 * @return {bool} returns true if at least 1 file was found
 */
func CreateDownloadList() bool {

	fileFound := false
	bar := pb.StartNew(len(threads))
	alreadyDownloaded := 0

	for _, t := range threads {
		thread, _ := api.GetThread(config.SearchBoard, t.Id()) // get whole thread with all replies

		for _, post := range thread.Posts {
			if post.File != nil {

				fileName := strconv.FormatInt(post.File.Id, 10) + post.File.Ext

				if(dnList.CheckForLocalFile(fileName)) {
					alreadyDownloaded++
				} else {
					dnList.Add(post.File)
					fileFound = true
				}
			}
		}

		bar.Increment()
	}

	statStr := strconv.Itoa(dnList.NumTotal)   + " file(s) foud | "  +
			   strconv.Itoa(alreadyDownloaded) + " files(s) skipped | " +
			   strconv.Itoa(dnList.SizeTotal / 1024) + " KB total";
	bar.FinishPrint("Done!, " + statStr)

	return fileFound
}

func PrintThreads() {

	if len(threads) > 0 {

		fmt.Printf("Found %v threads totally!\n", len(threads))
		fmt.Println("")

		for _, thread := range threads {

			fmt.Printf("URL: %v/%v/thread/%v \nName: %v\n\n", BASE_URL_THREAD, config.SearchBoard, thread.Id(), thread.Posts[0].Subject)
		}
	} else {

		fmt.Println("No threads were found.")
	}
}

func main() {

	InitVariables()

	if !ReadConfig("config.json") {
		os.Exit(EXIT_ERROR)
	}

	fmt.Println("Build local file DB...")
	dnList.ScanLocalFiles(config.FilePath);
	fmt.Printf("%d downloaded file(s) found!\n\n", len(dnList.LocalFiles))

	fmt.Println("Search for Threads...")
	if !SearchThreads(10) {
		fmt.Println("No threads found!")
		os.Exit(EXIT_SUCCESS)
	}
	fmt.Println("")

	PrintThreads()
	fmt.Println("")

	fmt.Println("Fetch Data to download...")
	if !CreateDownloadList() {
		fmt.Println("No files found!")
		os.Exit(EXIT_SUCCESS)
	}
	fmt.Println("")

	fmt.Println("Start downloading files...")
	dnList.DownloadAll()
	fmt.Println("")

}
