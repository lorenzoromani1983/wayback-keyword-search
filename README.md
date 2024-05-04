# wayback-keyword-search 

IMPORTANT NOTICE: WAYBACK IS NOW RATE LIMITING. TOOL IS STILL RUNNING, BUT SLOWLY (IT CAN'T LEVERAGE PARALLELISM ANYMORE). THE PYTHON VERSION IS UPDATED; THE GO VERSION IS NOT (YET).

This tools downloads each page from the Wayback Machine for a specific input domain and saves each page as a local .txt file, so that you can later search for keyword matches within the saved files.

Downloading is done with the "download" file; and searching with the "search" file.

You can download pages saved in specific years (i.e.: 2020), or years and months (i.e.: 202001), or years and months and days (i.e: 20200101), just specifying the date format in the prompt. If you want to download everything in the 2000's or 19**'s regardless the saved date, just type "2" (for the pages saved past 2000) or "1" (for the pages saved in the XXth century) in the prompt, and the tool will save each page matching that criteria. So, if you want to download a website that has been archived across 1999 and 2000, you will need to run the tool twice.

If you need to download big websites (thousands of saved pages), it may require quite a long time now. Still better than nothing. I advice using a VPN with auto switch, changing IP address every 30 minutes to avoid blocking.

Sometimes, when the Wayback API is down, you cannot fetch the entire list of URLs it has archived (this happens quite often based on recent experience); so be patient and retry.

There is a Python3 version and a Go version.

--------------------------

[*] Python usage:

> python3 download.py > specify your domain like: nytimes.com (no quotes!)

When the download is completed, a directory named as the domain will be saved in the local path.

So you can search for keyword matches within each file in the local dir using the "search.py" file:

> python3 search.py > specify your keyword (no quotes!).

--------------------------

[*] Go usage: [NOT WORKING NOW DUE TO ARCHIVE BLOCKING TOO MANY PARALLEL REQUESTS]

> go run download.go

and then:

> go run search.go

The best way to use the Go version is by running the compiled executables:

go build search.go

go build download.go

Notice that the Go version also features a download_channels.go version (thanks to Stephen Paulger for such improvement) which is a bit more efficient. Consider testing both!

