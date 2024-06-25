package net

import (
	"fmt"
	"log"
	"net"
	"time"
)

type client struct {
	Conn net.Conn
}

func Conn(netT, netW, addr string) (*client, error) {
	tcpAddr, err := net.ResolveTCPAddr(netT, addr)
	if err != nil {
		//panic(err)
		return nil, err
	}
	conn, err := net.DialTCP(netW, nil, tcpAddr) //链接
	if err != nil {
		//panic(err)
		return nil, err
	}
	go func() {
		for {
			beatch := make(chan []byte) // beatch没有缓冲区
			go heartBeat(conn, beatch, 15)
			//go writeBeat(conn, beatch)
			/*	fmt.Println("阻塞读取数据：")
				buf:=make([]byte,1024)
				n,_:=conn.Read(buf)//读取数据
				fmt.Println(string(buf[:n]))*/
		}
	}()
	return &client{conn}, nil
}

// HeartBeat 判断30秒内有没有产生通信
// 超过30秒退出
func heartBeat(conn net.Conn, heartchan chan []byte, timeout int) {
	fmt.Println(" HeartBeat start。。。。。")
	select { // 监听IO操作，当IO操作发生时，触发相应的动作
	case hc := <-heartchan:
		log.Println("read chan******》", string(hc))
		conn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second)) // conn会自己接收
		fmt.Println("conn deadline set finish")
	case <-time.After(time.Duration(timeout) * time.Second): // 这一段是没有用的
		fmt.Println("kkk")
		log.Println(conn.RemoteAddr(), "time out. will close conn") //客户端超时
		conn.Close()
		fmt.Println("ddd")
	}
	fmt.Println(" HeartBeat end。。。。。")
}

// Send 向服务器端发送消息
func (cl *client) Send() {

}

// Close 关闭与服务器端的连接
func (cl *client) Close() error {
	err := cl.Conn.Close()
	return err
}
