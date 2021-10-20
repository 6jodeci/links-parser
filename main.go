package main

import (
	"fmt"
	"net/http"
	"os"

	"golang.org/x/net/html"
)

func main() {
	//все параметры переданные в командную строку
	for _, url := range os.Args[1:] {
		//все полученные ссылки
		links, err := findLinks(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "parse: %v\n", err)
		}
		for _, link := range links {
			fmt.Println(link)
		}
	}
}

func findLinks(url string) ([]string, error) {
	//запись ответа
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	//проверка на ответ сервера
	if resp.StatusCode != http.StatusOK {
		//закрытие соединения
		resp.Body.Close()
		return nil, fmt.Errorf("getting %s: %s", url, resp.Status)
	}
	//возвраm документа
	doc, err := html.Parse(resp.Body)
	//закрытие докумета
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("parsing %s as HTML: %v", url, err)
	}
	return visit(nil, doc), nil
}

func visit(links []string, n *html.Node) []string {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			//если атрибурт == href - то это ссылка
			if a.Key == "href" {
				links = append(links, a.Val)
			}
		}
	}
	//перебор всего HTML
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = visit(links, c)
	}
	return links
}