// Package ccsv provides a "thread" safe way of writing to CSV files
package ccsv

import (
	"encoding/csv"
	"io"
	"os"
	"sync"
)

// CsvWriter holds pointers to a Mutex, csv.Writer and the underlying CSV file
type CsvWriter struct {
	mutex     *sync.Mutex
	csvWriter *csv.Writer
	f         *os.File
}

// NewCsvWriter creates a CSV file and returns a CsvWriter
func NewWriterToFile(fileName string) (*CsvWriter, error) {
	csvFile, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}
	w := csv.NewWriter(csvFile)
	return &CsvWriter{csvWriter: w, mutex: &sync.Mutex{}, f: csvFile}, nil
}

// NewCSVWriter returns new CSVWriter with JSONPointerStyle.
func NewWriter(w io.Writer) (*CsvWriter, error) {
	return &CsvWriter{csvWriter: csv.NewWriter(w), mutex: &sync.Mutex{}, f: nil}, nil
}

// Write a single row to a CSV file
func (w *CsvWriter) Write(row []string) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	return w.csvWriter.Write(row)
}

// WriteAll writes multiple rows to a CSV file
func (w *CsvWriter) WriteAll(records [][]string) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	return w.csvWriter.WriteAll(records)
}

func (w *CsvWriter) Error() error {

	err := w.csvWriter.Write(nil)
	return err
}

// Comma is the field delimiter, set to '.'
func (w *CsvWriter) Comma() rune {

	w.mutex.Lock()
	defer w.mutex.Unlock()
	return w.csvWriter.Comma
}

// SetComma takes the passed rune and uses it to set the field
// delimiter for CSV fields.
func (w *CsvWriter) SetComma(r rune) {

	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.csvWriter.Comma = r
}

// UseCRLF exposes the csv writer's UseCRLF field.
func (w *CsvWriter) UseCRLF() bool {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	return w.csvWriter.UseCRLF
}

// SetUseCRLF set's the csv'writer's UseCRLF field
func (w *CsvWriter) SetUseCRLF(b bool) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.csvWriter.UseCRLF = b
}

// Flush writes any pending rows
func (w *CsvWriter) Flush() error {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.csvWriter.Flush()
	return w.csvWriter.Error()
}

// Close CSV file for writing
// Implicitly calls Flush() before
func (w *CsvWriter) Close() error {
	err := w.Flush()
	if err != nil {
		return err
	}
	return w.f.Close()
}
