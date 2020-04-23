package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"./dl"
	"github.com/lucas-clemente/quic-go"
	"github.com/lucas-clemente/quic-go/h2quic"
)

var (
	url       string
	output    string
	chanSize  int
	protocol  string
	multipath bool
)

func init() {
	flag.StringVar(&url, "u", "", "M3U8 URL, required")
	flag.IntVar(&chanSize, "c", 25, "Maximum number of occurrences")
	flag.StringVar(&output, "o", "", "Output folder, required")
	flag.StringVar(&protocol, "p", "quic", "The protocol to be used")
	flag.BoolVar(&multipath, "m", false, "Enable multipath use for QUIC")

}

func main() {
	flag.Parse()
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("[error]", r)
			os.Exit(-1)
		}
	}()
	if url == "" {
		panicParameter("u")
	}
	if output == "" {
		panicParameter("o")
	}
	if chanSize <= 0 {
		panic("parameter 'c' must be greater than 0")
	}
	c := http.Client{
		Timeout: time.Duration(60) * time.Second,
	}
	if protocol == "quic" {
		fmt.Println("Using QUIC")
		quicConfig := &quic.Config{CreatePaths: multipath}
		tlsConfig := &tls.Config{InsecureSkipVerify: true}
		c.Transport = &h2quic.RoundTripper{QuicConfig: quicConfig, TLSClientConfig: tlsConfig}
	}

	downloader, err := dl.NewTask(output, url, c)
	if err != nil {
		panic(err)
	}
	if err := downloader.Start(chanSize); err != nil {
		panic(err)
	}
	fmt.Println("Done!")
}

func panicParameter(name string) {
	panic("parameter '" + name + "' is required")
}
