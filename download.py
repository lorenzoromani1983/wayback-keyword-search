import requests
import os
import sys
from multiprocessing import Pool
from functools import partial
import time


API_URL = "https://web.archive.org/cdx/search/cdx?url=*."
BASE_URL = "http://web.archive.org/web/"


def getUrls(data, domain, timeframe):
    """
    Return a set of wayback URLs
    """
    wayback_urls = set()
    for record in data:
        items = record.split(' ')
        savedpage = items[0].split(')/')[1]
        url = domain + "/" + savedpage
        timestamp = items[1]
        if str(timestamp.strip()).startswith(timeframe):
            wayback_url = BASE_URL + timestamp + "/" + url
            wayback_urls.add(wayback_url)
    return wayback_urls


def download(savePath, url):
    noSlash = url.rstrip('/').replace('/', '£').replace(":", "!!!").replace("?", "§§")
    if not noSlash.endswith('.txt'):
        output = os.path.join(savePath, noSlash)+".txt"
    else:
        output = os.path.join(savePath, noSlash)
    while True:
        if len(output) < 255:
            try:
                response = requests.get(url)
            except Exception:
                print("CANNOT RETRIEVE URL: ",url)
                break
            if response.status_code == 200:
                print("Writing to file:",url)
                file = open(output, "w+")
                data = response.text
                file.write(data)
                file.close()
                break
        if len(output) > 255:
            print("Skipping url: ",url)
            break
        break


def main():
    domain = input("Type the target domain: ")
    timeStamp = str(input("Specify a timestamp, ex: yyyymmdd, but also: yyyymm. To download everything just type '2'): "))

    localDir = path = os.path.dirname(os.path.abspath(__file__))
    savePath = os.path.join(localDir, domain)

    try:
        os.mkdir(savePath)
    except FileExistsError:
        print("An output directory for the given domain already exists.")
        print("Quitting to avoid over-writing your data")
        sys.exit(1)

    history = requests.get(API_URL + domain).text.splitlines()

    waybackurls = getUrls(history, domain, timeStamp)

    print("Downloading {} pages".format(str(len(waybackurls))))
    time.sleep(2)

    p = Pool(10)
    p.map(partial(download, savePath), waybackurls)
    p.terminate()
    p.join()


if __name__ == "__main__":
    main()
