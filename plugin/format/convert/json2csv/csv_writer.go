package json2csv

import (
	"io"
	"sort"
	"strconv"
	"sync"

	// "github.com/k0kubun/pp"
	"github.com/sniperkit/xutil/plugin/format/convert/json2csv/jsonpointer"
	csv "github.com/sniperkit/xutil/plugin/format/csv/concurrent-writer"
)

// KeyStyle represents the specific style of the key.
type KeyStyle uint

// Header style
const (
	// "/foo/bar/0/baz"
	JSONPointerStyle KeyStyle = iota

	// "foo/bar/0/baz"
	SlashStyle

	// "foo.bar.0.baz"
	DotNotationStyle

	// "foo.bar[0].baz"
	DotBracketStyle
)

// CSVWriter writes CSV data.
type CSVWriter struct {
	*csv.CsvWriter
	mutex       *sync.Mutex
	HeaderStyle KeyStyle
	Transpose   bool
	hNone       bool
	hDone       bool
	hMap        map[string]string
	rows        int
	cols        int
}

// NewCSVWriter returns new CSVWriter with DotBracketStyle.
func NewCSVWriter(w io.Writer) (*CSVWriter, error) {
	writer, err := csv.NewWriter(w)
	if err != nil {
		return nil, err
	}
	return &CSVWriter{
		writer,
		&sync.Mutex{},
		DotBracketStyle,
		false,
		false,
		false,
		nil,
		0,
		0,
	}, nil
}

// NewCSVWriter returns new CSVWriter with DotBracketStyle.
func NewCSVWriterToFile(filename string) (*CSVWriter, error) {
	writer, err := csv.NewWriterToFile(filename)
	if err != nil {
		return nil, err
	}
	return &CSVWriter{
		writer,
		&sync.Mutex{},
		DotBracketStyle,
		false,
		false,
		false,
		nil,
		0,
		0,
	}, nil
}

func (w *CSVWriter) SetHeaders(hMap map[string]string) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	w.hMap = hMap
	return nil
}

// WriteHeaders as the first row of the CSV outputfile.
func (w *CSVWriter) NoHeaders(status bool) bool {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	w.hNone = status
	return w.hNone
}

// WriteCSV writes CSV data.
func (w *CSVWriter) WriteCSV(results []KeyValue) error {
	if w.Transpose {
		return w.writeTransposedCSV(results)
	}
	return w.writeCSV(results)
}

// WriteCSV writes CSV data.
func (w *CSVWriter) WriteListCSV(results []KeyValue) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	return w.writeListCSV(results)
}

// WriteCSV writes CSV data.
func (w *CSVWriter) writeCSV(results []KeyValue) error {
	// w.mutex.Lock()
	// defer w.mutex.Unlock()

	pts, err := allPointers(results)
	if err != nil {
		return err
	}
	sort.Sort(pts)
	keys := pts.Strings()
	header := w.getHeader(pts)

	if !w.hDone {
		w.hDone = true
		if err := w.Write(header); err != nil {
			return err
		}
	}

	for _, result := range results {
		record := toRecord(result, keys)
		if err := w.Write(record); err != nil {
			return err
		}
	}

	// w.Flush()
	if err := w.Error(); err != nil {
		return err
	}

	return nil
}

// WriteCSV writes CSV data which is transposed rows and columns.
func (w *CSVWriter) writeTransposedCSV(results []KeyValue) error {
	pts, err := allPointers(results)
	if err != nil {
		return err
	}
	sort.Sort(pts)
	keys := pts.Strings()
	header := w.getHeader(pts)

	for i, key := range keys {
		record := toTransposedRecord(results, key, header[i])
		if err := w.Write(record); err != nil {
			return err
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return err
	}

	return nil
}

// WriteCSV writes CSV data which is transposed rows and columns.
func (w *CSVWriter) writeListCSV(results []KeyValue) error {
	// w.mutex.Lock()
	// defer w.mutex.Unlock()

	pts, err := allPointers(results)
	if err != nil {
		return err
	}
	sort.Sort(pts)
	keys := pts.Strings()
	// header := w.getHeader(pts)
	extra := []string{"remote_id", "processed_at"}

	var line []string
	for _, key := range keys {
		record := toTransposedList(results, key)
		line = append(line, record...)
	}
	line = append(line, extra...)
	if err := w.Write(line); err != nil {
		return err
	}

	// log.Fatal("test...")
	w.Flush()
	if err := w.Error(); err != nil {
		return err
	}
	return nil
}

func toTransposedList(results []KeyValue, key string) (record []string) {
	for _, result := range results {
		if value, ok := result[key]; ok {
			record = append(record, toString(value))
		}
	}
	return record
}

func toTransposedList2(results KeyValue, extra ...string) []string {
	record := make([]string, 0, 1+len(extra))
	record = append(record, extra...)
	for _, result := range results {
		record = append(record, toString(result))
	}
	return record
}

func allPointers(results []KeyValue) (pointers pointers, err error) {
	set := make(map[string]bool, 0)
	for _, result := range results {
		for _, key := range result.Keys() {
			// log.Println("key=", key)
			if !set[key] {
				set[key] = true
				pointer, err := jsonpointer.New(key)
				if err != nil {
					return nil, err
				}
				pointers = append(pointers, pointer)
			}
		}
	}
	return
}

func (w *CSVWriter) getHeader(pointers pointers) []string {
	switch w.HeaderStyle {
	case JSONPointerStyle:
		return pointers.Strings()
	case SlashStyle:
		return pointers.Slashes()
	case DotNotationStyle:
		return pointers.DotNotations(false)
	case DotBracketStyle:
		return pointers.DotNotations(true)
	default:
		return pointers.Strings()
	}
}

func printRow(w *CSVWriter, keys []string, d map[string]interface{}) error {
	var record []string
	for _, k := range keys {
		switch f := d[k].(type) {
		case string:
			record = append(record, f)
		case float64:
			record = append(record, strconv.FormatFloat(f, 'f', -1, 64))
		case bool:
			if f {
				record = append(record, "true")
			} else {
				record = append(record, "false")
			}
		default:
			log.Fatalf("Unsupported type %T. Aborting.\n", f)
		}
	}
	return w.Write(record)
}

func toRecord(kv KeyValue, keys []string) []string {
	record := make([]string, 0, len(keys))
	for _, key := range keys {
		if value, ok := kv[key]; ok {
			record = append(record, toString(value))
		} else {
			record = append(record, "")
		}
	}
	return record
}

func toTransposedRecord(results []KeyValue, key string, header string) []string {
	record := make([]string, 0, len(results)+1)
	record = append(record, header)
	for _, result := range results {
		if value, ok := result[key]; ok {
			record = append(record, toString(value))
		} else {
			record = append(record, "")
		}
	}
	return record
}
