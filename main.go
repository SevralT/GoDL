package main

// Import required packages
import (
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	logic "github.com/SevralT/GoDL/func"
	"github.com/jeandeaual/go-locale"
)

func main() {
	// Get OS language
	userLanguage, _ := locale.GetLanguage()

	// Locatization
	var usage, finished, check_internet_connection, connected string
	if userLanguage == "ru" {
		usage = "Использование: godl [-n filename] [-p] [-stdout] [-sha2sum] URL..."
		logic.File_download = "Загрузка файла"
		finished = "окончена!"
		check_internet_connection = "Ошибка: Проверьте подключение к интернету!"
		connected = "Соединение установлено!"
	} else if userLanguage == "uk" {
		usage = "Використання: godl [-n filename] [-p] [-stdout] [-sha2sum] URL..."
		logic.File_download = "Завантаження файлу"
		finished = "закінчено!"
		check_internet_connection = "Помилка: Перевірте підключення до інтернету!"
		connected = "З'єднання встановлено!"
	} else {
		usage = "Usage: godl [-n filename] [-p] [-stdout] [-sha2sum] URL..."
		logic.File_download = "Downloading file"
		finished = "finished!"
		check_internet_connection = "Error: Check internet connection!"
		connected = "Connection established!"
	}

	// Check for args amount
	if len(os.Args) < 2 {
		fmt.Println(usage)
		os.Exit(0)
	}

	// Set flags and args
	var filename string
	var stdout, sha2sum bool
	flag.StringVar(&filename, "n", "", "Указать пользовательское название файла")
	flag.BoolVar(&logic.Progress, "p", false, "Использовать прогресс-бар")
	flag.BoolVar(&stdout, "stdout", false, "Вывести содержимое через stdout")
	flag.BoolVar(&sha2sum, "sha2sum", false, "Посчитать sha2sum для выходного файла")
	flag.Parse()

	logic.FileUrl = flag.Arg(0)

	// Set filename
	if filename == "" {
		r, _ := http.NewRequest("GET", logic.FileUrl, nil)
		logic.FileName = path.Base(r.URL.Path)
	} else {
		logic.FileName = filename
	}

	// Check internet connection (ping url)
	if !logic.Connected() {
		fmt.Println(check_internet_connection)
		os.Exit(0)
	} else {
		println(connected)
	}

	// Use stdout or download file
	if stdout == true {
		logic.Stdout(logic.FileUrl)
	} else {
		logic.DownloadFile(logic.FileUrl, logic.FileName)
		fmt.Println(logic.File_download, logic.FileName, finished)
		// Optionally calculate sha2sum
		if sha2sum == true {
			hasher := sha256.New()
			s, _ := ioutil.ReadFile(logic.FileName)
			hasher.Write(s)
			fmt.Println("SHA2SUM:", hex.EncodeToString(hasher.Sum(nil)))
		}
	}
}
