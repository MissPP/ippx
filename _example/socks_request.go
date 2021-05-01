package main

import (
	"context"
	"fmt"
	"golang.org/x/net/proxy"
	_ "io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
)

func main() {
	// create a socks5 dialer
	dialer, err := proxy.SOCKS5("tcp", "127.0.0.1:7555", &proxy.Auth{"user", "pwd"}, proxy.Direct)
	if err != nil {
		fmt.Fprintln(os.Stderr, "can't connect to the proxy:", err)
		os.Exit(1)
	}
	httpTransport := &http.Transport{}
	httpClient := &http.Client{Transport: httpTransport}
	dc := dialer.(interface {
		DialContext(ctx context.Context, network, addr string) (net.Conn, error)
	})
	httpTransport.DialContext = dc.DialContext
	if resp, err := httpClient.Get("http://yourdomain.com"); err != nil {
		log.Fatalln(err)
	} else {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("%s\n", body)
	}
}
