package main

import (
	"flag"
	"fmt"
	"github.com/anaskhan96/soup"
	"github.com/djimenez/iconv-go"
	"io/ioutil"
	"net/http"
	"strings"
)

var client = &http.Client{}

var firstURL string

func init() {
	flag.StringVar(&firstURL, "url", "https://m.rzlib.net/b/52/52352/23791490.html", "the first page")
	flag.Parse()
}
func main() {
	link := printContent(firstURL)
	for link != "" {
		url := "https://m.rzlib.net" + link
		link = printContent(url)
	}
}

func getHtml(url string) []byte {
	req, err := http.NewRequest("GET", url, nil)
	req.Close = true
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:66.0) Gecko/20100101 Firefox/66.0")

	response, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer response.Body.Close()

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(content)
	}

	out := make([]byte, len(content)*2)
	out = out[:]
	iconv.Convert(content, out, "gb2312", "utf-8")
	return out
}

func printContent(baseURL string) string {
	out := getHtml(baseURL)
	doc := soup.HTMLParse(string(out))

	chaptername := doc.Find("h1", "id", "chaptername").Text()
	fmt.Print(chaptername)

	//root := doc.Find("div", "id", "txt")
	//for index, c := range root.Children() {
	//	if index > 2 && c.NodeValue != "br" && c.NodeValue != "script"{
	//		print(strings.TrimSpace(c.NodeValue))
	//	}
	//	print(c.Text())
	//}
	textURL := getTextURL(baseURL)
	//getHtml(textURL)
	strs := strings.Split(string(getHtml(textURL)), "\n")
	//println(strs[0])
	content := strs[0]
	for i := 1; i < 5; i++ {
		regParam := strings.Split(strs[i], ",")
		arg1 := strings.ReplaceAll(regParam[0], "cctxt=cctxt.replace(/", "")
		arg1 = strings.ReplaceAll(arg1, "/g", "")
		arg2 := strings.ReplaceAll(regParam[1], "'", "")
		arg2 = strings.ReplaceAll(arg2, ");\r", "")
		//println(arg1)
		//println(arg2)
		content = strings.ReplaceAll(content, arg1, arg2)
	}
	content = strings.ReplaceAll(content, "var cctxt='", "")
	content = strings.ReplaceAll(content, "<br />", "\n")
	content = strings.ReplaceAll(content, "&nbsp;", "")
	content = strings.ReplaceAll(content, "';", "")
	fmt.Print(content)
	fmt.Println()

	link := doc.Find("a", "id", "pb_next").Attrs()["href"]
	return link
}

func getTextURL(htmlURL string) string {
	//https://m.rzlib.net/b/52/52352/23791490.html
	textURL := strings.ReplaceAll(htmlURL, ".html", ".txt")
	return strings.ReplaceAll(textURL, "https://m.rzlib.net/b/103", "https://www.rzlib.net/b/txtg333")
}
