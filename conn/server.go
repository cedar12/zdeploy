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
	"golang.org/x/text/encoding/simplifiedchinese"
	"log"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
	. "zdeploy/common"
	"zdeploy/config"
	"zdeploy/encrypt"
	"zdeploy/files"
)

func Server(parse config.IniParser) {
	host := parse.GetString(ServerStr, HostStr)
	port := parse.GetString(ServerStr, PortStr)
	listen, err := net.Listen(Network, host+":"+port)
	if err != nil {
		log.Println("Service failed to start ", err.Error())
		return
	}
	log.Println("The service starts successfully，listening port: ", port)
	defer listen.Close()

	for {
		conn, err := listen.Accept()
		if err != nil {
			continue
		}

		conn.Write([]byte(Version))

		err = parse.Reload()
		if err != nil {
			log.Println("Failed to reload ", parse.FileName())
		}
		log.Println(conn.RemoteAddr().String(), " join")
		if ipFilter(conn, parse) {
			go accept(conn, parse)
		} else {
			log.Println(conn.RemoteAddr().String(), " intercepted")
			conn.Close()
		}

	}
}

func ipFilter(conn net.Conn, parse config.IniParser) bool {

	addr := conn.RemoteAddr().String()
	ip := strings.Split(addr, ":")[0]
	whites := parse.GetString(ServerStr, "white")
	if whites == EmptyStr {
		return true
	}
	whiteList := strings.Split(whites, ",")
	for i := range whiteList {
		if IpWhiteMatch(whiteList[i], ip) {
			return true
		}
	}
	return false
}

func accept(conn net.Conn, parse config.IniParser) {
	addr := conn.RemoteAddr().String()

	defer func() {
		log.Println(addr + " exit")
		err := conn.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}()
	passStr := parse.GetString(ServerStr, "pass")

	pass, err := ReadStr(conn)
	if err != nil {
		log.Println(err.Error())
	}
	s := fmt.Sprintf("%x", md5.Sum([]byte(passStr)))
	if encrypt.Decode(pass) != s {
		log.Println(addr, " Authentication failed")
		conn.Write([]byte{Err})
		return
	}
	log.Println(addr, " Certification passed")
	conn.Write([]byte{Ok})

	fileName, err := ReadStr(conn)
	if err != nil {
		log.Println(err.Error())
		return
	}
	conn.Write([]byte{Ok})

	fileSize, err := ReadInt(conn)
	if err != nil {
		log.Println(err.Error())
		return
	}

	if FileExist(fileName) {
		nowTime := time.Now().Format(TimeFormat)
		log.Println("rename ", fileName, " -> ", nowTime+fileName)
		os.Rename(fileName, nowTime+fileName)
	}
	log.Println("create " + fileName)
	f, err := os.Create(fileName)
	if err != nil {
		log.Println("create file error ", err.Error())
	}
	fileBufSize := parse.GetInt64(ServerStr, "buf")
	if fileBufSize <= 0 {
		fileBufSize = 2048
	}
	conn.Write(Int64ToBytes(fileBufSize))

	log.Println("file size ", fileSize)
	var count int
	for {
		buf := make([]byte, fileBufSize)
		n, _ := conn.Read(buf)
		count += n
		f.Write(buf[:n])
		if int64(count) >= fileSize {
			log.Println("File reception completed")
			break
		}
	}

	fileInfo, err := f.Stat()
	if err != nil {
		log.Println(err.Error())
		conn.Write([]byte{Err})
		return
	}
	if fileInfo.Size() != fileSize {
		log.Println("file size mismatch")
		conn.Write([]byte{Err})
		return
	}
	f.Close()
	conn.Write([]byte{Ok})

	cmdCount, err := ReadInt(conn)
	if err != nil {
		log.Println(err.Error())
		return
	}

	conn.Write([]byte{Ok})

	for i := 0; i < int(cmdCount); i++ {
		err = execArgs(conn, fileName, parse)
		if err != nil {
			break
		}
	}

}

func execArgs(conn net.Conn, fileName string, parse config.IniParser) error {
	cmd, err := ReadStr(conn)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	log.Println("receive cmd ", cmd)
	if cmd == "unzip" {
		err = files.Unzip(fileName, "")
		if err != nil {
			log.Println("Unzip error ", err.Error())
			conn.Write([]byte(Fail + Delimiter + err.Error() + "\n"))
			return err
		}
		conn.Write([]byte(Success + Delimiter + "Decompression succeeded\n"))
	} else if cmd == "file" {

	} else {
		arg := parse.GetString("cmd", cmd)
		if arg == "" {
			conn.Write([]byte(Fail + Delimiter + "illegal cmd"))
			return nil
		}
		log.Println("run cmd ", cmd, "["+arg+"]")
		res, err := command(arg)
		if err != nil {
			conn.Write([]byte(Fail + Delimiter + err.Error()))
			return err
		}

		conn.Write([]byte(Success + Delimiter + res))
	}

	return nil
}

func command(arg string) (string, error) {
	args := []string{"/C"}
	cmdName := "cmd"
	if runtime.GOOS != Windows {
		cmdName = "bash"
		args[0] = "-c"
	}
	args = append(args, arg)
	cmd := exec.Command(cmdName, args...)
	var output = make([]byte, 1024)
	var err error
	if output, err = cmd.CombinedOutput(); err != nil {
		log.Print(err)
		return "", err
	}
	result := string(output)
	if runtime.GOOS == Windows {
		result = ConvertByte2String(output, GB18030)
	}
	log.Print(result)
	return result, nil
}

type Charset string

const (
	UTF8    = Charset("UTF-8")
	GB18030 = Charset("GB18030")
)

func ConvertByte2String(byte []byte, charset Charset) string {
	var str string
	switch charset {
	case GB18030:
		var decodeBytes, _ = simplifiedchinese.GB18030.NewDecoder().Bytes(byte)
		str = string(decodeBytes)
	case UTF8:
		fallthrough
	default:
		str = string(byte)
	}
	return str
}

func IpWhiteMatch(ip string, targetIp string) bool {
	ipParts := strings.Split(ip, Dot)
	targetIpParts := strings.Split(targetIp, Dot)

	for i, part := range targetIpParts {
		if ipParts[i] != MatchAll && ipParts[i] != part {
			return false
		}
	}
	return true
}
