// delta.go - Highlight timestamp gaps.
//
// delta reads lines of text from log files or stdin, tries to find timestamps 
// in those lines of text and calculates the difference of the timestamps
// between subsecquent lines. If this delta is larger then a certain limit, an
// extra line of text that visually seperates those two lines is inserted.
//
// TODO: automatic duration limit

package main

import "bufio"
import "flag"
import "fmt"
import "io"
import "os"
import "regexp"
import "time"

// Two subsequent lines with timestamp differences larger then
// timestampDifferenceLimit will get seperated.
var timestampDifferenceLimit time.Duration

// TODO doc
var comparisionFunction func(t time.Time) bool = largeTimestampDifference

// Name of the input file ("" if stdin is used).
var inputFileName string

// Holds the timestamp of the previous line.
var previousTimestamp time.Time

// Type for timestamp formats. The fields definition and example are used to
// specify a timestamp. The field compiled holds the timestamp as compiled
// regular expression.
type TimestampFormat struct {
	definition string
	example    string
	compiled   regexp.Regexp
}

// Optional custom format from the command line.
var customFormat TimestampFormat

// All specified timestamp formats.
var timestampFormats []TimestampFormat

// Prepare predefined or custom timestamp formats.
func prepareTimestampFormats() {
	if customFormat.definition == "" {
		// hh:mm:ss.mmmuuu (glog)
		timestampFormats = append(timestampFormats, TimestampFormat{
			definition: "(?P<time>[0-9]{2}:[0-9]{2}:[0-9]{2}.[0-9]{6})",
			example:    "15:04:05.000000"})
		// hh:mm:ss.mmm
		timestampFormats = append(timestampFormats, TimestampFormat{
			definition: "(?P<time>[0-9]{2}:[0-9]{2}:[0-9]{2}.[0-9]{3})",
			example:    "15:04:05.000"})
		// hh:mm:ss
		timestampFormats = append(timestampFormats, TimestampFormat{
			definition: "(?P<time>[0-9]{2}:[0-9]{2}:[0-9]{2})",
			example:    "15:04:05"})
		// hh:mm
		timestampFormats = append(timestampFormats, TimestampFormat{
			definition: "(?P<time>[0-9]{2}:[0-9]{2})",
			example:    "15:04"})
	} else {
		timestampFormats = append(timestampFormats, TimestampFormat{
			definition: customFormat.definition,
			example:    customFormat.example})
	}
	// Compile regular expressions.
	for i := 0; i < len(timestampFormats); i++ {
		compiled, err := regexp.Compile(timestampFormats[i].definition)
		if err != nil {
			panic(err)
		}
		timestampFormats[i].compiled = *compiled
	}
}

// Holds the seperator configuration.
var seperator struct {
	pattern string
	reps    int
	line    string
}

// Prepare the seperator line. We only want to do this once.
func prepareSeperator() {
	for i := 0; i < seperator.reps; i++ {
		seperator.line += seperator.pattern
	}
}

// Check if we have a large timestamp difference.
func largeTimestampDifference(t time.Time) bool {
	diff := -previousTimestamp.Sub(t)
	previousTimestamp = t
	if diff > timestampDifferenceLimit {
		return true
	}
	return false
}

var first int64 = 0
var mean int64

// todo: doesn't work very good yet!
func dynamicDifference(t time.Time) bool {
	if first == 0 {
		previousTimestamp = t
		first = 1
		return false
	} else if first == 1 {
		diff := -previousTimestamp.Sub(t)
		mean = int64(diff)
		first = 2
		return false
	} else if first == 2 {
		diff := -previousTimestamp.Sub(t)
		previousTimestamp = t
		eta := 0.1
		mean += int64(float64((int64(diff) - mean)) * eta)
		fmt.Printf("diff = %12d, mean = %12d |", diff/1000000, mean/1000000)
		if int64(diff) > mean {
			return false // true
		}
		return false
	}
	return false
}

// Analyze a single line.
func analyzeLine(line []byte) {
	// Check if any of the known timestamp formats fits.
	for i := 0; i < len(timestampFormats); i++ {
		regexp := timestampFormats[i].compiled
		tuple := regexp.FindSubmatchIndex(line)
		if tuple != nil {
			start := tuple[0]
			end := tuple[1]
			raw := line[start:end]
			parsed, err := time.Parse(timestampFormats[i].example, string(raw))
			if err != nil {
				continue
			}
			if comparisionFunction(parsed) {
				fmt.Println(seperator.line)
			}
			break
		}
	}
	fmt.Println(string(line))
}

// Analyze proper.
func analyze(reader *bufio.Reader) {
	// Main loop.
	for {
		line, err := reader.ReadString('\n')
		if err == nil {
			analyzeLine([]byte(line[:len(line)-1]))
		} else if err == io.EOF {
			return
		} else {
			panic(err)
		}
	}
}

// Read from file.
func analyzeFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	analyze(reader)
}

// Read from stdin.
func analyzeStdin() {
	reader := bufio.NewReader(os.Stdin)
	analyze(reader)
}

// Print usage message.
func usage() {
	fmt.Println(`Usage: delta <[FILE] >[FILE]

	tail -f /var/log/messages | delta

delta - highlight timestamp gaps.

It reads from stdin, tries to find timestamps and calculates the timestamp
delta between subsequent lines. If this delta is larger then a certain limit,
an extra line of ASCII decoration that visually seperates those two lines is 
inserted.
	
Options:`)
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	// Care about command line flags.
	dynamic := false
	duration := "100ms"
	flag.Usage = usage
	flag.StringVar(&inputFileName, "f", "", "Read from this file")
	flag.StringVar(&customFormat.definition, "c", "",
		"Use a custom timestamp format instead of the predefined ones. "+
			"If used, an example has to provided with the -e switch")
	flag.StringVar(&customFormat.example, "e", "",
		"Example for the custom timestamp format")
	flag.StringVar(&duration, "d", duration,
		"Duration limit with unit suffix, e.g. 250ms, 1h45m. Valid time "+
			"units are ns, us, ms, s, m, h")
	flag.StringVar(&seperator.pattern, "p", "-",
		"Defines a custom seperator pattern")
	flag.IntVar(&seperator.reps, "r", 80,
		"Defines how often the seperator pattern will be repeated")
	flag.BoolVar(&dynamic, "y", dynamic, "Use dynamic duration limit")
	flag.Parse()
	if dynamic {
		comparisionFunction = dynamicDifference
	}
	d, err := time.ParseDuration(duration)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		timestampDifferenceLimit = d
	}

	// Compile regular expressions and prepare seperator.
	prepareTimestampFormats()
	prepareSeperator()

	// Do work.
	if inputFileName == "" {
		analyzeStdin()
	} else {
		analyzeFile(inputFileName)
	}
}

