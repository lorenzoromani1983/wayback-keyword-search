import os
import re

folder = input("Specify the domain: ")

kw = input("Specify your search string: ")

localpath = path = os.path.dirname(os.path.abspath(__file__))

files = os.listdir(localpath+"/"+folder)

print("\nMatches: \n")

for file in files:
     try:
          data = open(localpath+"/"+folder+"/"+file).read()
     except Exception as e:
          print(e)
          continue
     if re.search(kw, data, flags= re.IGNORECASE):
          if not file.endswith('robots.txt'):
               print(file.replace('£','/').replace('.txt',''))
          else:
               print(file.replace('£','/'))
 
