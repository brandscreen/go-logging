// Copyright 2013, Cong Ding. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// author: Cong Ding <dinggnu@gmail.com>

package logging

import (
	"path"
	"runtime"
)

// This variable maps fields in recordArgs to relevent function signatures
var setRequestFields = map[string]func(*request){
	"filename": (*request).setFilename, // source filename of the caller
	"pathname": (*request).setPathname, // filename with path
	"lineno":   (*request).setLineno,   // line number in source code
	"funcname": (*request).setFuncname, // function name of the caller
}

var runtimeFields = map[string]bool{
	"name":      false,
	"seqid":     false,
	"levelno":   false,
	"levelname": false,
	"created":   false,
	"nsecs":     false,
	"time":      false,
	"timestamp": false,
	"rtime":     false,
	"filename":  true,
	"pathname":  true,
	"module":    false,
	"lineno":    true,
	"funcname":  true,
	"process":   false,
	"message":   false,
}

// request struct stores the logger request
type request struct {
	level  Level
	format string

	// runtime fields, empty unless specified
	pathname string
	filename string
	lineno   int
	funcname string

	v []interface{}
}

// If it fails to get some fields with string type, these fields are set to
// errString value.
const errString = "???"

// genRuntime generates the runtime information, including pathname, function
// name, filename, line number.
func genRuntime(r *request) {
	calldepth := 5
	pc, file, line, ok := runtime.Caller(calldepth)
	if ok {
		// Generate short function name
		fname := runtime.FuncForPC(pc).Name()
		fshort := fname
		for i := len(fname) - 1; i > 0; i-- {
			if fname[i] == '.' {
				fshort = fname[i+1:]
				break
			}
		}

		r.pathname = file
		r.funcname = fshort
		r.filename = path.Base(file)
		r.lineno = line
	} else {
		r.pathname = errString
		r.funcname = errString
		r.filename = errString
		// Here we uses -1 rather than 0, because the default value in
		// golang is 0 and we should know the value is uninitialized
		// or failed to get
		r.lineno = -1
	}
}

// Construct new request, evaluating any runtime arguments
// The runtime arguments depend on the call stack, so must be evaluated
// explicitly in the main context, rather than lazily (asynchronously) in the
// writer context. Thus we pre-format any runtime arguments in this constructor
// and simply save the rest.
func NewRequest(logger *Logger, level Level, format string, v []interface{}) *request {
	r := new(request)
	r.level = level
	r.format = format
	r.v = v
	// Find runtime arguments and explicitly evaluate them
	for _, v := range logger.recordArgs {
		if runtimeFields[v] {
			setRequestFields[v](r)
		}
	}
	return r
}

// File name of calling logger, with whole path
func (r *request) setPathname() {
	if r.pathname == "" {
		genRuntime(r)
	}
}

// File name of calling logger
func (r *request) setFilename() {
	if r.filename == "" {
		genRuntime(r)
	}
}

// Line number
func (r *request) setLineno() {
	if r.lineno == 0 {
		genRuntime(r)
	}
}

// Function name
func (r *request) setFuncname() {
	if r.funcname == "" {
		genRuntime(r)
	}
}

