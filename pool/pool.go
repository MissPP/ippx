package pool

import (
	"errors"
	"fmt"
	"ippx/config"
	s5 "ippx/socks5"
	"log"
	"net"
	"strconv"
	"sync"
	"time"
)

var (
	poolType = "tcp"
	timeOut  = 60 * 60
)

type Pool interface {
	Get(string) *s5.SocksClientConn
	NewConn(string) *s5.SocksClientConn
	Put(*s5.SocksClientConn) (int, error)
	Release(*s5.SocksClientConn)
	CheckExpire()
	Close()
}

type SocksPool struct {
	IdleConn    chan *s5.SocksClientConn
	MaxLimit    int
	IdleCount   int
	BusingCount int
	TimeLimit   time.Duration
	Source      map[string]string
	Mu          *sync.Mutex
	Closed      bool
}

func NewSocksPool(user string, password string, proxyAddr string, maxLimit int, expire time.Duration) Pool {
	var p Pool = &SocksPool{
		IdleConn:    make(chan *s5.SocksClientConn, maxLimit),
		MaxLimit:    maxLimit,
		IdleCount:   0,
		BusingCount: 0,
		TimeLimit:   time.Duration(expire),
		Source:      map[string]string{`username`: user, `password`: password, `proxyAddr`: proxyAddr},
		Mu:          &sync.Mutex{},
		Closed:      false,
	}
	go func() {
		for {
			p.CheckExpire()
			time.Sleep(60 * time.Second)
		}
	}()
	return p
}

func (t *SocksPool) isConnExpire(conn *s5.SocksClientConn) bool {
	if conn.CreatedAt.Before(time.Now().Add(-t.TimeLimit * time.Second)) {
		fmt.Println("expire out of date")
		conn.Conn.Close()
		return true
	}
	return false
}

func (t *SocksPool) CheckExpire() {
	pp(t, "checkexpire")
	t.Mu.Lock()
	channelLen := len(t.IdleConn)
	tmp := make(chan *s5.SocksClientConn, t.MaxLimit)
	for i := 0; i < channelLen; i++ {
		s := <-t.IdleConn
		if s.CreatedAt.Before(time.Now().Add(-t.TimeLimit * time.Second)) {
			s.Conn.Close()
		} else {
			tmp <- s
		}
	}
	t.IdleConn = tmp
	t.IdleCount = len(tmp)
	t.Mu.Unlock()
	pp(t, "-- end checkexpire")

}

func (t *SocksPool) Close() {
	t.Mu.Lock()
	l := len(t.IdleConn)
	close(t.IdleConn)
	for i := 0; i < l; i++ {
		select {
		case a := <-t.IdleConn:
			a.Conn.Close()
		}
	}
	t.Closed = true
	t.Mu.Unlock()
}

func (t *SocksPool) NewConn(target string) *s5.SocksClientConn {
	if t.Closed {
		return nil
	}
	var d net.Dialer
	if config.Config.KeepAlive {
		d = net.Dialer{Timeout: time.Duration(timeOut) * time.Second, KeepAlive: config.Config.KeepAliveTimeout * time.Second}
	} else {
		d = net.Dialer{Timeout: time.Duration(timeOut) * time.Second}
	}
	tcpConn, err := d.Dial(poolType, t.Source["proxyAddr"])
	if err != nil {
		fmt.Println(err)
		panic("net dial fatal error")
	}
	err = tcpConn.SetDeadline(time.Now().Add(60 * 3 * time.Second))
	if err != nil {
		return nil
	}
	s := &s5.SocksClientConn{tcpConn, time.Now(), &sync.Mutex{}, target}
	err = s.SendVersion()
	if err != nil {
		return nil
	}
	err = s.SendAuth(t.Source["username"], t.Source["password"])
	if err != nil {
		return nil
	}
	err = s.SendAddr(target)
	if err != nil {
		return nil
	}
	if err != nil {
		return nil
	}
	return s
}

func (t *SocksPool) Get(target string) (res *s5.SocksClientConn) {
	if t.Closed {
		return nil
	}
	t.Mu.Lock()
	if (t.IdleCount + t.BusingCount) < t.MaxLimit {
		if (t.IdleCount + t.BusingCount) < t.MaxLimit/2 {
			res = t.NewConn(target)
		} else {
			for i := 0; i < t.IdleCount; i++ {
				v := <-t.IdleConn
				if v.TargetAddr == target {
					log.Println("get old conn")
					t.IdleCount--
					if t.isConnExpire(v) {
						res = t.NewConn(target)
					} else {
						res = v
					}
					break
				} else {
					t.IdleConn <- v
				}
			}
			res = t.NewConn(target)
		}
		t.BusingCount++
		t.Mu.Unlock()
	} else {
		fmt.Println("get need wait----")
		t.Mu.Unlock()
		select {
		case a := <-t.IdleConn:
			t.BusingCount++
			if a.TargetAddr == target {
				t.IdleCount--
				if t.isConnExpire(a) {
					res = t.NewConn(target)
				} else {
					res = a
				}
			} else {
				t.IdleCount--
				//debug
				a.Conn.Close()
				res = t.NewConn(target)
			}
		}
	}
	return res
}

func (t *SocksPool) Put(s *s5.SocksClientConn) (int, error) {
	if t.Closed {
		return 0, nil
	}
	t.Mu.Lock()
	defer t.Mu.Unlock()
	existCount := t.IdleCount + t.BusingCount
	if existCount >= t.MaxLimit {
		fmt.Println(" over limit ")
		t.BusingCount--
		s.Conn.Close()
		return 0, errors.New(" over limit ")
	}
	t.IdleConn <- s
	if t.BusingCount > 0 {
		t.BusingCount = t.BusingCount - 1
	}
	t.IdleCount++
	return existCount, nil
}

func (t *SocksPool) Release(s *s5.SocksClientConn) {
	if t.Closed {
		return
	}
	t.Mu.Lock()
	defer t.Mu.Unlock()
	s.Conn.Close()
	t.BusingCount--
}

func pp(t *SocksPool, s string) {
	log.Println("debug in function ---- " + s)
	log.Println("exist idle length --- " + strconv.Itoa(len(t.IdleConn)))
	log.Println("exist idle count --- " + strconv.Itoa(t.IdleCount))
	log.Println("exist busing count --- " + strconv.Itoa(t.BusingCount))
}
