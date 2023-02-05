import os
import re
import time

folder = input("Specify the domain: ")

kw = input("Specify your search string: ")

localpath = os.path.dirname(os.path.abspath(__name__))

files = os.listdir(localpath + "/" + folder)

print("\nMatches: \n")

for file in files:
    try:
        data = open(localpath + "/" + folder + "/" + file).read()
    except Exception as e:
        continue
    if re.search(kw, data, flags=re.IGNORECASE):
        if not file.endswith("robots.txt"):
            print(file.replace("£", "/").replace(".txt", "").replace("!!!", ":").replace("§§", "?"))
        else:
            print(file.replace("£", "/").replace("!!!", ":").replace("§§", "?"))

time.sleep(10000)
