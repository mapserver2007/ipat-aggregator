package file_gateway

import (
	"os"
	"path/filepath"
)

type PathOptimizer interface {
	GetProjectRoot() (string, error)
	// TODO 相対パス渡したら絶対パス返すメソッド
}

const (
	targetFileName = "go.mod"
)

func NewPathOptimizer() PathOptimizer {
	return &pathOptimizer{}
}

type pathOptimizer struct{}

func (p *pathOptimizer) GetProjectRoot() (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", err
	}

	execDir := filepath.Dir(execPath)

	for {
		modPath := filepath.Join(execDir, targetFileName)
		if _, err := os.Stat(modPath); err == nil {
			return execDir, nil
		}

		parentDir := filepath.Dir(execDir)
		if parentDir == execDir {
			break
		}

		execDir = parentDir
	}

	rootPath, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return rootPath, nil
}
