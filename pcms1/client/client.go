package main

import (
	"fmt"
	"net"
	"os"
)

func loop() {
	conn, err := net.Dial("tcp", "localhost:9999")
	if err != nil {
		fmt.Println("连接服务器失败:", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("")
	fmt.Println("[个人通讯录系统]")
	fmt.Println("1. 查看联系人信息")
	fmt.Println("2. 添加新联系人")
	fmt.Println("3. 修改联系人信息")
	fmt.Println("4. 删除联系人")
	fmt.Println("5. 退出")

	fmt.Print("请选择操作：")
	var choice int
	_, err = fmt.Scanln(&choice)

	if err != nil {
		fmt.Println("输入无效，请重新输入")
		fmt.Printf("Error: %s", err)
		return
	}

	switch choice {
	case 1:
		conn.Write([]byte("查看联系人信息"))
		readResponse(conn)
	case 2:
		addContact(conn)
		readResponse(conn)
	case 3:
		editContact(conn)
		readResponse(conn)
	case 4:
		deleteContact(conn)
		readResponse(conn)
	case 5:
		fmt.Println("谢谢使用，再见！")
		os.Exit(0)
	default:
		fmt.Println("无效的选择，请重新输入")
	}

}

func main() {
	for {
		loop()
	}
}

func readResponse(conn net.Conn) {
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("读取服务器响应错误:", err)
		return
	}
	fmt.Println("服务器响应:", string(buf[:n]))
}

func addContact(conn net.Conn) {
	var name, address, phone string
	fmt.Print("请输入联系人姓名：")
	fmt.Scanln(&name)
	fmt.Print("请输入联系人地址：")
	fmt.Scanln(&address)
	fmt.Print("请输入联系人电话：")
	fmt.Scanln(&phone)
	data := fmt.Sprintf("添加新联系人:%s:%s:%s", name, address, phone)
	conn.Write([]byte(data))
}

func editContact(conn net.Conn) {
	var index int
	fmt.Print("请输入要修改的联系人编号：")
	_, err := fmt.Scanln(&index)
	if err != nil {
		fmt.Println("无效的编号")
		return
	}

	if index < 1 {
		fmt.Println("无效的编号")
		return
	}

	var name, address, phone string
	fmt.Print("请输入新的联系人姓名：")
	fmt.Scanln(&name)
	fmt.Print("请输入新的联系人地址：")
	fmt.Scanln(&address)
	fmt.Print("请输入新的联系人电话：")
	fmt.Scanln(&phone)

	conn.Write([]byte(fmt.Sprintf("修改联系人信息:%d:%s:%s:%s", index, name, address, phone)))
}

func deleteContact(conn net.Conn) {
	var index int
	fmt.Print("请输入要删除的联系人编号：")
	_, err := fmt.Scanln(&index)
	if err != nil {
		fmt.Println("无效的编号")
		return
	}

	if index < 1 {
		fmt.Println("无效的编号")
		return
	}

	conn.Write([]byte(fmt.Sprintf("删除联系人:%d", index)))
}
