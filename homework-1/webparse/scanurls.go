package webparse

import "fmt"

// GetUrls попросит ввести строку c urls (разделитель запятая без пробелов)
func GetUrls( urls *string ) {
    fmt.Println("Введите строку со списком urls (разделитель: запятая без пробелов)")
    fmt.Scanln(urls)
}
