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
MAX_RETRIES = 3


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


def safe_filename(url):
    """
    Generate a filesystem safe filename from the URL.
    """
    output = url.rstrip("/").replace("/", "£")
    if not output.endswith(".txt"):
        output += ".txt"
    return output


def download(savePath, url):
    """
    Saves the given url into a provided directory with a
    safe filename
    """
    output = safe_filename(url)
    if len(output) >= 255:
        print("Skipping url: ", url)
        return

    for _ in range(MAX_RETRIES):
        try:
            response = requests.get(url, timeout=5)
            response.raise_for_status()
        except requests.exceptions.ReadTimeout:
            continue
        except requests.exceptions.HTTPError:
            continue
        else:
            break
    else:
        print(f"Failed to download {url} after {MAX_RETRIES} retries")
        return

    print("Writing to file:", url)
    with open(os.path.join(savePath, output), "w+") as outfile:
        data = response.text
        outfile.write(data)


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
