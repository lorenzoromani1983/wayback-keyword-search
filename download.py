import requests
import os
from multiprocessing import Pool
import time

domain = input("Type the target domain: ")

api_url = "https://web.archive.org/cdx/search/cdx?url=*."
base_url = "http://web.archive.org/web/"

localDir = path = os.path.dirname(os.path.abspath(__file__))
savePath = os.path.join(localDir, domain)
os.mkdir(savePath)

history = requests.get(api_url+domain).text.splitlines()

waybackurls = []

def getUrls(data):
    for record in data:
        if "text/html" in record:
            items = record.split(' ')
            savedpage = items[0].split(')')[1]
            url = domain+savedpage
            timestamp = items[1]
            wayback_url = base_url+timestamp+"/"+url
            waybackurls.append(wayback_url)

def download(url):
    noSlash = url.rstrip('/').replace('/', '£')
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
    
getUrls(history)

print("Downloading {} pages".format(str(len(set(waybackurls)))))
time.sleep(2)

p = Pool(10)
p.map(download, waybackurls)
p.terminate()
p.join()
