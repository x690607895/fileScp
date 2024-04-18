package main

import (
	"flag"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"time"
)

func main() {
	flag.String("help", "", `
	send 发送文件
	recv 接收文件
	help 帮助
	`)
	flag.Parse()
	switch flag.Arg(0) {
	case "send":
		send()
		break
	case "recv":
		recv()
		break
	default:
		break
	}
}

var (
	letterRunes       = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	short_uuid_length = 7
	udpBoardcastPort  = 5174
	tcpServerPort     = 5175
)

func GenerateShortUUID() string {
	b := make([]rune, short_uuid_length)
	for i := range b {
		if i == 3 {
			b[i] = '-'
		} else {
			b[i] = letterRunes[rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(letterRunes))]
		}
	}
	return string(b)
}

func send() {
	fileName := flag.Arg(1)
	if _, err := os.Lstat(fileName); err != nil {
		log.Fatal(err)
	}
	uuid := GenerateShortUUID()
	log.Println("请在其他机器上运行下面命令：")
	log.Println("go run main.go recv", uuid)
	udpBoardcastClient, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(255, 255, 255, 255),
		Port: udpBoardcastPort,
	})
	checkErr(err)
	tcpServer, err := net.ListenTCP("tcp", &net.TCPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: tcpServerPort,
	})
	checkErr(err)
	log.Println("等待客户端连入")
	var conn net.Conn
	for {
		_, err := udpBoardcastClient.Write([]byte(uuid))
		checkErr(err)
		tcpServer.SetDeadline(time.Now().Add(time.Second))
		conn, err = tcpServer.Accept()
		if err != nil {
			continue
		}
		log.Println("有新的tcp连接加入")
		udpBoardcastClient.Close()
		break
	}
	log.Println("准备发送数据")
	fileBaseName := filepath.Base(fileName)
	conn.Write([]byte(fileBaseName))
	file, err := os.Open(fileName)
	checkErr(err)
	io.Copy(conn, file)
	log.Println("发送成功")
}

func recv() {
	uuid := flag.Arg(1)
	if uuid == "" {
		log.Fatal("uuid不能为空")
	}
	udpBoardcastServer, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: udpBoardcastPort,
	})
	checkErr(err)

	var addr *net.UDPAddr
	for {
		readData := make([]byte, short_uuid_length)
		_, addr, err = udpBoardcastServer.ReadFromUDP(readData)
		checkErr(err)
		if uuid == string(readData) {
			log.Println("发现目标")
			udpBoardcastServer.Close()
			break
		}
	}
	addr.Port = tcpServerPort
	tcpClient, err := net.Dial("tcp", addr.String())
	checkErr(err)
	fileName := make([]byte, 1024)
	n, err := tcpClient.Read(fileName)
	log.Println("接收到的文件：", string(fileName[:n]))
	file, err := os.Create(string(fileName[:n]))
	checkErr(err)
	io.Copy(file, tcpClient)
	log.Println("接收成功")
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
