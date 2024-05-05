package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/prometheus/common/version"
)

type myCollector struct {
	fileCount *prometheus.Desc
}

func newMyCollector() *myCollector {
	return &myCollector{
		fileCount: prometheus.NewDesc("file_count",
			"Count files in a folder",
			[]string{"path"}, nil,
		),
	}
}

func (collector *myCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.fileCount
}

func (collector *myCollector) Collect(ch chan<- prometheus.Metric) {

	for _, path := range paths {
		fileCnt := walkDir(path)
		// 'path' is also metrics label
		m := prometheus.MustNewConstMetric(collector.fileCount, prometheus.GaugeValue, float64(fileCnt), path)

		ch <- m
	}
}

// we will be monitoring these paths(folders)
var (
	paths []string
)

func main() {

	// listener endpoint
	endpoint := flag.String("endpoint", "localhost:8800", "The address to listen on for HTTP requests <IP>:<PORT>")
	flag.Parse()

	log.Println("Build context", "build_context", version.BuildContext())
	log.Println("endpoint:", *endpoint)

	if len(flag.Args()) < 1 {
		flag.Usage()
		log.Fatal("missing paths")
	}

	paths = append(paths, flag.Args()...)

	myColl := newMyCollector()
	prometheus.MustRegister(myColl)

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*endpoint, nil))
}

func init() {

	// Define a custom usage message
	// ... --base-path [] --exclude-path [] --path [single path]
	flag.Usage = func() {
		fmt.Println("\nUsage:")
		fmt.Printf("  %s [options] [path1 path2 ...]\n", os.Args[0])
		fmt.Println("\nOptions:")
		fmt.Println("\nlist of paths")
		flag.PrintDefaults()
	}

}

func walkDir(path string) int {
	i := 0
	files, err := os.ReadDir(path)
	if err != nil {
		log.Printf("Error reading directory: %v", err)
		return i
	}

	for _, file := range files {
		if file.Type().IsRegular() {
			i++
		}
	}
	return i
}
