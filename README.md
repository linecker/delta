# Delta
Highlight lines with large timestamp delta.

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
Usage: delta [flags]
  -c="": Use a custom timestamp format instead of the predefined ones. The format has 
         to be specified as a regular expression, e.g. "([0-9]{2}:[0-9]{2}:[0-9]{2})". 
         If used, a layout example has to provided as well using the -e switch. The 
         layout example has to show how the reference time, defined to be 
         15:04:05.000000000 would be displayed
  -d="100ms": Duration limit with unit suffix, e.g. 250ms, 1h45m. Valid time units 
         are ns, us, ms, s, m, h
  -e="": Example for the custom timestamp format, see -c
  -f="": Provides a file to read from
  -p="~": Defines a custom seperator pattern
  -r=120: Defines how often the seperator pattern should be repeated to form the 
         seperator line
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


