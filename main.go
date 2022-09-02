package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/tidwall/gjson"
)

func main() {
	// 引数を定義する
	url := flag.String("url", "http://localhost:8888", "EPGStation の ホスト:ポート を指定します。(例: http://your.server:8888)")
	flag.Parse()

	// EPGStation への接続性を確認する
	_, err := http.Get(*url)
	if err != nil {
		log.Fatal("[ERROR] EPGStation への接続に失敗しました。")
		os.Exit(1)
	}

	if check(*url) {
		log.Println("[INFO] EPGStation は正常に Mirakurun と接続しています。")
	} else {
		log.Println("[INFO] EPGStation は Mirakurun に接続できていません。")
	}

}

func check(base_url string) bool {
	url := base_url + "/api/streams/live/" + get_channel(base_url) + "/m2ts?mode=2"
	client := http.Client{
		Timeout: 2 * time.Second,
	}
	res, err := client.Get(url)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	buf := bytes.NewBuffer(body)
	mimeType := http.DetectContentType(buf.Bytes())

	if mimeType == "application/octet-stream" {
		return true
	} else {
		return false
	}
}

// 適当なチャンネルを自動的に取得する
func get_channel(base_url string) string {
	url := base_url + "/api/channels"
	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	buf := bytes.NewBuffer(body)

	value := gjson.Get(buf.String(), "0.id")

	return value.String()
}
