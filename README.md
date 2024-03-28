
This is Prometheus exporter that shows the current number of files per directory.


```bash
$ ./gowatch -h

Usage:
  ./gowatch [options] [path1 path2 ...]

Options:

list of paths
  -endpoint string
    	The address to listen on for HTTP requests <IP>:<PORT> (default "localhost:8800")
```

Exporter takes a list of path to be monitored.

Files are counted using __"fsnotify"__ package, so even large folders can be tracked easily.
Go package __"fsnotify"__ uses __"inotify"__ API on Linux, which is basically relying on Linux kernel to send us notifications on file creation/deletion.


