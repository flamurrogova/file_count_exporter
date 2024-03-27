package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/prometheus/common/version"
)

func main() {

	// listening endpoint
	endpoint := flag.String("endpoint", "localhost:8800", "The address to listen on for HTTP requests <IP>:<PORT>")
	flag.Parse()

	log.Println("Build context", "build_context", version.BuildContext())

	log.Println("endpoint:", *endpoint)
	for _, s := range flag.Args() {
		fmt.Println("arg: ", s)
	}

	if len(flag.Args()) < 1 {
		log.Fatal("missing paths")
	}

	// prometheus: create a gauge with one label named ("path")
	fileCount := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "file_count_on_path",
			Help: "Count number of files on the path.",
		},
		// The label name by which to split the metric.
		[]string{"path"},
	)

	// prometheus: create a non-global registry
	reg := prometheus.NewRegistry()
	reg.MustRegister(collectors.NewBuildInfoCollector(), fileCount)

	// fsnotify: create new watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// fsnotify: add paths passed on command line
	for _, path := range flag.Args() {
		err = watcher.Add(path)
		if err != nil {
			log.Fatal(err, " : ", path)
		}
	}

	// fsnotify: start listening for events
	// we are interested only on file create/remove events
	go func() {
		for {
			select {

			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				if event.Has(fsnotify.Create) {
					name := filepath.Dir(event.Name)
					fileCount.WithLabelValues(name).Inc()
				}

				if event.Has(fsnotify.Remove) {
					name := filepath.Dir(event.Name)
					fileCount.WithLabelValues(name).Dec()
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	// prometheus: expose the registered metrics via HTTP.
	http.Handle("/metrics", promhttp.HandlerFor(
		reg,
		promhttp.HandlerOpts{
			// Pass custom registry
			Registry: reg,
		},
	))
	log.Fatal(http.ListenAndServe(*endpoint, nil))

}

func init() {

	// Define a custom usage message
	flag.Usage = func() {
		fmt.Println("\nUsage:")
		fmt.Printf("  %s [options] [path1 path2 ...]\n", os.Args[0])
		fmt.Println("\nOptions:")
		fmt.Println("\nlist of paths")
		flag.PrintDefaults()
	}

}
