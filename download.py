import argparse
import os
import sqlite3
import time
import requests

API_URL = "https://web.archive.org/cdx/search/cdx?url=*."
BASE_URL = "http://web.archive.org/web/"

# Timing controls
batch_timer = 5
divisore = 20

# Simple, fixed user agent to avoid external deps
DEFAULT_UA = (
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) "
    "AppleWebKit/537.36 (KHTML, like Gecko) "
    "Chrome/120.0.0.0 Safari/537.36"
)


def ensure_db(conn):
    cur = conn.cursor()
    cur.execute(
        """
        CREATE TABLE IF NOT EXISTS snapshots (
            url TEXT PRIMARY KEY,
            content TEXT
        )
        """
    )
    conn.commit()


def getUrls(data, timeframe):
    """
    Return a set of Wayback snapshot URLs filtered by timeframe prefix.
    """
    wayback_urls = set()
    for record in data:
        items = record.split(" ")
        if len(items) < 3:
            continue
        timestamp = items[1].strip()
        savedpage = items[2]
        if str(timestamp).startswith(timeframe):
            wayback_url = f"{BASE_URL}{timestamp}/{savedpage}"
            wayback_urls.add(wayback_url)
    return wayback_urls


def batch_existing(conn, urls, batch_size=500):
    """Return a set of URLs that already exist in DB, checked in batches."""
    existing = set()
    cur = conn.cursor()
    urls_list = list(urls)
    for i in range(0, len(urls_list), batch_size):
        chunk = urls_list[i : i + batch_size]
        placeholders = ",".join(["?"] * len(chunk))
        cur.execute(f"SELECT url FROM snapshots WHERE url IN ({placeholders})", chunk)
        existing.update(row[0] for row in cur.fetchall())
    return existing


def fetch_and_store(session, conn, url, intercall_timer=0):
    time.sleep(intercall_timer)
    resp = session.get(url, timeout=30)

    status_code = resp.status_code

    if status_code != 200:
        print(f"Failed to download URL: {url} with status code: {resp.status_code}")
        return False

    # Accumulate body as bytes then decode as latin-1 to preserve bytes mapping
    resp.encoding = resp.encoding or resp.apparent_encoding or "utf-8"
    content_text = resp.text

    try:
        conn.execute(
            "INSERT OR IGNORE INTO snapshots(url, content) VALUES(?, ?)",
            (url, content_text),
        )
        conn.commit()
        print("Downloaded:", url)
        return True
    except sqlite3.Error as e:
        print("DB error while saving", url, e)
        return False

def main():
    parser = argparse.ArgumentParser(
        description="Download archived pages from Wayback Machine for a given domain."
    )
    parser.add_argument(
        "-d", "--domain",
        required=True,
        help="Target domain (e.g. example.com)"
    )
    parser.add_argument(
        "-s", "--sleep",
        type=float,
        default=1,
        help="Sleep time in seconds between each request (default: 1)"
    )
    parser.add_argument(
        "-t", "--timestamp",
        default="",
        help="Timestamp filter (e.g. yyyymmdd, yyyymm, yyyy). "
             "Use '2' for all 2*** years, '1' for all 19** years, or leave empty for all."
    )
    args = parser.parse_args()

    domain = args.domain
    timeStamp = args.timestamp
    intercall_timer = args.sleep

    localDir = os.path.dirname(os.path.abspath(__file__))
    db_path = os.path.join(localDir, "archive.db")

    conn = sqlite3.connect(db_path)
    ensure_db(conn)

    # Build candidate URL list from Wayback CDX
    history = requests.get(API_URL + domain, timeout=30).text.splitlines()
    waybackurls = getUrls(history, timeStamp)

    # Check resume state via DB
    existing = batch_existing(conn, waybackurls)
    to_fetch = [u for u in waybackurls if u not in existing]
    print(f"Preparing to download {len(to_fetch)} pages (skipping {len(existing)} already saved)")

    with requests.Session() as session:
        session.headers = {"user-agent": DEFAULT_UA}

        downloaded = 0
        for _, url in enumerate(to_fetch, start=1):
            if downloaded % divisore == 0 and downloaded > 0:
                print(f"Sleeping {batch_timer} seconds...ZzZz")
                time.sleep(batch_timer)
            try:
                if fetch_and_store(session, conn, url, intercall_timer):
                    downloaded += 1
            except requests.exceptions.ConnectionError as e:
                cause = e.__cause__ or e.__context__
                if isinstance(cause, OSError) and getattr(cause, "errno", None) in {61, 111}:
                    print(f"Connection refused when reaching {url}")
                else:
                    print("Blocked. Sleeping for 5 minutes.")
                    time.sleep(300)
                continue
            except Exception as e:
                print("Error downloading:", url, e)
                continue

    conn.close()

if __name__ == "__main__":
    main()
