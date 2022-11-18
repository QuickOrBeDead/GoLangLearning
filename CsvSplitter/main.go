package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {
	csvPath := "./files/source.csv"
	err := createCsv(csvPath)
	if err != nil {
		fmt.Println("source csv create error: ", err.Error())
		return
	}

	file, err := os.Open(csvPath)
	if err != nil {
		fmt.Println("source csv open error: ", err.Error())
		return
	}

	defer file.Close()

	reader := bufio.NewReader(file)
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Println("error reading csv line", err.Error())
			break
		}

		fmt.Printf("%s \n", line)
	}
}

func createCsv(path string) error {
	source, err := os.Create(path)
	if err != nil {
		return err
	}

	defer source.Close()

	_, err = source.WriteString("name,email,phone number,address\n")
	if err != nil {
		return err
	}

	for i := 0; i < 50; i++ {
		source.WriteString(fmt.Sprintf("Example%d,example%d@example.com,222-22-22,Example Address%d\n", i, i, i))
	}

	return nil
}
