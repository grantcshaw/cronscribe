package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/user/cronscribe/internal/cli"
)

func main() {
	var (
		previewCount = flag.Int("n", 5, "number of next-run previews to display")
		verbose      = flag.Bool("v", false, "show human-readable description of the cron expression")
	)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: cronscribe [options] <schedule>\n\n")
		fmt.Fprintf(os.Stderr, "Translates a human-readable schedule into a cron expression.\n\n")
		fmt.Fprintf(os.Stderr, "Examples:\n")
		fmt.Fprintf(os.Stderr, "  cronscribe \"every day at noon\"\n")
		fmt.Fprintf(os.Stderr, "  cronscribe -v -n 3 \"every monday at 9am\"\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	input := strings.Join(args, " ")
	cfg := cli.Config{
		PreviewCount: *previewCount,
		Verbose:      *verbose,
	}

	if err := cli.Run(input, cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
