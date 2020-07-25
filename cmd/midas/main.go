package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

func main() {
	apikey := "<APIKEY>"

	readstocks("../../stocklist.dat", apikey)
}

func readstocks(fn string, apikey string) {
	file, err := os.Open(fn)

	defer file.Close()

	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(file)

	var line string
	for {
		line, err = reader.ReadString('\n')
		lineslice := strings.Split(line, ",")
		symbolstr := lineslice[0]
		datestr := lineslice[1]

		trimmedsymbol := strings.ToUpper(strings.TrimSpace(symbolstr))
		trimmeddate := strings.TrimSpace(datestr)
		date, _ := time.Parse("2-1-2006", trimmeddate)
		testhandler(trimmedsymbol, date, apikey)

		if err != nil {
			break
		}
	}

	if err != io.EOF {
		fmt.Printf(" > Failed!: %v\n", err)
	}
}
