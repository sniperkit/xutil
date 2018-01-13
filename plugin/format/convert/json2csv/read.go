package json2csv

import (
	"fmt"
	"io"
	"math"
	"strings"

	csv "github.com/sniperkit/xutil/plugin/format/csv/concurrent-writer"
	json "github.com/sniperkit/xutil/plugin/format/json"
)

type LineReader interface {
	ReadBytes(delim byte) (line []byte, err error)
}

func getValue(data map[string]interface{}, keyparts []string) string {
	if len(keyparts) > 1 {
		subdata, _ := data[keyparts[0]].(map[string]interface{})
		return getValue(subdata, keyparts[1:])
	} else if v, ok := data[keyparts[0]]; ok {
		switch v.(type) {
		case nil:
			return ""
		case float64:
			f, _ := v.(float64)
			if math.Mod(f, 1.0) == 0.0 {
				return fmt.Sprintf("%d", int(f))
			} else {
				return fmt.Sprintf("%f", f)
			}
		default:
			return fmt.Sprintf("%+v", v)
		}
	}
	return ""
}

func lineReader(r LineReader, w *csv.CsvWriter, keys []string, printHeader bool) {
	var line []byte
	var err error
	line_count := 0

	var expanded_keys [][]string
	for _, key := range keys {
		expanded_keys = append(expanded_keys, strings.Split(key, "."))
	}

	for {
		if err == io.EOF {
			return
		}
		line, err = r.ReadBytes('\n')
		if err != nil {
			if err != io.EOF {
				log.Printf("Input ERROR: %s", err)
				break
			}
		}
		line_count++
		if len(line) == 0 {
			continue
		}

		if printHeader {
			w.Write(keys)
			w.Flush()
			printHeader = false
		}

		var data map[string]interface{}
		err = json.Unmarshal(line, &data)
		if err != nil {
			log.Printf("ERROR Decoding JSON at line %d: %s\n%s", line_count, err, line)
			continue
		}

		var record []string
		for _, expanded_key := range expanded_keys {
			record = append(record, getValue(data, expanded_key))
		}

		w.Write(record)
		w.Flush()
	}
}
