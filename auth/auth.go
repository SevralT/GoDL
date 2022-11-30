package auth

import (
	"encoding/base64"
	"net/http"

	"github.com/SevralT/GoDL/vars"
)

// Function for http auth
func GetAuth() http.Header {
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
