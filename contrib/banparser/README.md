# banparser.py

That script takes buscarron's log output from stdin, parses it and prints hashes on bad actors and the amount of times they were banned.
That script intended to be used as semi-automatic "data provider" for `BUSCARRON_BAN_LIST` var.

Example:

```bash
$ journalctl -o cat -u buscarron | python3 ./banparser.py
111111111 (7 bans)
222222222 (6 bans)
333333333 (5 bans)
444444444 (4 bans)
555555555 (4 bans)
666666666 (3 bans)
777777777 (3 bans)
888888888 (3 bans)
```

You may adjust minimal amount of bans required in the script itself (`minBans` var).
Keep in mind that you will have a lot of entries of random scanners and bots, so setting `minBans` to lower values is bad idea.
