package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

const version = "v0.3.7"

func main() {
	// ignore error
	os.Remove("tables.go")

	f, err := os.Create("tables.go")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	if _, err := f.WriteString("// Code generated by internal/main.go DO NOT EDIT\n\n"); err != nil {
		log.Fatal(err)
	}

	{
		resp, err := http.Get(fmt.Sprintf("https://raw.githubusercontent.com/golang/text/%s/LICENSE", version))
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		br := bufio.NewScanner(resp.Body)
		for br.Scan() {
			line := br.Text()
			if _, err := f.WriteString("// " + line + "\n"); err != nil {
				log.Fatal(err)
			}
		}
	}

	{
		resp, err := http.Get(fmt.Sprintf("https://raw.githubusercontent.com/golang/text/%s/encoding/japanese/tables.go", version))
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		br := bufio.NewScanner(resp.Body)
		for br.Scan() {
			line := br.Text()
			if strings.HasPrefix(line, "// Package") {
				continue
			}
			line = strings.ReplaceAll(line, "package japanese", "package sjisreplace")
			if _, err := f.WriteString(line + "\n"); err != nil {
				log.Fatal(err)
			}
		}
	}
}
