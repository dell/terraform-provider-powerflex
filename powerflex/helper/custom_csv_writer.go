package helper

import (
	"bufio"
	"io"
	"strings"
)

// CustomCSVWriter desfines the srtuct for the CSV Writer
type CustomCSVWriter struct {
	writer io.Writer
	buf    *bufio.Writer
}

// NewCustomCSVWriter function for the Custom CSV Writer
func NewCustomCSVWriter(w io.Writer) *CustomCSVWriter {
	return &CustomCSVWriter{
		writer: w,
		buf:    bufio.NewWriter(w),
	}
}

// Write function allows the CustomCSVWriter to write CSV records to a specified writer (e.g., file, buffer) one row at a time.
// The function takes care of concatenating the record's values and adding new lines between rows while handling any potential errors that might occur during the write process
func (c *CustomCSVWriter) Write(record []string) error {
	_, err := c.writer.Write([]byte(strings.Join(record, ",")))
	if err != nil {
		return err
	}
	_, err = c.writer.Write([]byte("\n"))
	return err
}

// Flush flushes any buffered data to the underlying writer
func (cw *CustomCSVWriter) Flush() error {
	return cw.buf.Flush()
}
