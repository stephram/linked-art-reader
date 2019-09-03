# linked-art-reader

**Note: Does not yet do anything useful**

Reads an activity-stream containing references to linked-art objects.

On macOS build with
```
./build.sh
```
which will create the `reader` and `importer` binaries.

Run with `-h` to see the various command line options.

```
$ ./reader -h
Usage of ./reader:
  -end int
    	stop at page (default -1)
  -host string
    	activity stream host
  -path string
    	path to the activity stream (default "activity-stream")
  -pretty
    	pretty print JSON output
  -prof
    	enable the pprof package. Listening on port 8080
  -scheme string
    	http(s) (default "http")
  -start int
    	start at page (default 1)
```

Run it with the following to see some output.

`./reader -host <host> -pretty`
