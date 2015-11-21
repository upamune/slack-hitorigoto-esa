# slack-hitorigoto-esa

## Description

ここの分報を1日分まとめて[esa.io](https://esa.io/) に日報として投げるやつです.
![スクリーンショット]() みたいな感じになります. cronとかで毎日実行してあげると良さそうです.

FYI: [Slackで簡単に「日報」ならぬ「分報」をチームで実現する3ステップ 〜 Problemが10分で解決するチャットを作ろう](http://c16e.com/1511101558/)
## Usage

```bash
$ slack-hitorigoto-esa -s SLACK_TOKEN -e ESA_TOKEN
```

設定ファイルを読みこませることもできます. 設定ファイルはTOML形式で `c` オプションで設定ファイルの名前を指定します.

```toml:config.toml
SLACK_TOKEN = "SLACK_TOKEN"
ESA_TOKEN = "ESA_TOKEN"
```

```bash
$ slack-hitorigoto-esa -c ./config.toml
```

設定ファイルとオプションでTOKENを同時に指定した場合はオプションで指定した方が優先されます.

```toml:config.toml
SLACK_TOKEN = "SLACK_TOKEN2"
ESA_TOKEN = "ESA_TOKEN2"
```

```bash
$ slack-hitorigoto-esa -c ./config.toml -e "ESA_TOKEN1"
# SLACK_TOKEN => "SLACK_TOKEN2"
# ESA_TOKEN => "ESA_TOKEN1"
```

## Install

```bash
$ go get -d github.com/upamune/slack-hitorigoto-esa
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
