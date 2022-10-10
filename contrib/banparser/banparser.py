#!/usr/bin/env python

from sys import stdin
from re import search
import operator

minBans = 3

def summarize_bans(data:str)->None:
    lines = data.split("\n")
    matches = {}
    for line in lines:
        if "banned" in line:
            match = search(r"(\d+) has been banned$", line)
            if match:
                id = match.group(1)
                if id not in matches.keys():
                    matches[id] = 1
                else:
                    matches[id] += 1
    prettyprint(matches)


def prettyprint(data: dict) -> None:
    data = sorted(data.items(),key=operator.itemgetter(1),reverse=True)
    for item in data:
        if item[1] >= minBans:
            print("{} ({} bans)".format(item[0], item[1]))

if __name__ == "__main__":
    summarize_bans(stdin.read())
