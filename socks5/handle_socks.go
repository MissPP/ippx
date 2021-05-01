package socks5

import (
	"errors"
	"fmt"
	"ippx/conn"
	"log"
	"net"
	"runtime/debug"
	"time"
)

func (s *SocksServer) HandleSocks(IsEncrypted bool) {
	s.Mu.Lock()
	defer func() {
		s.Mu.Unlock()
		if err := recover(); err != nil {
			log.Println(err)
			debug.PrintStack()
		}
	}()

	if s.Conn == nil {
		return
	}

	if err := s.checkVersion(); err != nil {
		return
	}
	if err := s.checkAuth(); err != nil {
		return
	}
	hostStr, portStr, err := s.getHostPort()
	if err != nil {
		return
	}
	addr := hostStr + ":" + portStr

	server, err := Dial("tcp", addr)
	if err != nil {
		log.Println("dial dst server fail:" + err.Error())
		return
	}
	//debug helpful
	if IsEncrypted {
		go conn.CopyWithEncode(s.Conn, server, 6)
		conn.CopyWithEncode(server, s.Conn, 7)
	} else {
		go conn.Copy(s.Conn, server, 60)
		conn.Copy(server, s.Conn, 60)
	}
	//log.Println("server tcp connection:", server.LocalAddr().String(), "->", server.RemoteAddr().String())
	//for {
	//	var tcpReqInfo [10240]byte
	//	reqLength, err := s.Conn.Read(tcpReqInfo[:])
	//	err = s.dialSever(server, tcpReqInfo[:], reqLength)
	//	if err != nil {
	//		//retry once
	//		server, err = reDial(addr)
	//		err = s.dialSever(server, tcpReqInfo[:], reqLength)
	//		if err != nil {
	//			log.Println(`return`)
	//			break
	//		} else {
	//			break
	//		}
	//	}
	//}
	return
}

func (s *SocksServer) dialSever(server net.Conn, tcpInfo []byte, n int) error {
	_, err := server.Write(tcpInfo[:n])
	if err != nil {
		return err
	}
	log.Println(fmt.Sprintf("dial id : %v", s.ClientId))
	for {
		var buf [1024]byte
		n, err := server.Read(buf[:])
		if n > 0 {
			if _, err := s.Conn.Write(buf[0:n]); err != nil {
				break
			}
		}
		if err != nil {
			log.Println(fmt.Sprintf("errrrr%v", err))
			break
		}
	}
	return nil
}

func reDial(addr string) (net.Conn, error) {
	server, err := Dial("tcp", addr)
	if err != nil {
		log.Println(err)
		return nil, errors.New("redial fail")
	}
	server.SetDeadline(time.Now().Add(50 * time.Second))
	return server, nil
}
