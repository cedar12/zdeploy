// Copyright 2022 cedar12, cedar12.zxd@qq.com
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

package conn

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	. "zdeploy/common"
	"zdeploy/config"
	"zdeploy/progress"
)

func Client(parser config.IniParser) {
	host := parser.GetString("server", "host")
	port := parser.GetString("server", "port")
	src := parser.GetString("deploy", "src")
	f, err := os.Open(src)
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
	defer f.Close()

	fi, err := os.Stat(src)
	if err != nil {
		log.Fatalln(err.Error())
		return
	}

	conn, err := net.Dial("tcp", host+":"+port)
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Println("connect ", conn.RemoteAddr().String())
	defer conn.Close()

	version, err := ReadStr(conn)
	if err != nil {
		log.Println(err.Error())
		return
	}
	if version != Version {
		log.Println("Version mismatch\ncurrent version：", Version, "\ntarget version：", version)
		return
	}

	s := fmt.Sprintf("%x", md5.Sum([]byte(parser.GetString("server", "pass"))))
	conn.Write([]byte(s))

	buf := make([]byte, 1)
	_, err = conn.Read(buf)
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
	if buf[0] != Ok {
		log.Fatalln("Authentication failed")
		return
	}
	log.Println("Certification passed")
	dist := parser.GetString("deploy", "dist")
	if dist == "" {
		log.Fatalln("Configuration error [deploy]->dist")
		return
	}
	conn.Write([]byte(dist))
	log.Println("send dist path", dist)
	buf = make([]byte, 1)
	_, err = conn.Read(buf)
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
	if buf[0] != Ok {
		log.Fatalln("server not ok")
		return
	}

	conn.Write(Int64ToBytes(fi.Size()))

	buf = make([]byte, 8)
	n, err := conn.Read(buf)
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
	fileBufSize := BytesToInt64(buf[:n])

	log.Println("file ready buf ", fileBufSize)

	p := progress.NewProgress(0, fi.Size())
	for {
		buf = make([]byte, fileBufSize)
		n, err = f.Read(buf)
		if err != nil && io.EOF == err {
			log.Println("File sending completed, waiting for receiving to complete")
			break
		}
		n, err = conn.Write(buf[:n])

		p.Add(int64(n))

	}

	buf = make([]byte, 1)
	_, err = conn.Read(buf)
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
	if buf[0] != Ok {
		log.Fatalln("server file not ok")
		return
	}

	cmdPath := parser.GetString("cmd", "path")
	log.Println("send cmd ", cmdPath)
	if cmdPath == "" {
		return
	}

	cmds := strings.Split(cmdPath, ",")

	conn.Write(Int64ToBytes(int64(len(cmds))))

	isOk := ReadBool(conn)
	if !isOk {
		log.Println("server not cmd count")
		return
	}

	for _, cmd := range cmds {
		conn.Write([]byte(cmd))

		cmdResult, err := ReadStr(conn)
		if err != nil {
			log.Println(err.Error())
			return
		}
		log.Print(cmd, " -> ", strings.Replace(strings.Replace(cmdResult, Fail+Delimiter, EmptyStr, 1), Success+Delimiter, EmptyStr, 1))
		if strings.Index(cmdResult, Fail+Delimiter) == 0 {
			return
		}

	}

	log.Println("finish")
}
