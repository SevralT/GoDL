package main

import (
	"fmt"
	"net/http"
	"os"
	"path"

	dl "github.com/SevralT/GoDL/func"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Использование: godl https://example.com filename.png")
		os.Exit(1)
	}

	dl.FileUrl = os.Args[1]

	if len(os.Args) == 3 {
		dl.FileName = os.Args[2]
	} else {
		r, _ := http.NewRequest("GET", dl.FileUrl, nil)
		dl.FileName = path.Base(r.URL.Path)
	}

	fmt.Print("Загрузка файла ", dl.FileName, "...")

	err := dl.DownloadFile(dl.FileUrl, dl.FileName)
	if err != nil {
		panic(err)
	}
	fmt.Println("\rЗагрузка файла", dl.FileName, "окончена!")
}
