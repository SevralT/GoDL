package checks

import (
	"net"
	"net/http"
	"net/url"
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
func FileExist() bool {
	_, err := os.Stat(dl.FileName)
	if os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

// Detect content type
func GetFileContentType() (filetype string) {
	client := http.Client{}
	req, _ := http.NewRequest("HEAD", dl.FileUrl, nil)
	req.Header.Set("Accept", "*/*")
	resp, _ := client.Do(req)
	return resp.Header.Get("Content-Type")
}

func GetDomain() (domain string) {
	u, _ := url.Parse(dl.FileUrl)
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
