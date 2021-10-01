package main

import (
	"fmt"
	"os"

	dl "github.com/SevralT/GoDL/func"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Использование: https://example.com filename.png")
		os.Exit(1)
	}

	dl.FileUrl = os.Args[1]
	dl.FileName = os.Args[2]

	fmt.Print("Загрузка файла ", dl.FileName, "...")

	err := dl.DownloadFile(dl.FileUrl, dl.FileName)
	if err != nil {
		panic(err)
	}
	fmt.Println("\rЗагрузка файла", dl.FileName, "окончена!")
}
