package sonnarplugins

import (
	"fmt"

	"github.com/diegocsandrim/engagement-by-quality/cmd"
)

type GolangCILintPlugin struct {
	cmdFactory *cmd.CmdFactory
}

func NewGolanCILintPlugin(cmdFactory *cmd.CmdFactory) *GolangCILintPlugin {
	return &GolangCILintPlugin{
		cmdFactory: cmdFactory,
	}
}

func (p *GolangCILintPlugin) Setup() error {
	_, err := p.cmdFactory.ExecF("mkdir -p ./72bd2909-e59c-459a-a739-1a7660e9d67d")
	if err != nil {
		return fmt.Errorf("failed to create report path for plugin golangci-lint: %w", err)
	}

	return err
}

func (p *GolangCILintPlugin) Run() error {
	output, err := p.cmdFactory.ExecF(`docker run --rm -v $(pwd):/app -w /app golangci/golangci-lint:v1.23.3 \
golangci-lint run \
--out-format checkstyle \
--max-issues-per-linter=0 \
--timeout=5m \
--max-same-issues=0 \
--uniq-by-line=false \
--issues-exit-code=0 \
> ./72bd2909-e59c-459a-a739-1a7660e9d67d/report.xml`)
	if err != nil {
		if !p.canIgnoreGolangciError(output, err) {
			return fmt.Errorf("failed to run golangci-lint: %w", err)
		}
	}
	return nil
}

func (g *GolangCILintPlugin) canIgnoreGolangciError(reportError string, err error) bool {
	return true
}

func (g *GolangCILintPlugin) ReportFlag() string {
	return "-D sonar.go.golangci-lint.reportPaths=/root/src/72bd2909-e59c-459a-a739-1a7660e9d67d/report.xml"
}
