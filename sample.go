package main

import (
	"rand"
	"bufio"
	"io"
	"flag"
	"fmt"
	"os"
	"time"
)

var (
	samplingProbability float64
	reservoirSize       uint
)

func main() {
	flag.Float64Var(&samplingProbability, "p", 0.0, "Sample lines at probability p in [0, 1].")
	flag.UintVar(&reservoirSize, "r", 0, "Sample k lines.")
	flag.Parse()

	rand.Seed(time.Nanoseconds())

	if reservoirSize > 0 && samplingProbability > 0.0 {
		fmt.Println("You can specify either -p or -r, but not both.")
		flag.PrintDefaults()
		return
	}

	input := inputReader(flag.Args())

	if samplingProbability > 0.0 {
		sampleAtRate(input, samplingProbability, os.Stdout)
	} else if reservoirSize > 0 {
		lines := reservoirSample(input, reservoirSize)

		for _, line := range lines {
			fmt.Print(line)
		}
	}
}

func inputReader(args []string) (reader io.Reader) {
	readers := make([]io.Reader, 0)

	for _, arg := range args {
		file, err := os.Open(arg)
		//defer file.Close()

		if arg == "-" {
			readers = append(readers, os.Stdin)
		}

		if file != nil && err == nil {
			readers = append(readers, file)
		}
	}

	if len(readers) == 0 {
		return os.Stdin
	}

	return io.MultiReader(readers...)
}

func sampleAtRate(reader io.Reader, rate float64, writer io.Writer) {
	bufReader := bufio.NewReader(reader)

	for {
		line, err := bufReader.ReadBytes('\n')
		if err != nil {
			break
		}

		if rand.Float64() < rate {
			writer.Write(line)
			writer.Write([]uint8{'\n'})
		}
	}
}

func reservoirSample(reader io.Reader, numLines uint) (reservoir []string) {
	reservoir = make([]string, numLines)
	bufReader := bufio.NewReader(reader)
	var seen uint = 0

	for {
		line, err := bufReader.ReadString('\n')
		if err != nil {
			break
		}
		if seen < numLines {
			reservoir[seen] = line
		} else if rand.Float64() < float64(numLines)/float64(seen) {
			reservoir[rand.Intn(int(numLines))] = line
		}

		seen++
	}

	return

}
