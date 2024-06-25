package net

import (
	"fmt"
	"log"
	"net"
	"time"
)

// HeartBeat 判断30秒内有没有产生通信
// 超过30秒退出
func HeartBeat(conn net.Conn, heartchan chan byte, timeout int) {
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

// HeartChanHandler 处理心跳的channel
func HeartChanHandler(n []byte, beatCh chan byte) {
	fmt.Println(" HeartChanHandler", len(n))
	for _, v := range n {
		fmt.Println("put *******> chan", string(v))
		beatCh <- v
	}
	close(beatCh) //关闭管道
	//fmt.Println(" HeartChanHandler end, close chan")
}

func MsgHandler(conn net.Conn) {
	//flag := true
	buf := make([]byte, 1024)
	defer conn.Close()
	for {
		fmt.Println("阻塞等待数据读取：")
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("when read conn, conn closed")
			//flag = false
			break
		}

		msg := buf[:n]
		/*
		 * do something
		 */
		//conn.Write([]byte("收到数据:"+string(buf[1:n])+"\n"))
		fmt.Println("收到数据: ", string(buf))
		beatch := make(chan byte) // beatch没有缓冲区
		go HeartBeat(conn, beatch, 30)
		go HeartChanHandler(msg[:1], beatch) // beatch中的数据必须要及时读取完，否则会一直阻塞。导致该协程不能推出， 因此只msg的第一个数组作为心跳
	}
	fmt.Println("当前处理数据的协程结束。。。。")
}

// Start 启动服务端
func Start(netw, addr string) {
	listener, err := net.Listen(netw, addr)
	if err != nil {
		panic(err) //处理错误
	}
	defer listener.Close() //延迟关闭
	for {
		newConn, err := listener.Accept() //接收消息
		if err != nil {
			panic(err) //处理错误
		}
		go MsgHandler(newConn) //处理客户端消息
	}
}
