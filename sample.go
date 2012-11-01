package main

import (
	"math/rand"
	"bufio"
	"io"
	"flag"
	"fmt"
	"os"
	"time"
	"strconv"
	"strings"
)

var (
	samplingProbability float64
	reservoirSize       uint
	splitStr            string
	outputFilename      string
)

func main() {
	flag.Float64Var(&samplingProbability, "p", 0.0, "Sample lines at probability p in [0, 1].")
	flag.UintVar(&reservoirSize, "r", 0, "Sample exactly k lines.")
	flag.StringVar(&splitStr, "s", "", "Split the input into n random subsets. Delimit weights with commas, e.g. 8,1,1. Weights will be normalized to sum to 1.")
	flag.StringVar(&outputFilename, "o", "", "Output file(s). If not specified, write to stdout. Required and used as a basename for splitting. The basename will be appended with _1, _2, etc.")
	flag.Parse()

	rand.Seed(time.Now().UnixNano())

	if reservoirSize > 0 && samplingProbability > 0.0 {
		fmt.Println("You can specify either -p or -r, but not both.")
		flag.PrintDefaults()
		return
	}


	input := inputReader(flag.Args())

	if splitStr != "" {
		if (outputFilename == "") {
			fmt.Println("To split input, you must specify an output base filename with -o.")	
			flag.PrintDefaults()
			return
		}
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

func Split(input io.Reader) {
	weights := parseSplitWeights()
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

func parseSplitWeights() (weights []float64) {
	weights = make([]float64, 0)
	sum := 0.0
	for _, weightStr := range strings.Split(splitStr, ",") {
		w, err := strconv.ParseFloat(weightStr, 64)
		if err != nil {
			panic(fmt.Sprintf("Bad weight string: %s", weightStr))
		}

		weights = append(weights, w)
		sum += w
	}
	for i, w := range weights {
		weights[i] = w / sum
	}
	return weights
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
		}
	}
}

func reservoirSample(reader io.Reader, k uint) (reservoir []string) {
	reservoir = make([]string, k)
	bufReader := bufio.NewReader(reader)

	var n uint = 0

	for {
		n++
		line, err := bufReader.ReadString('\n')
		if err != nil {
			break
		}
		if n <= k {
			reservoir[n-1] = line

		// Replace a random element with probability k/n
		} else if rand.Float64() < float64(k)/float64(n) {
			reservoir[rand.Intn(int(k))] = line
		}
	}

	return reservoir
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
