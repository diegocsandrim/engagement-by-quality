package git

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

const GitBaseDir string = "/src/github.com"

func EnsureGitBaseDirExists() error {
	cmd := exec.Command("mkdir", "-p", GitBaseDir)
	_, err := execWithLog(cmd)
	return err
}

func Clone(namespace string, project string) error {
	cmd := exec.Command("bash", "-c", fmt.Sprintf("mkdir -p %s", namespaceDir(namespace)))
	_, err := execWithLog(cmd)
	if err != nil {
		return err
	}

	cmd = exec.Command("bash", "-c", fmt.Sprintf("rm -rf %s", projectDir(namespace, project)))
	_, err = execWithLog(cmd)
	if err != nil {
		return err
	}

	cmd = exec.Command("git", "clone", fmt.Sprintf("https://github.com/%s/%s.git", namespace, project))
	cmd.Dir = namespaceDir(namespace)
	_, err = execWithLog(cmd)
	return err
}

func GetContributors(namespace string, project string) ([]*Contributor, error) {
	cmd := exec.Command("bash", "-c", "git log --all --no-merges --format='%aN <%aE>' | sort | uniq")
	cmd.Dir = projectDir(namespace, project)
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

func GetContributorFirstCommit(namespace string, project string) ([]*Contributor, error) {
	cmd := exec.Command("bash", "-c", "git log --all --no-merges --format='%aN <%aE>' | sort | uniq")
	cmd.Dir = projectDir(namespace, project)
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

func projectDir(namespace string, project string) string {
	return fmt.Sprintf("%s/%s/%s", GitBaseDir, namespace, project)
}

func namespaceDir(namespace string) string {
	return fmt.Sprintf("%s/%s", GitBaseDir, namespace)
}
