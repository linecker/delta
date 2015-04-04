# Delta

Highlight timestamp gaps.

```
$ cat example.log
12:00:00 A first line ...
12:00:01 ... immediately followed by a second one ...
12:00:02 ... and a third one ...
12:00:02 ... and another one.
12:01:00 Some time later.
```
This simplistic text file contains five lines of text, each accompanied by a timestamp (like it's the case with a typical log file). The first four lines have timestamps close to each other while the last line is a bit off. It takes quite a while for a human beeing to scan through those timestamps and spot this circumstance.
```
$ cat example.log | delta -d 10s
12:00:00 A first line ...
12:00:01 ... immediately followed by a second one ...
12:00:02 ... and a third one ...
12:00:02 ... and another one.
--------------------------------------------------------------------------------
12:01:00 Some time later.
```
When using delta, a seperation line (simple ASCII, nothing else) gets drawn between lines with timestamps that differ more then 10 seconds (in the given example).

## Usage
```
./delta -h
Usage: delta <[FILE] >[FILE]

	tail -f /var/log/messages | delta

delta - highlight timestamp gaps.

It reads from stdin, tries to find timestamps and calculates the timestamp
delta between subsequent lines. If this delta is larger then a certain limit,
an extra line of ASCII decoration that visually seperates those two lines is 
inserted.
	
Options:
  -c="": Use a custom timestamp format instead of the predefined ones. If used, 
  an example has to provided with the -e switch
  -d="100ms": Duration limit with unit suffix, e.g. 250ms, 1h45m. Valid time 
  units are ns, us, ms, s, m, h
  -e="": Example for the custom timestamp format
  -f="": Read from this file
  -p="-": Defines a custom seperator pattern
  -r=80: Defines how often the seperator pattern will be repeated

```
Examples:
```
delta -d 250ms -p "#" -r 80 -f /var/log/messages
tail -f /var/log/messages | delta
```

## Building
Assuming you have git and go installed, building delta is straightforward:
```
git clone http://github.com/linecker/delta
cd delta/
go build delta.go 
./delta -f sample.log 
```


