package builder

import (
	"os/exec"
	"os"
	"path/filepath"
	"errors"
	"strings"
	"fmt"
	"path"
)

type GoTools struct {
	goTool string
	toolName string
	envMap map[string]string
	buildFlags []string
	packageNames []string
	workingDir string
}

const (
	GoBuild = iota
	GoGet
	GoRun
	GoTest
)

func NewGoTools(t int, goPath string) (*GoTools, error) {
	var err error
	goTools := &GoTools{envMap:make(map[string]string)}
	if goTools.goTool, err = FindGoTool(); err != nil {
		return nil, err
	}
	goTools.SetGoPath(goPath)
	switch t {
	case GoBuild:
		goTools.toolName = "build"
	case GoGet:
		goTools.toolName = "get"
		if gitTool, err := FindGitTool(); err == nil {
			goTools.AddEnvVar("PATH", path.Dir(gitTool))
		} else {
			return nil, err
		}
	case GoRun:
		goTools.toolName = "run"
	case GoTest:
		goTools.toolName = "test"
	default:
		return nil, errors.New("tool not found")
	}
	return goTools, nil
}

func (gt *GoTools) AddEnvVar(key, val string) {
	key = strings.ToUpper(key)
	if v, ok := gt.envMap[key]; ok {
		gt.envMap[key] = strings.Join([]string{v, val}, string(os.PathListSeparator))
	} else {
		gt.envMap[key] = val
	}
}

func (gt *GoTools) envVarStr() ([]string) {
	envs := []string{}
	for k, v := range gt.envMap {
		envs = append(envs, fmt.Sprintf("%s=%s", k, v))
	}
	return envs
}

func (gt *GoTools) SetGoPath(path string) {
	gt.envMap["GOPATH"] = path
}

func (gt *GoTools) SetWorkingDir(dir string) {
	gt.workingDir = dir
}

func (gt *GoTools) AddBuildFlags(flags ...string) {
	gt.buildFlags = append(gt.buildFlags, flags...)
}

func (gt *GoTools) AddPackageNames(names ...string) {
	gt.packageNames = append(gt.packageNames, names...)
}

func (gt *GoTools) Run() (string, bool) {
	args := gt.args()
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = gt.workingDir
	cmd.Env = gt.envVarStr()
	if out, err := cmd.CombinedOutput(); err != nil {
		return string(out), false
	} else {
		return string(out), true
	}
}

func (gt *GoTools) String() string {
	return strings.Join(gt.args(), " ")
}

func (gt *GoTools) args() []string {
	args := []string{gt.goTool, gt.toolName}
	for _, buildFlag := range gt.buildFlags {
		args = append(args, buildFlag)
	}
	for _, packageName := range gt.packageNames {
		args = append(args, packageName)
	}
	return args
}

func FindGoTool() (string, error) {
	if goTool, err := exec.LookPath("go"); err == nil {
		return goTool, nil
	}
	if goBin, ok := os.LookupEnv("GOBIN"); ok {
		return filepath.Join(goBin, "go"), nil
	}
	if goRoot, ok := os.LookupEnv("GOROOT"); ok {
		return filepath.Join(goRoot, "bin", "go"), nil
	}
	return "", errors.New("environment variable $GOROOT not found")
}

func FindGitTool() (string, error) {
	if gitTool, err := exec.LookPath("git"); err == nil {
		return gitTool, nil
	}
	if gitHome, ok := os.LookupEnv("GIT_HOME"); ok {
		return filepath.Join(gitHome, "bin", "git"), nil
	}
	return "", errors.New("git not found")
}