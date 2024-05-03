import requests
import os
import sys
import time
from contextlib import closing
from fake_useragent import UserAgent

API_URL = "https://web.archive.org/cdx/search/cdx?url=*."
BASE_URL = "http://web.archive.org/web/"
useragent = UserAgent()

def getUrls(data, domain, timeframe):
    """
    Return a set of wayback URLs
    """
    wayback_urls = set()
    for record in data:
        items = record.split(' ')
        try:
            savedpage = items[2]
        except Exception:
            continue
        url = savedpage
        timestamp = items[1]
        if str(timestamp.strip()).startswith(timeframe):
            wayback_url = BASE_URL + timestamp + "/" + url
            wayback_urls.add(wayback_url)
    return wayback_urls

def download(session, savePath, url):
    """
    Download the webpage using requests with a single session and save it in the specified path.
    """
    noSlash = url.rstrip('/').replace('/', '£').replace(":", "!!!").replace("?", "§§")
    filename = noSlash + ".txt" if not noSlash.endswith('.txt') else noSlash
    output = os.path.join(savePath, filename)
     
    if len(filename) <= 255 and not os.path.exists(os.path.join(savePath, filename)):
        with closing(session.get(url, stream=True)) as response:
            if response.status_code == 200:
                with open(output, 'wb') as f:
                    for chunk in response.iter_content(chunk_size=1024):
                        f.write(chunk)
                    print("Downloaded:",url)
            else:
                print(f"Failed to download URL: {url} with status code: {response.status_code}")
    else:
        if len(filename) > 255:
            print("Skipped - filename too long", url)
        if os.path.exists(os.path.join(savePath, filename)):
            print("Skipped (file already existing):",url)

def main():
    domain = input("Type the target domain: ")
    timeStamp = str(input("Specify a timestamp, ex: yyyymmdd, but also: yyyymm. Type '2' or '1' > to download everything for the years past 20** or 19**): "))

    localDir = path = os.path.dirname(os.path.abspath(__file__))
    savePath = os.path.join(localDir, domain)

    try:
        os.mkdir(savePath)
    except FileExistsError:
        print("Resuming download")
        pass

    history = requests.get(API_URL + domain).text.splitlines()
    waybackurls = getUrls(history, domain, timeStamp)

    print("Preparing to download {} pages".format(str(len(waybackurls))))
    time.sleep(2) #try to reduce (or eliminate) this sleep time, but at your own risk of being blocked.

    with requests.Session() as session:
        session.headers = {'user-agent':useragent.chrome}
        for url in waybackurls:
            try:
                download(session, savePath, url)
            except Exception as e:
                print("Error downloading:",url)
                continue

if __name__ == "__main__":
    main()
