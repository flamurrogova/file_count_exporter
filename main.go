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

var (

	// Create a gauge with one label named ("path").
	fileCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "file_count_on_path",
			Help: "Count number of files on the path.",
		},
		// The label name by which to split the metric.
		[]string{"path"},
	)
)

func main() {

	endpoint := flag.String("endpoint", "localhost:8800", "The address to listen on for HTTP requests <IP>:<PORT>")
	flag.Parse()

	log.Println("msg", "Build context", "build_context", version.BuildContext())

	fmt.Println("endpoint:", *endpoint)
	for _, s := range flag.Args() {
		fmt.Println("arg: ", s)
	}

	if len(flag.Args()) < 1 {
		log.Fatal("missing paths")
	}

	// Create a non-global registry.
	reg := prometheus.NewRegistry()
	reg.MustRegister(collectors.NewBuildInfoCollector(), fileCount)

	// Create new watcher.
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// get paths from arg line
	for _, v := range flag.Args() {
		err = watcher.Add(v)
		if err != nil {
			log.Fatal(err, " : ", v)
		}
	}

	// Start listening for events.
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

	// Expose the registered metrics via HTTP.
	http.Handle("/metrics", promhttp.HandlerFor(
		reg,
		promhttp.HandlerOpts{
			// Opt into OpenMetrics to support exemplars.
			EnableOpenMetrics: true,
			// Pass custom registry
			Registry: reg,
		},
	))
	log.Fatal(http.ListenAndServe(*endpoint, nil))

	// Block main goroutine forever.
	//<-make(chan struct{})
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
