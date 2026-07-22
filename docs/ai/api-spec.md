# APIサマリ

正は `docs/api/openapi.yaml`。本ファイルはフロントエンド実装時に参照しやすいよう、エンドポイント一覧とレスポンス形状を要約したもの。
詳細な型・バリデーション・エラーレスポンスは必ず openapi.yaml 側を確認すること。

ベースURL: ローカル開発では `http://localhost:8080/api/v1`

---

## users

| メソッド・パス | 概要 | 成功レスポンス | エラー |
|---|---|---|---|
| `POST /users` | ユーザー登録（`email`） | `201 User` | `400` / `409`（email重複） |
| `GET /users/{userId}` | ユーザー取得 | `200 User` | `404` |

## artists

| メソッド・パス | 概要 | 成功レスポンス | エラー |
|---|---|---|---|
| `GET /artists` | 一覧（`q`部分一致、`limit`/`offset`） | `200 ArtistList` | - |
| `POST /artists` | 追加（管理者用） | `201 Artist` | `400` / `409`（`official_site_url`重複） |
| `GET /artists/{artistId}` | 取得 | `200 Artist` | `404` |
| `PUT /artists/{artistId}` | 更新（管理者用、部分更新） | `200 Artist` | `400` / `404` |
| `DELETE /artists/{artistId}` | 削除（管理者用） | `204` | `404` |

## follows

| メソッド・パス | 概要 | 成功レスポンス | エラー |
|---|---|---|---|
| `GET /users/{userId}/follows` | フォロー中アーティスト一覧 | `200 ArtistList` | `404`（user不在） |
| `POST /users/{userId}/follows` | フォロー登録（`artist_id`） | `201` | `404` / `409`（重複フォロー） |
| `DELETE /users/{userId}/follows/{artistId}` | フォロー解除 | `204` | `404` |

## concerts

| メソッド・パス | 概要 | 成功レスポンス | エラー |
|---|---|---|---|
| `GET /concerts` | 一覧（`artist_id`/`from`/`to`/`is_festival`/`limit`/`offset`） | `200 ConcertList` | - |
| `GET /concerts/{concertId}` | 詳細 | `200 Concert` | `404` |

## dashboard

| メソッド・パス | 概要 | 成功レスポンス | エラー |
|---|---|---|---|
| `GET /users/{userId}/dashboard` | フォロー中アーティストの今後のライブ一覧（`date >= 今日`、`date`昇順、`is_new`付与） | `200 DashboardResponse` | `404`（user不在） |

`DashboardResponse.items` は `Concert` のフィールド全部 + `is_new: boolean`（`discovered_at` が過去7日以内なら`true`）。

## admin

| メソッド・パス | 概要 | 成功レスポンス | エラー |
|---|---|---|---|
| `POST /admin/scrape` | 手動スクレイピング実行（`artist_ids`省略時は全件） | `202 TriggerScrapeResponse`（`job_ids`） | - |
| `GET /admin/scrape-jobs` | ジョブ一覧（`artist_id`/`status`/`limit`/`offset`） | `200 ScrapeJobList` | - |
| `GET /admin/scrape-jobs/{jobId}` | ジョブ詳細 | `200 ScrapeJob` | `404` |

---

## 共通スキーマの形

### 一覧レスポンス（`*List`）
すべて `{ items: T[], total: number }` の形。`total` はフィルタ適用後・ページング前の全件数。

### エラーレスポンス
```json
{ "code": "NOT_FOUND", "message": "指定されたリソースが見つかりません" }
```
`code` は `NOT_FOUND` 等のスネークケース定数（フロントの分岐に使える想定）。

### ページング
`limit`（デフォルト20・最大100）・`offset`（デフォルト0）はクエリパラメータで共通。

---

## フロントエンド実装時の注意

- `Concert.artist_name` はJOIN済みの表示用フィールド。アーティスト名表示のために別途 `GET /artists/{artistId}` を呼ぶ必要はない
- `ScrapeJob` の `started_at`/`finished_at`/`error_message` は `pending` 状態では `null`
- `admin` タグの操作はUIから叩く場合、認証・権限チェックは初期スコープ外（`docs/spec.md` 6.）だが将来的にJWT等で保護される前提でUI設計すること
