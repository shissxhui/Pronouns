package f2

import (
	"fmt"
	"net/http"
	"time"

	"github.com/urfave/cli/v2"
)

func init() {
	// Override the default help template
	cli.AppHelpTemplate = `DESCRIPTION:
	{{.Usage}}

USAGE:
   {{.HelpName}} {{if .UsageText}}{{ .UsageText }}{{end}}
{{if len .Authors}}
AUTHOR:
   {{range .Authors}}{{ . }}{{end}}{{end}}
{{if .Version}}
VERSION:
	 {{.Version}}{{end}}
{{if .VisibleFlags}}
FLAGS:{{range .VisibleFlags}}
	 {{.}}{{end}}{{end}}

DOCUMENTATION:
	https://github.com/ayoisaiah/f2#examples

WEBSITE:
	https://github.com/ayoisaiah/f2
`

	// Override the default version printer
	oldVersionPrinter := cli.VersionPrinter
	cli.VersionPrinter = func(c *cli.Context) {
		oldVersionPrinter(c)
		checkForUpdates(GetApp())
	}
}

func checkForUpdates(app *cli.App) {
	fmt.Println("Checking for updates...")

	c := http.Client{Timeout: 20 * time.Second}
	resp, err := c.Get("https://github.com/ayoisaiah/f2/releases/latest")
	if err != nil {
		fmt.Println("HTTP Error: Failed to check for update")
		return
	}

	defer resp.Body.Close()

	var version string
	_, err = fmt.Sscanf(
		resp.Request.URL.String(),
		"https://github.com/ayoisaiah/f2/releases/tag/%s",
		&version,
	)
	if err != nil {
		fmt.Println("Failed to get latest version")
		return
	}

	if version == app.Version {
		fmt.Printf(
			"Congratulations, you are using the latest version of %s\n",
			app.Name,
		)
	} else {
		fmt.Printf("%s: %s at %s\n", green("Update available"), version, resp.Request.URL.String())
	}
}

// GetApp retrieves the f2 app instance
func GetApp() *cli.App {
	return &cli.App{
		Name: "F2",
		Authors: []*cli.Author{
			{
				Name:  "Ayooluwa Isaiah",
				Email: "ayo@freshman.tech",
			},
		},
		Usage:                "F2 is a command-line tool for batch renaming multiple files and directories quickly and safely",
		UsageText:            "FLAGS [OPTIONS] [PATHS...]",
		Version:              "v1.3.0",
		EnableBashCompletion: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "find",
				Aliases: []string{"f"},
				Usage:   "Search `<pattern>`. Treated as a regular expression by default. Use -s or --string-mode to opt out",
			},
			&cli.StringFlag{
				Name:    "replace",
				Aliases: []string{"r"},
				Usage:   "Replacement `<string>`. If omitted, defaults to an empty string. Supports built-in and regex capture variables",
			},
			&cli.IntFlag{
				Name:        "start-num",
				Aliases:     []string{"n"},
				Usage:       "When using an auto incrementing number in the replacement string such as %03d, start the count from `<number>`",
				Value:       1,
				DefaultText: "1",
			},
			&cli.StringSliceFlag{
				Name:    "exclude",
				Aliases: []string{"E"},
				Usage:   "Exclude files/directories that match the given find pattern. Treated as a regular expression. Multiple exclude `<pattern>`s can be specified.",
			},
			&cli.StringFlag{
				Name:    "output-file",
				Aliases: []string{"o"},
				Usage:   "Output a map `<file>` for the current operation",
			},
			&cli.BoolFlag{
				Name:    "exec",
				Aliases: []string{"x"},
				Usage:   "Execute the batch renaming operation",
			},
			&cli.BoolFlag{
				Name:    "recursive",
				Aliases: []string{"R"},
				Usage:   "Rename files recursively",
			},
			&cli.IntFlag{
				Name:        "max-depth",
				Aliases:     []string{"m"},
				Usage:       "positive `<integer>` indicating the maximum depth for a recursive search (set to 0 for no limit)",
				Value:       0,
				DefaultText: "0",
			},
			&cli.StringFlag{
				Name:      "undo",
				Aliases:   []string{"u"},
				TakesFile: true,
				Usage:     "Undo a successful operation using a previously created map `file`",
			},
			&cli.BoolFlag{
				Name:    "ignore-case",
				Aliases: []string{"i"},
				Usage:   "Ignore case",
			},
			&cli.BoolFlag{
				Name:    "ignore-ext",
				Aliases: []string{"e"},
				Usage:   "Ignore extension",
			},
			&cli.BoolFlag{
				Name:    "include-dir",
				Aliases: []string{"d"},
				Usage:   "Include directories",
			},
			&cli.BoolFlag{
				Name:    "only-dir",
				Aliases: []string{"D"},
				Usage:   "Rename only directories (implies include-dir)",
			},
			&cli.BoolFlag{
				Name:    "hidden",
				Aliases: []string{"H"},
				Usage:   "Include hidden files and directories",
			},
			&cli.BoolFlag{
				Name:    "fix-conflicts",
				Aliases: []string{"F"},
				Usage:   "Fix any detected conflicts with auto indexing",
			},
			&cli.BoolFlag{
				Name:    "string-mode",
				Aliases: []string{"s"},
				Usage:   "Opt into string literal mode by treating find expressions as non-regex strings",
			},
		},
		UseShortOptionHandling: true,
		Action: func(c *cli.Context) error {
			op, err := newOperation(c)
			if err != nil {
				return err
			}

			return op.run()
		},
	}
}
