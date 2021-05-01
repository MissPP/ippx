package main

import (
	"ippx/config"
	s5 "ippx/socks5"
	"log"
	"math/rand"
	"net"
	"sync"
	"time"
)

var conf = *config.Config

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Printf(`panic in main: %v`, err)
		}
	}()
	l, err := net.Listen("tcp", ":"+conf.ServerPort)
	if err != nil {
		log.Panic(err)
	}

	for {
		//set socks proxy conn id
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		randId := r.Intn(10000)
		tcpConn, err := l.Accept()
		var s = &s5.SocksServer{tcpConn, randId, &sync.Mutex{}}
		if err != nil {
			log.Panic(err)
		}
		if conf.SocksSwitch {
			//socks
			go s.HandleSocks(conf.IsEncrypted)
		} else {
			//tcp
			go s.Handle()
		}
	}
}
