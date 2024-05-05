
This is Prometheus exporter that shows the current number of files per directory, written in Go.


```bash
$ ./file_count_exporter -h

Usage:
  ./file_count_exporter [options] [path1 path2 ...]

Options:

list of paths
  -endpoint string
    	The address to listen on for HTTP requests <IP>:<PORT> (default "localhost:8800")
```

Exporter takes a list of paths which should be be monitored.

File are counted using os.ReadDir().



