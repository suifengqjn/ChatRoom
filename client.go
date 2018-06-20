package main

import (
	"net"
	"goExt/goExt"
	"fmt"
	"os"
	"ketang/netWork/0604_Socket/Tool"
	"encoding/json"
)

const(
	singleChat = 1
	roomChat = 2
)

var chatType int

type Mess struct {
	code int	`code`
	mess string `message`
}


func main()  {

	conn, err :=net.Dial("tcp", Tool.GetLocalIp()+ ":12345")
	if goExt.CheckErr(err) {
		return
	}
	var ty string


	start: for {
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


		} else if ty == "2"{
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
	go func(){

		buf := make([]byte, 1024)
		for{
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
	for {
		n, err := os.Stdin.Read(str)
		if err != nil {
			fmt.Println("input err")
			return
		}

		var m = make(map[string]string)

		mes := str[:n]
		if string(mes) == "userLists" {
			 m["act"] = "userLists"
		} else if string(mes) == "allUsers" {
			m["act"] = "allUsers"
		} else {
			m["act"] = "message"
			m["message"] = string(mes)
		}
		byteMes, _ := json.Marshal(m)
		_, err = conn.Write(byteMes)
		if  err != nil{
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
	for{

		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read server data err:", err.Error())
			return false, err.Error()
		} else {
			//M := Mess{}
			//fmt.Println(string(buf[:n]))
			//err = json.Unmarshal(buf[:n], &M)
			//fmt.Println(M)
			//var mes = M.mess
			//return true, mes
			M := make(map[string]interface{})
			err = json.Unmarshal(buf[:n], &M)
			var mes = M["message"]
			res := false
			if  M["code"] == 1{
				res = true
			}
			return res, mes.(string)
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
	for{

		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read server data err:", err.Error())
			return false, err.Error()
		} else {
			M := make(map[string]interface{})
			err = json.Unmarshal(buf[:n], &M)
			var mes = M["message"]
			res := false
			if  M["code"] == 1{
				res = true
			}
			return res, mes.(string)
		}
	}

}
