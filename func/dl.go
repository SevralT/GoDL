package dl

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/schollz/progressbar/v3"
)

var FileName, FileUrl, File_download string
var Progress bool

func DownloadFile(url string, filepath string) error {
	resp, _ := http.Get(FileUrl)
	defer resp.Body.Close()

	temp, _ := os.Create(filepath + ".godl")
	defer temp.Close()

	if Progress == false {
		total := progressbar.DefaultBytes(
			-1,
			File_download,
		)
		io.Copy(io.MultiWriter(temp, total), resp.Body)
	} else if Progress == true {
		bar := progressbar.DefaultBytes(
			resp.ContentLength,
			File_download,
		)
		io.Copy(io.MultiWriter(temp, bar), resp.Body)
	}

	os.Rename(FileName+".godl", FileName)
	return nil
}

func Connected() (ok bool) {
	_, err := http.Get(FileUrl)
	if err != nil {
		return false
	}
	return true
}

func Stdout(url string) error {
	resp, _ := http.Get(FileUrl)
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(body)

	fmt.Println(bodyString)
	return nil
}
