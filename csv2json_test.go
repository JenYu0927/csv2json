package main

import (
	"bytes"
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
	/*
			expected := `[
		    {
		        "id": "1",
		        "first_name": "John",
		        "last_name": "Doe",
		        "email": "john.doe@example.com",
		        "description": "John is a software engineer with 5 years of experience.",
		        "role": "Software Engineer",
		        "phone": "555-555-5555"
		    },
		    {
		        "id": "2",
		        "first_name": "Jane",
		        "last_name": "Smith",
		        "email": "jane.smith@example.com",
		        "description": "Jane is a project manager with 10 years of experience.",
		        "role": "Project Manager",
		        "phone": "555-555-5555"
		    }
			]`
	*/

	serializer := JSONSerializer{}
	buffer := new(bytes.Buffer)
	err := serializer.Serialize(buffer, employees)

	assert.NoError(t, err)
	assert.NotEmpty(t, buffer.Bytes())

}
