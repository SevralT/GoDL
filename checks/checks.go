package checks

import (
	"net/http"
	"os"

	"github.com/SevralT/GoDL/dl"
)

// Function for check internet connections
func Connected() (ok bool) {
	// "Download file"
	_, err := http.Get(dl.FileUrl)
	if err != nil {
		return false
	}
	return true
}

// Check if file exists
func FileExist(filename string) error {
	_, err := os.Stat(dl.FileName)
	if os.IsNotExist(err) {
		dl.Exists = false
	} else {
		dl.Exists = true
	}
	return nil
}

// Detect content type
func GetFileContentType() (filetype string) {
	client := http.Client{}
	req, _ := http.NewRequest("HEAD", dl.FileUrl, nil)
	req.Header.Set("Accept", "*/*")
	resp, _ := client.Do(req)
	return resp.Header.Get("Content-Type")
}
