package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	samplingProbability float64
	reservoirSize       uint
	splitStr            string
	outputFilename      string
)

func main() {
	flag.Float64Var(&samplingProbability, "p", 0.0, "Sample lines at probability p in [0, 1].")
	flag.UintVar(&reservoirSize, "r", 0, "Sample k lines.")
	flag.StringVar(&splitStr, "s", "", "Split ")
	flag.StringVar(&outputFilename, "o", "", "Output file(s)")
	flag.Parse()

	rand.Seed(time.Now().UnixNano())

	if reservoirSize > 0 && samplingProbability > 0.0 {
		fmt.Println("You can specify either -p or -r, but not both.")
		flag.PrintDefaults()
		return
	}

	input := inputReader(flag.Args())
	lines := make(chan string, 100)

	go func() {
		readLines(input, lines)
	}()

	if splitStr != "" {
		Split(input)
	} else if samplingProbability > 0.0 {
		sampleAtRate(input, samplingProbability, os.Stdout)
	} else if reservoirSize > 0 {
		lines := reservoirSample(input, reservoirSize)

		for _, line := range lines {
			fmt.Print(line)
		}
	}
}

func readLines(input bufio.Reader, lines chan<- string) err {
	for {
		line, err := input.ReadString('\n')
		if err != nil {
			return
		}

		lines <- line
	}
}

func Split(input io.Reader) {
	weights := parseSplitWeights()
	fmt.Println("Weights:", weights)
	bufReader := bufio.NewReader(input)
	writers := make(map[int]io.Writer)
	for i, _ := range weights {
		outName := fmt.Sprintf("%s_%d", outputFilename, i)
		file, err := os.Create(outName)
		if err != nil {
			panic(fmt.Sprintf("Couldn't open %s", outName))
		}
		defer file.Close()

		writers[i] = bufio.NewWriter(file)
	}

	for {
		line, err := bufReader.ReadBytes('\n')
		if err != nil {
			break
		}
		r := rand.Float64()
		for i, weight := range weights {
			if r < weight {
				_, err := writers[i].Write(line)
				if err != nil {
					panic("writing err")
				}
			}
		}
	}
	return
}

func parseSplitWeights() (splitWeights []float64) {
	splitWeights = make([]float64, 0)
	weightSum := 0.0
	fmt.Println("splitStr:", splitStr, ".")
	for _, weightStr := range strings.Split(splitStr, ",") {
		weight, err := strconv.ParseFloat(weightStr, 64)
		if err != nil {
			panic(fmt.Sprintf("Bad weight string: %s", weightStr))
		}

		splitWeights = append(splitWeights, weight)
		weightSum += weight
	}
	splitWeights = append(splitWeights, 1.-weightSum)
	return
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

func lineCount(reader io.Reader) (count uint) {
	count = 0

	bufReader := bufio.NewReader(reader)

	for {
		c, err := bufReader.ReadByte()
		if err != nil {
			return
		}
		if c == '\n' {
			count++
		}
	}
	return
}
