twi-meteor
===

24時間以上経過したツイートを削除するツール

# Getting Started

twi-meteorは[Heroku](https://jp.heroku.com/)を利用してworker dynoとしてデプロイする事を想定しています。
worker dynoに対して、[Heroku Scheduler](https://devcenter.heroku.com/articles/scheduler)などを利用して定期的に実行してあげることで、
24時間以上経過したツイートの定期的な削除を可能にします。

oauth1.0aでの接続となるため以下の情報を必要とします。

```
CONSUMER_KEY
CONSUMER_SECRET
ACCESS_TOKEN
ACCESS_TOKEN_SECRET
TWITTER_ID
```

## Herokuへのデプロイ方法

[![Deploy](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy)

## 手元での実行 (Docker)

ローカルでDockerを利用して実行する際の手順です。

1. リポジトリのクローン
2. .env.sampleファイルをコピーして.envファイルを作成し、必要な情報を入力する。
3. リポジトリのルートで `docker build -t twi-meteor .` をしてコンテナのビルドをする。
4. リポジトリのルートで `docker run -d  --env-file .env twi-meteor` をしてコンテナの実行をする。

## 手元での実行　(Go)

ローカルでGoを利用して実行する際の手順です。Goは1.17を想定しています。

1. リポジトリのクローン
2. .env.sampleファイルをコピーして.envファイルを作成し、必要な情報を入力する。
3. リポジトリのルートで　`go run main.go` を実行する。


