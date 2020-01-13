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
	g.cmdFactory = cmd.NewCmdFactory(g.ProjectDir())

	return &g
}

func (g *GitRepo) ForceClone() error {
	_, err := cmd.NewCmdFactory("/").ExecF("mkdir -p %s", g.namespaceDir())
	if err != nil {
		return err
	}

	_, err = cmd.NewCmdFactory("/").ExecF("rm -rf %s", g.ProjectDir())
	if err != nil {
		return err
	}

	_, err = cmd.NewCmdFactory(g.namespaceDir()).ExecF("git clone https://github.com/%s/%s.git", g.namespace, g.project)
	return err
}

func (g *GitRepo) LoadCommits() error {
	g.commits = make(map[string]*Commit)
	g.contributors = make(map[string]*Contributor)

	commitsLog, err := g.cmdFactory.ExecF("git log --all --format='%%H %%at %%aE %%P'")
	if err != nil {
		return err
	}

	commitsLogLines := strings.Split(commitsLog, "\n")
	for _, commitLog := range commitsLogLines {
		if commitLog == "" {
			continue
		}
		splittedLogLine := strings.Split(commitLog, " ")
		if len(splittedLogLine) < 4 {
			log.Panicf("commits log line is in a bad format: '%s'", commitLog)
		}

		commitHash := splittedLogLine[0]
		commitTimestampString := splittedLogLine[1]
		contributorId := splittedLogLine[2]
		parentCommitHashs := splittedLogLine[3:]

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

		commit := NewCommit(commitHash, parentCommitHashs[0], commitTimestamp, contributor)
		contributor.AddCommit(commit)

		g.commits[commit.Hash] = commit
	}
	return nil
}

func (g *GitRepo) ProjectDir() string {
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

func (g *GitRepo) ContributorAttractorCommits() []*ContributorAttractorCommit {
	contributorAttractorCommitsByCommitHash := make(map[string]*ContributorAttractorCommit)

	for _, contributor := range g.contributors {
		contributorFirstCommit := contributor.FirstCommit()
		if contributorFirstCommit.ParentHash == "" {
			continue
		}
		parentCommit := g.commits[contributorFirstCommit.ParentHash]
		if parentCommit == nil {
			log.Println("fail!")
		}
		contributorAttractorCommit, exists := contributorAttractorCommitsByCommitHash[parentCommit.Hash]
		if !exists {
			contributorAttractorCommit = NewContributorAttractorCommit(parentCommit)
			contributorAttractorCommitsByCommitHash[parentCommit.Hash] = contributorAttractorCommit
		}

		contributorAttractorCommit.AddAttractedContributor(contributor)
	}

	contributorAttractorCommits := make([]*ContributorAttractorCommit, 0, len(contributorAttractorCommitsByCommitHash))
	for _, contributorAttractorCommit := range contributorAttractorCommitsByCommitHash {
		contributorAttractorCommits = append(contributorAttractorCommits, contributorAttractorCommit)
	}

	return contributorAttractorCommits
}

func (g *GitRepo) Checkout(ref string) error {
	_, err := g.cmdFactory.ExecF("git checkout %s", ref)
	return err
}
