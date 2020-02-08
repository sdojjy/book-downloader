package main

import (
	"flag"
	"fmt"
	"github.com/anaskhan96/soup"
	"github.com/djimenez/iconv-go"
	js "github.com/dop251/goja"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

var client = &http.Client{}

var firstURL string

var vm = js.New()

func init() {
	flag.StringVar(&firstURL, "url", "https://m.rzlib.net/b/103/103300/48601197.html", "the first page")
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
	for {
		req, err := http.NewRequest("GET", url, nil)
		req.Close = true
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:66.0) Gecko/20100101 Firefox/66.0")

		response, err := client.Do(req)
		if err != nil {
			print(err)
			continue
		}
		content, err := ioutil.ReadAll(response.Body)
		if err != nil {
			print(err)
			deferClose(response.Body)
			continue
		}

		out := make([]byte, len(content)*2)
		out = out[:]
		_, w, err := iconv.Convert(content, out, "gb2312", "utf-8")
		if err != nil {
			continue
		}
		deferClose(response.Body)
		return out[:w]
	}
}

func printContent(baseURL string) string {
	out := getHtml(baseURL)
	doc := soup.HTMLParse(string(out))

	chaptername := doc.Find("h1", "id", "chaptername").Text()
	textURL := getTextURL(baseURL)
	//strs := strings.Split(string(getHtml(textURL)), "\n")
	//content := fmt.Sprintf("%s\n%s", chaptername, strs[0])
	//for i := 1; i < len(strs)-1; i++ {
	//	regParam := strings.Split(strs[i], ",")
	//	arg1 := strings.ReplaceAll(regParam[0], "cctxt=cctxt.replace(/", "")
	//	arg1 = strings.ReplaceAll(arg1, "/g", "")
	//	arg2 := strings.ReplaceAll(regParam[1], "'", "")
	//	arg2 = strings.ReplaceAll(arg2, ");\r", "")
	//	content = strings.ReplaceAll(content, arg1, arg2)
	//}
	//content = strings.ReplaceAll(content, "var cctxt='", "")
	_, err := vm.RunString(string(getHtml(textURL)))
	if err != nil {
		panic(err)
	}
	eval := vm.Get("cctxt").String()

	eval = strings.ReplaceAll(eval, "<br />", "\n")
	eval = strings.ReplaceAll(eval, "&nbsp;", "")
	eval = strings.ReplaceAll(eval, "';", "")
	fmt.Printf("%s\n%s", chaptername, eval)
	fmt.Println()

	link := doc.Find("a", "id", "pb_next").Attrs()["href"]
	return link
}

func getTextURL(htmlURL string) string {
	//https://m.rzlib.net/b/52/52352/23791490.html
	textURL := strings.ReplaceAll(htmlURL, ".html", ".txt")
	//index := strings.Index(htmlURL, ".html")
	return strings.ReplaceAll(textURL, "https://m.rzlib.net/b/103", "https://www.rzlib.net/b/txtg333")
}

func deferClose(c io.Closer) {
	if err := c.Close(); err != nil {
		print(err)
	}
}
