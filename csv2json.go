package main

import (
	"fmt"
	"os"
	"strconv"

	"encoding/csv"
	"encoding/json"

	"github.com/urfave/cli"
)

type Employee struct {
	ID          int
	FirstName   string `json:"FirstName"`
	LastName    string `json:"LastName"`
	Email       string `json:"Email"`
	Description string `json:"Description"`
	Role        string `json:"Role"`
	Phone       string `json:"Phone"`
}

func readCSV(csvFileName string) ([]Employee, error) {
	csvFile, err := os.Open(csvFileName)
	defer csvFile.Close()
	if err != nil {
		return nil, err
	}

	csvReader := csv.NewReader(csvFile)
	lines, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	var employees []Employee = []Employee{}
	for index, line := range lines {
		if index != 0 { // skip first line which is the title of below lines
			id, err := strconv.Atoi(line[0])
			if err != nil {
				return nil, err
			}
			employee := Employee{ID: id, FirstName: line[1], LastName: line[2],
				Email: line[3], Description: line[4], Role: line[5], Phone: line[6]}
			employees = append(employees, employee)
		}
	}

	return employees, nil
}

func writeToJsonFile(employees []Employee, outputFileName string) error {
	employeesByte, err := json.Marshal(employees)
	if err != nil {
		return err
	}
	//fmt.Println(string(employeesByte))

	outputFile, err := os.OpenFile(outputFileName, os.O_WRONLY|os.O_CREATE, 0664)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	_, err = outputFile.Write(employeesByte)
	if err != nil {
		return err
	}

	return nil
}

/*
func writeToYaml([]Employee) error {
	return nil
}*/

func main() {
	var outputFileName string
	var outputFormat string
	outputFormatFunctions := map[string]func([]Employee, string) error{"json": writeToJsonFile} // determine by variable outputFormat

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
		//println(outputFileName, outputFormat, outputFormatFunctions[outputFormat], reflect.TypeOf(outputFormatFunctions[outputFormat]))
		if c.NArg() < 1 {
			fmt.Println("Designate a file to parse")
		} else if c.NArg() > 1 {
			fmt.Println("This tool require only one parameter as target file")
		} else if _, err := os.Stat(c.Args().Get(0)); os.IsNotExist(err) {
			fmt.Println("Designated file dose not exist")
			return err
		} else if _, ok := outputFormatFunctions[outputFormat]; !ok { // Non supported output format
			fmt.Println("Wrong output format")
		} else {
			employees, err := readCSV(c.Args().Get(0)) // read CSV into array of Employee
			//fmt.Println(employees)
			if err != nil {
				fmt.Println("Parsing csv file failed. Error message:", err)
				return err
			} else if len(employees) == 0 {
				fmt.Println("Can't parse any employee from target file")
			}
			//fmt.Println(employees, reflect.TypeOf(employees))

			err = outputFormatFunctions[outputFormat](employees, outputFileName) // call decode function
			if err != nil {
				return err
			}
		}
		return nil
	}
	//fmt.Println(os.Args)
	app.Run(os.Args)
}
