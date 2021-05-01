package main

import (
	"bytes"
	"fmt"
	"ippx/config"
	"ippx/conn"
	"ippx/pool"
	"log"
	"net"
	"net/url"
	"strings"
	"time"
)

var poolObj pool.Pool
var conf *config.ConfigData

func main() {
	conf = config.Config
	//get user pwd
	poolObj = pool.NewSocksPool(conf.Username, conf.Password, conf.ServerHost+":"+conf.ServerPort, conf.MaxPoolConn, time.Second*60*60*24)
	l, _ := net.Listen("tcp", ":"+conf.ClientPort)
	for {
		log.Println("accept start=================")
		customerConn, _ := l.Accept()
		go handleConn(customerConn)
	}
}

func handleConn(sourceConn net.Conn) {
	var b [1024]byte
	_, err := sourceConn.Read(b[:])
	var method, host, address string
	splitIndex := bytes.IndexByte(b[:], '\n')
	if err != nil || splitIndex == -1 {
		log.Println(err)
		return
	}
	//get tcp source info
	fmt.Sscanf(string(b[:splitIndex]), "%s%s", &method, &host)
	hostPortURL, err := url.Parse(host)
	if err != nil {
		log.Println(err)
		return
	}

	if hostPortURL.Opaque == "443" {
		address = hostPortURL.Scheme + ":443"
	} else {
		if strings.Index(hostPortURL.Host, ":") == -1 { //default port 80
			address = hostPortURL.Host + ":80"
		} else {
			address = hostPortURL.Host
		}
	}
	//get the connection from pool
	s := poolObj.Get(address)
	log.Println("address:" + address)

	if method == "CONNECT" {
		fmt.Fprint(sourceConn, "HTTP/1.1 200 Connection established\r\n\r\n")
	} else {
		s.Conn.Write(b[:])
	}

	defer func() {
		sourceConn.Close()
		if err := recover(); err != nil {
			log.Printf("panic in defer: %v", err)
		} else {
			if conf.KeepAlive {
				poolObj.Put(s)
			} else {
				poolObj.Release(s)
			}
		}
	}()
	if conf.SocksSwitch && conf.IsEncrypted {
		go conn.CopyWithEncode(sourceConn, s.Conn, 4)
		conn.CopyWithEncode(s.Conn, sourceConn, 5)
	} else {
		go conn.Copy(sourceConn, s.Conn, 60)
		conn.Copy(s.Conn, sourceConn, 60)
	}

}
