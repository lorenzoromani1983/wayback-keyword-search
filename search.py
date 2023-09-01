import os
import re
import time
import sys

folder = input("Specify the domain: ")

kw = input("Specify your search string: ")

if getattr(sys, 'frozen', False):
    localpath = os.path.dirname(sys.executable)
else:
    localpath = os.path.dirname(os.path.abspath(__file__))

repo = os.listdir(localpath + "/" + folder)

print("\nMatches: \n")

for file in repo:
    try:
        data = open(localpath + "/" + folder + "/" + file, encoding = 'Latin1').read()
    except Exception as e:
        print("Error reading file, check manually:",file)
        continue
    if re.search(kw, data, flags=re.IGNORECASE):
        if not file.endswith("robots.txt"):
            print(file.replace("£", "/").replace(".txt", "").replace("!!!", ":").replace("§§", "?"))
        else:
            print(file.replace("£", "/").replace("!!!", ":").replace("§§", "?"))

input("Search Finished. Press Enter to exit...")
