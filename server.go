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

//广播消息
var message = make(chan  string)
//聊天室消息
var chatRoomMessage = make(chan string)

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
//聊天室消息

func ChatRoomMessager()  {

	for  {
		chatRoomMsg := <- chatRoomMessage

		for _,client := range onlineMap  {
			fmt.Println(">>>>>>>>>>>>>>>>>>")
			if client.IsInChatRoom == true {
				client.ChatroomC <- chatRoomMsg
			}

		}
	}

}


func writeMsgToClient(cli server.Client, conn net.Conn)  {

	for msg := range cli.C {
		conn.Write([]byte(msg))
	}


}

func writeChatroomMes(cli server.Client, conn net.Conn)  {
	for roomMsg := range  cli.ChatroomC {
		conn.Write([]byte(roomMsg))
	}
}

//处理用户连接
func handleConn(conn net.Conn)  {

	defer conn.Close()
	var cli = server.Client{}
	cli.C = make(chan string)
	cli.ChatroomC = make(chan  string)
	addr := conn.RemoteAddr().String()
	cli.Address = addr

	// 新开协程
	go writeMsgToClient(cli, conn)

	go writeChatroomMes(cli, conn)


	//用户是否主动退出
	isQuit := make(chan bool)

	//超时处理
	hasData := make(chan  bool)

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

				success :=server.DealLogin(conn, cli, registerMap,onlineMap)
				//广播某个人上线
				if success {
					message <- makeMsg(cli, "login")
				}


			case "register":
				cli.Username = Map["username"]
				cli.Password = Map["password"]
				success := server.DealResgister(conn, cli,registerMap, onlineMap)
				fmt.Println("当前用户列表：", registerMap)
				//广播某个人上线
				if success {
					message <- makeMsg(cli, "register")
				}
			case "message":
				var mes = Map["message"]
				if cli.IsInChatRoom == false {
					server.SendMessageToPerson(mes,onlineMap, conn, cli)
				} else {
					server.SendMessageToChatRoom(mes, cli.Username, chatRoomMessage)

				}
				//message <- makeMsg(cli, mes)

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
			case "chatRoom":
				cli.IsInChatRoom = !cli.IsInChatRoom
				onlineMap[cli.Username] = cli
				if cli.IsInChatRoom == true {
					conn.Write([]byte("你已经进入聊天室\n"))
					server.SendMessageToChatRoom("进入聊天室", cli.Username, chatRoomMessage)
				} else {
					conn.Write([]byte("你已经离开聊天室\n"))
					server.SendMessageToChatRoom("离开聊天室", cli.Username, chatRoomMessage)
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


	// 第三个go 转发聊天室消息
	go ChatRoomMessager()

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
