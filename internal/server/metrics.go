package server

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	FileModificationTime = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "file_modification_time",
			Help: "Gouge to show the file modification time.",
		},
		[]string{"filepath"},
	)

	FilesCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "files_count",
			Help: "Gouge to show the number of files in a directory.",
		},
		[]string{"directory"},
	)
)

func RegisterMetrics() {
	prometheus.MustRegister(FileModificationTime)
	prometheus.MustRegister(FilesCount)
}
