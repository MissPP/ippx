package socks5

import (
	"errors"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

type SocksClientConn struct {
	Conn       net.Conn
	CreatedAt  time.Time
	Mu         *sync.Mutex
	TargetAddr string
}

func (s *SocksClientConn) SendVersion() error {
	s.Conn.Write([]byte{5, 1})
	var a [1024]byte
	s.Conn.Read(a[:])
	if a[0] != 5 || a[1] != 2 {
		return errors.New("check proxy version fail")
	}
	return nil
}

func (s *SocksClientConn) SendAuth(user string, password string) error {
	var a [1024]byte
	userLength := len(user)
	pwdLength := len(password)
	resp := []byte{1}
	resp = append(resp, uint8(userLength))
	resp = append(resp, []byte(user)...)
	resp = append(resp, uint8(pwdLength))
	resp = append(resp, []byte(password)...)
	s.Conn.Write(resp)
	s.Conn.Read(a[:])
	if a[0] != 1 || a[1] != 0 {
		return errors.New("check proxy auth  fail")
	}
	return nil
}

func (s *SocksClientConn) SendAddr(addr string) error {
	var a [1024]byte
	addrArr := strings.Split(addr, ":")
	resp := []byte{5, 1, 1, 3}
	resp = append(resp, uint8(len(addrArr[0])))
	resp = append(resp, []byte(addrArr[0])...)
	port, _ := strconv.Atoi(addrArr[1])
	resp = append(resp, []byte{uint8(port >> 8), uint8(port & 255)}...)
	s.Conn.Write(resp)
	s.Conn.Read(a[:])
	if a[0] != 5 {
		return errors.New("proxy send addr fail")
	}
	return nil
}
