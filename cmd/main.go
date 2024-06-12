package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	jira "github.com/andygrunwald/go-jira"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
		return
	}

	var arguments = map[int]time.Time{
		0: time.Now(),
	}

	for i, arg := range os.Args[1:] {
		arguments[i], err = time.Parse("2006-01-02", arg)

		if err != nil {
			log.Fatal("Date should be formatted as YYYY-MM-DD")
			return
		}

		if i == 1 && arguments[i].Before(arguments[0]) {
			log.Fatal("Dates must be in ascending order")
			return
		}
	}

	apiUrls := strings.Split(os.Getenv("API_URL"), ",")
	apiKey := os.Getenv("JIRA_TOKEN")
	user := os.Getenv("API_USER")
	workLogUserId := os.Getenv("WORKLOG_USER_ID")

	tp := jira.BasicAuthTransport{
		Username: user,
		Password: apiKey,
	}

	search := fmt.Sprintf(
		"worklogAuthor = %s and worklogDate >= %s",
		workLogUserId, arguments[0].Format("2006-01-02"),
	)

	fmt.Println("WORKLOG SINCE:", arguments[0].Format(time.RFC1123))

	if endDate, ok := arguments[1]; ok {
		search = search + fmt.Sprintf(" and worklogDate <= %s", endDate.Format("2006-01-02"))
		fmt.Println("UNTIL:", endDate.Format(time.RFC1123))
	}

	sumInSeconds := 0.00

	for _, url := range apiUrls {
		client, err := jira.NewClient(tp.Client(), url)

		if err != nil {
			log.Fatal("Failed create a client")
			return
		}

		issues, _, err := client.Issue.Search(search, &jira.SearchOptions{})

		if err != nil {
			log.Fatal("Failed fetch a board", err.Error())
			return
		}

		length := 0

		for _, issue := range issues {
			worklog, _, _ := client.Issue.GetWorklogs(issue.ID)

			for _, wl := range worklog.Worklogs {
				if wl.Author.AccountID != workLogUserId {
					continue
				}

				started := time.Time(*wl.Started)
				endDate, endDateDefined := arguments[1]

				isBefore := true

				if endDateDefined {
					isBefore = started.Before(endDate)
				}

				if started.After(arguments[0]) && isBefore { // logged for the date (not when you logged, but for the date you logged)
					length++
					fmt.Println(issue.Key, time.Time(*wl.Started).Local().Format(time.RFC1123), "Author:", wl.Author.DisplayName)
					fmt.Println(wl.Comment)
					fmt.Println("TIME SPENT:", wl.TimeSpent)
					fmt.Println("______________________________________________________")

					sumInSeconds += float64(wl.TimeSpentSeconds)
				}

			}
		}
	}

	duration := time.Duration(sumInSeconds * float64(time.Second))

	fmt.Printf("\nTotal: %.0fh %.0fm \n", duration.Hours(), duration.Hours()*60-duration.Minutes())
}
