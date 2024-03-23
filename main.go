package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/manifoldco/promptui"
)

const statusFileName = "status.toml"

var hiddenFilesIgnore bool = false

const tomlHeader = `
[NOT WORKING]

[COMPLETED]

[PROGRESS]

[PLAYING]

[NOT PLAYED]
`

func overwriteStatusAsk() bool {
	prompt := promptui.Select{
		Label: "Rewrite an existing file? [y/n]",
		Items: []string{"Yes", "No"},
	}
	_, result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}

	return result == "Yes"
}

func writeStatus(fileName string) {
	file, err := os.OpenFile(fileName, os.O_CREATE | os.O_WRONLY | os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	file.WriteString(tomlHeader)

	files, err :=  os.ReadDir(".")
    if err != nil {
        log.Fatal(err)
    }

	for _, fileInfo := range files {
		name := fileInfo.Name()

		if name == statusFileName { continue }
		if !hiddenFilesIgnore && name[0] == '.' { continue }

		fmt.Println("Writing file in status:", name)
		file.WriteString(name + "\n")
	}
}

func main() {
	hiddenFiles := flag.Bool("hidden", false, "Does not include hidden files in the status")
	flag.Parse()

	hiddenFilesIgnore = *hiddenFiles

	if _, err := os.Stat(statusFileName); !os.IsNotExist(err) {
		if overwriteStatusAsk() {
			writeStatus(statusFileName)
		}
	} else {
		writeStatus(statusFileName)
	}
}
