package main

// Import required packages
import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	logic "github.com/SevralT/GoDL/func"
	"github.com/jeandeaual/go-locale"
	flag "github.com/spf13/pflag"
)

func main() {
	// Get OS language
	userLanguage, _ := locale.GetLanguage()

	// Locatization
	var usage, finished, check_internet_connection, connected, override string
	if userLanguage == "ru" {
		usage = "Использование: godl [ПАРАМЕТРЫ]... [URL]...\n"
		logic.File_download = "Загрузка файла"
		finished = "окончена!"
		check_internet_connection = "Ошибка: Проверьте подключение к интернету!"
		connected = "Соединение установлено!"
		override = "Ошибка: Такой файл уже существует! Если вы хотите его перезаписать, используйте -override.\n"
	} else if userLanguage == "uk" {
		usage = "Використання: godl [ПАРАМЕТРИ]... [URL]..."
		logic.File_download = "Завантаження файлу"
		finished = "закінчено!"
		check_internet_connection = "Помилка: Перевірте підключення до інтернету!"
		connected = "З'єднання встановлено!"
		override = "Помилка: Такий файл вже існує! Якщо ви бажаєте його перезаписати, використовуйте -override.\n"
	} else {
		usage = "Usage: godl [OPTIONS]... [URL]..."
		logic.File_download = "Downloading file"
		finished = "finished!"
		check_internet_connection = "Error: Check internet connection!"
		connected = "Connection established!"
		override = "Error: Such file already exists! If you want to overwrite it, use --override or -o."
	}

	// Set custom message on error
	flag.Usage = func() {
		fmt.Printf(usage)
		flag.PrintDefaults()
		os.Exit(0)
	}

	// Set flags and args
	var filename string
	var stdout, hash, override_bool bool
	flag.StringVarP(&filename, "name", "n", "", "Указать пользовательское название файла")
	flag.BoolVarP(&logic.Progress, "progress", "p", false, "Использовать прогресс-бар")
	flag.BoolVarP(&stdout, "stdout", "s", false, "Вывести содержимое через stdout")
	flag.BoolVarP(&hash, "hash", "h", false, "Посчитать sha2sum для выходного файла")
	flag.BoolVarP(&override_bool, "override", "o", false, "Перезаписывать существующий файл")

	flag.Parse()

	logic.FileUrl = flag.Arg(0)

	// Check for args amount
	if len(os.Args) < 2 {
		fmt.Printf(usage)
		flag.PrintDefaults()
		os.Exit(0)
	}

	// Set filename
	if filename == "" {
		r, _ := http.NewRequest("GET", logic.FileUrl, nil)
		logic.FileName = path.Base(r.URL.Path)
	} else {
		logic.FileName = filename
	}

	// Check if exists
	logic.FileExist(logic.FileName)
	if logic.Exists == true && override_bool == false {
		fmt.Println(override)
		os.Exit(0)
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
		if hash == true {
			hasher := sha256.New()
			s, _ := ioutil.ReadFile(logic.FileName)
			hasher.Write(s)
			fmt.Println("SHA2SUM:", hex.EncodeToString(hasher.Sum(nil)))
		}
	}
}
