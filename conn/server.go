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
	"log"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
	. "zdeploy/common"
	"zdeploy/config"
	"zdeploy/files"
)

func Server(parse config.IniParser) {
	host := parse.GetString(ServerStr, "host")
	port := parse.GetString(ServerStr, "port")
	listen, err := net.Listen("tcp", host+":"+port)
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

	conn.Write([]byte(Version))

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

func ip2binary(ip string) string {
	str := strings.Split(ip, Dot)
	var ipstr string
	for _, s := range str {
		i, err := strconv.ParseUint(s, 10, 8)
		if err != nil {
			fmt.Println(err)
		}
		ipstr = ipstr + fmt.Sprintf("%08b", i)
	}
	return ipstr
}

func ipMatch(ip, ipRange string) bool {
	ipb := ip2binary(ip)
	ipr := strings.Split(ipRange, "/")
	masklen, err := strconv.ParseUint(ipr[1], 10, 32)
	if err != nil {
		fmt.Println(err)
		return false
	}
	iprb := ip2binary(ipr[0])
	return strings.EqualFold(ipb[0:masklen], iprb[0:masklen])
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
	passSrc := parse.GetString(ServerStr, "pass")

	pass, err := ReadStr(conn)
	if err != nil {
		log.Println(err.Error())
	}
	s := fmt.Sprintf("%x", md5.Sum([]byte(passSrc)))
	if pass != s {
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

	buf := make([]byte, 8)
	n, err := conn.Read(buf)
	if err != nil {
		log.Println(err.Error())
		return
	}
	fileSize := BytesToInt64(buf[:n])

	if FileExist(fileName) {
		log.Println("rename ", fileName)
		nowtime := time.Now().Format("20060102150405")
		os.Rename(fileName, nowtime+fileName)
	}
	log.Println("create " + fileName)
	f, err := os.Create(fileName)
	if err != nil {
		log.Println("create file error ", err.Error())
	}
	fileBufSize := parse.GetInt64("server", "buf")
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
	if runtime.GOOS != "windows" {
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
	log.Print(result)
	return result, nil
}

func CheckIp(ip string) bool {
	parts := strings.Split(ip, Dot)
	for _, part := range parts {
		intPart, err := strconv.Atoi(part)
		if err != nil || intPart < 0 || intPart > 255 {
			return false
		}
	}
	return true
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
