package server

import (
	"strings"
	"net"
	"fmt"
)

func SendMessageToPerson(message string, onLineLists map[string]Client, conn net.Conn, fromcli Client)  {

	if strings.Contains(message, ":") == false{
		conn.Write([]byte("格式不正确"))
		return
	}
	username := strings.Split(message,":")[0]
	if len(username) == 0 {
		conn.Write([]byte("格式不正确"))
		return
	}
	mes := strings.Split(message,":")[1]
	Tocli := onLineLists[username]

	sendMes := fmt.Sprintf("%v发来消息：%v", fromcli.Username, mes)

	Tocli.C <- sendMes

}

func SendMessageToChatRoom(message string, from string, ChatRoomC chan string)  {

	sendMes := fmt.Sprintf("聊天室：%v发来消息：%v", from, message)
	ChatRoomC <- sendMes

	fmt.Println("-----聊天室消息")
}
