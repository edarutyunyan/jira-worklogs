package main

import (
	"fmt"
	"jira-worklogs/pkg/services"
	"log"
	"os"
	"strings"
	"time"

	jira "github.com/andygrunwald/go-jira"
	"github.com/joho/godotenv"
)

var ResultsPerPage int = 50
var TotalCount int = -1

func Seconds2hm(seconds int) string {
	minutes := seconds / 60
	return fmt.Sprintf("%dh %dm", minutes/60, minutes%60)
}

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
		return
	}

	apiUrls := strings.Split(os.Getenv("API_URL"), ",")
	apiKey := os.Getenv("JIRA_TOKEN")
	user := os.Getenv("API_USER")
	workLogUserId := os.Getenv("WORKLOG_USER_ID")

	tp := jira.BasicAuthTransport{
		Username: user,
		Password: apiKey,
	}

	srv := services.NewServices(*services.NewInputService(), *services.NewOutputService())

	arguments, err := srv.Input.GetArgs()

	if err != nil {
		log.Fatal(err.Error())
	}

	search := fmt.Sprintf(
		"worklogAuthor = %s and worklogDate >= %s",
		workLogUserId, arguments[0].Format("2006-01-02"),
	)

	fmt.Printf("WORKLOG SINCE: %s \n\n", arguments[0].Format(time.RFC1123))

	if endDate, ok := arguments[1]; ok {
		search = search + fmt.Sprintf(" and worklogDate <= %s", endDate.Format("2006-01-02"))
		fmt.Println("UNTIL:", endDate.Format(time.RFC1123))
	}

	var sumInSeconds int

	for _, url := range apiUrls {
		client, err := jira.NewClient(tp.Client(), url)

		if err != nil {
			log.Fatal("Failed create a client")
			return
		}
		lastIssue := 0

	lpGetIssues:
		for {
			issues, resp, err := client.Issue.Search(search, &jira.SearchOptions{
				StartAt:    lastIssue,
				MaxResults: ResultsPerPage,
			})

			if err != nil {
				log.Fatal("Failed fetch a board", err.Error())
				return
			}
			if resp.Total == 0 {
				break
			}
			if TotalCount < 0 || TotalCount > resp.Total {
				TotalCount = resp.Total
			}
			lastIssue = resp.StartAt

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
						fmt.Print(srv.Output.FormattedStdOutString(issue, wl))
						sumInSeconds += wl.TimeSpentSeconds
					}
				}

				lastIssue++
				// to see pages on stderr:
				// fmt.Fprintf(os.Stderr, "%d/%d\n", lastIssue, TotalCount)
				if lastIssue >= TotalCount {
					break lpGetIssues
				}
			}
		}
	}

	fmt.Printf("\nTotal: %s\n", Seconds2hm(sumInSeconds))
}
