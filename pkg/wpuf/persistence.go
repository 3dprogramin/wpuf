package wpuf

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// ReadPasswords from file (wordlist)
func ReadPasswords(path string) []string {
	inFile, err := os.Open(path)
	if err != nil {
		fmt.Println(err.Error() + `: ` + path)
		CheckError(err)
	}
	defer func(inFile *os.File) {
		err := inFile.Close()
		CheckError(err)
	}(inFile)

	var lines []string
	scanner := bufio.NewScanner(inFile)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "# ") {
			continue
		}
		if strings.HasPrefix(line, "##########																																																																																																		") {
			continue
		}
		lines = append(lines, line)
	}
	return lines
}
