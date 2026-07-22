# backend

Go (net/http) によるREST API / スクレイピングバッチ。

## パッケージ構成方針

- `cmd/server`: Cloud Run Service相当。常駐HTTPサーバーのエントリポイント
- `cmd/scraper`: Cloud Run Job相当。バッチ実行のエントリポイント（`internal/domain`/`internal/repository`をserverと共有）
- `internal/domain`: ドメインモデル・ビジネスルール。外部パッケージに依存しない
- `internal/oas`: **[自動生成／編集禁止]** `docs/api/openapi.yaml` から [ogen](https://github.com/ogen-go/ogen) で生成したサーバーコード（型・`Handler`インターフェース・ルーター）
- `internal/handler`: `internal/oas.Handler` の実装（`APIHandler`）。リソースごとにファイルを分ける（`users.go`, `artists.go`, ...）。未実装の操作は `oas.UnimplementedHandler` にフォールバックする
- `internal/service`: ユースケース層。handlerとrepositoryの間を仲介するビジネスロジック
- `internal/repository`: DB・外部APIアクセス層。[sqlc](https://sqlc.dev/) で生成した型安全なクエリを利用する
- `internal/scraping`: scraper専用。HTML取得・前処理
- `internal/ai`: scraper専用。Claude APIクライアント・構造化抽出

依存方向は `handler → service → repository` / `domain` を一方向とし、逆方向の依存を作らない。

## コード生成

```bash
cd backend
go generate ./...   # internal/oas（ogen）・internal/repository配下のsqlc生成物を再生成
```

- OpenAPI仕様（`docs/api/openapi.yaml`）を変更したら `internal/oas` を再生成する
- SQLスキーマ・クエリ（`internal/repository/queries/*.sql`）を変更したら sqlc生成物を再生成する
- 生成物はリポジトリにコミットする（`// Code generated ... DO NOT EDIT.` が付与されたファイルは直接編集しない）

## DBマイグレーション

マイグレーションファイルは `backend/migrations/`（`NNNN_name.up.sql` / `.down.sql`）に配置する。
sqlcはこのディレクトリをスキーマ定義として読み込むため、テーブル定義の変更は必ずここに追加すること。

```bash
docker compose up -d db
docker exec -i live-concert-notifier-db-1 psql -U app -d live_concert_notifier < backend/migrations/0001_init.up.sql
```

テーブル定義の詳細は [`docs/ai/database-schema.md`](../docs/ai/database-schema.md) を参照。

## テストの実行方法

```bash
cd backend
go test ./...
```
