# wayback-keyword-search

This tools downloads each text/html page from the Wayback Machine for a specific input domain and saves each page as a local .txt file, so that you can later search for keyword matches within the saved files.

There is a Python3 version and a Go version.

[*] Python usage:

> python3 download.py > specify your domain like: nytimes.com (no quotes!)

When the download is completed, a directory named as the domain will be saved in the local path.

So you can search for keyword matches within each file in the local dir using the "search.py" file:

> python3 search.py > specify your keyword (no quotes!).

[*] Go usage:

> go run download.go

and then:

> go run search.go

The Go version has an additional feature compared to the Python version, as it resumes a previously aborted download.

The best way to use the Go version is by running the compiled executables:

go build search.go
go build download.go
