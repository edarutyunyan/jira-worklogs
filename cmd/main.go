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

var START_DATE = time.Date(2024, 05, 1, 0, 0, 0, 0, time.Now().Location())

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

	sumInHours := 0.00

	fmt.Println("WORKLOG SINCE:", START_DATE.Format(time.RFC1123))

	for _, v := range apiUrls {
		client, err := jira.NewClient(tp.Client(), v)

		if err != nil {
			log.Fatal("Failed create a client")
			return
		}

		issues, _, err := client.Issue.Search(fmt.Sprintf("worklogAuthor = %s and worklogDate >= %s", workLogUserId, START_DATE.Format("2006-01-02")), &jira.SearchOptions{})

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

				if time.Time(*wl.Updated).After(START_DATE) {
					length++
					fmt.Println(issue.Key, time.Time(*wl.Updated).Format(time.RFC1123), "Author:", wl.Author.DisplayName)
					fmt.Println(wl.Comment)
					fmt.Println("TIME SPENT:", wl.TimeSpent)
					fmt.Println("______________________________________________________")

					sumInHours += float64(wl.TimeSpentSeconds) / 60 / 60
				}

			}
		}
	}

	fmt.Printf("\nTotal: %dh %.0fm \n", int(sumInHours), 60*(sumInHours-float64(int(sumInHours))))
}
