
This is Prometheus exporter that shows the current number of files per directory, written in Go.


```bash
$ ./gowatch -h

Usage:
  ./gowatch [options] [path1 path2 ...]

Options:

list of paths
  -endpoint string
    	The address to listen on for HTTP requests <IP>:<PORT> (default "localhost:8800")
```

Exporter takes a list of path which should be be monitored.

Files are counted using __"fsnotify"__ Go package, so even large folders can be tracked easily.


Go package __"fsnotify"__ uses __"inotify"__ on Linux.
__"inotify"__ is a Linux kernel subsystem that provides a mechanism to efficently monitor filesystem events such as :
- creation, deletion, modification, accessing

of files or directories.



