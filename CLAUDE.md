# live-concert-notifier

ライブ・コンサート情報を収集し、ユーザーに通知するAI活用モノリポ。

- **Backend**: Go (net/http)
- **Frontend**: Next.js
- **インフラ**: Docker / Docker Compose

---

## ディレクトリ構成

```
live-concert-notifier/
├── CLAUDE.md                        # ← このファイル（AI向けルートコンテキスト）
│
├── docs/
│   ├── ai/                          # AI（Claude等）向けコンテキストドキュメント
│   │   ├── architecture.md          # システム全体のアーキテクチャ図・説明
│   │   ├── domain.md                # ドメイン知識・用語集（Concert / Venue / Artist 等）
│   │   ├── api-spec.md              # REST APIエンドポイント一覧と仕様
│   │   └── database-schema.md       # テーブル定義・ER図の説明
│   └── adr/                         # Architecture Decision Records
│       └── 0001-use-go-net-http.md  # 技術選定の意思決定記録
│
├── backend/                         # Go net/http API
│   ├── CLAUDE.md                    # バックエンド固有のAI向けコンテキスト
│   ├── cmd/
│   │   └── server/
│   │       └── main.go
│   ├── internal/
│   │   ├── domain/                  # ドメインモデル・ビジネスルール
│   │   ├── handler/                 # HTTPハンドラ
│   │   ├── service/                 # ユースケース層
│   │   ├── repository/              # DB・外部API アクセス層
│   │   └── ai/                      # AI機能（スクレイピング補助・要約等）
│   ├── go.mod
│   ├── go.sum
│   └── Dockerfile
│
├── frontend/                        # Next.js
│   ├── CLAUDE.md                    # フロントエンド固有のAI向けコンテキスト
│   ├── src/
│   │   ├── app/                     # App Router
│   │   ├── components/
│   │   └── lib/                     # APIクライアント・ユーティリティ
│   ├── package.json
│   ├── next.config.ts
│   └── Dockerfile
│
└── docker-compose.yml
```

---

## AI向けドキュメント（`docs/ai/`）の種類と役割

| ファイル | 役割 | 更新タイミング |
|---|---|---|
| `architecture.md` | サービス間の依存関係・データフロー。新機能追加時の影響範囲把握に使う | アーキテクチャ変更時 |
| `domain.md` | ドメイン用語の定義（例: "セトリ" = セットリスト）。AIが命名やコメントで迷わないための辞書 | 用語追加・変更時 |
| `api-spec.md` | エンドポイント・リクエスト/レスポンス形式。フロントエンド実装時の参照先 | API変更時 |
| `database-schema.md` | テーブル定義と制約。クエリ生成やマイグレーション作成時の参照先 | スキーマ変更時 |

`docs/adr/` には「なぜそう決めたか」の記録を残す。AIは現在のコードだけでは意図を読み取れないため、ADRがあると不必要なリファクタを防げる。
adrは必要最低限で構わない。gitのcommit idと変更理由のみ

---

## 各サブディレクトリの `CLAUDE.md`

ルートの `CLAUDE.md`（このファイル）はプロジェクト全体の概要のみ記載する。
詳細なコンテキスト（テスト方針・コーディング規約・よく使うコマンド）は各サブディレクトリの `CLAUDE.md` に分割して置く。

- `backend/CLAUDE.md`: Goのパッケージ構成方針、DBマイグレーション手順、テストの実行方法
- `frontend/CLAUDE.md`: コンポーネント設計方針、状態管理の考え方、スタイリングルール

---

## よく使うコマンド

```bash
# 全サービス起動
docker compose up

# バックエンドのみ起動
docker compose up backend

# フロントエンドのみ起動
docker compose up frontend
```
