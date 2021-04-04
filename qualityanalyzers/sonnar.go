package qualityanalyzers

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/diegocsandrim/engagement-by-quality/cmd"
	"github.com/diegocsandrim/engagement-by-quality/qualityanalyzers/sonnarplugins"
)

type Sonnar struct {
	projectKey    string
	sonarLogin    string
	sonnarHostUrl string
	sonnarPlugins []sonnarplugins.SonnarPlugin
	projectDir    string
	cmdFactory    *cmd.CmdFactory
}

func CreateSonnarAnalyser(projectKey string, sonarLogin string, sonnarHostUrl string, sonnarPlugins []string, projectDir string) (*Sonnar, error) {
	analyser := Sonnar{
		projectKey:    projectKey,
		sonarLogin:    sonarLogin,
		sonnarHostUrl: sonnarHostUrl,
		projectDir:    projectDir,
		cmdFactory:    cmd.NewCmdFactory(projectDir),
	}

	plugins, err := sonnarplugins.CreatePlugins(analyser.cmdFactory, sonnarPlugins)
	if err != nil {
		return nil, err
	}
	analyser.sonnarPlugins = plugins

	return &analyser, nil
}

func (s *Sonnar) Run(projectVersion string, date time.Time, attractedContributors int) error {
	reportFlags := make([]string, 0, len(s.sonnarPlugins))
	for _, plugin := range s.sonnarPlugins {
		err := plugin.Run()
		if err != nil {
			return err
		}
		reportFlags = append(reportFlags, plugin.ReportFlag())
	}
	reportFlagsParam := strings.Join(reportFlags, " \\\n")

	projectDate := date.UTC().Format("2006-01-02")
	output, err := s.cmdFactory.ExecF(`docker run --name sonar-scanner --network host -dit -v %s:/root/src diegocsandrim/sonar-scanner:4.2 \
	-D sonar.scm.disabled=True \
    -D sonar.host.url=%s \
    -D sonar.projectKey=%s \
    -D sonar.projectBaseDir=/root/src \
	-D sonar.login=%s \
	-D sonar.projectVersion=%s \
	-D sonar.projectDate=%s \
	-D sonar.analysis.contributorGain=%d \
	-D sonar.exclusions=**/vendor/**,**/*.pb.go,**/*generated*.go,**/*.cs,**/*.css,**/*.less,**/*.scss,**/*as,**/*.html,**/*.xhtml,**/*.cshtml,**/*.vbhtml,**/*.aspx,**/*.ascx,**/*.rhtml,**/*.erb,**/*.shtm,**/*.shtml,**/*.jsp,**/*.jspf,**/*.jspx,**/*.java,**/*.jav,**/*.js,**/*.jsx,**/*.vue,**/*.kt,**/*php,**/*php3,**/*php4,**/*php5,**/*phtml,**/*inc,**/*py,**/*.rb,**/*.scala,**/*.ts,**/*.tsx,**/*.vb,**/*.xml,**/*.xsd,**/*.xsl,**/*_gen.go \
	%s \
	`, s.projectDir, s.sonnarHostUrl, s.projectKey, s.sonarLogin, projectVersion, projectDate, attractedContributors, reportFlagsParam)

	if err != nil {
		return fmt.Errorf("failed to start scanner: %w", err)
	}

	output, err = s.cmdFactory.ExecF("docker wait sonar-scanner")
	if err != nil {
		return fmt.Errorf("failed waiting scanner to finish: %w", err)
	}
	defer func() {
		output, err = s.cmdFactory.ExecF("docker rm sonar-scanner")
		if err != nil {
			log.Printf("failed to remove scanner: %s", err.Error())
			log.Printf("output: %s", output)
		}
	}()

	exitsCode := strings.Split(output, "\n")[0]
	if exitsCode != "0" {
		logs, err := s.cmdFactory.ExecF("docker logs sonar-scanner")
		if err != nil {
			return fmt.Errorf("%s: %s: %w", "sonar analyser has failed, but we could not get the logs", output, err)
		}

		return fmt.Errorf("%s: %s", "sonar analyser has failed", logs)
	}

	return err
}

func FormatProjectKey(parts ...string) string {
	return strings.Join(parts, ":")
}

func CollectMetrics() {
	// TODO: use the SQL to colect the metrics after all analysis has run
	// SELECT projects.name, metrics.name, project_measures.value
	// FROM
	// 	project_measures
	// 	inner join metrics on project_measures.metric_id=metrics.id
	// 	inner join projects on project_measures.component_uuid = projects.uuid
	// where metrics.name like 'cognitive_complexity'
	// order by project_measures.id;
}

func (s *Sonnar) Close() {
	s.cleanTempDirs()
}

func (s *Sonnar) cleanTempDirs() {
	output, err := s.cmdFactory.ExecF(`docker run -i --rm --network host -v %s:/root/src --entrypoint='' diegocsandrim/sonar-scanner:4.2 \
	rm -rf /root/src/.scannerwork \
	`, s.projectDir)
	if err != nil {
		log.Printf("failed to cleanup scanner: %w: %s", err, output)
		return
	}
}
