# データベーススキーマ

PostgreSQL。定義は `backend/migrations/0001_init.up.sql` を正とする（本ファイルは要約）。

## ER図

```
users ──< user_artists >── artists ──< concerts
                                    └─< scrape_jobs
```

## テーブル定義

### users

| カラム | 型 | 制約 |
|---|---|---|
| id | uuid | PK |
| email | text | UNIQUE, NOT NULL |
| created_at | timestamptz | NOT NULL, DEFAULT now() |

### artists

| カラム | 型 | 制約 |
|---|---|---|
| id | uuid | PK |
| name | text | NOT NULL |
| official_site_url | text | UNIQUE, NOT NULL |
| created_at | timestamptz | NOT NULL, DEFAULT now() |

### user_artists（フォロー関係）

| カラム | 型 | 制約 |
|---|---|---|
| user_id | uuid | PK（複合）, FK → users.id ON DELETE CASCADE |
| artist_id | uuid | PK（複合）, FK → artists.id ON DELETE CASCADE |
| followed_at | timestamptz | NOT NULL, DEFAULT now() |

### concerts

| カラム | 型 | 制約 |
|---|---|---|
| id | uuid | PK |
| artist_id | uuid | FK → artists.id ON DELETE CASCADE |
| title | text | NOT NULL DEFAULT '' |
| venue_name | text | NOT NULL |
| venue_location | text | NOT NULL |
| date | date | NOT NULL |
| co_performers | text[] | NOT NULL DEFAULT '{}' |
| is_festival | boolean | NOT NULL DEFAULT false |
| source_url | text | NOT NULL |
| raw_text | text | NOT NULL DEFAULT '' |
| discovered_at | timestamptz | NOT NULL, DEFAULT now() |
| created_at | timestamptz | NOT NULL, DEFAULT now() |

- `UNIQUE (artist_id, date, venue_name)` — スクレイピングバッチの重複登録防止（`docs/spec.md` 6章）
- `INDEX (date)` — ダッシュボード・一覧取得の `date >= 今日` 絞り込み用

### scrape_jobs

| カラム | 型 | 制約 |
|---|---|---|
| id | uuid | PK |
| artist_id | uuid | FK → artists.id ON DELETE CASCADE |
| status | text | NOT NULL, CHECK IN ('pending','running','succeeded','failed') |
| started_at | timestamptz | NULL可 |
| finished_at | timestamptz | NULL可 |
| error_message | text | NULL可 |

- `INDEX (artist_id)`, `INDEX (status)` — 管理画面の絞り込み用

## マイグレーション適用

```bash
docker compose up -d db
docker exec -i live-concert-notifier-db-1 psql -U app -d live_concert_notifier < backend/migrations/0001_init.up.sql
```

## クエリ生成（sqlc）

`backend/internal/repository/queries/*.sql` を編集したら以下で再生成する。

```bash
cd backend
go tool sqlc generate
```

生成物は `backend/internal/repository/sqlcgen/`（コミット対象、直接編集禁止）。
