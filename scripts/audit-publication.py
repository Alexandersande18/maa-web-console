#!/usr/bin/env python3
"""Check publishable files and optionally every reachable Git blob. Never print secrets."""
import argparse
import ipaddress
from pathlib import Path
import re
import subprocess
import sys

ROOT = Path(__file__).resolve().parents[1]
PRIVATE_PATH = re.compile(
    r"(^|/)(?:\.env(?:\..*)?|certs|secrets|\.runtime|\.logs|debug|\.kilo)(/|$)"
    r"|(?:^|/)config(?:-[^/]*)?\.ya?ml$|^Golang/conf\.yaml$"
    r"|\.(?:pem|key|p12|pfx|crt|cer|db|sqlite3?|dump|bundle|exe|dll|so|dylib|log)$"
    r"|^check_(?:status|detailed)\.py$|^docs/roguelike-home-entry-.*\.md$",
    re.I,
)
PATTERNS = {
    "private-key": re.compile(r"-----BEGIN (?:RSA |EC |OPENSSH )?PRIVATE KEY-----"),
    "database-password": re.compile(r"mysql:[^:\s\"']+:([^@\s\"']+)@"),
    "url-credential": re.compile(r"https?://[^/\s:@]+:([^/\s@]+)@"),
    "local-user-directory": re.compile(r"[A-Za-z]:[/\\]+Users[/\\]+[^/\\\s\"'`]+", re.I),
    "embedded-jwt-secret": re.compile(r'JwtSecretKey\s*=\s*"([^"\n]+)"'),
}
PLACEHOLDERS = {"CHANGE_ME", "REDACTED", "<password>", "YOUR_PASSWORD"}


def git(*args):
    return subprocess.check_output(["git", *args], cwd=ROOT)


def check(path, blob, scope):
    findings = []
    if PRIVATE_PATH.search(path):
        findings.append((scope, path, 0, "private-file"))
    if b"\0" in blob:
        return findings
    text = blob.decode("utf-8", errors="replace")
    for rule, pattern in PATTERNS.items():
        for match in pattern.finditer(text):
            if match.lastindex and match.group(1) in PLACEHOLDERS:
                continue
            findings.append((scope, path, text.count("\n", 0, match.start()) + 1, rule))
    for match in re.finditer(r"\b(?:\d{1,3}\.){3}\d{1,3}\b", text):
        try:
            address = ipaddress.ip_address(match.group())
        except ValueError:
            continue
        if address.is_loopback or address.is_unspecified or any(
            address in ipaddress.ip_network(net)
            for net in ("192.0.2.0/24", "198.51.100.0/24", "203.0.113.0/24")
        ):
            continue
        findings.append((scope, path, text.count("\n", 0, match.start()) + 1, "deployment-ip"))
    return findings


def history_blobs():
    rows = git("rev-list", "--objects", "--all").decode("utf-8").splitlines()
    process = subprocess.Popen(
        ["git", "cat-file", "--batch"], cwd=ROOT,
        stdin=subprocess.PIPE, stdout=subprocess.PIPE,
    )
    try:
        for row in rows:
            oid, _, path = row.partition(" ")
            process.stdin.write((oid + "\n").encode())
            process.stdin.flush()
            header = process.stdout.readline().split()
            if len(header) != 3:
                raise RuntimeError("Could not inspect Git object")
            blob = process.stdout.read(int(header[2]))
            process.stdout.read(1)
            if header[1] == b"blob":
                yield path, blob, oid[:12]
    finally:
        process.stdin.close()
        process.wait()


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--history", action="store_true")
    args = parser.parse_args()
    findings = []
    count = 0
    names = git("ls-files", "-z", "--cached", "--others", "--exclude-standard")
    for name in sorted(set(names.decode("utf-8").split("\0")) - {""}):
        path = ROOT / name
        if path.is_file():
            count += 1
            findings.extend(check(name, path.read_bytes(), "worktree"))
    if args.history:
        for path, blob, scope in history_blobs():
            count += 1
            findings.extend(check(path, blob, scope))
    for scope, path, line, rule in sorted(set(findings)):
        print(f"{scope} {path}:{line}: {rule} (value hidden)")
    print(f"Checked {count} files/blobs; {len(set(findings))} findings.")
    print("This project-specific check complements a general secret scanner and manual review.")
    return 1 if findings else 0


if __name__ == "__main__":
    sys.exit(main())
