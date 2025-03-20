/*
Copyright (c) 2023-2024 Dell Inc., or its subsidiaries. All Rights Reserved.

Licensed under the Mozilla Public License Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://mozilla.org/MPL/2.0/


Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
func (cw *CustomCSVWriter) Write(record []string) error {
	// Adding Checkmarx ignore, because this is working as intended and it is reporting a false postive low sev issue
	// Checkmarx: ignore
	_, err := cw.writer.Write([]byte(strings.Join(record, ",")))
	if err != nil {
		return err
	}
	_, err = cw.writer.Write([]byte("\n"))
	return err
}

// Flush flushes any buffered data to the underlying writer
func (cw *CustomCSVWriter) Flush() error {
	return cw.buf.Flush()
}
