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

package frpc

import (
	"fmt"

	"github.com/fatedier/frp/libs/frpc/core"
)

type CallbackBridge interface {
	OnEvent(action string, path string)
}

//export InitializeBridge
func InitializeBridge(
	to string,
	level string,
	maxDays int,
	cb CallbackBridge,
) {
	core.Initialize(
		to,
		level,
		maxDays,
		func(action core.Action, path string) {
			cb.OnEvent(string(action), path)
		},
	)
}

//export VerifyBridge
func VerifyBridge(path string) string {

	err := core.Verify(path)

	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}

	return "OK"
}

//export StartBridge
func StartBridge(path string) bool {

	err := core.Start(path)

	if err != nil {
		fmt.Printf("start frp client config has error %v \n", err)
		return false
	}

	return true
}

//export StopBridge
func StopBridge(path string) bool {

	err := core.Stop(path)

	if err != nil {
		fmt.Printf("stop frp client config has error %v \n", err)
		return false
	}

	return true
}
