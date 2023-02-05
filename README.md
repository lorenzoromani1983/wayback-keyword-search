# wayback-keyword-search
A quick and dirty but useful tool to download each text/html page from the wayback machine for a specific domain and search for keywords within the saved content

This tool downloads EACH "text/html" page for a specific domain from the Way Back Machine and stores the content of each retrieved page in the local directory as a .txt file.

> python3 download.py > specify your domain like: nytimes.com (no quotes!)

When the download is finished, a directory named as the domain will be saved in the local path.

So you can search for keyword matches within each file in the local dir using the "search.py" file:

> python3 search.py > specify your keyword (no quotes!).

BE CAREFUL: BIG DOMAINS MAY REQUIRE A LONG TIME TO DOWNLOAD! 

The tool is available in both Python3 and Go versions. Go is currently linux-only version. Will add Windows version soon.

IF YOU CHOOSE TO RUN THE GO VERSION, YOU BEST COMPILE IT (go build download.go; go build search.go). 
OTHERWISE REMEMBER TO RUN IT IN THE TERMINAL FROM THE SCRIPT'S DIRECTORY.

FEB. 2023: Added download file for Windows OS. Compile with Go build path_to_file and then use the search_windows.py file in order to get matches. No Golang search tool for the time being, just use the Python version.
