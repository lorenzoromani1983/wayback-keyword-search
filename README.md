# wayback-keyword-search
A quick and dirty but useful tool to download each text/html page from the wayback machine for a specific domain and search for keywords within the saved content

This tool downloads EACH "text/html" page for a specific domain from the Way Back Machine and stores the content of each retrieved page in the local directory as a .txt file.

> python3 download.py > specify your domain like: nytimes.com (no quotes!)

When the download is finished, a directory named as the domain will be saved in the local path.

So you can search for keyword matches within each file in the local dir using the "search.py" file:

> python3 search.py > specify your keyword (no quotes!).

BE CAREFUL: BIG DOMAINS MAY REQUIRE A LONG TIME TO DOWNLOAD! 
