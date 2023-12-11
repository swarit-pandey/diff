package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
)

type EditAction interface{}

type Keep struct {
	line string
}

type Insert struct {
	line string
}

type Remove struct {
	line string
}

type Frontier struct {
	x       int
	history []EditAction
}

// Myers diff algorithm
// Amazing read: blog.jcoglan.com/2017/02/12/the-myers-diff-algorithm-part-1/ and other parts
// Few more reads: [https://www.nathaniel.ai/myers-diff/] [https://epxx.co/artigos/diff_en.html]
// Complexity: O((N+M)D) where N and M are the lengths of the sequences and D is the number of edits
// Space: O(N+M)
// Reference implementation in Python: https://gist.github.com/adamnew123456/37923cf53f51d6b9af32a539cdfa7cc4
func myersDiff(aLines, bLines []string) []EditAction {
	frontier := make(map[int]Frontier)
	frontier[1] = Frontier{0, []EditAction{}}

	aMax := len(aLines)
	bMax := len(bLines)
	for d := 0; d <= aMax+bMax; d++ {
		for k := -d; k <= d; k += 2 {
			goDown := k == -d || (k != d && frontier[k-1].x < frontier[k+1].x)

			var oldX int
			var history []EditAction

			if goDown {
				oldX = frontier[k+1].x
				history = append([]EditAction{}, frontier[k+1].history...)
			} else {
				oldX = frontier[k-1].x + 1
				history = append([]EditAction{}, frontier[k-1].history...)
			}

			y := oldX - k

			if 1 <= y && y <= bMax && goDown {
				history = append(history, Insert{bLines[y-1]})
			} else if 1 <= oldX && oldX <= aMax {
				history = append(history, Remove{aLines[oldX-1]})
			}

			for oldX < aMax && y < bMax && aLines[oldX] == bLines[y] {
				history = append(history, Keep{aLines[oldX]})
				oldX++
				y++
			}

			if oldX >= aMax && y >= bMax {
				return history
			} else {
				frontier[k] = Frontier{oldX, history}
			}
		}
	}

	panic("No file found")
}

func printHorizontal(aLines, bLines []string, diff []EditAction) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Original", "Change", "Updated"})

	aIndex, bIndex := 0, 0
	for _, action := range diff {
		var row []string
		switch e := action.(type) {
		case Keep:
			row = []string{aLines[aIndex], " ", bLines[bIndex]}
			aIndex++
			bIndex++
		case Insert:
			row = []string{" ", "\033[32m+\033[0m", e.line}
			bIndex++
		case Remove:
			row = []string{e.line, "\033[31m-\033[0m", " "}
			aIndex++
		}

		table.Append(row)
	}

	table.Render()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage:", os.Args[0], "getdiff <string1> <string2> or <file1>.txt <file2>.txt")
		os.Exit(1)
	}

	mode := os.Args[1]

	switch mode {
	case "getdiff": // CL args
		if len(os.Args) != 4 {
			fmt.Println("Usage: getdiff <string1> <string2>")
			os.Exit(1)
		}
		aLines := strings.Split(os.Args[2], "")
		bLines := strings.Split(os.Args[3], "")
		diff := myersDiff(aLines, bLines)
		// printDiff(aLines, bLines, diff)
		printHorizontal(aLines, bLines, diff)
	default: // File args
		if len(os.Args) != 3 {
			fmt.Println("Usage:", os.Args[0], "<file1> <file2>")
			os.Exit(1)
		}
		aLines := readLines(os.Args[1])
		bLines := readLines(os.Args[2])
		diff := myersDiff(aLines, bLines)
		// printDiff(aLines, bLines, diff)
		printHorizontal(aLines, bLines, diff)
	}
}

func readLines(filePath string) []string {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, strings.TrimSpace(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return lines
}

// No table print (a bit confusing)
/* func printDiff(aLines, bLines []string, diff []EditAction) {
	aIndex, bIndex := 0, 0

	for _, action := range diff {
		var original, updated, change string

		switch e := action.(type) {
		case Keep:
			original = aLines[aIndex]
			updated = bLines[bIndex]
			change = " "
			aIndex++
			bIndex++
		case Insert:
			original = ""
			updated = e.line
			change = "\033[32m+\033[0m" // Green plus
			bIndex++
		case Remove:
			original = e.line
			updated = ""
			change = "\033[31m-\033[0m" // Red minus
			aIndex++
		}

		fmt.Printf("%-15s %-15s %-3s\n", original, updated, change)
	}
} */
