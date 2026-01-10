import os
import sqlite3
import sys

domain = input("Specify the domain to filter (e.g. example.com): ")
kw = input("Specify your search string: ")

if getattr(sys, "frozen", False):
    localpath = os.path.dirname(sys.executable)
else:
    localpath = os.path.dirname(os.path.abspath(__file__))

db_path = os.path.join(localpath, "archive.db")

if not os.path.exists(db_path):
    print("archive.db not found in current directory.")
    sys.exit(1)

conn = sqlite3.connect(db_path)
cur = conn.cursor()

print("\nMatches:\n")

try:
    # Case-insensitive search using LIKE; filter by domain occurrence in URL
    cur.execute(
        """
        SELECT url
        FROM snapshots
        WHERE url LIKE ? AND content LIKE ?
        """,
        (f"%{domain}%", f"%{kw}%"),
    )
    rows = cur.fetchall()
    for (url,) in rows:
        print(url)
except sqlite3.Error as e:
    print("Search error:", e)

conn.close()

input("Search Finished. Press Enter to exit...")
