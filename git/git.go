package git

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"../cmd"
)

const GitBaseDir string = "/src/github.com"

type GitRepo struct {
	namespace    string
	project      string
	cmdFactory   *cmd.CmdFactory
	commits      map[string]*Commit
	contributors map[string]*Contributor
}

func NewGitRepo(namespace string, project string) *GitRepo {
	g := GitRepo{
		namespace: namespace,
		project:   project,
	}
	g.cmdFactory = cmd.NewCmdFactory(g.projectDir())

	return &g
}

func (g *GitRepo) ForceClone() error {
	_, err := cmd.NewCmdFactory("/").ExecF("mkdir -p %s", g.namespaceDir())
	if err != nil {
		return err
	}

	_, err = cmd.NewCmdFactory("/").ExecF("rm -rf %s", g.projectDir())
	if err != nil {
		return err
	}

	_, err = cmd.NewCmdFactory(g.namespaceDir()).ExecF("git clone https://github.com/%s/%s.git", g.namespace, g.project)
	return err
}

func (g *GitRepo) LoadCommits() error {
	g.commits = make(map[string]*Commit)
	g.contributors = make(map[string]*Contributor)

	commitsLog, err := g.cmdFactory.ExecF("git log --all --no-merges --format='%%H %%P %%at %%aE'")
	if err != nil {
		return err
	}

	commitsLogLines := strings.Split(commitsLog, "\n")
	for _, commitLog := range commitsLogLines {
		if commitLog == "" {
			continue
		}
		splittedLogLine := strings.Split(commitLog, " ")
		if len(splittedLogLine) != 4 {
			log.Panicf("commits log line is in a bad format: '%s'", commitLog)
		}

		commitHash := splittedLogLine[0]
		parentCommitHash := splittedLogLine[1]
		commitTimestampString := splittedLogLine[2]
		contributorId := splittedLogLine[3]

		contributor, contributorExists := g.contributors[contributorId]
		if !contributorExists {
			contributor = NewContributor(contributorId)
			g.contributors[contributor.Id] = contributor
		}

		commitTimestampInt, err := strconv.ParseInt(commitTimestampString, 10, 64)
		if err != nil {
			log.Panicf("commits log timestamp is in a bad format: '%s'", commitTimestampString)
		}

		commitTimestamp := time.Unix(commitTimestampInt, 0)

		commit := NewCommit(commitHash, parentCommitHash, commitTimestamp, contributor)
		contributor.AddCommit(commit)

		g.commits[commit.Hash] = commit
	}
	return nil
}

func (g *GitRepo) projectDir() string {
	return fmt.Sprintf("%s/%s/%s", GitBaseDir, g.namespace, g.project)
}

func (g *GitRepo) namespaceDir() string {
	return fmt.Sprintf("%s/%s", GitBaseDir, g.namespace)
}

func (g *GitRepo) Contributors() []*Contributor {
	contributors := make([]*Contributor, 0, len(g.contributors))

	for _, contributor := range g.contributors {
		contributors = append(contributors, contributor)
	}
	return contributors
}
