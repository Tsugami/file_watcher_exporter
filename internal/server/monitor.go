package server

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Monitor struct {
	Config *Config
}

type MonitorPath struct {
	Path       string
	AgeSeconds float64
}

func NewMonitorPath(path string, ageSeconds float64) MonitorPath {
	return MonitorPath{
		Path:       path,
		AgeSeconds: ageSeconds,
	}
}

func NewMonitor(config *Config) Monitor {
	return Monitor{
		Config: config,
	}
}

func (m Monitor) IsExcluded(path string) bool {
	for _, excludePattern := range m.Config.Exclude {
		matched, _ := filepath.Match(excludePattern, path)
		if matched {
			return true
		}
	}
	return false
}

func (m Monitor) IsIncluded(path string) bool {
	for _, includePattern := range m.Config.Include {
		matched, _ := filepath.Match(includePattern, path)
		if matched {
			return true
		}
	}
	return false
}

func (m Monitor) IsValid(path string) bool {
	return m.IsIncluded(path) && !m.IsExcluded(path)
}

func GetDir(path string, d os.DirEntry) string {
	if d.IsDir() {
		return path
	}

	return filepath.Dir(path)
}

func (m Monitor) CheckFilesRecursively() []MonitorPath {
	FileModificationTime.Reset()
	results := []MonitorPath{}
	dirCounter := map[string]int{}

	err := filepath.WalkDir(m.Config.Dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		valid := m.IsValid(path)
		if !valid {

			return nil
		}

		fmt.Printf("valid: %t path %q\n", valid, path)
		info, err := os.Stat(path)
		if err != nil {
			fmt.Printf("Erro ao obter informações do arquivo %q: %v\n", path, err)
		}

		dir := GetDir(path, d)
		if _, ok := dirCounter[dir]; !ok {
			dirCounter[dir] = 0
		}

		if d.IsDir() {
			return nil
		}

		dirCounter[dir]++
		fileAge := time.Since(info.ModTime())
		FileModificationTime.WithLabelValues(path).Set(fileAge.Seconds())
		results = append(results, NewMonitorPath(path, fileAge.Seconds()))

		return nil
	})

	for dir, count := range dirCounter {
		FilesCount.WithLabelValues(dir).Set(float64(count))
	}

	fmt.Printf("dirCounter: %+v\n", dirCounter)
	if err != nil {
		fmt.Printf("Erro ao percorrer o diretório %q: %v\n", m.Config.Dir, err)
	}

	return results
}
