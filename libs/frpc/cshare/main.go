// Copyright 2021 The frp Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

/*
#include <stdlib.h>

typedef void (*ClientChangedCallback)(const char* action, const char* path);

static inline void call(ClientChangedCallback cb, const char* action, const char* path) {
	if (cb != NULL) {
		cb(action, path);
	}
}
*/
import "C"

import (
	"fmt"
	"unsafe"

	"github.com/fatedier/frp/libs/frpc/core"
)

//export InitializeBridge
func InitializeBridge(
	to *C.char,
	level *C.char,
	maxDays C.int,
	cb C.ClientChangedCallback,
) {
	core.Initialize(
		C.GoString(to),
		C.GoString(level),
		int(maxDays),
		func(action, path string) {
			caction := C.CString(action)
			cpath := C.CString(path)
			defer C.free(unsafe.Pointer(caction))
			defer C.free(unsafe.Pointer(cpath))
			C.call(cb, caction, cpath)
		},
	)

	fmt.Printf("initialize finished \n")
}

//export VerifyBridge
func VerifyBridge(path *C.char) C.int {

	err := core.Verify(C.GoString(path))

	if err != nil {
		fmt.Printf("verify frp client config has error %v \n", err)
		return 0
	}

	return 1
}

//export StartBridge
func StartBridge(path *C.char) C.int {

	err := core.Start(C.GoString(path))

	if err != nil {
		fmt.Printf("start frp client config has error %v \n", err)
		return 0
	}

	return 1
}

//export StopBridge
func StopBridge(path *C.char) C.int {

	err := core.Stop(C.GoString(path))

	if err != nil {
		fmt.Printf("stop frp client config has error %v \n", err)
		return 0
	}

	return 1
}

func main() {}
