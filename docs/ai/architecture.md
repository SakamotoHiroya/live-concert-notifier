# アーキテクチャ（GCPデプロイ前提）

`docs/spec.md`（仕様）・`docs/api/openapi.yaml`（API仕様）を踏まえた、本番環境（GCP）でのシステム構成。
ローカル開発は引き続き `docker compose` を使うが、本番は各コンポーネントを独立したGCPリソースとしてデプロイする。

---

## 1. コンポーネント一覧

| コンポーネント | 実体 | GCPリソース |
|---|---|---|
| Frontend | Next.js | Cloud Run Service |
| Backend API | Go net/http | Cloud Run Service |
| Scraper | Go（`cmd/scraper`） | Cloud Run Job |
| DB | PostgreSQL | Cloud SQL |
| LLM | Claude API | 外部API（Anthropic） |
| 定期実行トリガー | — | Cloud Scheduler |
| シークレット管理 | — | Secret Manager |
| コンテナイメージ | — | Artifact Registry |
| 監視・ログ | — | Cloud Logging / Cloud Monitoring / Error Reporting |

呼び出し関係:
- Frontend → Backend API（HTTPS/JSON）
- Backend API → Cloud SQL（Cloud SQL Connector）
- Backend API → Scraper Job（Cloud Run Admin APIの `jobs.run` を呼び出して手動起動）
- Cloud Scheduler → Scraper Job（OIDC認証付きHTTPで `jobs.run` を呼び出して定期起動）
- Scraper Job → Cloud SQL（Cloud SQL Connector）
- Scraper Job → Claude API（HTTPS）

---

## 2. 各ノードの責務分割

スクレイピング〜LLM構造化〜DB書き込みまでを **Backend APIとは別の実行単位（Cloud Run Job）** に分離するのが最大のポイント。
理由は [ADR 0001](../adr/0001-scraping-batch-on-cloud-run-jobs.md) を参照。

| ノード | 責務 | 責務でないこと |
|---|---|---|
| **Frontend** | ダッシュボードUI描画、アーティスト検索・フォロー操作、Backend APIの呼び出し | DBへの直接アクセス、LLM呼び出し |
| **Backend API** | ユーザー・アーティスト・フォロー・コンサートのCRUD、ダッシュボード集約、`scrape_jobs`のpendingレコード作成、Scraper Jobの起動 | スクレイピング処理そのもの、LLM呼び出し |
| **Scraper Job** | 対象アーティストのHTML取得→前処理→Claude APIで構造化→重複チェック→`concerts`へINSERT→`scrape_jobs`のステータス更新 | HTTPリクエストの直接受付（常駐サーバーではない）、フロント向けAPI提供 |
| **Cloud SQL (PostgreSQL)** | 唯一の永続化層。Backend API / Scraper Job の両方から接続 | ビジネスロジック |
| **Claude API** | スクレイピング済みテキストからConcert情報をJSON構造化 | データ永続化 |
| **Cloud Scheduler** | 毎朝5時（Asia/Tokyo）にScraper Jobを起動 | ジョブの実処理 |
| **Secret Manager** | DB認証情報・Claude APIキーの管理と各Cloud Run実行時への注入 | — |

Backend APIとScraper Jobは `internal/domain`・`internal/repository` を共有するが、**デプロイ単位（コンテナイメージ・Cloud Runリソース種別）は別**にする。
Backend APIは常駐HTTPサーバー（Cloud Run **Service**）、Scraper Jobは起動→処理→終了のバッチ（Cloud Run **Job**）という性質の違いに対応させるため。

---

## 3. GCPリソースマッピング（詳細）

| コンポーネント | GCPリソース | 備考 |
|---|---|---|
| Frontend | Cloud Run Service | `frontend/Dockerfile` からビルド。min instances 0〜1 |
| Backend API | Cloud Run Service | `backend/Dockerfile`（`cmd/server`）からビルド |
| Scraper | Cloud Run Job | `backend/Dockerfile.scraper`（`cmd/scraper`）からビルド。1日1回 + 手動実行 |
| DB | Cloud SQL for PostgreSQL | Backend API / Scraper から Cloud SQL Connector（unixソケット）経由で接続 |
| 定期実行 | Cloud Scheduler | OIDCトークン付きHTTPで Cloud Run Jobs の `:run` エンドポイントを叩く |
| シークレット | Secret Manager | `DB_PASSWORD`, `ANTHROPIC_API_KEY` 等。Cloud Runの環境変数にシークレット参照として注入 |
| コンテナイメージ | Artifact Registry | frontend / backend / scraper の3イメージ |
| ログ・監視 | Cloud Logging / Cloud Monitoring / Error Reporting | `scrape_jobs.status = failed` が続く場合のアラート等 |
| CI/CD | Cloud Build または GitHub Actions | push時にビルド→Artifact Registryへpush→`gcloud run deploy` / `gcloud run jobs deploy` |

---

## 4. データフロー

### 4-1. 定期スクレイピングバッチ（Cloud Scheduler起点）

1. Cloud Scheduler が毎朝5:00（Asia/Tokyo）にOIDC付きHTTP POSTで Cloud Run Jobs API（`projects/.../jobs/scraper:run`）を呼び出す
2. Scraper Job が起動し、`artists` テーブルを全件取得
3. 各アーティストごとに `scrape_jobs` を `status=running` で記録
4. HTML取得 → 前処理 → Claude APIで構造化 → `Concert[]` を抽出
5. `(artist_id, date, venue_name)` で既存DBと重複チェック
6. 新規のみ `concerts` へINSERT（`discovered_at = now`）
7. `scrape_jobs.status` を `succeeded` / `failed` に更新（失敗時は `error_message` も記録）

