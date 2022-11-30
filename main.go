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

	"github.com/SevralT/GoDL/checks"
	"github.com/SevralT/GoDL/dl"
	"github.com/SevralT/GoDL/upload"
	"github.com/SevralT/GoDL/vars"
	"github.com/jeandeaual/go-locale"
	flag "github.com/spf13/pflag"
)

func main() {
	// Get OS language
	userLanguage, _ := locale.GetLanguage()

	// Locatization
	var usage, finished, check_internet_connection1, check_internet_connection2, connected, override, ver, c_type, file_download, error_url, file_not_available, file_upload string
	var versia = "1.2.0"
	if userLanguage == "ru" {
		usage = "Использование: godl [ПАРАМЕТРЫ]... [URL]... [ИМЯ ФАЙЛА (НЕОБЯЗАТЕЛЬНО)]\n"
		file_download = "\nЗагрузка файла"
		vars.File_download = "Загрузка файла"
		finished = "окончена!"
		check_internet_connection1 = "Ошибка: Не удалось поключится к "
		check_internet_connection2 = "Проверьте подключение к интернету!"
		connected = "Соединение установлено!"
		override = "Ошибка: Такой файл уже существует! Если вы хотите его перезаписать, используйте --override или -o.\n"
		ver = "Версия программы -"
		c_type = "Тип контента:"
		error_url = "Ошибка: Введён неверный URL-адрес"
		file_not_available = "Ошибка: Файл недоступен"
		file_upload = "Выгрузка файла"
	} else if userLanguage == "uk" {
		usage = "Використання: godl [ПАРАМЕТРИ]... [URL]... [ІМЯ ФАЙЛУ (НЕОБОВЯЗКО)]\n"
		vars.File_download = "Завантаження файлу"
		file_download = "\rЗавантаження файлу"
		finished = "закінчено!"
		check_internet_connection1 = "Помилка: Не вдалося підключитися до "
		check_internet_connection2 = "Перевірте підключення до інтернету!"
		connected = "З'єднання встановлено!"
		override = "Помилка: Такий файл вже існує! Якщо ви бажаєте його перезаписати, використовуйте --override чи -o.\n"
		ver = "Версія програми –"
		c_type = "Тип контенту:"
		error_url = "Помилка: Введено неправильну URL-адресу"
		file_not_available = "Помилка: Файл недоступний"
		file_upload = "Вивантаження файлу"
	} else {
		usage = "Usage: godl [OPTIONS]... [URL]... [FILE NAME (OPTIONAL)]\n"
		vars.File_download = "Downloading file"
		file_download = "\rDownloading file"
		finished = "finished!"
		check_internet_connection1 = "Error: Could not connect to "
		check_internet_connection2 = "Please check your internet connection!"
		connected = "Connection established!"
		override = "Error: Such file already exists! If you want to overwrite it, use --override or -o."
		ver = "Program version -"
		c_type = "Content Type:"
		error_url = "Error: Invalid URL entered"
		file_not_available = "Error: File not available"
		file_upload = "Uploading file"
	}

	// Set custom message on error
	flag.Usage = func() {
		fmt.Print(usage)
		flag.PrintDefaults()
		os.Exit(0)
	}

	// Set flags and args
	var filename string
	var stdout, hash, override_bool, version, upload_check bool
	flag.StringVarP(&filename, "name", "n", "", "Указать пользовательское название файла")
	flag.BoolVarP(&vars.Progress, "progress", "p", false, "Использовать прогресс-бар")
	flag.BoolVarP(&stdout, "stdout", "s", false, "Вывести содержимое через stdout")
	flag.BoolVarP(&hash, "hash", "h", false, "Посчитать sha2sum для выходного файла")
	flag.BoolVarP(&override_bool, "override", "o", false, "Перезаписывать существующий файл")
	flag.BoolVarP(&vars.QuiteMode, "quite", "q", false, "Тихий режим")
	flag.StringVarP(&vars.Username, "http-user", "", "", "Имя пользователя для HTTP-протокола")
	flag.StringVarP(&vars.Password, "http-pass", "", "", "Пароль для HTTP-протокола")
	flag.StringVarP(&vars.ProxyServer, "proxy", "", "", "Прокси-сервер")
	flag.StringVarP(&vars.ProxyPort, "proxy-port", "", "", "Порт прокси-сервера")
	flag.BoolVarP(&upload_check, "upload", "u", false, "Загрузить файл на сервер")
	flag.BoolVarP(&version, "version", "v", false, "Номер версии программы")

	flag.Parse()

	vars.FileUrl = flag.Arg(0)

	if vars.FileUrl == "" && !version {
		fmt.Print(usage)
		flag.PrintDefaults()
		os.Exit(0)
	}

	// Version number
	if version && len(os.Args) == 2 {
		fmt.Println(ver, versia)
		os.Exit(0)
	}

	// Check for args amount
	if len(os.Args) < 2 {
		fmt.Print(usage)
		flag.PrintDefaults()
		os.Exit(0)
	}

	// Set filename
	if flag.Arg(1) != "" {
		vars.FileName = flag.Arg(1)
	} else if filename == "" {
		r, _ := http.NewRequest("GET", vars.FileUrl, nil)
		vars.FileName = path.Base(r.URL.Path)
	} else {
		vars.FileName = filename
	}

	if vars.FileName == "." || vars.FileName == "" {
		fmt.Println(error_url)
		os.Exit(0)
	}

	// Check if exists
	if checks.FileExist() && !override_bool && !upload_check {
		fmt.Println(override)
		os.Exit(0)
	}

	// Check internet connection (ping url)
	if !checks.Connected() {
		fmt.Print(check_internet_connection1, checks.GetDomain(), "\n", check_internet_connection2, "\n")
		os.Exit(0)
	} else if !vars.QuiteMode && !stdout {
		println(connected)
	}

	// Check file availibility on remote server
	if !checks.CheckFileAvailability() && !vars.QuiteMode {
		fmt.Print(file_not_available)
		os.Exit(0)
	}

	// Use stdout or download file
	if !upload_check {
		if stdout {
			dl.Stdout(vars.FileUrl)
		} else {
			if !vars.QuiteMode {
				fmt.Println("IPv4:", checks.GetIP())
				contentType := checks.GetFileContentType()
				fmt.Println(c_type, contentType)
				fmt.Println()
			}
			dl.DownloadFile(vars.FileUrl, vars.FileName)
			if !vars.QuiteMode {
				fmt.Println(file_download, vars.FileName, finished)
			}
			// Optionally calculate sha2sum
			if hash {
				hasher := sha256.New()
				s, _ := ioutil.ReadFile(vars.FileName)
				hasher.Write(s)
				fmt.Println("SHA2SUM:", hex.EncodeToString(hasher.Sum(nil)))
			}
		}
	} else if upload_check {
		if !vars.QuiteMode {
			fmt.Println("IPv4:", checks.GetIP())
			fmt.Println()
		}
		upload.UploadFile(vars.FileUrl, vars.FileName)
		if !vars.QuiteMode {
			fmt.Println(file_upload, vars.FileName, finished)
		}
	}

}
