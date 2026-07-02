# Git Rules

## Branch

形式

```text
<type>/<description>
```

例

```text
feature/user-login
fix/token-refresh
setup/docker
docs/readme
refactor/auth
```

### Branch Types

| Type       | 用途             |
| ---------- | -------------- |
| `feature`  | 新機能            |
| `fix`      | バグ修正           |
| `refactor` | リファクタリング       |
| `docs`     | ドキュメント         |
| `test`     | テスト            |
| `setup`    | 初期設定・環境構築      |
| `chore`    | その他（設定・依存関係など） |
| `style`    | フォーマットのみ       |
| `perf`     | パフォーマンス改善      |
| `ci`       | CI/CD          |
| `release`  | リリース           |

---

## Commit

形式

```text
<gitmoji> <type>: <description>
```

例

```text
✨ feature: add login page
🐛 fix: fix token refresh
📝 docs: update README
♻️ refactor: simplify auth service
```

### Gitmoji

| Emoji | Code                    | 用途         |
| ----- | ----------------------- | ---------- |
| ✨     | `:sparkles:`            | 新機能        |
| 🐛    | `:bug:`                 | バグ修正       |
| ♻️    | `:recycle:`             | リファクタリング   |
| 📝    | `:memo:`                | ドキュメント     |
| ✅     | `:white_check_mark:`    | テスト        |
| 🎨    | `:art:`                 | コード整形・スタイル |
| ⚡️    | `:zap:`                 | パフォーマンス改善  |
| 🔧    | `:wrench:`              | 設定変更       |
| ⬆️    | `:arrow_up:`            | 依存関係更新     |
| 👷    | `:construction_worker:` | CI/CD      |
| 🚀    | `:rocket:`              | リリース       |
| ⏪️    | `:rewind:`              | Revert     |

---

## Pull Request

形式

```text
<gitmoji> <type>: <description>
```

例

```text
✨ feature: add login page
🐛 fix: resolve token refresh issue
📝 docs: update setup guide
```

---

## Issue

形式

```text
<gitmoji> <type>: <description>
```

例

```text
✨ feature: add push notification for new concerts
🐛 fix: fix duplicate notification on scraping retry
📝 docs: add API spec for concerts endpoint
```

### Rules

* Issueのtype・gitmojiはBranch/Commitと同じ表を使う
* 1つのIssueでは1つの目的（機能・バグ・タスク）のみ扱う
* 本文には背景・完了条件（Acceptance Criteria）を書く
* 対応するPRはIssueをクローズするキーワード（`Closes #<番号>`等）を含める

---

## Rules

* ブランチ名・コミット名・PR名・Issue名は英語で記述する
* `description` は短く内容が分かる名前にする
* 1つのブランチ・PR・Issueでは1つの目的のみ扱う
* Gitmojiを必ず付ける
