package dl

// Import required packages
import (
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/schollz/progressbar/v3"
)

// Set global variables
var FileName, FileUrl, File_download string
var Username, Password string
var Progress, Exists, QuiteMode bool

// Function for file download
func DownloadFile(url string, filepath string) error {
	// Download file
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header = getAuth()
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create temp file for download
	temp, _ := os.Create(filepath + ".tmp")
	defer temp.Close()

	// Use progress bar or no
	if !Progress {
		if !QuiteMode {
			total := progressbar.DefaultBytes(
				-1,
				File_download,
			)
			io.Copy(io.MultiWriter(temp, total), resp.Body)
		} else {
			io.Copy(io.MultiWriter(temp), resp.Body)
		}
	} else if Progress {
		bar := progressbar.DefaultBytes(
			resp.ContentLength,
			File_download,
		)
		io.Copy(io.MultiWriter(temp, bar), resp.Body)
	}

	// Rename temp file
	os.Rename(FileName+".tmp", FileName)
	return nil
}

// Function for stdout
func Stdout(url string) error {
	// "Download" file
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header = getAuth()
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	// Read file
	defer resp.Body.Close()

	// Convert response in string
	body, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(body)

	// Output file content
	fmt.Println(bodyString)
	return nil
}

// Function for http auth
func getAuth() http.Header {
	// Check if username and password are set
	if Username != "" && Password != "" {
		// Create auth header
		auth := Username + ":" + Password
		authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
		header := http.Header{}
		header.Add("Authorization", authHeader)
		return header
	} else {
		return nil
	}
}