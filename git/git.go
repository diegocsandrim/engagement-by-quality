package git

import (
	"fmt"
	"strings"

	"../cmd"
)

const GitBaseDir string = "/src/github.com"

type GitRepo struct {
	namespace  string
	project    string
	cmdFactory *cmd.CmdFactory
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

func (g *GitRepo) GetContributors() ([]*Contributor, error) {
	out, err := g.cmdFactory.ExecF("git log --all --no-merges --format='%%aN <%%aE>' | sort | uniq")
	if err != nil {
		return nil, err
	}

	contributorIds := strings.Split(out, "\n")
	contributors := make([]*Contributor, 0, len(contributorIds))
	for _, contributorId := range contributorIds {
		if contributorId == "" {
			continue
		}
		contributor := NewContributor(contributorId)
		contributors = append(contributors, contributor)
	}

	return contributors, nil
}

func (g *GitRepo) GetContributorFirstCommit() ([]*Contributor, error) {
	out, err := g.cmdFactory.ExecF("git log --all --no-merges --format='%%aN <%%aE>' | sort | uniq")
	if err != nil {
		return nil, err
	}

	contributorIds := strings.Split(out, "\n")
	contributors := make([]*Contributor, 0, len(contributorIds))
	for _, contributorId := range contributorIds {
		contributor := NewContributor(contributorId)
		contributors = append(contributors, contributor)
	}

	return contributors, nil
}

func (g *GitRepo) projectDir() string {
	return fmt.Sprintf("%s/%s/%s", GitBaseDir, g.namespace, g.project)
}

func (g *GitRepo) namespaceDir() string {
	return fmt.Sprintf("%s/%s", GitBaseDir, g.namespace)
}
