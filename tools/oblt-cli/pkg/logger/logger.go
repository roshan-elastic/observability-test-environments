// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http:// www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/fatih/color"
)

// Verbose enable/disable verbose logs output.
var Verbose bool

var Debug = log.New(&logWriter{Color: *DebugColor, Level: "debug", CheckVerbose: true}, "", 0)
var Info = log.New(&logWriter{Color: *InfoColor, Level: "info", CheckVerbose: false}, "", 0)
var Warn = log.New(&logWriter{Color: *WarnColor, Level: "warn", CheckVerbose: false}, "", 0)
var Error = log.New(&logWriter{Color: *ErrColor, Level: "error", CheckVerbose: false}, "", 0)

type logWriter struct {
	Color        color.Color
	Level        string
	CheckVerbose bool
}

// Write implements io.Writer interface
// it writes the log message to stderr
// it adds a timestamp and the log level
// it returns the number of bytes written and any error
// it checks if verbose is enabled before writing
func (writer logWriter) Write(bytes []byte) (num int, err error) {
	now := time.Now().UTC().Format("2006-01-02T15:04:05.999Z")
	prefix := writer.Color.Sprintf("[%s]:[%s]", now, writer.Level)
	if writer.CheckVerbose && !Verbose {
		num, err = 0, nil
	} else {
		num, err = fmt.Fprintf(os.Stderr, "%s:%s", prefix, string(bytes))
	}
	return num, err
}

// Debugf it prints a debug message only if verbose is enabled
func Debugf(format string, v ...any) {
	Debug.Printf(format, v...)
}

// Infof it prints a info message
func Infof(format string, v ...any) {
	Info.Printf(format, v...)
}

// Warnf it prints warning message
func Warnf(format string, v ...any) {
	Warn.Printf(format, v...)
}

// Errorf it prints error message
func Errorf(format string, v ...any) {
	Error.Printf(format, v...)
}

// Fatalf it prints error message and exit
func Fatalf(format string, v ...any) {
	Error.Printf(format, v...)
	os.Exit(1)
}

// Fatalf it prints error message and exit
func Fatal(v ...any) {
	Error.Print(v...)
	os.Exit(1)
}

// LogError it prints error message and return the error
// It logs the context and the error
func LogError(context string, err error) (string, error) {
	if err != nil {
		Error.Print(context)
		_, filename, line, _ := runtime.Caller(2)
		Error.Printf("%s:%d %v", filename, line, err)
		_, filename, line, _ = runtime.Caller(1)
		Error.Printf("%s:%d %v", filename, line, err)
	} else {
		Debugf("git.log.out: %s", context)
	}
	return context, err
}

// SetOutput sets the output destination for the logger.
func SetOutput(out io.Writer) {
	Debug.SetOutput(out)
	Info.SetOutput(out)
	Warn.SetOutput(out)
	Error.SetOutput(out)
}
