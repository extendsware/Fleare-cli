package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

type Executor struct {
	config *Config
}

func NewExecutor(config *Config) *Executor {
	return &Executor{config: config}
}

func (e *Executor) BuildTarget(targetName string) error {
	var target BuildTarget
	var ok bool

	if targetName == "default" {
		target = e.config.Build.Default
	} else {
		target, ok = e.config.Build.Targets[targetName]
		if !ok {
			return fmt.Errorf("unknown build target: %s", targetName)
		}
	}

	outputPath, err := e.renderTemplate(target.Output, map[string]string{
		"project": e.config.Project,
		"version": e.config.Version,
	})
	if err != nil {
		return fmt.Errorf("error rendering output path: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("error creating output directory: %w", err)
	}

	cmd := exec.Command("go", "build")
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("GOOS=%s", target.OS),
		fmt.Sprintf("GOARCH=%s", target.Arch),
	)
	if target.Cgo {
		cmd.Env = append(cmd.Env, "CGO_ENABLED=1")
	} else {
		cmd.Env = append(cmd.Env, "CGO_ENABLED=0")
	}

	if target.Ldflags != "" {
		ldflags, err := e.renderTemplate(target.Ldflags, map[string]string{
			"version": e.config.Version,
		})
		if err != nil {
			return fmt.Errorf("error rendering ldflags: %w", err)
		}
		cmd.Args = append(cmd.Args, "-ldflags", ldflags)
	}

	if len(target.Tags) > 0 {
		cmd.Args = append(cmd.Args, "-tags", strings.Join(target.Tags, ","))
	}

	cmd.Args = append(cmd.Args, "-o", outputPath, ".")

	fmt.Printf("Building for %s/%s to %s...\n", target.OS, target.Arch, outputPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("build failed: %w", err)
	}

	fmt.Printf("Build successful: %s\n", outputPath)
	return nil
}

func (e *Executor) RunTasks(taskName string) error {
	tasks, ok := e.config.Tasks[taskName]
	if !ok {
		return fmt.Errorf("unknown task: %s", taskName)
	}

	for _, task := range tasks {
		expandedTask, err := e.renderTemplate(task, map[string]string{
			"project": e.config.Project,
			"version": e.config.Version,
		})
		if err != nil {
			return fmt.Errorf("error rendering task: %w", err)
		}

		fmt.Printf("Running task: %s\n", expandedTask)
		cmd := exec.Command("sh", "-c", expandedTask)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("task failed: %w", err)
		}
	}
	return nil
}

func (e *Executor) RunPreBuildTasks() error {
	if _, ok := e.config.Tasks["pre-build"]; ok {
		return e.RunTasks("pre-build")
	}
	return nil
}

func (e *Executor) RunPostBuildTasks() error {
	if _, ok := e.config.Tasks["post-build"]; ok {
		return e.RunTasks("post-build")
	}
	return nil
}

func (e *Executor) renderTemplate(templateStr string, data map[string]string) (string, error) {
	tmpl, err := template.New("").Parse(templateStr)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
