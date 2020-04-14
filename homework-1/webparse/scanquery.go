package webparse

import "fmt"

// GetQuery попросит ввести строку поиска
func GetQuery( query *string ) {
    fmt.Println("Введите строку поиска")
    fmt.Scanln(query)
}
