# cronscribe

Translates human-readable schedules into cron expressions with validation and next-run previews.

---

## Installation

```bash
go get github.com/yourname/cronscribe
```

Or install the CLI tool:

```bash
go install github.com/yourname/cronscribe/cmd/cronscribe@latest
```

---

## Usage

### As a library

```go
package main

import (
    "fmt"
    "github.com/yourname/cronscribe"
)

func main() {
    expr, err := cronscribe.Translate("every day at 9am")
    if err != nil {
        panic(err)
    }

    fmt.Println(expr.Cron)        // 0 9 * * *
    fmt.Println(expr.NextRuns(3)) // next 3 scheduled times
}
```

### As a CLI

```bash
$ cronscribe "every Monday at noon"
Cron:      0 12 * * 1
Next runs:
  - Mon, 14 Jul 2025 12:00:00 UTC
  - Mon, 21 Jul 2025 12:00:00 UTC
  - Mon, 28 Jul 2025 12:00:00 UTC
```

---

## Supported Phrases

| Input | Cron |
|---|---|
| `every hour` | `0 * * * *` |
| `every day at 9am` | `0 9 * * *` |
| `every Monday at noon` | `0 12 * * 1` |
| `every weekday at 8:30am` | `30 8 * * 1-5` |

---

## License

MIT © [yourname](https://github.com/yourname)