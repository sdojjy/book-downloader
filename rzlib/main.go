package main

import (
	"compress/gzip"
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
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
		req.Header.Set("Accept-Encoding", "gzip, deflate, br")
		req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7,zh-TW;q=0.6,fr;q=0.5")

		response, err := client.Do(req)
		body := response.Body
		if response.Header.Get("Content-Encoding") == "gzip" {
			body, err = gzip.NewReader(response.Body)
			if err != nil {
				fmt.Println("http resp unzip is failed,err: ", err)
			}
		}

		if err != nil {
			print(err)
			continue
		}

		//reader := transform.NewReader(body, simplifiedchinese.HZGB2312.NewEncoder())
		//d, e := ioutil.ReadAll(reader)
		//if e != nil {
		//	panic(d)
		//}
		//deferClose(body)
		//return d

		r, err := iconv.NewReader(body, "gb2312", "utf-8")
		if err != nil {
			deferClose(body)
			panic(err)
		}
		out, err := ioutil.ReadAll(r)
		if err != nil {
			deferClose(body)
			panic(err)
		}
		deferClose(body)
		return out
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
