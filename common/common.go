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

package common

import (
	"encoding/binary"
	"log"
	"net"
	"os"
)

const (
	Version   string = "0.1"
	Success   string = "success"
	Fail      string = "fail"
	Ok        byte   = 1
	Err       byte   = 0
	EmptyStr  string = ""
	ServerStr string = "server"
	Delimiter string = "|"
	Dot       string = "."
	MatchAll  string = "*"
)

func Int64ToBytes(num int64) []byte {
	buf := make([]byte, binary.MaxVarintLen64)
	n := binary.PutVarint(buf, num)
	b := buf[:n]
	return b
}

func BytesToInt64(b []byte) int64 {
	x, _ := binary.Varint(b)
	return x
}

func ReadStr(conn net.Conn) (string, error) {
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		return "", err
	}
	str := string(buf[:n])
	return str, nil
}

func ReadBool(conn net.Conn) bool {
	buf := make([]byte, 1)
	_, err := conn.Read(buf)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	if buf[0] != Ok {
		return false
	}
	return true
}

func ReadInt(conn net.Conn) (int64, error) {
	buf := make([]byte, 8)
	_, err := conn.Read(buf)
	if err != nil {
		return 0, err
	}
	return BytesToInt64(buf), nil
}

func FileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}
