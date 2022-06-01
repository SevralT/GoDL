package dl

// Import required packages
import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/schollz/progressbar/v3"
)

// Set global variables
var FileName, FileUrl, File_download string
var Progress, Exists, QuiteMode bool

// Function for file download
func DownloadFile(url string, filepath string) error {
	// Download file
	resp, _ := http.Get(FileUrl)
	defer resp.Body.Close()

	// Create temp file for download
	temp, _ := os.Create(filepath + ".godl")
	defer temp.Close()

	// Use progress bar or no
	if Progress == false {
		if QuiteMode == false {
			total := progressbar.DefaultBytes(
				-1,
				File_download,
			)
			io.Copy(io.MultiWriter(temp, total), resp.Body)
		} else {
			io.Copy(io.MultiWriter(temp), resp.Body)
		}
	} else if Progress == true {
		bar := progressbar.DefaultBytes(
			resp.ContentLength,
			File_download,
		)
		io.Copy(io.MultiWriter(temp, bar), resp.Body)
	}

	// Rename temp file
	os.Rename(FileName+".godl", FileName)
	return nil
}

// Function for stdout
func Stdout(url string) error {
	// "Download" file
	resp, _ := http.Get(FileUrl)
	defer resp.Body.Close()

	// Convert response in string
	body, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(body)

	// Output file content
	fmt.Println(bodyString)
	return nil
}
