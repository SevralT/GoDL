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

	"github.com/SevralT/GoDL/checks"
	"github.com/SevralT/GoDL/dl"
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
		dl.File_download = "Загрузка файла"
		finished = "окончена!"
		check_internet_connection = "Ошибка: Проверьте подключение к интернету!"
		connected = "Соединение установлено!"
		override = "Ошибка: Такой файл уже существует! Если вы хотите его перезаписать, используйте --override или -o.\n"
		ver = "Версия программы - 1.0.0"
		c_type = "Тип контента:"
	} else if userLanguage == "uk" {
		usage = "Використання: godl [ПАРАМЕТРИ]... [URL]..."
		dl.File_download = "Завантаження файлу"
		file_download = "\rЗавантаження файлу"
		finished = "закінчено!"
		check_internet_connection = "Помилка: Перевірте підключення до інтернету!"
		connected = "З'єднання встановлено!"
		override = "Помилка: Такий файл вже існує! Якщо ви бажаєте його перезаписати, використовуйте --override чи -o.\n"
		ver = "Версія програми – 1.0.0"
		c_type = "Тип контенту:"
	} else {
		usage = "Usage: godl [OPTIONS]... [URL]..."
		dl.File_download = "Downloading file"
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
	flag.BoolVarP(&dl.Progress, "progress", "p", false, "Использовать прогресс-бар")
	flag.BoolVarP(&stdout, "stdout", "s", false, "Вывести содержимое через stdout")
	flag.BoolVarP(&hash, "hash", "h", false, "Посчитать sha2sum для выходного файла")
	flag.BoolVarP(&override_bool, "override", "o", false, "Перезаписывать существующий файл")
	flag.BoolVarP(&dl.QuiteMode, "quite", "q", false, "Тихий режим")
	flag.BoolVarP(&version, "version", "v", false, "Номер версии программы")

	flag.Parse()

	dl.FileUrl = flag.Arg(0)

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
		r, _ := http.NewRequest("GET", dl.FileUrl, nil)
		dl.FileName = path.Base(r.URL.Path)
	} else {
		dl.FileName = filename
	}

	// Check if exists
	checks.FileExist(dl.FileName)
	if dl.Exists && !override_bool {
		fmt.Println(override)
		os.Exit(0)
	}

	// Check internet connection (ping url)
	if !checks.Connected() {
		fmt.Println(check_internet_connection)
		os.Exit(0)
	} else if !dl.QuiteMode {
		println(connected)
	}

	// Get website ip
	u, _ := url.Parse(dl.FileUrl)
	ips, _ := net.LookupIP(u.Hostname())

	// Use stdout or download file
	if stdout {
		dl.Stdout(dl.FileUrl)
	} else {
		if dl.QuiteMode == false {
			for _, ip := range ips {
				if ipv4 := ip.To4(); ipv4 != nil {
					fmt.Println("IPv4:", ipv4)
				}
			}
			contentType := checks.GetFileContentType()
			fmt.Println(c_type, contentType)
			fmt.Println()
		}
		dl.DownloadFile(dl.FileUrl, dl.FileName)
		if dl.QuiteMode == false {
			fmt.Println(file_download, dl.FileName, finished)
		}
		// Optionally calculate sha2sum
		if hash {
			hasher := sha256.New()
			s, _ := ioutil.ReadFile(dl.FileName)
			hasher.Write(s)
			fmt.Println("SHA2SUM:", hex.EncodeToString(hasher.Sum(nil)))
		}
	}
}
