package fileMgr

import (
	"io"
	"log"
	"os"
)

func check(err error) {
	if err != nil {
		log.Fatal("Error: ", err)
	}
}

func ReadFile(fileName string) string {
	input, err := os.Open(fileName)
	check(err)
	defer input.Close()

	buf := make([]byte, 1024)
	fileContent := ""
	for {
		count, err := input.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
		}
		fileContent = fileContent + string(buf[:count])
	}
	return clearCarriage(fileContent)
}

func clearCarriage(s string) string {
	result := ""
	for _, rn := range s {
		if rn != '\r' {
			result += string(rn)
		}
	}
	return result
}

func WriteFile(fileName, s string) {
	err := os.WriteFile(fileName, []byte(s), 0o644)
	check(err)
}
