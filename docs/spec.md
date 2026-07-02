# Live Concert Notifier — 仕様書

## 1. 概要

ユーザーが登録したアーティストの新着ライブ情報を自動収集し、ダッシュボードで確認できるWebアプリ。

毎朝5時にバッチ処理がアーティスト公式サイトをスクレイピングし、LLMでライブ情報を構造化してDBに格納する。

---

## 2. ユーザーストーリー

| # | ストーリー |
|---|---|
| 1 | ユーザーはアーティストを検索して「フォロー」登録できる |
| 2 | ダッシュボードでフォロー中アーティストの今後のライブ一覧を確認できる |
| 3 | 新たにスクレイピングで追加されたライブには「NEW」バッジが付く |
| 4 | 各ライブの詳細（会場・日付・共演者・フェス判定など）を確認できる |

---

## 3. ドメインモデル

### User（ユーザー）
| フィールド | 型 | 説明 |
|---|---|---|
| id | UUID | 主キー |
| email | string | メールアドレス（ユニーク） |
| created_at | datetime | 登録日時 |

### Artist（アーティスト）
| フィールド | 型 | 説明 |
|---|---|---|
| id | UUID | 主キー |
| name | string | アーティスト名 |
| official_site_url | string | 公式サイトURL（スクレイピング対象） |
| created_at | datetime | 登録日時 |

### UserArtist（フォロー関係）
| フィールド | 型 | 説明 |
|---|---|---|
| user_id | UUID | FK → User |
| artist_id | UUID | FK → Artist |
| followed_at | datetime | フォロー日時 |

### Concert（ライブ情報）
| フィールド | 型 | 説明 |
|---|---|---|
| id | UUID | 主キー |
| artist_id | UUID | FK → Artist |
| title | string | ライブタイトル（例: "ARENA TOUR 2026"） |
| venue_name | string | 会場名（例: "さいたまスーパーアリーナ"） |
| venue_location | string | 開催地（都道府県）（例: "埼玉県"） |
| date | date | 開催日 |
| co_performers | string[] | 共演者リスト |
| is_festival | boolean | フェス・イベント形式か否か |
| source_url | string | スクレイピング元URL |
| raw_text | string | LLMに渡した前処理済みテキスト（デバッグ用） |
| discovered_at | datetime | スクレイピングで発見した日時 |
| created_at | datetime | DB登録日時 |

### ScrapeJob（スクレイピングジョブ）
| フィールド | 型 | 説明 |
|---|---|---|
| id | UUID | 主キー |
| artist_id | UUID | FK → Artist |
| status | enum | `pending` / `running` / `succeeded` / `failed` |
| started_at | datetime | 開始日時 |
| finished_at | datetime | 完了日時 |
| error_message | string | エラー内容（失敗時） |

---

## 4. システムフロー

### 4-1. 毎朝5時バッチ（スクレイピング〜DB格納）

```
[Cron: 05:00]
    │
    ▼
全アーティストの official_site_url を取得
    │
    ├─ [Artist A] ──▶ HTMLフェッチ ──▶ 前処理（不要タグ除去・本文抽出） ──▶ LLM
    ├─ [Artist B] ──▶  ...                                                   │
    └─ [Artist C] ──▶  ...                                       ┌───────────┘
                                                                  │
                                                           Concert[] を抽出
                                                                  │
                                                   既存DBと重複チェック（artist_id + date + venue_name）
                                                                  │
                                                   ┌──────────────┴──────────────┐
                                                 新規のみ                     既存はスキップ
                                                   │
                                               DB INSERT
                                           (discovered_at = now)
```

### 4-2. LLMへの入力・出力仕様

**入力**（プロンプト）
```
以下はアーティスト「{artist_name}」の公式サイトから抽出したテキストです。
ライブ・コンサート情報を全て抽出し、JSONの配列で返してください。

{preprocessed_html_text}
```

**出力**（期待するJSONスキーマ）
```json
[
  {
    "title": "string",
    "venue_name": "string",
    "venue_location": "string",
    "date": "YYYY-MM-DD",
    "co_performers": ["string"],
    "is_festival": false
  }
]
```

### 4-3. ダッシュボード表示フロー

```
[ユーザーがダッシュボードを開く]
    │
    ▼
フォロー中アーティストの Concert を取得
（date >= 今日 でフィルタ、date 昇順）
    │
    ▼
discovered_at が過去7日以内のものに "NEW" フラグを付与
    │
    ▼
フロントエンドで一覧表示
```

---

## 5. 技術スタック

| レイヤー | 技術 |
|---|---|
| Backend API | Go (net/http) |
| Frontend | Next.js (App Router) |
| DB | PostgreSQL |
| LLM | Claude API (Anthropic) |
| バッチ実行 | Cron（コンテナ内または外部スケジューラー） |
| インフラ | Docker / Docker Compose |

---

## 6. 非機能要件（初期スコープ）

- スクレイピングは各アーティストを順次実行（並列数は設定可能）
- LLM呼び出し失敗時はリトライ2回、それ以降は `ScrapeJob.status = failed` に記録
- 同一ライブの重複登録は `(artist_id, date, venue_name)` の組み合わせでユニーク制約により防ぐ
- 認証は初期スコープ外（将来的にJWT）
