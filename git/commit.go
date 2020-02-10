package git

import "time"

type Commit struct {
	Hash        string
	ParentHash  string
	Date        time.Time
	Contributor *Contributor
	HasGoCode   bool
}

func NewCommit(hash string, parentHash string, date time.Time, contributor *Contributor, hasGoCode bool) *Commit {
	c := Commit{
		Hash:        hash,
		ParentHash:  parentHash,
		Date:        date,
		Contributor: contributor,
		HasGoCode:   hasGoCode,
	}
	return &c
}
