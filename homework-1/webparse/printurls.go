package webparse

import "fmt"

//SayUrls выведет на экран cписок страниц, на которых обнаружен поисковый запрос
func SayUrls(urls []string) {
	for _, url := range urls {
		fmt.Println("Результат:")
		fmt.Println(url)
	}
}
