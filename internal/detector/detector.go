package detector

import (
	"errors"
	"os"
	"path/filepath"
)

type ProjectType string

const (
	Go      ProjectType = "go"
	Python  ProjectType = "py"
	Node    ProjectType = "node"
	Unknown ProjectType = "unknown"
)

type rule struct {
	filename string
	project  ProjectType
}

var defaultRule = []rule{
	{filename: "go.mod", project: Go},
	{filename: "package.json", project: Node},
	{filename: "requirements.txt", project: Python},
	{filename: "pyproject.toml", project: Python},
}

func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {

		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}

		return false, err
	}
	return true, nil
}

func Detect(targetDir string) (ProjectType, error) {
	for _, r := range defaultRule {
		path := filepath.Join(targetDir, r.filename)
		exists, err := fileExists(path)
		if err != nil {
			return Unknown, err
		}

		if exists {
			return r.project, nil
		}
	}
	return Unknown, errors.New("file its not compatiable with the directory")
}
