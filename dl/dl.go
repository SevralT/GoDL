package dl

// Import required packages
import (
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"github.com/schollz/progressbar/v3"
	"github.com/SevralT/GoDL/vars"
)

// Function for file download
func DownloadFile(url string, filepath string) error {
	// Download file
	client := getClient()
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
	client := getClient()
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
	if vars.Username != "" && vars.Password != "" {
		// Create auth header
		auth := vars.Username + ":" + vars.Password
		authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
		header := http.Header{}
		header.Add("Authorization", authHeader)
		return header
	} else {
		return nil
	}
}

// Function for proxy server
func getClient() *http.Client {
	// Check if proxy server is setted
	if vars.ProxyServer != "" {
		// Check if proxy port is setted
		if vars.ProxyPort != "" {
			// Create proxy url
			proxyUrl, _ := url.Parse("http://" + vars.ProxyServer + ":" + vars.ProxyPort)
			// Create client
			client := &http.Client{
				Transport: &http.Transport{
					Proxy: http.ProxyURL(proxyUrl),
				},
			}
			return client
		} else {
			// Create proxy url
			proxyUrl, _ := url.Parse("http://" + vars.ProxyServer + ":8080")
			// Create client
			client := &http.Client{
				Transport: &http.Transport{
					Proxy: http.ProxyURL(proxyUrl),
				},
			}
			return client
		}
	} else {
		// Create client
		client := &http.Client{}
		return client
	}
}