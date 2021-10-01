package dl

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"github.com/dustin/go-humanize"
)

var FileName, FileUrl string

func DownloadFile(url string, filepath string) error {
	out, err := os.Create(filepath + ".godl")
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(FileUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	count := &WriteCounter{}
	_, err = io.Copy(out, io.TeeReader(resp.Body, count))
	if err != nil {
		return err
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	err = os.Rename(FileName+".godl", FileName)
	if err != nil {
		return err
	}
	return nil
}

type WriteCounter struct {
	Total uint64
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.PrintProgress()
	return n, nil
}

func (wc WriteCounter) PrintProgress() {
	fmt.Printf("\r%s", strings.Repeat(" ", 50))
	fmt.Printf("\rЗавершено %s", humanize.Bytes(wc.Total))
}
