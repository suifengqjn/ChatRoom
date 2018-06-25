package server

import (
	"net"
	json2 "encoding/json"
	"fmt"
)

type User struct {
	Username string
	Password string
	UserId int
}

type Client struct {
	C chan string  //用户发送数据的管道
	ChatroomC chan string //聊天室管道
	Address string  //ip +port
	IsInChatRoom bool // 是否在聊天室中
	User
}

func DealLogin(conn net.Conn, cli Client, lists map[string]Client, onlineMap map[string]Client) bool  {

	var code int
	var mes string
	value ,ok := lists[cli.Username]
	fmt.Println("login===",lists, cli, ok)
	if ok {
		if value.Password == cli.Password {
			code = 1
			mes = "登入成功"


			onlineMap[cli.Username] = cli




		} else {
			code = 0
			mes = "密码错误"
		}
	} else {
		code = 0
		mes = "用户不存在"
	}

	var mp = map[string]interface{}{"code":code,"message":mes}
	json ,err := json2.Marshal(mp)
	fmt.Println("dd", json, err)
	if err == nil {
		conn.Write(json)
	} else {
		conn.Close()
	}

	return code == 1

}

func DealResgister(conn net.Conn, cli Client, lists map[string]Client, onlineMap map[string]Client) bool {

	var code int
	var mes string
	if len(cli.Username) > 3 && len(cli.Password) > 3 {

		//判断用户名是否存在

		_, ok := lists[cli.Username]
		if ok {
			code = 0
			mes = "用户名已存在"
		} else {
			var username = cli.Username

			lists[username] = cli
			onlineMap[username] = cli
			code = 1
			mes = "注册成功"
		}
	} else {
		code = 0
		mes = "注册失败，用户名和密码长度必须大于3"
	}

	var mp = map[string]interface{}{"code":code, "message":mes}
	json ,err := json2.Marshal(mp)
	fmt.Println(string(json))
	if err == nil {
		conn.Write(json)
	} else {
		conn.Close()
	}
	return code == 1
}