package qualityanalyzers

import (
	"fmt"
	"strings"
	"time"

	"../cmd"
)

type Sonnar struct {
	projectKey    string
	sonarLogin    string
	sonnarHostUrl string
	projectDir    string
	cmdFactory    *cmd.CmdFactory
}

func NewSonnar(projectKey string, sonarLogin string, sonnarHostUrl string, projectDir string) *Sonnar {
	return &Sonnar{
		projectKey:    projectKey,
		sonarLogin:    sonarLogin,
		sonnarHostUrl: sonnarHostUrl,
		projectDir:    projectDir,
		cmdFactory:    cmd.NewCmdFactory(projectDir),
	}
}

func (s *Sonnar) Run(projectVersion string, date time.Time) error {
	projectDate := date.UTC().Format("2006-01-02")
	output, err := s.cmdFactory.ExecF(`docker run --name sonar-scanner --network host -dit -v %s:/root/src sonar-scanner:4 \
    -D sonar.host.url=%s \
    -D sonar.projectKey=%s \
    -D sonar.projectBaseDir=/root/src \
	-D sonar.login=%s \
	-D sonar.projectVersion=%s \
	-D sonar.projectDate=%s
	`, s.projectDir, s.sonnarHostUrl, s.projectKey, s.sonarLogin, projectVersion, projectDate)

	if err != nil {
		return err
	}

	output, err = s.cmdFactory.ExecF("docker wait sonar-scanner")
	if err != nil {
		return err
	}

	exitsCode := strings.Split(output, "\n")[0]
	if exitsCode != "0" {
		logs, err := s.cmdFactory.ExecF("docker wait sonar-scanner")
		if err != nil {
			return fmt.Errorf("%s: %s: %w", "sonar analyser has failed, but we could not get the logs", output, err)
		}

		return fmt.Errorf("%s: %s", "sonar analyser has failed", logs)
	}

	output, err = s.cmdFactory.ExecF("docker rm sonar-scanner")
	if err != nil {
		return err
	}

	return err
}

func FormatProjectKey(parts ...string) string {
	return strings.Join(parts, ":")
}
