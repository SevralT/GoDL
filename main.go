package main

// Import required packages
import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
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
	var usage, finished, check_internet_connection, connected, override, ver, c_type, file_download string
	if userLanguage == "ru" {
		usage = "Использование: godl [ПАРАМЕТРЫ]... [URL]...\n"
		file_download = "\nЗагрузка файла"
		logic.File_download = "Загрузка файла"
		finished = "окончена!"
		check_internet_connection = "Ошибка: Проверьте подключение к интернету!"
		connected = "Соединение установлено!"
		override = "Ошибка: Такой файл уже существует! Если вы хотите его перезаписать, используйте --override или -o.\n"
		ver = "Версия программы - 1.0.0"
		c_type = "Тип контента:"
	} else if userLanguage == "uk" {
		usage = "Використання: godl [ПАРАМЕТРИ]... [URL]..."
		logic.File_download = "Завантаження файлу"
		file_download = "\rЗавантаження файлу"
		finished = "закінчено!"
		check_internet_connection = "Помилка: Перевірте підключення до інтернету!"
		connected = "З'єднання встановлено!"
		override = "Помилка: Такий файл вже існує! Якщо ви бажаєте його перезаписати, використовуйте --override чи -o.\n"
		ver = "Версія програми – 1.0.0"
		c_type = "Тип контенту:"
	} else {
		usage = "Usage: godl [OPTIONS]... [URL]..."
		logic.File_download = "Downloading file"
		file_download = "\rDownloading file"
		finished = "finished!"
		check_internet_connection = "Error: Check internet connection!"
		connected = "Connection established!"
		override = "Error: Such file already exists! If you want to overwrite it, use --override or -o."
		ver = "Program version - 1.0.0"
		c_type = "Content Type:"
	}

	// Set custom message on error
	flag.Usage = func() {
		fmt.Print(usage)
		flag.PrintDefaults()
		os.Exit(0)
	}

	// Set flags and args
	var filename string
	var stdout, hash, override_bool, version bool
	flag.StringVarP(&filename, "name", "n", "", "Указать пользовательское название файла")
	flag.BoolVarP(&logic.Progress, "progress", "p", false, "Использовать прогресс-бар")
	flag.BoolVarP(&stdout, "stdout", "s", false, "Вывести содержимое через stdout")
	flag.BoolVarP(&hash, "hash", "h", false, "Посчитать sha2sum для выходного файла")
	flag.BoolVarP(&override_bool, "override", "o", false, "Перезаписывать существующий файл")
	flag.BoolVarP(&logic.QuiteMode, "quite", "q", false, "Тихий режим")
	flag.BoolVarP(&version, "version", "v", false, "Номер версии программы")

	flag.Parse()

	logic.FileUrl = flag.Arg(0)

	// Version number
	if version && len(os.Args) == 2 {
		fmt.Println(ver)
		os.Exit(0)
	}

	// Check for args amount
	if len(os.Args) < 2 {
		fmt.Print(usage)
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
	if logic.Exists && !override_bool {
		fmt.Println(override)
		os.Exit(0)
	}

	// Check internet connection (ping url)
	if !logic.Connected() {
		fmt.Println(check_internet_connection)
		os.Exit(0)
	} else if !logic.QuiteMode {
		println(connected)
	}

	// Get website ip
	u, _ := url.Parse(logic.FileUrl)
	ips, _ := net.LookupIP(u.Hostname())

	// Use stdout or download file
	if stdout {
		logic.Stdout(logic.FileUrl)
	} else {
		if logic.QuiteMode == false {
			for _, ip := range ips {
				if ipv4 := ip.To4(); ipv4 != nil {
					fmt.Println("IPv4:", ipv4)
				}
			}
			contentType := logic.GetFileContentType()
			fmt.Println(c_type, contentType)
			fmt.Println()
		}
		logic.DownloadFile(logic.FileUrl, logic.FileName)
		if logic.QuiteMode == false {
			fmt.Println(file_download, logic.FileName, finished)
		}
		// Optionally calculate sha2sum
		if hash {
			hasher := sha256.New()
			s, _ := ioutil.ReadFile(logic.FileName)
			hasher.Write(s)
			fmt.Println("SHA2SUM:", hex.EncodeToString(hasher.Sum(nil)))
		}
	}
}
