package main

import (
	"github.com/moshee/go-4chan-api/api"
	"gopkg.in/cheggaaa/pb.v1"
	"strconv"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
)

type Download_List struct {
	Files []*api.File
	LocalFiles []string
	SizeTotal int
	NumTotal int
}

/**
 * adds a file to the list and updates the statistics
 * @param  {api.File} file file to add
 */
func (d *Download_List) Add(file *api.File) {
    d.Files = append(d.Files, file)
    d.SizeTotal += file.Size
    d.NumTotal++
}

/**
 * scan for already downloaded files, save the list to "LocalFiles"
 * @type {string} download path
 */
func (d *Download_List) ScanLocalFiles(path string) {

	files, _ := ioutil.ReadDir(path)
	for _, f := range files {
		d.LocalFiles = append(d.LocalFiles, f.Name())
	}

	sort.Strings(d.LocalFiles)
}

/**
 * checks if a file-name is in the list of local files
 * @type {bool} true if it's already there, false if not
 */
func (d *Download_List) CheckForLocalFile(name string) bool {

	i := sort.SearchStrings(d.LocalFiles, name)
	if i < len(d.LocalFiles) && d.LocalFiles[i] == name {
	    return true
	}
	return false
}

/**
 * downloads a single file from an url to a local path
 * @type {string} url URL of the file
 * @type {string} path local filesystem path to save to
 * @return bool true on success, false on any error
 */
func (d *Download_List) DownloadFile(url string, path string) bool {

	response, e := http.Get(url)
	if e != nil {
		return false
	}
	defer response.Body.Close()

	file, err := os.Create(path)
	if err != nil {
		return false
	}

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return false
	}

	file.Close()

	return true
}

/**
 * downloads all files in "Files"
 */
func (d *Download_List) DownloadAll() {

	//bar := pb.StartNew(len(d.Files))
	bar := pb.StartNew(d.SizeTotal)
	bar.SetUnits(pb.U_BYTES)

	for _, file := range d.Files {

		fileName := strconv.FormatInt(file.Id, 10) + file.Ext
		url := BASE_URL_MEDIA + "/" + config.SearchBoard + "/" + fileName

		if(!d.DownloadFile(url, config.FilePath + "/" + fileName)) {
			//@TODO handle error
		}

		bar.Add(file.Size)
	}

	bar.FinishPrint("Done!")
}
