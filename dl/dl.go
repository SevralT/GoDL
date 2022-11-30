package dl

// Import required packages
import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/SevralT/GoDL/auth"
	"github.com/SevralT/GoDL/proxy"
	"github.com/SevralT/GoDL/vars"
	"github.com/schollz/progressbar/v3"
)

// Function for file download
func DownloadFile(url string, filepath string) error {
	// Download file
	client := proxy.GetProxyClient()
	req, err := http.NewRequest("GET", url, nil)
	req.Header = auth.GetAuth()
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create temp file for download
	temp, _ := os.Create(filepath + ".tmp")
	defer temp.Close()

	// Use progress bar or no
	if !vars.Progress {
		if !vars.QuiteMode {
			total := progressbar.DefaultBytes(-1, vars.File_download)
			io.Copy(io.MultiWriter(temp, total), resp.Body)
		} else {
			io.Copy(io.MultiWriter(temp), resp.Body)
		}
	} else if vars.Progress {
		bar := progressbar.DefaultBytes(resp.ContentLength, vars.File_download)
		io.Copy(io.MultiWriter(temp, bar), resp.Body)
	}

	// Rename temp file
	os.Rename(vars.FileName+".tmp", vars.FileName)
	return nil
}

// Function for stdout
func Stdout(url string) error {
	// "Download" file
	client := proxy.GetProxyClient()
	req, err := http.NewRequest("GET", url, nil)
	req.Header = auth.GetAuth()
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
