package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

var newLine []byte = []byte("\n")

func main() {
	csvPath := "./files/source.csv"
	err := createCsv(csvPath)
	if err != nil {
		fmt.Println("source csv create error: ", err.Error())
		return
	}

	err = splitCsv(csvPath, 12)
	if err != nil {
		fmt.Println("error spliting csv file: ", err.Error())
	}
}

func splitCsv(csvPath string, fileRowSize int) error {
	file, err := os.Open(csvPath)
	if err != nil {
		return err
	}

	defer file.Close()

	rowCounter := 0
	reader := bufio.NewReader(file)
	var header []byte = nil
	var outputFile *os.File
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		if header == nil {
			header = line
			continue
		}

		fileNumber := rowCounter / fileRowSize
		if rowCounter%fileRowSize == 0 {
			outputFile, err = os.Create(fmt.Sprintf("./files/output_%d.csv", fileNumber))
			if err != nil {
				return err
			}

			defer outputFile.Close()

			err = writeLine(outputFile, &header)
			if err != nil {
				return err
			}
		}

		err = writeLine(outputFile, &line)
		if err != nil {
			return err
		}

		rowCounter++
	}

	return nil
}

func writeLine(file *os.File, line *[]byte) error {
	_, err := file.Write(*line)
	if err != nil {
		return err
	}

	_, err = file.Write(newLine)
	if err != nil {
		return err
	}

	return nil
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
