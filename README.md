# wayback-keyword-search

IMPORTANT NOTICE: WAYBACK IS NOW RATE LIMITING. TOOL IS STILL RUNNING, BUT SLOWLY (IT CAN'T LEVERAGE PARALLELISM ANYMORE).

This tools downloads each page from the Wayback Machine for a specific input domain and saves each page as a local .txt file, so that you can later search for keyword matches within the saved files.

Downloading is done with the "download" file; and searching with the "search" file.

You can download pages saved in specific years (i.e.: 2020), or years and months (i.e.: 202001), or years and months and days (i.e: 20200101), just specifying the date format in the prompt. If you want to download everything in the 2000's or 19**'s regardless the saved date, just type "2" (for the pages saved past 2000) or "1" (for the pages saved in the XXth century) in the prompt, and the tool will save each page matching that criteria. So, if you want to download a website that has been archived across 1999 and 2000, you will need to run the tool twice.

If you need to download big websites (thousands of saved pages), it may require quite a long time now. Still better than nothing. I advice using a VPN with auto switch, changing IP address every 30 minutes to avoid blocking.

Sometimes, when the Wayback API is down, you cannot fetch the entire list of URLs it has archived (this happens quite often based on recent experience); so be patient and retry.

There is a Python3 version and a Go version.

--------------------------

[*] Python usage:

```bash
python3 download.py > specify your domain like: nytimes.com (no quotes!)
```

When the download is completed, a directory named as the domain will be saved in the local path.

So you can search for keyword matches within each file in the local dir using the "search.py" file:

```bash
python3 search.py > specify your keyword (no quotes!).
```

--------------------------

[*] Go usage: [FOLLOWING PULL REQUEST THE CODE HAS BEEN REFACTORED - remember that setting too many workers may get you blocked.]

If you need to build the binary, use the following command:

```bash
make run-downloader
```

When using the downloader, you can specify the following arguments for running it:

```bash
downloader --domain=YOUR_SITE --timeStamp=2023 --workers=10
```

where the parameters are:

* domain - specify the target domain (only lowercase)
* timeStamp - specify timestamp in the format:'yyyymmdd' (also: 'yyyy' > download only a specific year; 'yyyymm' > year and month; '2' or '1' > everything for the years past 20** or 19**
* workers - specify the max workers (default=10)

And then, to build the search utilities, run the following command:

```bash
make run-search
```

The best way to use the Go version is by running the compiled executables:

```bash
make builds
```

Notice that the Go version also features a `./cmd/downloader/main.go` version (thanks to Stephen Paulger for such improvement) which is a bit more efficient. Consider testing both!

