package main

import (
	"net"
	"goExt/goExt"
	"fmt"
	"io"
	"time"
	"encoding/json"
	"goDemo/ChatRoom/server"
	"goDemo/ChatRoom/config"
)

const (
	Register = 1
	Login = 2
	Message= 3
	UserList = 4

)



//保存线上用户的数据 key 为userName 用户名不允许重复
var onlineMap map[string]server.Client

//保存注册用户的数据
var registerMap map[string]server.Client

var message = make(chan  string)

//广播消息
func Manager()  {
	onlineMap = make(map[string]server.Client)
	registerMap = make(map[string]server.Client)

	for {
		msg := <- message//没有消息的时候，这里会阻塞
		// 遍历map 给每个成员发送消息
		for _,client := range onlineMap  {
			client.C <- msg
		}

	}

}

func writeMsgToClient(cli server.Client, conn net.Conn)  {

	for msg := range cli.C {
		conn.Write([]byte(msg))
	}
}

//处理用户连接
func handleConn(conn net.Conn)  {

	defer conn.Close()
	var cli = server.Client{}
	cli.C = make(chan string)
	addr := conn.RemoteAddr().String()
	cli.Address = addr

	// 新开协程
	go writeMsgToClient(cli, conn)

	//广播某个人上线
	//message <- makeMsg(cli, "login")

	//用户是否主动退出
	isQuit := make(chan bool)

	//超时处理
	hasData := make(chan  bool)

	//提示信息
	//cli.C <- makeMsg(cli, " i am here ")

	//新开go 接受用户发送的数据
	go func(){
		buf := make([]byte, 1024)

		for {
			n, err := conn.Read(buf)

			if err == io.EOF {
				isQuit <- true
				return
			}
			if n== 0 {
				fmt.Println("server read err,", err.Error())
				return
			}

			// 服务器转发数据 给所有用户
			msg := buf[:n]
			fmt.Println("receive message form client:", string(msg))
			//fmt.Printf("%d-%v--\n", len(msg), msg)
			Map := make(map[string]string)
			err = json.Unmarshal(msg, &Map)

			var act = Map["act"]
			switch act {
			case "login":

				cli.Username = Map["username"]
				cli.Password = Map["password"]

				server.DealLogin(conn, cli, registerMap,onlineMap)

			case "register":
				cli.Username = Map["username"]
				cli.Password = Map["password"]
				server.DealResgister(conn, cli,registerMap, onlineMap)
				fmt.Println("当前用户列表：", registerMap)
			case "message":
				var mes = Map["message"]
				message <- makeMsg(cli, mes)
			case "userLists":
				conn.Write([]byte("user list:\n"))
				for _, v := range onlineMap {
					conn.Write([]byte(v.Address +"   " + v.Username+"\n"+"-----------------\n"))
				}
			case "allUsers":
				conn.Write([]byte("user list:\n"))
				for _, v := range registerMap {
					conn.Write([]byte(v.Address +"   "+ v.Username+"\n"+"-----------------\n"))
				}
			default:

			}

			hasData <- true

		}

	}()

	for {

		select {
		case <- isQuit:
			delete(onlineMap, cli.Address)
			message <- makeMsg(cli,"logout")
			return
		case <- hasData:
		case <- time.After(time.Second * config.LimitTimeout):
			delete(onlineMap, cli.Address)
			message <- makeMsg(cli,"time out,请重新登入")
			return 

		}

	}

}

func makeMsg(cli server.Client, msg string) string  {
	buf := "[" + cli.Username + "]: " + msg + "\n"
	return buf

}

func main() {

	listner,err := net.Listen("tcp",":"+ config.NetPort)
	if goExt.CheckErr(err) {
		return
	}

	defer listner.Close()

	// 第二个go 转发消息，只要有消息，给map每个成员发送消息
	go Manager()

	for {

		conn, err:=listner.Accept()
		fmt.Println("等待用户连接")
		if goExt.CheckErr(err) {
			continue
		}
		//第一个go 处理用户连接
		go handleConn(conn)

	}

}
