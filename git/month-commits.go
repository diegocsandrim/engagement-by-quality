package git

import "time"

type MonthCommits struct {
	Month   MonthYear
	Commits []*Commit
}

type MonthYear struct {
	Month time.Month
	Year  int
}
