# EPGSWatcher
EPGStation と Mirakurun (mirakc) 間の接続性を監視

## 機能
* cron 形式の間隔指定を解釈し、それにしたがって監視する
* Discord への Webhook を経由した通知

## 引数 / 環境変数について
* コマンドライン引数は、以下のように指定します。
    ```shell
    $ ./EPGSWatcher -<key> <value>    
    ```

    | 環境変数 | コマンドライン引数 | 説明 |
    | - | - | - |
    | `EPGS_URL` | `url` | EPGStation への URL <br>例 / 既定値: `http://localhost:8888` |
    | `CRON` | `cron` | cron 形式の間隔指定 <br>例: `@daily`, `0 30 * * * *`, ... <br> 既定値: `@every 15s`|
    | `DISCORD_URL` | `discord_url` | Discord の Webhook URL を指定 <br> 既定値: (空欄) <br> 空文字の場合、通知処理は行われません。 |
    | `MIRAKURUN_MSG` | `mirakurun_msg` | Mirakurun と接続できないときの警告メッセージを指定 <br> 既定値: `:warning: EPGStation が Mirakurun (mirakc) バックエンドと接続できていません！` |
    | `EPGS_MSG` | `epgs_msg` | EPGStation と接続できないときの警告メッセージを指定 <br> 既定値: `:warning: EPGStation に接続できません！` |

## セットアップ
### バイナリを直接実行する場合
```shell
# ビルド
$ go build .
# 実行
$ ./EPGSWatcher -url http://your.server:port -discord https://discord.com/api/webhooks/xxx
```

### Docker Compose を使用する場合
* 以下の内容を `docker-compose.yaml` として保存し、環境変数を編集してください。
```yaml
version: '3.9'

services:
  epgswatcher:
    image: ghcr.io/yude/epgswatcher:master
    build: .
    restart: unless-stopped
    environment:
      EPGS_URL: ""
      CRON: ""
      DISCORD_URL: ""
      MIRAKURUN_MSG: ""
      EPGS_MSG: ""
```
* `docker compose up -d` で起動します。

## ライセンス
MIT
