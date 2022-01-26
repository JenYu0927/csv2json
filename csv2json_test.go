package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
)

func TestValidReadCSV(t *testing.T) {
	testFileName := "csv/input.csv"
	employees, err := readCSV(testFileName)
	assert.Equal(t, err, nil)

	expectedEmployees := []Employee{}
	expectedEmployees = append(expectedEmployees, Employee{ID: 1, FirstName: "Marc", LastName: "Smith", Email: "marc@glasnostic.com",
		Description: "Writer of Java", Role: "Dev", Phone: "541-754-3010"})

	expectedEmployees = append(expectedEmployees, Employee{ID: 2, FirstName: "John", LastName: "Young", Email: "john@glasnostic.com",
		Description: "Interested in MHW", Role: "HR", Phone: "541-75-3010"})

	expectedEmployees = append(expectedEmployees, Employee{ID: 3, FirstName: "Peter", LastName: "Scott", Email: "peter@glasnostic.com",
		Description: "amateur boxer", Role: "Dev", Phone: "541-754-3010"})

	assert.Equal(t, employees, expectedEmployees)

}

func TestInvalidReadCSV(t *testing.T) {
	testFileName := "csv/error.csv"
	employees, err := readCSV(testFileName)
	expectedEmployees := []Employee{}

	assert.Equal(t, err, nil)
	assert.Equal(t, employees, expectedEmployees)
}

func TestEmptyReadCSV(t *testing.T) {
	testFileName := "csv/empty.csv"
	employees, err := readCSV(testFileName)
	expectedEmployees := []Employee{}

	assert.Equal(t, err, nil)
	assert.Equal(t, employees, expectedEmployees)
}

func TestWriteToJsonFile(t *testing.T) {
	var employees []Employee = []Employee{}
	employees = append(employees, Employee{ID: 1, FirstName: "Marc", LastName: "Smith", Email: "marc@glasnostic.com",
		Description: "Writer of Java", Role: "Dev", Phone: "541-754-3010"})

	employees = append(employees, Employee{ID: 2, FirstName: "John", LastName: "Young", Email: "john@glasnostic.com",
		Description: "Interested in MHW", Role: "HR", Phone: "541-75-3010"})

	employees = append(employees, Employee{ID: 3, FirstName: "Peter", LastName: "Scott", Email: "peter@glasnostic.com",
		Description: "amateur boxer", Role: "Dev", Phone: "541-754-3010"})

	outputFileName := "unit_test.json"
	err := writeToJsonFile(employees, outputFileName)
	assert.Equal(t, err, nil)

	_, err = os.Stat(outputFileName)
	defer os.Remove(outputFileName)
	assert.Equal(t, err, nil)

}

func TestValidCommand(t *testing.T) {
	os.Args = []string{"./csv2json", "-o", "unit_test.json", "csv/input.csv"} // mock CLI for unit test

	var outputFileName string
	var outputFormat string
	outputFormatFunctions := map[string]func([]Employee, string) error{"json": writeToJsonFile} // determine by variable output_format

	app := cli.NewApp()
	app.Name = "csv2json"
	app.Usage = "convert csv file to json file"

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "output,o",
			Usage:       "Specify the output file name",
			Value:       "Employees",
			Destination: &outputFileName,
		},
		&cli.StringFlag{
			Name:        "format,f",
			Usage:       "Specify format of output file",
			Value:       "json",
			Destination: &outputFormat,
		},
	}
	app.Action = func(c *cli.Context) error {
		if c.NArg() < 1 {
			fmt.Println("Designate a file to parse")
		} else if c.NArg() > 1 {
			fmt.Println("This tool require only one parameter as target file")
			return fmt.Errorf("This tool require only one parameter as target file")
		} else if _, err := os.Stat(c.Args().Get(0)); os.IsNotExist(err) { // Lack input file
			fmt.Println("Designated file dose not exist")
			return err
		} else if _, ok := outputFormatFunctions[outputFormat]; !ok { // Non supported output format
			fmt.Println("Wrong output format")
		} else {
			employees, err := readCSV(c.Args().Get(0)) // read CSV into array of Employee
			if err != nil {
				fmt.Println("Parsing csv file failed. Error message:", err)
				return err
			} else if len(employees) == 0 {
				fmt.Println("Can't parse any employee from target file")
			}

			err = outputFormatFunctions[outputFormat](employees, outputFileName) // call decode function
			if err != nil {
				return err
			}
		}
		return nil
	}
	app.Run(os.Args)

	err := app.Run(os.Args)
	assert.Equal(t, err, nil)
	os.Remove("unit_test.json")
}
