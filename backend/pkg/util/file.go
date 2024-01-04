package util

import (
	"bufio"
	"io"
	"os"
	"strings"
)

// ReadAllLines 读取文件中的所有行
func ReadAllLines(infile string) []string {
	lines := make([]string, 0)
	f, err := os.Open(infile)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadString('\n')
		lines = append(lines, strings.TrimSpace(line))
		if err != nil {
			if err == io.EOF {
				break
			} else {
				panic(err)
			}
		}
	}
	return lines
}
