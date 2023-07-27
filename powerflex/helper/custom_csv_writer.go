package helper

import (
	"bufio"
	"io"
	"strings"
)

type CustomCSVWriter struct {
	writer io.Writer
	buf    *bufio.Writer
}

func NewCustomCSVWriter(w io.Writer) *CustomCSVWriter {
	return &CustomCSVWriter{
		writer: w,
		buf:    bufio.NewWriter(w),
	}
}

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
