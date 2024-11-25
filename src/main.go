package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("Missing filename, please include file name for the parsing")
		os.Exit(1)
	}
	data, err := read(args[0])
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("Parsed Data: %v\n", data)
	}
}

func read(file string) (interface{}, error) {
	// Read file
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	return JSONParser(string(data))
}

func JSONParser(item string) (interface{}, error) {
	// Remove newlines and whitespace
	item = strings.TrimSpace(item)
	if len(item) == 0 {
		return nil, errors.New("empty input")
	}

	firstChar := item[0]
	lastChar := item[len(item)-1]

	// Object parsing
	if firstChar == '{' && lastChar == '}' {
		if len(item) == 2 {
			return make(map[string]interface{}), nil
		}

		var res = make(map[string]interface{})
		objectPairs, err := ArrayAndObjectSplitter(item[1:len(item)-1], ',')
		if err != nil {
			return nil, err
		}

		for _, pair := range objectPairs {
			keyValue, err := ArrayAndObjectSplitter(pair, ':')
			if err != nil || len(keyValue) != 2 {
				return nil, fmt.Errorf("invalid key-value pair: %s", pair)
			}
			key, err := JSONParser(strings.TrimSpace(keyValue[0]))
			if err != nil {
				return nil, err
			}
			value, err := JSONParser(strings.TrimSpace(keyValue[1]))
			if err != nil {
				return nil, err
			}
			res[fmt.Sprint(key)] = value
		}
		return res, nil

		// Array parsing
	} else if firstChar == '[' && lastChar == ']' {
		if len(item) == 2 {
			return []interface{}{}, nil
		}

		elements, err := ArrayAndObjectSplitter(item[1:len(item)-1], ',')
		if err != nil {
			return nil, err
		}

		var res []interface{}
		for _, element := range elements {
			parsed, err := JSONParser(strings.TrimSpace(element))
			if err != nil {
				return nil, err
			}
			res = append(res, parsed)
		}
		return res, nil

		// String parsing
	} else if firstChar == '"' && lastChar == '"' {
		return item[1 : len(item)-1], nil

		// Boolean parsing
	} else if item == "true" {
		return true, nil
	} else if item == "false" {
		return false, nil

		// Null parsing
	} else if item == "null" {
		return nil, nil

		// Number parsing
	} else {
		number, err := strconv.ParseFloat(item, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid number: %s", item)
		}
		return number, nil
	}
}

func ArrayAndObjectSplitter(item string, splitterChar rune) ([]string, error) {
	var res []string
	level := 0
	startIndex := 0

	for i, char := range item {
		if char == '{' || char == '[' {
			level++
		} else if char == '}' || char == ']' {
			level--
		} else if char == splitterChar && level == 0 {
			res = append(res, item[startIndex:i])
			startIndex = i + 1
		}
	}

	if level != 0 {
		return nil, errors.New("mismatched brackets in JSON")
	}

	res = append(res, item[startIndex:])
	return res, nil
}
