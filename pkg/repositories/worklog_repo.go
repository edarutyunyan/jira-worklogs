package repositories

type WorklogRecord struct {
	Id               string
	IssueId          string
	ProjectUrl       string
	AuthorFullName   string
	AuthorAccountId  string
	TimeSpentSeconds int
	TimeSpent        string // 1h 15m
	Comment          string
}
