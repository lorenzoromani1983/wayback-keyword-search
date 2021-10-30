import requests
import os
import sys
from multiprocessing import Pool
from functools import partial
import time
from dataclasses import dataclass
from datetime import datetime


API_URL = "https://web.archive.org/cdx/search/cdx?url=*."
BASE_URL = "http://web.archive.org/web/"


@dataclass
class WaybackRecord:
    urlkey: str
    timestamp: str
    original: str
    mimetype: str
    statuscode: str
    digest: str
    length: int

    @property
    def path(self):
        return self.urlkey.split(")")[1]


def parse_wayback_record(record_str):
    (
        urlkey,
        timestamp,
        original,
        mimetype,
        statuscode,
        digest,
        length,
    ) = record_str.split()
    return WaybackRecord(
        urlkey=urlkey,
        timestamp=timestamp,
        original=original,
        mimetype=mimetype,
        statuscode=statuscode,
        digest=digest,
        length=int(length)
    )


def getUrls(domain):
    """
    Return a set of wayback URLs
    """
    wayback_urls = set()
    history = requests.get(API_URL + domain).text.splitlines()
    for line in history:
        record = parse_wayback_record(line)
        if record.mimetype == "text/html":
            url = domain + record.path
            wayback_url = BASE_URL + record.timestamp + "/" + url
            wayback_urls.add(wayback_url)
    return wayback_urls


def download(savePath, url):
    noSlash = url.rstrip("/").replace("/", "£")
    if not noSlash.endswith(".txt"):
        output = os.path.join(savePath, noSlash) + ".txt"
    else:
        output = os.path.join(savePath, noSlash)
    while True:
        if len(output) < 255:
            try:
                response = requests.get(url)
            except Exception:
                print("CANNOT RETRIEVE URL: ", url)
                break
            if response.status_code == 200:
                print("Writing to file:", url)
                with open(output, "w+") as outfile:
                    data = response.text
                    outfile.write(data)
                break
        if len(output) > 255:
            print("Skipping url: ", url)
            break
        break


def main():
    domain = input("Type the target domain: ")

    localDir = os.path.dirname(os.path.abspath(__file__))
    savePath = os.path.join(localDir, domain)

    try:
        os.mkdir(savePath)
    except FileExistsError:
        print("An output directory for the given domain already exists.")
        print("Quitting to avoid over-writing your data")
        sys.exit(1)

    waybackurls = getUrls(domain)

    print("Downloading {} pages".format(str(len(waybackurls))))

    with Pool(10) as p:
        p.map(partial(download, savePath), waybackurls)


if __name__ == "__main__":
    main()
