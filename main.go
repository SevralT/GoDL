package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path"

	dl "github.com/SevralT/GoDL/func"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Использование: godl [-u url] [-n filename]")
		os.Exit(0)
	}

	var filename string
	flag.StringVar(&dl.FileUrl, "u", "nil", "Указать URL-адрес")
	flag.StringVar(&filename, "n", "nil", "Указать пользовательское название файла")
	flag.Parse()

	if dl.FileUrl == "nil" && filename == "nil" {
		dl.FileUrl = os.Args[1]
		r, _ := http.NewRequest("GET", dl.FileUrl, nil)
		dl.FileName = path.Base(r.URL.Path)
	}

	if dl.FileUrl == "nil" {
		fmt.Println("Использование: godl [-u url] [-n filename]")
		os.Exit(0)
	}

	if filename == "nil" {
		r, _ := http.NewRequest("GET", dl.FileUrl, nil)
		dl.FileName = path.Base(r.URL.Path)
	} else {
		dl.FileName = filename
	}

	fmt.Print("Загрузка файла ", dl.FileName, "...")

	err := dl.DownloadFile(dl.FileUrl, dl.FileName)
	if err != nil {
		panic(err)
	}
	fmt.Println("\rЗагрузка файла", dl.FileName, "окончена!")
}
