package main

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"
)

func checkError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	log.Println("Begin conversion...")
	startTime := time.Now()
	files, err := ioutil.ReadDir(".")
	checkError(err)
	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".json") {
			continue
		}
		processJSON(file.Name())
	}
	endTime := time.Now()
	log.Println("Conversion done!", endTime.Sub(startTime))
}

var escapeRegex = regexp.MustCompile(`[0-9](%)`)

func processJSON(jsonFile string) {
	log.Println("Processing", jsonFile)

	// Read json
	jsonBytes, err := ioutil.ReadFile(jsonFile)
	checkError(err)

	// Convert format
	var jsonLang map[string]string
	json.Unmarshal(jsonBytes, &jsonLang)
	var rows = []string{}
	for key, val := range jsonLang {
		v := val
		v = strings.ReplaceAll(val, "\n", "\\n")
		matches := escapeRegex.FindAllStringSubmatchIndex(v, -1)
		for i, match := range matches {
			v = v[:match[2]+i] + "%" + v[match[2]+i:]
		}
		row := key + "=" + v
		rows = append(rows, row)
	}
	sort.Strings(rows)

	// Write to file
	fileName := strings.Split(jsonFile, ".")
	file, err := os.OpenFile(fileName[0]+".lang", os.O_CREATE|os.O_WRONLY, 0644)
	checkError(err)
	defer file.Close()

	datawriter := bufio.NewWriter(file)
	for _, data := range rows {
		_, _ = datawriter.WriteString(data + "\n")
	}
	datawriter.Flush()
}
