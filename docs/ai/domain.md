# ドメイン用語集

`docs/spec.md`「3. ドメインモデル」の実体を、AIがコード生成・命名・コメントを書く際に迷わないよう補足する用語辞典。
フィールドの型・制約そのものは `docs/ai/database-schema.md` を正とする。

---

## エンティティ

### Artist（アーティスト）
スクレイピング対象となるアーティスト本体。`official_site_url` が公式サイトのスケジュール/ライブ情報ページを指す想定。

### Concert（ライブ・コンサート情報）
1公演を表す最小単位。`docs/spec.md` の用語では「ライブ」と同義。フェスも1つの `Concert` として扱う（`is_festival=true` の場合、`title` にフェス名が入る）。

### ScrapeJob（スクレイピングジョブ）
1アーティスト・1回のスクレイピング実行を表す記録。`docs/ai/architecture.md` 4-1（定期バッチ）・4-2（手動トリガー）のどちらから作られても同じテーブルに記録される。

---

## 用語

| 用語 | 意味 |
|---|---|
| **フェス判定** | `Concert.is_festival`。単独公演ではなく複数アーティストが出演するフェス・イベント形式かどうかのフラグ。LLMへの構造化抽出プロンプト（`docs/spec.md` 4-2）が直接この値を出力する |
| **共演者 (co_performers)** | 同じ公演に出演する他アーティスト。フェスの場合は出演アーティスト一覧、単独公演の場合はゲスト出演者を指すことが多い |
| **NEW バッジ / is_new** | ダッシュボード上で新着ライブを示すフラグ。`Concert.discovered_at` が現在時刻から過去7日以内なら `true`（`Concert.IsNew` メソッド、境界値はちょうど7日前を含む＝`<=`） |
| **discovered_at** | スクレイピングでDBに初めて登録された日時。ライブ自体の開催日（`date`）とは無関係 |
| **重複チェック** | `(artist_id, date, venue_name)` の組み合わせで同一ライブを識別する。3項目が一致する場合は新規INSERTせず静かにスキップする（`ConcertRepository.Create` の `inserted=false`） |
| **raw_text** | スクレイピングしたHTMLを前処理（不要タグ除去・本文抽出）したテキスト。LLMへの入力そのものであり、抽出結果のデバッグ用に保存する |
| **手動スクレイピングトリガー** | 管理者が `POST /admin/scrape` でオンデマンドにScraper Jobを起動する操作。定期バッチ（毎朝5時）とは独立した経路だが、生成される `ScrapeJob` の扱いは同じ |
| **ScraperTrigger** | Backend APIからScraper Jobを起動する処理を抽象化したインターフェース（`internal/service/scraper_trigger.go`）。本番はCloud Run Admin API、ローカルはログ出力のみのスタブ（`LogScraperTrigger`） |

---

## 状態遷移

### ScrapeJobStatus

```
pending ──▶ running ──┬─▶ succeeded
                       └─▶ failed
```

- `pending`: 作成直後（`POST /admin/scrape` 実行時、または定期バッチのキュー投入時）
- `running`: Scraper Jobが処理を開始
- `succeeded` / `failed`: 処理完了。`failed` の場合は `error_message` に理由を記録
- LLM呼び出し失敗時はリトライ2回まで行い、それでも失敗した場合に `failed` へ遷移する（`docs/spec.md` 6. 非機能要件）

---

## 責務境界（迷ったときの判断基準）

- **フォロー中アーティストの絞り込み・NEW判定**はBackend API（`internal/service`）の責務。Scraper Jobは関与しない
- **HTML取得・LLM構造化・重複チェック・DB書き込み**はScraper Job（`internal/scraping`, `internal/ai`）の責務。Backend APIはトリガーするのみ
- 迷ったら `docs/ai/architecture.md` 2.「各ノードの責務分割」の表を参照する
