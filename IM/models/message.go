package models

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"gopkg.in/fatih/set.v0"
	"net"
	"net/http"
	"strconv"
	"sync"
)

type Message struct {
	Id       int    `gorm:"primary_key;AUTO_INCREMENT"`
	FromId   int64  //发送者
	TargetId int64  //接收者
	Type     int    //消息类型 私聊 群聊 公告
	Media    int    //消息类型 文字 图片 视频
	Content  string //消息内容
	Pic      string //图片相关
	Url      string //url相关
	Desc     string //描述
	Amount   int    //数字统计
}

func (table *Message) TableName() string {
	return "message"
}

type Node struct {
	Conn      *websocket.Conn
	DataQueue chan []byte
	GroupSets set.Interface
}

// 映射关系
var clientMap map[int64]*Node = make(map[int64]*Node, 0)

// 读写锁
var rwLocker sync.RWMutex

func Chat(writer http.ResponseWriter, request *http.Request) {
	//1.检验token等合法性
	query := request.URL.Query()
	Id := query.Get("userId")
	userId, _ := strconv.ParseInt(Id, 10, 64)
	//targetId := query.Get("targetId")
	//context := query.Get("context")
	//msgType := query.Get("type")
	isvalid := true //应该check的,先这样
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return isvalid
		},
	}).Upgrade(writer, request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	//2.获取链接
	node := &Node{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),
		GroupSets: set.New(set.ThreadSafe),
	}
	//3.用户关系

	//4.userid跟node绑定并加锁
	rwLocker.Lock()
	clientMap[userId] = node
	rwLocker.Unlock()

	//5.完成发送逻辑
	go sendProc(node)
	//6.完成接受逻辑
	go recvProc(node)
	sendMsg(userId, []byte("welcome to Chat system"))

}
func sendProc(node *Node) {
	for {
		select {
		case data := <-node.DataQueue:
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}
func recvProc(node *Node) {
	for {

		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		broadMsg(data)
		fmt.Println("[ws]<<<<", data)

	}
}

var udpsendChan chan []byte = make(chan []byte, 1024)

func broadMsg(data []byte) {
	udpsendChan <- data
}

func init() {
	go udpSendProc()
	go udpRecvProc()
}

// 完成udp数据发送协程
func udpSendProc() {
	con, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(172, 25, 0, 255),
		Port: 8000,
	})
	defer con.Close()
	if err != nil {
		fmt.Println(err)
	}
	for {
		select {
		case data := <-udpsendChan:
			_, err := con.Write(data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

// 完成udp数据接受协程
func udpRecvProc() {
	con, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 8080,
	})
	if err != nil {
		fmt.Println(err)
	}
	defer con.Close()
	for {
		var buf [512]byte
		n, err := con.Read(buf[0:])
		if err != nil {
			fmt.Println(err)
			return
		}
		dispatch(buf[0:n])
	}
}

// 后端调度处理逻辑
func dispatch(data []byte) {
	msg := Message{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	switch msg.Type {
	case 1:
		sendMsg(msg.TargetId, data) //私信

	}
}

func sendMsg(userId int64, msg []byte) {
	rwLocker.RLock()
	node, ok := clientMap[userId]
	rwLocker.RUnlock()
	if ok {
		node.DataQueue <- msg
	}
}
