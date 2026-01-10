import argparse
import csv
import os
import re
import sqlite3
import sys
from urllib.parse import urlparse


def list_domains(cur):
    """Extract and list unique domains from saved URLs."""
    cur.execute("SELECT url FROM snapshots")
    domains = set()
    for (url,) in cur.fetchall():
        # Wayback URL format: http://web.archive.org/web/TIMESTAMP/ORIGINAL_URL
        parts = url.split("/", 5)
        if len(parts) >= 6:
            original_url = parts[5]
            try:
                parsed = urlparse(original_url if "://" in original_url else f"http://{original_url}")
                if parsed.netloc:
                    domains.add(parsed.netloc)
            except Exception:
                continue
    return sorted(domains)


def export_csv(filename, keyword, urls):
    """Export results to CSV file."""
    with open(filename, "w", newline="", encoding="utf-8") as f:
        writer = csv.writer(f)
        writer.writerow(["keyword", "wayback_url"])
        for url in urls:
            writer.writerow([keyword, url])


def search_like(cur, domain, keyword):
    """Search using SQL LIKE (case-insensitive substring match)."""
    cur.execute(
        """
        SELECT url, content
        FROM snapshots
        WHERE url LIKE ? AND content LIKE ?
        """,
        (f"%{domain}%", f"%{keyword}%"),
    )
    return cur.fetchall()


def search_regex(cur, domain, pattern, ignore_case=False):
    """Search using Python regex on content."""
    cur.execute(
        """
        SELECT url, content
        FROM snapshots
        WHERE url LIKE ?
        """,
        (f"%{domain}%",),
    )
    flags = re.IGNORECASE if ignore_case else 0
    try:
        compiled = re.compile(pattern, flags)
    except re.error as e:
        print(f"Invalid regex pattern: {e}")
        sys.exit(1)

    results = []
    for url, content in cur.fetchall():
        if compiled.search(content):
            results.append((url, content))
    return results


def main():
    parser = argparse.ArgumentParser(
        description="Search archived pages in the local SQLite database."
    )
    parser.add_argument(
        "-d", "--domain",
        default="",
        help="Filter by domain in URL (e.g. example.com). Use 'all' to search across all saved domains."
    )
    parser.add_argument(
        "-k", "--keyword",
        help="Search string or regex pattern"
    )
    parser.add_argument(
        "-r", "--regex",
        action="store_true",
        help="Interpret keyword as a regular expression"
    )
    parser.add_argument(
        "-i", "--ignore-case",
        action="store_true",
        help="Case-insensitive regex search (only with -r)"
    )
    parser.add_argument(
        "-o", "--output",
        help="Export results to CSV file (specify filename)"
    )
    parser.add_argument(
        "-l", "--list-domains",
        action="store_true",
        help="List all domains saved in the database"
    )
    args = parser.parse_args()

    if not args.list_domains and not args.keyword:
        parser.error("-k/--keyword is required unless using -l/--list-domains")

    # Treat 'all' as no domain filter
    if args.domain and args.domain.lower() == "all":
        args.domain = ""

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

    try:
        if args.list_domains:
            domains = list_domains(cur)
            print(f"\nDomains saved ({len(domains)}):\n")
            for domain in domains:
                print(domain)
            conn.close()
            return

        if args.regex:
            rows = search_regex(cur, args.domain, args.keyword, args.ignore_case)
        else:
            rows = search_like(cur, args.domain, args.keyword)

        urls = [url for url, _ in rows]

        print(f"\nMatches ({len(urls)}):\n")
        for url in urls:
            print(url)

        if args.output and urls:
            output_path = os.path.join(localpath, args.output)
            export_csv(output_path, args.keyword, urls)
            print(f"\nResults exported to: {output_path}")

    except sqlite3.Error as e:
        print("Search error:", e)

    conn.close()


if __name__ == "__main__":
    main()