### 4-2. 手動スクレイピングトリガー（`POST /admin/scrape`）

1. Frontend/管理者が Backend API に `POST /admin/scrape { artist_ids? }` を送信
2. Backend API が対象アーティストごとに `scrape_jobs` を `status=pending` で作成
3. Backend API が Cloud Run Admin API を呼び出し、Scraper Job を対象 `artist_ids` 付きで起動（`overrides.containers[].env` で `TARGET_ARTIST_IDS` を渡す）
4. Backend API が `202 Accepted { job_ids: [...] }` を返す
5. 以降は 4-1 の3〜7と同じ流れでScraper Jobが処理

Backend APIのサービスアカウントには、Scraper Jobを起動するための最小権限（対象Jobリソースへの `run.jobs.run`）のみ付与する。

### 4-3. ダッシュボード表示

1. Frontend が Backend API に `GET /users/{userId}/dashboard` を送信
2. Backend API が Cloud SQL からフォロー中アーティストの `concerts` を取得（`date >= 今日`、`discovered_at` が7日以内なら `is_new=true`）
3. Frontend で一覧表示

---

## 5. ディレクトリ構成（更新版）

ルート `CLAUDE.md` に記載の構成をベースに、GCPデプロイに必要な要素を追加する。

```
live-concert-notifier/
├── CLAUDE.md
│
├── docs/
│   ├── ai/
│   │   ├── architecture.md          # 本ファイル
│   │   ├── domain.md
│   │   ├── api-spec.md
│   │   └── database-schema.md
│   ├── adr/
│   │   └── 0001-scraping-batch-on-cloud-run-jobs.md
│   ├── api/
│   │   └── openapi.yaml
│   └── spec.md
│
├── backend/
│   ├── CLAUDE.md
│   ├── cmd/
│   │   ├── server/                  # Cloud Run Service（REST API）
│   │   │   └── main.go
│   │   └── scraper/                 # Cloud Run Job（定期バッチ）※新規
│   │       └── main.go
│   ├── internal/
│   │   ├── domain/                  # server / scraper 共有
│   │   ├── handler/                 # server専用（HTTPハンドラ）
│   │   ├── service/                 # server / scraper 共有ユースケース
│   │   ├── repository/              # server / scraper 共有（DB・Cloud SQL接続）
│   │   ├── scraping/                # scraper専用：HTML取得・前処理 ※新規
│   │   │   ├── fetcher.go
│   │   │   └── extractor.go
│   │   └── ai/                      # scraper専用：Claude APIクライアント・構造化
│   │       └── claude_client.go
│   ├── Dockerfile                   # server用イメージ
│   ├── Dockerfile.scraper           # scraper用イメージ ※新規
│   ├── go.mod
│   └── go.sum
│
├── frontend/
│   ├── CLAUDE.md
│   ├── src/
│   │   ├── app/
│   │   ├── components/
│   │   └── lib/
│   ├── package.json
│   ├── next.config.ts
│   └── Dockerfile
│
├── infra/                           # GCPインフラ定義 ※新規
│   ├── terraform/
│   │   ├── main.tf
│   │   ├── cloud_run_service_backend.tf
│   │   ├── cloud_run_service_frontend.tf
│   │   ├── cloud_run_job_scraper.tf
│   │   ├── cloud_sql.tf
│   │   ├── cloud_scheduler.tf
│   │   ├── secret_manager.tf
│   │   ├── artifact_registry.tf
│   │   ├── iam.tf
│   │   └── variables.tf
│   └── README.md
│
├── .github/
│   └── workflows/
│       └── deploy.yml               # CI/CD ※新規
│
└── docker-compose.yml               # ローカル開発専用
```

`backend/cmd/scraper` は `internal/domain` / `internal/repository` を `cmd/server` と共有しつつ、独立したエントリポイント・独立したコンテナイメージとしてビルドする。
これによりCloud Run Service（常駐）とCloud Run Job（都度実行）というライフサイクルの違いを自然に表現できる。

---

## 6. IAM / サービスアカウント設計（概要）

| サービスアカウント | 付与する権限 | 用途 |
|---|---|---|
| `backend-api-sa` | Cloud SQL Client, Secret Manager Secret Accessor, Scraper Jobへの `run.jobs.run` | Backend APIの実行、手動スクレイピングトリガー |
| `scraper-job-sa` | Cloud SQL Client, Secret Manager Secret Accessor | Scraper Jobの実行 |
| `scheduler-invoker-sa` | Scraper Jobへの `run.jobs.run`（Cloud Run Invoker相当） | Cloud SchedulerからJobを起動 |
| `frontend-sa` | （Backend API呼び出しのみ、GCPリソースへの直接アクセスなし） | Frontendの実行 |

いずれも最小権限の原則に従い、Jobリソース単位で権限を絞る（プロジェクト全体への `roles/run.developer` 等は付与しない）。

---

## 7. ローカル開発との差分

| 項目 | ローカル (docker compose) | 本番 (GCP) |
|---|---|---|
| Backend API | `docker compose up backend` | Cloud Run Service |
| Scraper | ローカルcronまたは手動実行のコンテナ | Cloud Run Job（Cloud Scheduler起点） |
| DB | ローカルPostgreSQLコンテナ | Cloud SQL for PostgreSQL |
| シークレット | `.env` | Secret Manager |
| イメージ | ローカルビルド | Artifact Registry |

ローカルでは `cmd/scraper` を単発の `docker compose run scraper` 相当で実行し、本番同様のバッチ処理を検証できるようにする。
