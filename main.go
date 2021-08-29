package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/mmcdole/gofeed"
)

func truncateToLength(str string, n int) string {
	if len(str) <= n {
		return str
	}
	return string(str[0:n-3]) + "..."
}

func Notify(title string, messages []string, appIcon string) error {
	var length int
	switch string(os.PathSeparator) {
	case "\\":
		length = 45
	default:
		/*
			I don't know what the notification length on Linux
			If you find it, feel free to start an issue on GitHub
			and I'll fix it.
		*/
		length = 50
	}
	title = truncateToLength(title, length)
	processed_msg := ""
	for _, message := range messages {
		processed_msg += truncateToLength(message, length) + "\n"
	}
	return beeep.Notify(title, processed_msg, "")
}

func main() {

	const DateFormat = "02-01-2006 15:04:05"
	const LogFilePath = "./log.txt"

	// TODO: Create log file if it doesn't exist.

	logFile, err := os.Open(LogFilePath)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()
	scanner := bufio.NewScanner(logFile)
	var lastRun time.Time
	lines := 0
	for scanner.Scan() {
		lastRun, err = time.Parse(DateFormat, scanner.Text())
		if err != nil {
			panic(err)
		}
		lines += 1
		/*
			We only need the first line.
			An if block is used here to avoid Go raising
			an warning that there is an unconditional break.
		*/
		if true {
			break
		}
	}
	if lines == 0 {
		lastRun = time.Now().Local()
	}

	logFile.Close()

	// Empty file
	err = os.Remove(LogFilePath)
	if err != nil {
		panic(err)
	}

	logFile, err = os.Create(LogFilePath)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	dt := time.Now().Local()
	// Write last run date to file.
	_, err = logFile.WriteString(dt.Format(DateFormat))
	if err != nil {
		panic(err)
	}

	// Modified from https://stackoverflow.com/questions/8757389/
	file, err := os.Open("./list.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var feeds = []*gofeed.Feed{}
	fp := gofeed.NewParser()
	scanner = bufio.NewScanner(file)
	// For each line in "list.txt"
	for scanner.Scan() {
		feed, err := fp.ParseURL(scanner.Text())
		if err != nil {
			panic(err)
		}
		feeds = append(feeds, feed)
	}

	var missed_items = []*gofeed.Item{}
	for _, feed := range feeds {
		for i := 0; i < len(feed.Items); i++ {
			item := feed.Items[i]
			date := item.PublishedParsed
			if date == nil {
				continue
			}
			if date.After(lastRun) {
				missed_items = append(missed_items, item)
			}
		}
	}

	if len(missed_items) == 0 {
		// Do nothing
	} else if len(missed_items) <= 3 {
		messages := []string{}
		for i, item := range missed_items {
			messages = append(messages, fmt.Sprint(fmt.Sprint(i+1)+".", item.Title))
		}
		Notify(fmt.Sprint("You have missed ", len(missed_items), " items."), messages, "")
	} else {
		messages := []string{}
		for i, item := range missed_items[0:2] {
			messages = append(messages, fmt.Sprint(fmt.Sprint(i+1)+".", item.Title, "-", item.Author))
		}
		messages = append(messages, "and "+fmt.Sprint(len(missed_items)-2)+" more...")
		Notify(fmt.Sprint("You have missed ", len(missed_items), " items."), messages, "")
	}

}
