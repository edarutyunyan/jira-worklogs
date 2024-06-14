package services

import (
	"fmt"
	"strings"
	"time"

	"github.com/andygrunwald/go-jira"
)

type OutputService struct{}

func NewOutputService() *OutputService {
	return &OutputService{}
}

const WORKLOG_STARTED_DATE_FORMAT = "02 Jan 2006 15:04 (Timezone: MST)"

func (output *OutputService) FormattedStdOutString(issue jira.Issue, wl jira.WorklogRecord) string {
	outputStringsArray := []string{
		fmt.Sprintf("Ticket: %s, Date: %s, Author: %s", issue.Key, time.Time(*wl.Started).Local().Format(WORKLOG_STARTED_DATE_FORMAT), wl.Author.DisplayName),
		fmt.Sprintf("COMMENT: %s", wl.Comment),
		fmt.Sprintf("TIME SPENT: %s", wl.TimeSpent),
		"______________________________________________________\n\n",
	}

	return strings.Join(outputStringsArray, "\n")
}
