package main

import (
	"fmt"
	"math"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/jinzhu/now"
	. "github.com/logrusorgru/aurora"
	"github.com/urfave/cli"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknow"
)

var dateRe = regexp.MustCompile("[^0-9]")

func filenameWithoutExtension(fn string) string {
	return strings.TrimSuffix(fn, path.Ext(fn))
}

func main() {
	app := cli.NewApp()
	app.Name = "gfs-cleaner"
	app.Usage = "GFS Cleaner"
	app.Version = fmt.Sprintf("%s (Git: %s) %s", version, commit, date)

	app.Commands = []cli.Command{
		{
			Name:  "clean",
			Usage: "Clean a target directory based on GFS retention.",
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "daily, d",
					Value: 7,
					Usage: "Number of daily backups",
				},
				cli.IntFlag{
					Name:  "weekly, w",
					Value: 4,
					Usage: "Number of weekly backups",
				},
				cli.IntFlag{
					Name:  "monthly, m",
					Value: 12,
					Usage: "Number of monthly backups",
				},
				cli.IntFlag{
					Name:  "yearly, y",
					Value: 10,
					Usage: "Number of yearly backups",
				},
				cli.BoolFlag{
					Name:  "dry",
					Usage: "Dry run (don't remove files/folders)",
				},
			},
			Action: clean,
		},
		{
			Name:  "generate",
			Usage: "Generate test files.",
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "count, c",
					Value: 50,
					Usage: "Number of generated files/folders",
				},
				cli.StringFlag{
					Name:  "mode, m",
					Value: "files",
					Usage: "Generate 'files' or 'folders'",
				},
			},
			Action: generate,
		},
	}

	app.Run(os.Args)
}

func clean(c *cli.Context) error {
	if len(c.Args()) < 1 {
		cli.ShowCommandHelp(c, "clean")
		return cli.NewExitError("Please set working directory", 2)
	}

	numOfDaily := c.Int("daily")
	numOfWeekly := c.Int("weekly")
	numOfMonthly := c.Int("monthly")
	numOfYearly := c.Int("yearly")

	fmt.Printf("Start cleaning. Retention: %d daily, %d weekly, %d monthly and %d yearly", numOfDaily, numOfWeekly, numOfMonthly, numOfYearly)
	workDir := c.Args()[0]
	absWorkDir, _ := filepath.Abs(workDir)
	fmt.Printf("\nWorking directory: %s\n", absWorkDir)

	items := []string{}

	filepath.Walk(absWorkDir, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			if f != nil {
				fmt.Printf("Unable to traverse directory.")
			} else {
				fmt.Printf("Unable to traverse directory: %s", path)
			}
			return err
		}

		if f != nil {
			if path == absWorkDir {
				return nil
			}

			items = append(items, filepath.Base(path))

			if f.IsDir() {
				return filepath.SkipDir
			}
		}
		return nil
	})

	removableItems := []string{}

	if len(items) == 0 {
		fmt.Println("Not found backup files.")
		return nil
	}

	sort.Strings(items)
	fmt.Printf("Found %d files.\n", len(items))

	lastDay := getLastDay(items)
	fmt.Printf("Last backup day: %s\n", lastDay.Format("2006-01-02"))

	for _, item := range items {
		fmt.Printf("Check '%s'...", item)

		/*stat, err := os.Stat(p)
		if err != nil {
			fmt.Printf("Error to get file info: %s. %q\n", p, err)
			continue
		}*/

		//isDir := stat.IsDir()
		str := dateRe.ReplaceAllString(item, "")

		t, err := time.Parse("20060102", str)
		if err != nil {
			fmt.Printf("%s: Invalid date format: %s. %q\n", Red("SKIP").Bold(), str, err)
			continue
		}

		now.WeekStartDay = time.Monday
		status := "daily"
		shouldRemove := false

		// Check "daily"
		shouldRemove = dayDuration(lastDay, t) >= numOfDaily

		// Check "weekly"
		if t.Weekday() == time.Monday {
			status = "weekly"
			shouldRemove = weekDuration(lastDay, t) > numOfWeekly
		}

		// Check "monthly"
		if t.Day() == 1 {
			status = "monthly"
			shouldRemove = monthDuration(lastDay, t) > numOfMonthly
		}

		// Check "yearly"
		if t.Month() == 1 && t.Day() == 1 {
			status = "yearly"
			shouldRemove = lastDay.Year()-t.Year() >= numOfYearly
		}

		//now.BeginningOfWeek()
		if shouldRemove {
			fmt.Printf("%s %s (Day: %s)\n", Brown("DEL").Bold(), Cyan(status), t.Format("2006-01-02"))
			removableItems = append(removableItems, item)
		} else {
			fmt.Printf("%s %s (Day: %s)\n", Green("OK").Bold(), Cyan(status), t.Format("2006-01-02"))
		}
	}
	fmt.Println("")

	if !c.Bool("dry") {
		if len(removableItems) > 0 {
			fmt.Println("Cleaning folder...")
			for _, item := range removableItems {
				absPath := filepath.Join(absWorkDir, item)

				err := os.RemoveAll(absPath)
				if err != nil {
					fmt.Printf("Unable to remove %s file: %q\n", absPath, err)
					continue
				}

			}
			fmt.Printf("Removed %d old backup files.\n", len(removableItems))
		} else {
			fmt.Println("No need to remove backup files.")
		}
	} else {
		fmt.Printf("Dry run. Stay all files.")
	}

	return nil
}

func generate(c *cli.Context) error {
	day := time.Now()

	count := c.Int("count")
	mode := c.String("mode")

	fmt.Printf("Generate %d %s...\n", count, mode)
	os.Mkdir("test", 0755)

	for i := 0; i < count; i++ {
		if mode == "folders" {
			os.MkdirAll(path.Join("test", day.Format("2006-01-02")), 0755)
		} else {
			os.Create(path.Join("test", "backup_"+day.Format("2006-01-02")+".zip"))
		}
		day = day.Add(-time.Hour * 24)
	}

	return nil
}

func getLastDay(items []string) time.Time {

	var lastDay time.Time

	for _, item := range items {
		str := dateRe.ReplaceAllString(item, "")

		t, err := time.Parse("20060102", str)
		if err != nil {
			continue
		}

		if t.After(lastDay) {
			lastDay = t
		}
	}
	return lastDay
}

func monthDuration(now time.Time, past time.Time) int {
	m1 := now.Year()*12 + int(now.Month())
	m2 := past.Year()*12 + int(past.Month())

	return m1 - m2
}

func weekDuration(now time.Time, past time.Time) int {
	diff := now.Sub(past)
	return int(math.Ceil(diff.Hours() / 24 / 7))
}

func dayDuration(now time.Time, past time.Time) int {
	diff := now.Sub(past)
	return int(math.Ceil(diff.Hours() / 24))
}
