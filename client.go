package main

import (
	"net"
	"goExt/goExt"
	"fmt"
	"os"
	"ketang/netWork/0604_Socket/Tool"
	"encoding/json"
	"goDemo/ChatRoom/config"
)

//const (
//	systemMes  = 1
//	singleChat = 2
//	roomChat   = 3
//)

var chatType int

type Mess struct {
	Code    int    `code`
	Message string `message`
	Type    int    `type` // 消息类型  1：系统消息 2： 私聊信息 3：聊天室消息
}

func main() {

	conn, err := net.Dial("tcp", Tool.GetLocalIp()+":"+config.NetPort)
	if goExt.CheckErr(err) {
		return
	}
	var ty string

start:
	for {
		fmt.Println("1:登入  2:注册")
		fmt.Scanln(&ty)
		if ty == "1" {

			v, msg := login(conn)
			fmt.Println(msg)
			if v {
				break
			} else {
				goto start
			}

		} else if ty == "2" {
			v, msg := register(conn)
			fmt.Println(msg)
			if v {
				break
			} else {
				goto start
			}
		} else {
			fmt.Println("请重新输入")
		}
	}

	defer conn.Close()
	//接受服务器数据
	go func() {

		buf := make([]byte, 1024)
		for {
			n, err := conn.Read(buf)
			if err != nil {
				fmt.Println("read server data err:", err.Error())
				return
			} else {
				fmt.Println(string(buf[:n]))
			}
		}
	}()

	//send message
	str := make([]byte, 1024)
	fmt.Println("输入 '/' 进入聊天室和退出聊天室")
	fmt.Println("输入 'userLists' 查看在线人列表")
	fmt.Println("输入 'allUsers' 查看所有注册用户")
	fmt.Println("私聊格式 'userName:消息'")
	for {
		n, err := os.Stdin.Read(str)
		if err != nil {
			fmt.Println("input err")
			return
		}

		var m = make(map[string]string)

		mes := str[:n]
		if string(mes) == "userLists\n" {
			m["act"] = "userLists"
		} else if string(mes) == "allUsers\n" {
			m["act"] = "allUsers"
		} else if string(mes) =="/\n" {
			m["act"]= "chatRoom"
		} else {
			m["act"] = "message"
			m["message"] = string(mes)
		}
		byteMes, _ := json.Marshal(m)
		_, err = conn.Write(byteMes)
		if err != nil {
			fmt.Println("send message err:", err.Error())
		}
	}

}

// 登入并且接受服务器返回的结果
func login(conn net.Conn) (bool, string) {

	var m = make(map[string]string)
	m["act"] = "login"
	var username string
	var password string

	fmt.Println("输入用户名：")
	fmt.Scanln(&username)

	fmt.Println("输入密码")
	fmt.Scanln(&password)

	m["username"] = username
	m["password"] = password

	str, err := json.Marshal(m)
	if err != nil {

		return false, err.Error()
	}

	conn.Write(str)

	buf := make([]byte, 1024)
	for {

		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read server data err:", err.Error())
			return false, err.Error()
		} else {
			M := Mess{}
			fmt.Println(string(buf[:n]))
			err = json.Unmarshal(buf[:n], &M)

			res := false
			if M.Code == 1 {
				res = true
			}
			return res, M.Message

		}
	}

}

func register(conn net.Conn) (bool, string) {
	var m = make(map[string]string)
	m["act"] = "register"
	var username string
	var password string

	fmt.Println("输入用户名：")
	fmt.Scanln(&username)

	fmt.Println("输入密码")
	fmt.Scanln(&password)

	m["username"] = username
	m["password"] = password

	str, err := json.Marshal(m)
	if err != nil {

		return false, err.Error()
	}

	conn.Write(str)

	buf := make([]byte, 1024)
	for {

		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read server data err:", err.Error())
			return false, err.Error()
		} else {
			M := Mess{}
			fmt.Println(string(buf[:n]))
			err = json.Unmarshal(buf[:n], &M)
			res := false
			if M.Code == 1 {
				res = true
			}
			return res, M.Message
		}
	}

}
