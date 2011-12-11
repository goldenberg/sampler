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

var samplingProbability float64 
var reservoirSize uint

func main() {
	flag.Float64Var(&samplingProbability, "p", 0.5, "Sample lines at probability p in [0, 1].")
	flag.UintVar(&reservoirSize, "r", 10, "Sample k lines.")
	flag.Parse()

	rand.Seed(time.Nanoseconds())

	if flag.NArg() > 0 {
		file, _ := os.Open(flag.Arg(0))
		defer file.Close()

		// sampleAtRate(file, samplingProbability, os.Stdout)
		lines := reservoirSample(file, reservoirSize)
		for _, line := range lines {
			fmt.Print(line)
		}

	}
}

func sampleAtRate(reader io.Reader, rate float64, writer io.Writer) {
	bufReader := bufio.NewReader(reader)

	for  {
		line, err := bufReader.ReadBytes('\n');
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
		} else if rand.Float64() < float64(numLines) / float64(seen) {
			reservoir[rand.Intn(int(numLines))] = line
		}

		seen++
	}

	return


}