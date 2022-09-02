package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gtuk/discordwebhook"
	"github.com/robfig/cron"
	"github.com/tidwall/gjson"
)

func main() {
	log.SetOutput(os.Stdout)

	// 引数を定義する
	epgs_url := flag.String("url", "http://localhost:8888", "EPGStation の ホスト:ポート を指定します。(例: http://your.server:8888)")
	cron_string := flag.String("cron", "@every 15s", "どのような間隔で確認するかを指定します。cron 形式を使用できます。")
	discord_webhook_url := flag.String("discord_url", "", "Discord 上の Webhook 向け URL を指定します。")
	mirakurun_msg := flag.String("mirakurun_msg", ":warning: EPGStation が Mirakurun (mirakc) バックエンドと接続できていません！", "Mirakurun (mirakc) バックエンドと接続できないときのメッセージを指定します。")
	epgs_msg := flag.String("epgs_msg", ":warning: EPGStation に接続できません！", "EPGStation と接続できないときのメッセージを指定します。")
	flag.Parse()

	// 環境変数が指定されていれば、その値を優先する
	if os.Getenv("EPGS_URL") != "" {
		*epgs_url = os.Getenv("EPGS_URL")
	}
	if os.Getenv("CRON") != "" {
		*cron_string = os.Getenv("CRON")
	}
	if os.Getenv("DISCORD_URL") != "" {
		*discord_webhook_url = os.Getenv("DISCORD_URL")
	}
	if os.Getenv("MIRAKURUN_MSG") != "" {
		*mirakurun_msg = os.Getenv("MIRAKURUN_MSG")
	}
	if os.Getenv("EPGS_MSG") != "" {
		*epgs_msg = os.Getenv("EPGS_MSG")
	}

	log.Println("[INFO] 監視を開始します。")
	log.Println("[INFO] EPGStation の宛先は " + *epgs_url + " です。")
	log.Println("[INFO] cron の設定は " + *cron_string + " です。")

	// 確認を定期実行する
	c := cron.New()
	c.AddFunc(*cron_string, func() { call_check(*epgs_url, *discord_webhook_url, *mirakurun_msg, *epgs_msg) })
	c.Start()

	// 永眠
	select {}
}

func call_check(epgs_url string, discord_webhook_url string, mirakurun_msg string, epgs_msg string) {
	// EPGStation への接続性を確認する
	_, err := http.Get(epgs_url)
	if err != nil {
		log.Println("[ERROR] EPGStation への接続に失敗しました。")
		if discord_webhook_url != "" {
			discord_epgs(discord_webhook_url, epgs_msg)
		}
	} else {
		if check(epgs_url) {
			log.Println("[INFO] EPGStation は正常に Mirakurun と接続しています。")
		} else {
			log.Println("[INFO] EPGStation は Mirakurun に接続できていません。")
			if discord_webhook_url != "" {
				discord_mirakurun(discord_webhook_url, mirakurun_msg)
			}
		}
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

func discord_mirakurun(discord_webhook_url string, content string) {
	username := "EPGSWatch"
	message := discordwebhook.Message{
		Username: &username,
		Content:  &content,
	}
	err := discordwebhook.SendMessage(discord_webhook_url, message)
	if err != nil {
		log.Println(err)
	}
}

func discord_epgs(discord_webhook_url string, content string) {
	username := "EPGSWatch"
	message := discordwebhook.Message{
		Username: &username,
		Content:  &content,
	}
	err := discordwebhook.SendMessage(discord_webhook_url, message)
	if err != nil {
		log.Println(err)
	}
}
