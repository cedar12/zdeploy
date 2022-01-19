// Copyright 2018 cedar12, cedar12.zxd@qq.com
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

import (
	"fmt"
	"log"
	"os"
	"zdeploy/common"
	"zdeploy/config"
	"zdeploy/conn"
)

func main() {
	parser := config.IniParser{}
	if len(os.Args) == 2 {
		parser.Load(os.Args[1])
		conn.Client(parser)
	} else if len(os.Args) > 2 {
		if os.Args[1] == "-s" {
			parser.Load(os.Args[2])
			conn.Server(parser)
		} else if os.Args[1] == "-c" {
			parser.Load(os.Args[2])
			conn.Client(parser)
		} else {
			log.Fatalln("invalid option " + os.Args[1])
		}
	} else {
		fmt.Println("zdeploy is Deployment file tool\nversion: ", common.Version)
		fmt.Println("Usage:\n\tzdeploy [options] [config file]")
		fmt.Println("Options:\n\t-s\tstart the server\n\t-c\tstart the client")
		fmt.Println("Source: https://github.com/cedar12/zdeploy")
		fmt.Println("Author: cedar12.zxd@qq.com")
	}
}
