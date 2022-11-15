package main

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
)

func main() {
	filename := "birthday_001.txt"
	regex, err := regexp.Compile(`(.*)_([0-9]+).txt`)
	if err != nil {
		fmt.Println("regex compile error: ", err.Error())
	}

	matches := regex.FindStringSubmatch(filename)

	n, _ := strconv.Atoi(matches[2])
	total := int(math.Pow10(len(matches[2])))
	fmt.Printf("%s (%d of %d).txt", matches[1], n, total)
}
