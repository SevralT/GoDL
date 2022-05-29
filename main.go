package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path"

	dl "github.com/SevralT/GoDL/func"
	"github.com/jeandeaual/go-locale"
)

func main() {
	userLanguage, _ := locale.GetLanguage()

	var usage, finished, check_internet_connection, connected string
	if userLanguage == "ru" {
		usage = "Использование: godl [-n filename] [-p] [-stdout] URL..."
		dl.File_download = "Загрузка файла"
		finished = "окончена!"
		check_internet_connection = "Ошибка: Проверьте подключение к интернету!"
		connected = "Соединение установлено!"
	} else if userLanguage == "uk" {
		usage = "Використання: godl [-n filename] [-p] [-stdout] URL..."
		dl.File_download = "Завантаження файлу"
		finished = "закінчено!"
		check_internet_connection = "Помилка: Перевірте підключення до інтернету!"
		connected = "З'єднання встановлено!"
	} else {
		usage = "Usage: godl [-n filename] [-p] [-stdout] URL..."
		dl.File_download = "Downloading file"
		finished = "finished!"
		check_internet_connection = "Error: Check internet connection!"
		connected = "Connection established!"
	}

	if len(os.Args) < 2 {
		fmt.Println(usage)
		os.Exit(0)
	}

	var filename string
	var stdout bool
	flag.StringVar(&filename, "n", "", "Указать пользовательское название файла")
	flag.BoolVar(&dl.Progress, "p", false, "Использовать прогресс-бар")
	flag.BoolVar(&stdout, "stdout", false, "Вывести содержимое через stdout")
	flag.Parse()

	dl.FileUrl = flag.Arg(0)

	if filename == "" {
		r, _ := http.NewRequest("GET", dl.FileUrl, nil)
		dl.FileName = path.Base(r.URL.Path)
	} else {
		dl.FileName = filename
	}

	if !dl.Connected() {
		fmt.Println(check_internet_connection)
		os.Exit(0)
	} else {
		println(connected)
	}
	if stdout == true {
		dl.Stdout(dl.FileUrl)
	} else {
		dl.DownloadFile(dl.FileUrl, dl.FileName)
		fmt.Println(dl.File_download, dl.FileName, finished)
	}
}
