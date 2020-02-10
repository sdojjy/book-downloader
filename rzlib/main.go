package main

import (
	"compress/gzip"
	"flag"
	"fmt"
	"github.com/anaskhan96/soup"
	js "github.com/dop251/goja"
	"golang.org/x/net/html/charset"
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
	for link != "" && strings.HasSuffix(link, ".html") {
		url := "https://m.rzlib.net" + link
		link = printContent(url)
	}
	panic("没有找到更多章节, 下载完成！")
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
		if err != nil {
			print(err)
			continue
		}

		body := response.Body
		if response.Header.Get("Content-Encoding") == "gzip" {
			body, err = gzip.NewReader(response.Body)
			if err != nil {
				fmt.Println("http resp unzip is failed,err: ", err)
			}
		}

		r, err := charset.NewReader(body, response.Header.Get("Content-Type"))
		//r, err := iconv.NewReader(body, "gb2312", "utf-8")
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
	textURL := getTextURL(baseURL)
	_, err := vm.RunString(string(getHtml(textURL)))
	if err != nil {
		panic(err)
	}
	eval := vm.Get("cctxt").String()

	eval = strings.ReplaceAll(eval, "<br />", "\n")
	eval = strings.ReplaceAll(eval, "&nbsp;", "")

	chapterName := doc.Find("h1", "id", "chaptername").Text()
	fmt.Printf("\n\n%s\n\n%s\n", chapterName, eval)

	println(chapterName)
	link := doc.Find("a", "id", "pb_next").Attrs()["href"]
	return link
}

func getTextURL(htmlURL string) string {
	//https://m.rzlib.net/b/52/52352/23791490.html
	index := strings.Index(htmlURL, ".html")
	htmlURL = htmlURL[0:index]
	ps := strings.Split(htmlURL, "/")
	c := ps[len(ps)-1]
	d := ps[len(ps)-2]
	return fmt.Sprintf("https://www.rzlib.net/b/txtg333/%s/%s.txt", d, c)
}

func deferClose(c io.Closer) {
	if err := c.Close(); err != nil {
		print(err)
	}
}
