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

type matcher struct {
	field int
	value string
}

func (m matcher) Match(row []string) bool {
	if len(row) <= m.field-1 {
		return false
	}
	return row[m.field-1] == m.value
}

type matchers []matcher

func (m matchers) Match(row []string) bool {
	for _, mm := range m {
		if !mm.Match(row) {
			return false
		}
	}
	return true
}

func main() {
	var fieldsOpt string
	var grepOpt string
	flag.StringVar(&fieldsOpt, "f", "", "indices on a field for output. 1-origin. Ex. -f 1,2 or -f 1-10")
	flag.StringVar(&grepOpt, "g", "", "filter row by field value. Ex. -g 10=foobar,11=boo")
	flag.Parse()

	if fieldsOpt == "" {
		log.Println("-f is required")
		os.Exit(1)
	}

	showAll := false
	fields := []int{}
	if fieldsOpt == "-" {
		showAll = true
	} else {
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
	}
	filter := matchers{}
	if grepOpt != "" {
		fs := strings.Split(grepOpt, ",")
		for _, _f := range fs {
			fv := strings.SplitN(_f, "=", 2)
			if len(fv) != 2 {
				log.Printf("invalid filter option: %s", _f)
				os.Exit(1)
			}
			fieldStr, value := fv[0], fv[1]
			field, err := strconv.Atoi(fieldStr)
			if err != nil {
				log.Printf("fail to parse field number of filter: %s", _f)
				os.Exit(1)
			}
			filter = append(filter, matcher{
				field: field,
				value: value,
			})
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
		if !filter.Match(ss) {
			continue
		}
		if showAll {
			fmt.Println(scanner.Text())
			continue
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
