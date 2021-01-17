package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

func help(){
	fmt.Println(`usage: "-listen port1 port2"`)
	fmt.Println(`       "-tran port1 ip:port2"`)
	fmt.Println(`       "-slave ip1:port1 ip2:port2"`)
}

func main(){
	if len(os.Args) < 3{
		help()
		return
	}
	switch os.Args[1]{
	case "-listen":
		p2p(os.Args[2],os.Args[3])
	case "-tran":
		host := strings.Split(os.Args[3],":")[0]
		port := strings.Split(os.Args[3],":")[1]
		p2h(os.Args[2],lookHostname(host)+":"+port)
	case "-slave":
		host1 := strings.Split(os.Args[3],":")[0]
		port1 := strings.Split(os.Args[3],":")[1]
		host2 := strings.Split(os.Args[4],":")[0]
		port2 := strings.Split(os.Args[4],":")[1]
		h2h(lookHostname(host1)+":"+port1,lookHostname(host2)+":"+port2)
	}
}

func p2p(p1,p2 string){
	listener1,_ := net.Listen("tcp", "0.0.0.0:"+p1)
	listener2,_ := net.Listen("tcp", "0.0.0.0:"+p2)
	for {
		conn1 := accept(listener1)
		conn2 := accept(listener2)
		forward(conn1, conn2)
	}
}

func p2h(p,h string){
	listener,_ := net.Listen("tcp", "0.0.0.0:"+p)
	for{
		conn1 := accept(listener)
		if conn1 == nil{
			continue
		}
		go func(host string) {
			conn2,err := net.Dial("tcp", h)
			if err != nil {
				conn1.Close()
				return
			}
			forward(conn1,conn2)
		}(h)
	}
}

func h2h(h1,h2 string){
	var conn1,conn2 net.Conn
	var err error
	for{
		for{
			conn1,err = net.Dial("tcp", h1)
			if err != nil{
				time.Sleep(1 * time.Second)
			}else{
				break
			}
		}
		for{
			conn2,err = net.Dial("tcp", h2)
			if err != nil{
				time.Sleep(1 * time.Second)
			}else{
				break
			}
		}
		forward(conn1,conn2)
	}
}

func accept(listener net.Listener) (conn net.Conn){
	conn,_ = listener.Accept()
	return
}

func forward(conn1,conn2 net.Conn){
	var wg sync.WaitGroup
	wg.Add(2)
	go connCopy(conn1,conn2,&wg)
	go connCopy(conn2,conn1,&wg)
	wg.Wait()
}

func connCopy(conn1 net.Conn, conn2 net.Conn, wg *sync.WaitGroup){
	io.Copy(conn1,conn2)
	conn1.Close()
	wg.Done()
}


func lookHostname(hostname string) (string) {
	ip, err := net.ResolveIPAddr("ip4", hostname)
	if err != nil {
		panic(err)
	}
	fmt.Println(ip.String())
	return ip.String()
}