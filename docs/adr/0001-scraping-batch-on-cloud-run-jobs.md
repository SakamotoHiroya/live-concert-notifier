# 0001. スクレイピングバッチをCloud Run Jobsで実装する

## ステータス

承認

## コンテキスト

`docs/spec.md` の通り、毎朝5時に全アーティストの公式サイトをスクレイピングし、LLM（Claude API）で構造化してDBに格納するバッチ処理が必要。
このバッチは以下の特性を持つ。

- アーティスト数に応じて実行時間が伸びる（HTML取得・LLM呼び出しを順次実行）
- 外部サイトへのHTTPリクエストやLLM呼び出しの失敗・リトライが発生する
- 1日1回の定期実行に加えて、管理者による手動トリガー（`POST /admin/scrape`）も必要
- Backend API（常時稼働のHTTPサーバー）とはライフサイクルが異なる（起動→処理→終了）

GCP上でこのバッチをどう実装するかを決定する必要がある。

## 検討した選択肢

### 1. Cloud Functions (2nd gen) + Cloud Scheduler

- 実装がシンプルだが、HTTPトリガー関数は最大タイムアウトの制約が厳しく、複数アーティストを順次処理してLLM呼び出しを挟む処理には不向き
- Backend APIと実行環境（ランタイム・依存関係の共有方法）が分かれ、Goの共通パッケージ（`internal/domain`等）の再利用がしづらい

### 2. GKE CronJob

- 柔軟だが、常時稼働するクラスタの運用コストとオペレーション負荷が、1日1回のバッチには見合わない
- 本プロジェクトの規模（モノリポ・小規模チーム）にはオーバースペック

### 3. Compute Engine + cron

- サーバーレスではなく、VMの起動・パッチ適用・スケーリングを自前で管理する必要がある
- アイドル時間分の課金が発生し、コスト効率が悪い

### 4. Cloud Run Jobs + Cloud Scheduler（採用）

- サーバーレスで、実行時間に応じた課金（アイドルコストなし）
- 最大24時間の実行時間が確保でき、複数アーティストの順次処理・リトライに十分
- 既存のDocker/Goベースの開発体験をそのまま踏襲でき、`internal/domain`・`internal/repository`をBackend APIと共有可能
- Cloud Schedulerからの定期起動、Backend APIからのCloud Run Admin API経由の手動起動の両方に対応
- Cloud Run Jobsのタスク分割機能を使えば、将来的にアーティストごとの並列実行にも拡張できる

## 決定

スクレイピング〜LLM構造化〜DB書き込みのバッチ処理は、Backend API（Cloud Run Service）とは別の **Cloud Run Job** として実装する。

- エントリポイントは `backend/cmd/scraper`（Backend APIの `cmd/server` とは別バイナリ・別コンテナイメージ）
- 定期実行: Cloud Scheduler（毎朝5:00 Asia/Tokyo）が OIDC 認証付きHTTPで Cloud Run Jobs の `:run` エンドポイントを呼び出す
- 手動実行: Backend APIが `POST /admin/scrape` を受けて `scrape_jobs` に `pending` レコードを作成した後、Cloud Run Admin API経由で同じJobを対象 `artist_ids` 付きで起動する
- Backend APIとScraper Jobは `internal/domain` / `internal/repository` を共有するが、デプロイ単位（コンテナイメージ・Cloud Runリソース種別）は分離する

詳細は [`docs/ai/architecture.md`](../ai/architecture.md) を参照。

## 影響

- `backend/cmd/scraper` および `backend/internal/scraping`（HTML取得・前処理）、`backend/internal/ai`（LLM構造化）を新規に追加する
- `backend/Dockerfile.scraper` を別途用意し、Artifact Registryに独立したイメージとしてpushする
- Backend APIのサービスアカウントに、Scraper Jobを起動するための最小権限（対象Jobリソースへの `run.jobs.run`）を付与する必要がある
- Scraper Jobは常駐サーバーではないため、外部から直接HTTPで呼び出すことはできない（Cloud Scheduler / Cloud Run Admin API経由のみ）
