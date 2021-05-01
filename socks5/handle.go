package socks5

import (
	"bytes"
	"fmt"
	"io"
	"ippx/config"
	"log"
	"net"
	"net/url"
	"runtime/debug"
	"strings"
	"time"
)

func (s *SocksServer) Handle() {
	s.Mu.Lock()
	defer func() {
		s.Mu.Unlock()
		s.Conn.Close()
		if err := recover(); err != nil {
			log.Println(err)
			debug.PrintStack()
		}
	}()
	if s.Conn == nil {
		return
	}
	//s.Conn.SetDeadline(time.Now().Add(time.Duration(100) * time.Second))
	defer s.Conn.Close()
	//start get remote addr
	var b [1024]byte
	n, err := s.Conn.Read(b[:])

	if err != nil || bytes.IndexByte(b[:], '\n') == -1 {
		log.Println(err)
		return
	}
	var method, host, address string
	fmt.Sscanf(string(b[:bytes.IndexByte(b[:], '\n')]), "%s%s", &method, &host)
	hostPortURL, err := url.Parse(host)
	if err != nil {
		log.Println(err)
		return
	}

	if hostPortURL.Opaque == "443" {
		address = hostPortURL.Scheme + ":443"
	} else {
		if strings.Index(hostPortURL.Host, ":") == -1 {
			address = hostPortURL.Host + ":80"
		} else {
			address = hostPortURL.Host
		}
	}
	server, err := Dial("tcp", address)
	if err != nil {
		log.Println(err)
		return
	}
	defer server.Close()
	log.Println("connection:", server.LocalAddr().String(), "->", server.RemoteAddr().String())
	// server.SetDeadline(time.Now().Add(time.Duration(10) * time.Second))
	if method == "CONNECT" {
		fmt.Fprint(s.Conn, "HTTP/1.1 200 Connection established\r\n\r\n")
	} else {
		log.Println("server write", method)
		server.Write(b[:n])
	}

	go func() {
		io.Copy(server, s.Conn)
	}()
	io.Copy(s.Conn, server)
}

//Dial
func Dial(network, addr string) (net.Conn, error) {
	var d net.Dialer
	if config.Config.KeepAlive {
		d = net.Dialer{KeepAlive: config.Config.KeepAliveTimeout * time.Second}
	} else {
		d = net.Dialer{}
	}
	return d.Dial(network, addr)
	//return net.Dial(network, addr)
}
