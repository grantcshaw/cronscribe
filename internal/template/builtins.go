package template

// builtins contains the default set of cron schedule templates.
var builtins = []Template{
	{
		Name:        "hourly",
		Description: "Run once every hour at the top of the hour",
		Expression:  "0 * * * *",
		Tags:        []string{"frequent", "hourly"},
	},
	{
		Name:        "daily",
		Description: "Run once every day at midnight",
		Expression:  "0 0 * * *",
		Tags:        []string{"daily", "midnight"},
	},
	{
		Name:        "weekly",
		Description: "Run once every week on Sunday at midnight",
		Expression:  "0 0 * * 0",
		Tags:        []string{"weekly", "sunday"},
	},
	{
		Name:        "monthly",
		Description: "Run on the first day of every month at midnight",
		Expression:  "0 0 1 * *",
		Tags:        []string{"monthly"},
	},
	{
		Name:        "yearly",
		Description: "Run once a year on January 1st at midnight",
		Expression:  "0 0 1 1 *",
		Tags:        []string{"yearly", "annual"},
	},
	{
		Name:        "every-5-minutes",
		Description: "Run every 5 minutes",
		Expression:  "*/5 * * * *",
		Tags:        []string{"frequent", "minutes"},
	},
	{
		Name:        "every-15-minutes",
		Description: "Run every 15 minutes",
		Expression:  "*/15 * * * *",
		Tags:        []string{"frequent", "minutes"},
	},
	{
		Name:        "weekdays",
		Description: "Run every weekday (Monday to Friday) at 9am",
		Expression:  "0 9 * * 1-5",
		Tags:        []string{"weekday", "business"},
	},
	{
		Name:        "midnight-weekdays",
		Description: "Run at midnight on weekdays",
		Expression:  "0 0 * * 1-5",
		Tags:        []string{"weekday", "midnight"},
	},
	{
		Name:        "noon-daily",
		Description: "Run every day at noon",
		Expression:  "0 12 * * *",
		Tags:        []string{"daily", "noon"},
	},
}
