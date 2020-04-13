package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/mattn/go-shellwords"
)

func main() {
	var fieldsOpt string
	flag.StringVar(&fieldsOpt, "f", "", "output fields indexes that is 1-origin. Ex. -f 1,2 or -f 1-10")
	flag.Parse()

	if fieldsOpt == "" {
		log.Println("-f is required")
		os.Exit(1)
	}

	fields := []int{}
	commaSep := strings.Split(fieldsOpt, ",")
	for _, fieldStr := range commaSep {
		if strings.Contains(fieldStr, "-") {
			hsep := strings.SplitN(fieldStr, "-", 2)
			if len(hsep) != 2 {
				log.Printf("invalid specified hyphen-fields: %s", fieldStr)
				os.Exit(1)
			}
			beginStr, endStr := hsep[0], hsep[1]
			begin, err := strconv.Atoi(beginStr)
			if err != nil {
				log.Printf("fail to parse integer: %s", fieldStr)
				os.Exit(1)
			}
			end, err := strconv.Atoi(endStr)
			if err != nil {
				log.Printf("fail to parse integer: %s", fieldStr)
				os.Exit(1)
			}
			for i := begin; i <= end; i++ {
				fields = append(fields, i)
			}
		} else {
			field, err := strconv.Atoi(fieldStr)
			if err != nil {
				log.Printf("fail to parse integer: %s", fieldStr)
				os.Exit(1)
			}
			fields = append(fields, field)
		}
	}

	scanner := bufio.NewScanner(os.Stdin)
	p := shellwords.NewParser()

	line := 0
	for scanner.Scan() {
		line++
		ss, err := p.Parse(scanner.Text())
		if err != nil {
			log.Printf("fail to parse line: line=%d, text=%s", line, scanner.Text())
			os.Exit(1)
		}
		outputs := make([]string, 0, len(fields))
		for _, field := range fields {
			if len(ss) <= field-1 {
				continue
			}
			outputs = append(outputs, ss[field-1])
		}
		outputsStr := strings.Join(outputs, " ")
		fmt.Println(outputsStr)
	}
}
