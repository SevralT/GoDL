package checks

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/SevralT/GoDL/http/auth"
	"github.com/SevralT/GoDL/vars"
)

// Function for check internet connections
func Connected() (ok bool) {
	// "Download file"
	_, err := http.Get(vars.FileUrl)
	if err != nil {
		return false
	}
	return true
}

// Check if file exists
func FileExist() bool {
	_, err := os.Stat(vars.FileName)
	if os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

// Detect content type
func GetFileContentType() (filetype string) {
	client := http.Client{}
	req, _ := http.NewRequest("HEAD", vars.FileUrl, nil)
	req.Header.Set("Accept", "*/*")
	resp, _ := client.Do(req)
	return resp.Header.Get("Content-Type")
}

func GetDomain() (domain string) {
	u, _ := url.Parse(vars.FileUrl)
	return u.Hostname()
}

func GetIP() (ip net.IP) {
	// Get website ip
	ips, _ := net.LookupIP(GetDomain())

	for _, ip := range ips {
		if ipv4 := ip.To4(); ipv4 != nil {
			return ipv4
		}
	}
	return nil
}

func CheckFileAvailability() (ok bool) {
	client := &http.Client{}
	req, _ := http.NewRequest("HEAD", vars.FileUrl, nil)
	req.Header.Set("Accept", "*/*")
	req.Header = auth.GetAuth()
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		return true
	} else {
		return false
	}
}

func GetRedirect() (redirect string) {
	client := &http.Client{}
	req, _ := http.NewRequest("HEAD", vars.FileUrl, nil)
	req.Header.Set("Accept", "*/*")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		return ""
	} else {
		dump, _ := httputil.DumpResponse(resp, true)
		return strings.Split(string(dump), "\n")[7]
	}
}
