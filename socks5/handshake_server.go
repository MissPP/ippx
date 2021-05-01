package socks5

import (
	"errors"
	"ippx/config"
	"log"
	"net"
	"strconv"
	"sync"
)

type SocksServer struct {
	Conn     net.Conn
	ClientId int
	Mu       *sync.Mutex
}

func (s *SocksServer) readInfo() ([]byte, int) {
	tcpInfo := make([]byte, 1024)
	n, err := s.Conn.Read(tcpInfo)
	if err != nil {
		log.Println("readInfo err")
		return nil, 0
	}
	return tcpInfo, n
}

func (s *SocksServer) checkVersion() (err error) {
	//check version
	tcpInfo, _ := s.readInfo()

	if tcpInfo[0] != 5 {
		if _, err = s.Conn.Write([]byte{5, 0}); err == nil {
			err = errors.New("version check fail")
		}
	}
	s.Conn.Write([]byte{5, 2})
	return
}

func (s *SocksServer) checkAuth() (err error) {
	//check auth
	conf := config.Config
	tcpInfo, _ := s.readInfo()
	usernameLen := tcpInfo[1]
	if usernameLen < 1 {
		if _, err = s.Conn.Write([]byte{5, 2}); err == nil {
			err = errors.New("auth check fail")
		}
		return
	}
	username := tcpInfo[2 : usernameLen+2]
	if string(username) != conf.Username {
		if _, err = s.Conn.Write([]byte{5, 2}); err == nil {
			err = errors.New("auth check fail")
		}
		return
	}
	passwordLen := tcpInfo[usernameLen+2]
	if passwordLen < 1 {
		if _, err = s.Conn.Write([]byte{5, 2}); err == nil {
			err = errors.New("auth check fail")
		}
		return
	}
	password := tcpInfo[usernameLen+3 : passwordLen+usernameLen+3]
	if string(password) != conf.Password {
		if _, err = s.Conn.Write([]byte{5, 2}); err == nil {
			err = errors.New("auth check fail")
		}
		return
	}
	s.Conn.Write([]byte{tcpInfo[0], 0})
	return
}

func (s *SocksServer) getHostPort() (string, string, error) {
	var host, port []byte
	var hostStr, portStr string
	var err error
	tcpInfo, infoLen := s.readInfo()
	if infoLen < 8 {
		if _, err = s.Conn.Write([]byte{5, 2}); err == nil {
			err = errors.New("unknown data")
		}
		return "", "", err
	}
	if tcpInfo[0] != 5 {
		if _, err = s.Conn.Write([]byte{5, 2}); err == nil {
			err = errors.New("unknown data")
		}
		return "", "", err
	}
	//ATYP  1 ip4  3 domin 4 ip6
	addrType := tcpInfo[3]
	host = tcpInfo[4 : infoLen-2]
	port = tcpInfo[infoLen-2 : infoLen]
	switch addrType {
	case byte(1):
		//ip4
		for k, v := range host {
			if k == 0 {
				hostStr += strconv.Itoa(int(v))
			} else {
				hostStr += "." + strconv.Itoa(int(v))
			}
		}
	case byte(3):
		hostStr = string(host[1:])
	case byte(4):
		//TODO ip6
	default:
		s.Conn.Write([]byte{5, 2})
		return "", "", errors.New("unknown data")
	}
	portStr = strconv.Itoa(int(tcpInfo[infoLen-2])<<8 | int(tcpInfo[infoLen-1]))
	resMsg := []byte{5, 0, 0}
	resMsg = append(resMsg, addrType)
	resMsg = append(resMsg, tcpInfo[4:infoLen-2]...)
	resMsg = append(resMsg, port...)
	s.Conn.Write(resMsg)
	return hostStr, portStr, nil
}
