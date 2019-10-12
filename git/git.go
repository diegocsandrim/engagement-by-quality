package git

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

const GitBaseDir string = "/src/github.com"

type GitRepo struct {
	namespace string
	project   string
}

func NewGitRepo(namespace string, project string) *GitRepo {
	g := GitRepo{
		namespace: namespace,
		project:   project,
	}

	return &g
}

func (g *GitRepo) ForceClone() error {
	cmd := exec.Command("bash", "-c", fmt.Sprintf("mkdir -p %s", g.namespaceDir()))
	_, err := execWithLog(cmd)
	if err != nil {
		return err
	}

	cmd = exec.Command("bash", "-c", fmt.Sprintf("rm -rf %s", g.projectDir()))
	_, err = execWithLog(cmd)
	if err != nil {
		return err
	}

	cmd = exec.Command("git", "clone", fmt.Sprintf("https://github.com/%s/%s.git", g.namespace, g.project))
	cmd.Dir = g.namespaceDir()
	_, err = execWithLog(cmd)
	return err
}

func (g *GitRepo) GetContributors() ([]*Contributor, error) {
	cmd := exec.Command("bash", "-c", "git log --all --no-merges --format='%aN <%aE>' | sort | uniq")
	cmd.Dir = g.projectDir()
	out, err := execWithLog(cmd)
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
	cmd := exec.Command("bash", "-c", "git log --all --no-merges --format='%aN <%aE>' | sort | uniq")
	cmd.Dir = g.projectDir()
	out, err := execWithLog(cmd)
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

func execWithLog(cmd *exec.Cmd) (string, error) {
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("command: %s", cmd.String())
		log.Printf("output: %s", out)
		return string(out), err
	}
	return string(out), nil
}

func (g *GitRepo) projectDir() string {
	return fmt.Sprintf("%s/%s/%s", GitBaseDir, g.namespace, g.project)
}

func (g *GitRepo) namespaceDir() string {
	return fmt.Sprintf("%s/%s", GitBaseDir, g.namespace)
}
