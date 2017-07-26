package assistdog

import (
	"strings"
	"testing"

	"github.com/DATA-DOG/godog/gherkin"
)

type person struct {
	Name   string
	Height int
}

func TestCreateInstance(t *testing.T) {
	t.Run("successfully", func(t *testing.T) {
		table := buildTable([][]string{
			{"Name", "John"},
			{"Height", "182"},
		})

		result, err := NewDefault().CreateInstance(new(person), table)
		if err != nil {
			t.Error(err)
			return
		}

		typed := result.(*person)
		if typed.Name != "John" {
			t.Error("expected Name to be John, but was", typed.Name)
		}

		if typed.Height != 182 {
			t.Error("expected Height to be 182, but was", typed.Height)
		}
	})

	t.Run("with extra field", func(t *testing.T) {
		table := buildTable([][]string{
			{"Name", "John"},
			{"Height", "182"},
			{"Age", "25"},
		})

		_, err := NewDefault().CreateInstance(new(person), table)
		if err == nil {
			t.Error("expected an error, but got none")
			return
		}

		expectedMessage := `Age: field not found`
		if !strings.Contains(err.Error(), expectedMessage) {
			t.Errorf(`expected error message to contain "%v", but was "%v"`, expectedMessage, err.Error())
		}
	})

	t.Run("with invalid integer", func(t *testing.T) {
		table := buildTable([][]string{
			{"Name", "John"},
			{"Height", "nono"},
		})

		_, err := NewDefault().CreateInstance(new(person), table)
		if err == nil {
			t.Error("expected an error, but got none")
			return
		}

		expectedMessage := `parsing "nono": invalid syntax`
		if !strings.Contains(err.Error(), expectedMessage) {
			t.Errorf(`expected error message to contain "%v", but was "%v"`, expectedMessage, err.Error())
		}
	})
}

func TestCompareInstance(t *testing.T) {
	t.Run("successfully", func(t *testing.T) {
		table := buildTable([][]string{
			{"Name", "John"},
			{"Height", "182"},
		})

		actual := &person{
			Name:   "John",
			Height: 182,
		}

		err := NewDefault().CompareToInstance(actual, table)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("with different value for int", func(t *testing.T) {
		table := buildTable([][]string{
			{"Name", "John"},
			{"Height", "900"},
		})

		actual := &person{
			Name:   "John",
			Height: 182,
		}

		err := NewDefault().CompareToInstance(actual, table)
		if err == nil {
			t.Error("expected an error, but got none")
			return
		}

		expectedMessage := `Height: expected 900, but got 182`
		if !strings.Contains(err.Error(), expectedMessage) {
			t.Errorf(`expected error message to contain "%v", but was "%v"`, expectedMessage, err.Error())
		}
	})

	t.Run("with different value for string", func(t *testing.T) {
		table := buildTable([][]string{
			{"Name", "Mary"},
			{"Height", "182"},
		})

		actual := &person{
			Name:   "John",
			Height: 182,
		}

		err := NewDefault().CompareToInstance(actual, table)
		if err == nil {
			t.Error("expected an error, but got none")
			return
		}

		expectedMessage := `Name: expected Mary, but got John`
		if !strings.Contains(err.Error(), expectedMessage) {
			t.Errorf(`expected error message to contain "%v", but was "%v"`, expectedMessage, err.Error())
		}
	})
}

func buildTable(src [][]string) *gherkin.DataTable {
	rows := make([]*gherkin.TableRow, len(src))
	for i, row := range src {
		cells := make([]*gherkin.TableCell, len(row))
		for j, value := range row {
			cells[j] = &gherkin.TableCell{Value: value}
		}

		rows[i] = &gherkin.TableRow{Cells: cells}
	}

	return &gherkin.DataTable{Rows: rows}
}