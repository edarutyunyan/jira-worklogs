package main

import (
	"fmt"
	"log"
	"os"
	"time"

	jira "github.com/andygrunwald/go-jira"
	"github.com/joho/godotenv"
)

var START_DATE = time.Date(2024, 05, 15, 0, 0, 0, 0, time.Now().Location())
var ASSIGNEE_ID = "5f964829048052006bd12869"

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
		return
	}

	apiUrl := os.Getenv("API_URL")
	apiKey := os.Getenv("JIRA_TOKEN")
	user := os.Getenv("API_USER")

	tp := jira.BasicAuthTransport{
		Username: user,
		Password: apiKey,
	}

	client, err := jira.NewClient(tp.Client(), apiUrl)

	if err != nil {
		log.Fatal("Failed create a client")
		return
	}

	issues, _, err := client.Issue.Search(fmt.Sprintf("assignee = %s and worklogDate >= %s", ASSIGNEE_ID, START_DATE.Format("2006-01-02")), &jira.SearchOptions{})

	if err != nil {
		log.Fatal("Failed fetch a board", err.Error())
		return
	}

	sumInHours := 0.00

	for _, issue := range issues {
		worklog, _, _ := client.Issue.GetWorklogs(issue.ID)

		for _, wl := range worklog.Worklogs {

			if time.Time(*wl.Updated).After(START_DATE) {
				fmt.Println(issue.Key, wl.TimeSpent)
				sumInHours += float64(wl.TimeSpentSeconds) / 60 / 60
			}

		}
	}

	fmt.Print(float32(sumInHours))
}

// my id = 5f964829048052006bd12869

// assignee = 5f964829048052006bd12869 and Sprint = 'RevMax Sprint 72' and worklogDate >= '2024/05/13'

// todo: need to research how the timezones works
