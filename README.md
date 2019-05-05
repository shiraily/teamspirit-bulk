# teamspirit-bulk

teamspirit-bulkは、IFTTTでGoogleスプレッドシートに記録した特定のWi-Fiへの接続・切断時刻を参照して、
一括でTeamSpiritに勤怠時刻を入力するためのプログラムです。

## 要件

- AndroidのIFTTTアプリ
  - iPhoneは未検証
- Chrome
- [ChromeDriver](https://qiita.com/tenten0213/items/1f897ff8a64bd8b5270c)
- Golang

## 機能

- GoogleスプレッドシートをAPI経由で読み取ってTeamSpirit用の入力データを生成
- ChromeDriverでTeamSpiritを操作し勤怠を入力

## 設定

### IFTTT

以下手順で出勤時刻・退勤時刻用のアクションをそれぞれ登録します。

New Applet>Android Device>Connects to a specific WiFi networkから

- this
  - Network name: 接続・切断を記録したいWi-FiのSSID
- that
  - Select action service: Google Sheets>Add row to spreadsheet
  - Formatted row: "in ||| OccurredAt"
    - 出勤の場合はin, 退勤の場合はout
    - OccurredAtはもともと入力してある白抜きのテキスト

<img src="https://user-images.githubusercontent.com/25303121/56885617-d52ba200-6aa7-11e9-9be8-5223ad86c372.png" width="400">

### Googleスプレッドシート

IFTTTで記録したデータの中から今月の最初のデータの行番号を取得するために以下設定が必要です。

- IFTTTで設定したスプレッドシートに任意の名前でシートを追加（あとでシート名をyamlのsheet_settingに記入します）
  - 事前にIFTTTを発火させる必要。
  （未検証）または自分で同様にスプレッドシートを作成
- そのシートのA1セルに以下数式を入力
  - `=match(choose(month(today()),"January","February","March","April","May","June","July","August","September","October","November","December")&"*",time!B1:B500218,0)`
  - ↑の"time"はIFTTTで設定したシート名

### スプレッドシートへのアクセス権限付与

- GCPプロジェクトを作成
- [APIとサービス](https://console.cloud.google.com/apis/dashboard)>APIとサービスを有効化
  - Google Sheets API
- [APIとサービス>認証情報](https://console.cloud.google.com/apis/credentials)>認証情報を作成>サービスアカウントキー
  - 以下画像のように設定

<img src="https://user-images.githubusercontent.com/25303121/56885439-69493980-6aa7-11e9-964e-50377871ebf7.png">

- jsonファイルをダウンロード後保管（あとで使います）
- スプレッドシートを開く>共有>サービスアカウントIDを共有先に設定
  - サービスアカウントIDは、キーの作成時または[GCP>IAMと管理](https://console.cloud.google.com/iam-admin/iam)で確認できます

### yamlファイル

config/sample.yamlを参考に設定ファイルを作成します。
- sheet_id: スプレッドシートのID（URL `https://docs.google.com/spreadsheets/d/xxx/`のxxx部分です）
- sheet_work_time: IFTTTで設定したシート名
- sheet_setting: 自分で追加した設定用のシート名
- client_secret: さきほどダウンロードしたjsonの値を貼り付けます。


## 使い方

以下を実行すると、実行した現在時刻の月分の勤怠が自動入力されます。

```
$ git clone git@github.com:shiraily/teamspirit-bulk.git
$ cd /path/to/teamspirit-bulk
$ go get ./...
$ go run main.go /path/to/my.yaml
```

または、main.goを参考にカスタマイズもできます。

## Future work

- csv出力に対応
- 休憩時間がない場合の警告ポップアップへの対応
- 休日対応
- IFTTT以外の安定的な出退勤情報の利用
- 先月分のデータ入力
- 出勤打刻は一括入力ではなく出社打刻ボタンを押す
