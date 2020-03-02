package git

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/diegocsandrim/engagement-by-quality/cmd"
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

func (g *GitRepo) Clone() error {
	clearCommand := `git reset --hard HEAD && 
git clean -f -d &&
git checkout master &&
git fetch origin master &&
git reset --hard origin/master &&
git pull`

	_, err := cmd.NewCmdFactory(g.ProjectDir()).ExecF(clearCommand)
	if err != nil {
		log.Printf("failed to clone: %s. Will try to force clone", err.Error())
		return g.ForceClone()
	}
	return nil
}

func (g *GitRepo) hasGoCode(fileNames []string) (bool, error) {
	for _, fileName := range fileNames {
		if strings.HasSuffix(fileName, ".go") {
			return true, nil
		}
	}

	return false, nil
}

func (g *GitRepo) LoadCommits() error {
	g.commits = make(map[string]*Commit)
	g.contributors = make(map[string]*Contributor)

	commitLinePrefix := "commit:"
	commitsLog, err := g.cmdFactory.ExecF("git log --all --format='%s%%H/%%at/%%aE/%%P' --reverse --name-only", commitLinePrefix)
	if err != nil {
		return err
	}

	commitId := 0
	commitsLogLines := strings.Split(commitsLog, "\n")

	for line := 0; line < len(commitsLogLines); line++ {
		commitLog := commitsLogLines[line]

		if commitLog == "" {
			continue
		}
		splittedLogLine := strings.Split(commitLog, "/")
		if len(splittedLogLine) != 4 {
			log.Panicf("commits log line is in a bad format: '%s'", commitLog)
		}

		commitHash := splittedLogLine[0][len(commitLinePrefix):]
		commitTimestampString := splittedLogLine[1]
		contributorId := splittedLogLine[2]
		parentCommitHashs := strings.Split(splittedLogLine[3], " ")

		commitFileNames := make([]string, 0)

		for {
			if line+1 == len(commitsLogLines) {
				break
			}
			nextLine := commitsLogLines[line+1]
			if nextLine == "" {
				line++
				continue
			}

			if strings.HasPrefix(nextLine, commitLinePrefix) {
				break
			}

			if strings.HasPrefix(nextLine, "warning: ") {
				line++
				continue
			}

			commitFileNames = append(commitFileNames, nextLine)
			line++
		}

		hasGoCode, err := g.hasGoCode(commitFileNames)
		if err != nil {
			return err
		}

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

		commit := NewCommit(commitId, commitHash, parentCommitHashs[0], commitTimestamp, contributor, hasGoCode)
		commitId++

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
		contributorFirstGoCommit := contributor.FirstGoCommit()
		if contributorFirstGoCommit == nil {
			continue
		}
		if contributorFirstGoCommit.ParentHash == "" {
			continue
		}

		parentCommit := g.commits[contributorFirstGoCommit.ParentHash]
		if parentCommit == nil {
			log.Printf("Missing required parent commit! parent hash: %s", contributorFirstGoCommit.ParentHash)
		}

		if !contributor.IsMainContributor() {
			//Find a parent commit that is authored by other contributor
			for parentCommit.Contributor == contributor {
				parentCommit = g.commits[parentCommit.ParentHash]
			}
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
	_, err := g.cmdFactory.ExecF("git checkout --force %s", ref)
	return err
}
