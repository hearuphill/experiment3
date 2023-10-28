package main

import (
	"fmt"
	"net"
	"strings"
)

type Contact struct {
	Name    string
	Address string
	Phone   string
}

var contacts []Contact

func main() {
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
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("读取数据错误:", err)
		return
	}

	clientData := string(buf)
	fmt.Println("客户端请求:", clientData)

	// 根据客户端请求执行不同操作
	switch {
	case strings.HasPrefix(clientData, "查看联系人信息"):
		viewContacts(conn)
	case strings.HasPrefix(clientData, "添加新联系人"):
		addContact(clientData)
		conn.Write([]byte("联系人添加成功！"))
	case strings.HasPrefix(clientData, "修改联系人信息"):
		editContact(clientData)
		conn.Write([]byte("联系人信息修改成功！"))
	case strings.HasPrefix(clientData, "删除联系人"):
		deleteContact(clientData)
		conn.Write([]byte("联系人已删除"))
	default:
		conn.Write([]byte("无效的请求"))
	}
}

func viewContacts(conn net.Conn) {
	if len(contacts) == 0 {
		conn.Write([]byte("通讯录为空"))
		return
	}

	response := "联系人信息：\n"
	for i, contact := range contacts {
		response += fmt.Sprintf("%d. 姓名：%s, 地址：%s, 电话：%s\n", i+1, contact.Name, contact.Address, contact.Phone)
	}

	conn.Write([]byte(response))
}

func addContact(clientData string) {
	// 解析客户端请求
	parts := strings.Split(clientData, ":")
	if len(parts) < 4 {
		return
	}
	name, address, phone := parts[1], parts[2], parts[3]
	contact := Contact{Name: name, Address: address, Phone: phone}
	contacts = append(contacts, contact)
}

func editContact(clientData string) {
	// 解析客户端请求
	parts := strings.Split(clientData, ":")
	if len(parts) < 5 {
		return
	}
	index := atoi(parts[1]) - 1
	name, address, phone := parts[2], parts[3], parts[4]

	if index < 0 || index >= len(contacts) {
		return
	}

	contacts[index] = Contact{Name: name, Address: address, Phone: phone}
}

func deleteContact(clientData string) {
	// 解析客户端请求
	parts := strings.Split(clientData, ":")
	if len(parts) < 2 {
		return
	}
	index := atoi(parts[1]) - 1

	if index < 0 || index >= len(contacts) {
		return
	}

	contacts = append(contacts[:index], contacts[index+1:]...)
}

func atoi(s string) int {
	n := 0
	for _, c := range s {
		if c < '0' || c > '9' {
			return -1
		}
		n = n*10 + int(c-'0')
	}
	return n
}
