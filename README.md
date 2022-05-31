# GoDL
Программа полностью написана на языке GoLang. С помощью неё вы легко сможете загружать любые файлы на ваш компьютер.

## Возможности
* Загрузка любых типов файлов
* Поддержка вывода прогресса загрузки
* Поддержка вывода содержимого файла в stdout
* Проверка существования файла перед загрузкой
* Подсчёт sha2sum после загрузки и вывод в консоль
* Выбор имени скаченного файла или автоматическое получение имени файла из URL

## Как запустить
Для запуска GoDL, вам необходимо ввести следующую команду:
```
go run main.go [-n filename] [-p] [-stdout] [-sha2sum] [-override] URL...
```
Но для удобства, вы можете скомпилировать программу для вашей ОС или скачат её о страницы с релизами:
```
go build main.go
```
Далее, вы можете поместить бинарный файл в папку **/bin** для оперативного доступа.

## Пример использования
Для загрузки файла, вам необходимо указать адрес сайта и необходимые флаги:
```
godl -n example.html -sha2sum https://google.com/index.html
```
