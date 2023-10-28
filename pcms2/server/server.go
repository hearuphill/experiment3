package main

import (
	"fmt"
	"net"
)

func main() {
	databaseConn, err := net.Dial("tcp", "localhost:9988")
	if err != nil {
		fmt.Println("无法连接数据库服务器:", err)
		return
	}
	defer databaseConn.Close()

	listen, err := net.Listen("tcp", "localhost:9999")
	if err != nil {
		fmt.Println("服务器启动失败:", err)
		return
	}
	defer listen.Close()
	fmt.Println("服务器已启动，等待客户端连接...")

	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("客户端连接失败:", err)
			continue
		}
		go handleClient(conn, databaseConn)
	}
}

func handleClient(conn net.Conn, databaseConn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("读取数据错误:", err)
		return
	}

	clientData := string(buf)
	fmt.Println("客户端请求:", clientData)

	// 向数据库服务器发送请求
	_, err = databaseConn.Write([]byte(clientData))
	if err != nil {
		fmt.Println("向数据库服务器发送请求错误:", err)
		conn.Write([]byte("无法处理请求"))
		return
	}

	// 从数据库服务器接收响应
	databaseResponse := make([]byte, 1024)
	n, err := databaseConn.Read(databaseResponse)
	if err != nil {
		fmt.Println("从数据库服务器接收响应错误:", err)
		conn.Write([]byte("无法处理请求"))
		return
	}

	// 将数据库服务器的响应发送给客户端
	conn.Write(databaseResponse[:n])
}
