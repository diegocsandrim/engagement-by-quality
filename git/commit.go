package git

import "time"

type Commit struct {
	Hash        string
	ParentHash  string
	Date        time.Time
	Contributor *Contributor
}

func NewCommit(hash string, parentHash string, date time.Time, contributor *Contributor) *Commit {
	c := Commit{
		Hash:        hash,
		ParentHash:  parentHash,
		Date:        date,
		Contributor: contributor,
	}
	return &c
}
