package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// DataFormat represents a supported data format.
type DataFormat int

const (
	Unknown DataFormat = iota
	CSV
	// Add additional data formats here...
)

// Employee represents an employee record.
type Employee struct {
	ID          string `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	Description string `json:"description"`
	Role        string `json:"role"`
	Phone       string `json:"phone"`
}

// Converter is responsible for converting data from one format to another.
type Deserializer interface {
	Deserialize(reader io.Reader) ([]Employee, error)
}

// CSVConverter is a Converter that can read data from CSV files.
type CSVDeserializer struct{}

// Convert reads data from a CSV file and returns a slice of Employee records.
func (d CSVDeserializer) Deserialize(reader io.Reader) ([]Employee, error) {
	csvReader := csv.NewReader(reader)
	csvReader.Read() // skip the first line

	var employees []Employee
	for {
		line, err := csvReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		employee := Employee{
			ID:          line[0],
			FirstName:   line[1],
			LastName:    line[2],
			Email:       line[3],
			Description: line[4],
			Role:        line[5],
			Phone:       line[6],
		}

		employees = append(employees, employee)
	}

	return employees, nil
}

// JSONConverter is a Converter that can write data to JSON files.
type JSONSerializer struct{}

// Convert writes data to a JSON file and returns an error if any.
func (s JSONSerializer) Serialize(writer io.Writer, employees []Employee) error {
	jsonData, err := json.MarshalIndent(employees, "", "    ")
	if err != nil {
		return err
	}

	_, err = writer.Write(jsonData)
	if err != nil {
		return err
	}

	return nil
}

// GetConverter returns a Converter that can handle the specified data format.
func NewDeserialize(format DataFormat) (Deserializer, error) {
	switch format {
	case CSV:
		return CSVDeserializer{}, nil
	// Add additional data formats here...
	default:
		return nil, errors.New("unsupported data format")
	}
}

// DetectDataFormat returns the data format of the specified file based on its file extension.
func DetectDataFormat(filename string) (DataFormat, error) {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".csv":
		return CSV, nil
	// Add additional file extensions here...
	default:
		return Unknown, errors.New("unknown file extension")
	}
}

func decode(inputFilePath string) ([]Employee, error) {
	format, err := DetectDataFormat(inputFilePath)
	if err != nil {
		return nil, err
	}

	deserializer, err := NewDeserialize(format)
	if err != nil {
		return nil, err
	}

	inputReader, err := os.Open(inputFilePath)
	if err != nil {
		return nil, err
	}
	defer inputReader.Close()

	employees, err := deserializer.Deserialize(inputReader)
	if err != nil {
		return nil, err
	}

	return employees, nil
}

func encode(employees []Employee, outputFilePath string) error {
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		return fmt.Errorf("error creating output file: %s", err)
	}
	defer outputFile.Close()

	serializer := JSONSerializer{}
	if err := serializer.Serialize(outputFile, employees); err != nil {
		return fmt.Errorf("error serializing output file: %s", err)
	}

	return nil
}

func main() {
	inputFile := flag.String("input", "", "input file path")
	outputFileName := flag.String("output", "", "output file name")
	flag.Parse()

	if *inputFile == "" || *outputFileName == "" {
		fileName := filepath.Base(os.Args[0])
		fmt.Printf("Usage: ./%s -input=input_file -output=output_file\n", fileName)
		os.Exit(1)
	}

	employees, err := decode(*inputFile)
	if employees == nil {
		log.Fatalf("Error decoding input file: %s", err)
	}

	err = encode(employees, *outputFileName)
	if err != nil {
		log.Fatalf("Error writing output file: %s", err)
	} else {
		fmt.Println("Convert csv to json sucessfully")
	}

}
