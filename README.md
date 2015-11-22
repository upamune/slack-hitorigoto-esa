# slack-hitorigoto-esa


## Description

ここの分報を1日分まとめて[esa.io](https://esa.io/) に日報として投げるやつです.
![スクリーンショット](https://i.gyazo.com/e63a9e911003d74703f60faca85f92ad.png) みたいな感じになります. cronとかで毎日実行してあげると良さそうです.

FYI: [Slackで簡単に「日報」ならぬ「分報」をチームで実現する3ステップ 〜 Problemが10分で解決するチャットを作ろう](http://c16e.com/1511101558/)
## Usage

```bash
$ slack-hitorigoto-esa -c config.toml
```

設定ファイルはTOML形式で `c` オプションで設定ファイルのパスを指定します.

```toml:config.toml
[slack]
token = "access_token"
channel = "hitorigoto"

[esa]
token = "access_token"
team = "team_name"
category = "日報"
```

## Install

```bash
$ go get github.com/upamune/slack-hitorigoto-esa
```

[Releases](https://github.com/upamune/slack-hitorigoto-esa/releases)

バイナリを直接ダウンロードすることもできます.

## Contribution

1. Fork ([https://github.com/upamune/slack-hitorigoto-esa/fork](https://github.com/upamune/slack-hitorigoto-esa/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test ./...` command and confirm that it passes
1. Run `gofmt -s`
1. Create a new Pull Request

## Author

[upamune](https://github.com/upamune)
