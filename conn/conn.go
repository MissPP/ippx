package conn

import (
	"net"
	"time"
)

func setReadTimeout(c net.Conn, timeout int) {
	c.SetReadDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
}

func Copy(src, dst net.Conn, timeout int) {
	//defer dst.Close()
	var buf [1024]byte
	for {
		if timeout != 0 {
			setReadTimeout(src, timeout)
		}
		n, err := src.Read(buf[:])
		if n > 0 {
			if _, err := dst.Write(buf[0:n]); err != nil {
				break
			}
		}
		if err != nil {
			break
		}
	}
	return
}

func CopyWithEncode(src, dst net.Conn, timeout int) {
	//defer dst.Close()
	var buf [1024]byte
	for {
		if timeout != 0 {
			setReadTimeout(src, timeout)
		}
		n, err := src.Read(buf[:])
		//if timeout==6{
		//	fmt.Println("dst server reading")
		//	fmt.Println(n)
		//}
		//if timeout==5{
		//	fmt.Println("client get response")
		//	fmt.Println(n)
		//}
		if n > 0 {
			writeStr := Rc4(buf[0:n])
			//if timeout == 4 {
			//	fmt.Println("client write rc4")
			//	fmt.Println(len(writeStr))
			//}
			//if timeout==7 {
			//	fmt.Println("dst server writing")
			//	fmt.Println(len(writeStr))
			//}
			if _, err := dst.Write(writeStr); err != nil {
				break
			}
		}
		if err != nil {
			break
		}
	}
	return
}
