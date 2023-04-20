package main

import (
	"bytes"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCSVDeserializer_Deserialize(t *testing.T) {
	csvData := `ID,FirstName,LastName,Email,Description,Role,Phone
1,John,Doe,john.doe@example.com,"John is a software engineer with 5 years of experience.",Software Engineer,555-555-5555
2,Jane,Smith,jane.smith@example.com,"Jane is a project manager with 10 years of experience.",Project Manager,555-555-5555`

	expected := []Employee{
		{
			ID:          "1",
			FirstName:   "John",
			LastName:    "Doe",
			Email:       "john.doe@example.com",
			Description: "John is a software engineer with 5 years of experience.",
			Role:        "Software Engineer",
			Phone:       "555-555-5555",
		},
		{
			ID:          "2",
			FirstName:   "Jane",
			LastName:    "Smith",
			Email:       "jane.smith@example.com",
			Description: "Jane is a project manager with 10 years of experience.",
			Role:        "Project Manager",
			Phone:       "555-555-5555",
		},
	}

	reader := strings.NewReader(csvData)
	deserializer := CSVDeserializer{}
	actual, err := deserializer.Deserialize(reader)
	if err != nil {
		t.Errorf("CSVDeserializer.Deserialize() returned an unexpected error: %v", err)
	}

	if len(actual) != len(expected) {
		t.Errorf("CSVDeserializer.Deserialize() returned a slice with length %d; expected %d", len(actual), len(expected))
	}

	for i := range expected {
		if actual[i] != expected[i] {
			t.Errorf("CSVDeserializer.Deserialize() returned an unexpected Employee record: got %v, want %v", actual[i], expected[i])
		}
	}
}

func TestDecode(t *testing.T) {
	// Create a temporary test CSV file.
	csvData := []byte(`ID,FirstName,LastName,Email,Description,Role,Phone
1,John,Doe,john.doe@example.com,"John is a software engineer with 5 years of experience.",Software Engineer,555-555-5555`)
	tempFile, err := os.Create("testInput.csv")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	if _, err := tempFile.Write(csvData); err != nil {
		t.Fatalf("Failed to write to temporary file: %v", err)
	}
	if err := tempFile.Close(); err != nil {
		t.Fatalf("Failed to close temporary file: %v", err)
	}

	// Call decode() and check the result.
	expected := []Employee{
		{
			ID:          "1",
			FirstName:   "John",
			LastName:    "Doe",
			Email:       "john.doe@example.com",
			Description: "John is a software engineer with 5 years of experience.",
			Role:        "Software Engineer",
			Phone:       "555-555-5555",
		},
	}
	actual, err := decode(tempFile.Name())
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Unexpected result. Got: %v, want: %v", actual, expected)
	}
}

func TestJSONSerializer_Serialize(t *testing.T) {
	employees := []Employee{
		{
			ID:          "1",
			FirstName:   "John",
			LastName:    "Doe",
			Email:       "john.doe@example.com",
			Description: "John is a software engineer with 5 years of experience.",
			Role:        "Software Engineer",
			Phone:       "555-555-5555",
		},
		{
			ID:          "2",
			FirstName:   "Jane",
			LastName:    "Smith",
			Email:       "jane.smith@example.com",
			Description: "Jane is a project manager with 10 years of experience.",
			Role:        "Project Manager",
			Phone:       "555-555-5555",
		},
	}

	serializer := JSONSerializer{}
	buffer := new(bytes.Buffer)
	err := serializer.Serialize(buffer, employees)

	assert.NoError(t, err)
	assert.NotEmpty(t, buffer.Bytes())

}

func TestEncode(t *testing.T) {
	// Create a temporary test CSV file.
	employees := []Employee{
		{
			ID:          "1",
			FirstName:   "John",
			LastName:    "Doe",
			Email:       "john.doe@example.com",
			Description: "John is a software engineer with 5 years of experience.",
			Role:        "Software Engineer",
			Phone:       "555-555-5555",
		},
	}
	tempFile, err := os.Create("testOutput.json")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Call encode() and check the result.

	if err := encode(employees, tempFile.Name()); err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	assert.FileExists(t, "testOutput.json")
}
