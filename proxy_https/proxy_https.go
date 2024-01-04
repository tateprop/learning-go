package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"gio.test/proxy_https/proxyutil"
	"github.com/AdguardTeam/gomitmproxy"
)

func main() {
	proxy := gomitmproxy.NewProxy(gomitmproxy.Config{
		ListenAddr: &net.TCPAddr{
			IP:   net.IPv4(0, 0, 0, 0),
			Port: 8080,
		},
		OnRequest: func(session *gomitmproxy.Session) (request *http.Request, response *http.Response) {
			req := session.Request()

			log.Printf("onRequest: %s %s", req.Method, req.URL.String())

			if req.URL.Host == "example.net" {
				body := strings.NewReader("<html><body><h1>Replaced response</h1></body></html>")
				res := proxyutil.NewResponse(http.StatusOK, body, req)
				res.Header.Set("Content-Type", "text/html")

				// Use session props to pass the information about request being blocked
				session.SetProp("blocked", true)
				return nil, res
			}

			return nil, nil
		},
		OnResponse: func(session *gomitmproxy.Session) *http.Response {
			log.Printf("onResponse: %s", session.Request().URL.String())

			if _, ok := session.GetProp("blocked"); ok {
				log.Printf("onResponse: was blocked")
			}

			res := session.Response()
			req := session.Request()

			if strings.Index(res.Header.Get("Content-Type"), "text/html") != 0 {
				// Do nothing with non-HTML responses
				return nil
			}

			b, err := proxyutil.ReadDecompressedBody(res)
			// Close the original body
			_ = res.Body.Close()
			if err != nil {
				return proxyutil.NewErrorResponse(req, err)
			}

			// Use latin1 before modifying the body
			// Using this 1-byte encoding will let us preserve all original characters
			// regardless of what exactly is the encoding
			body, err := proxyutil.DecodeLatin1(bytes.NewReader(b))
			if err != nil {
				return proxyutil.NewErrorResponse(session.Request(), err)
			}

			// Modifying the original body
			modifiedBody, err := proxyutil.EncodeLatin1(body + "<!-- EDITED -->")
			if err != nil {
				return proxyutil.NewErrorResponse(session.Request(), err)
			}

			res.Body = ioutil.NopCloser(bytes.NewReader(modifiedBody))
			res.Header.Del("Content-Encoding")
			res.ContentLength = int64(len(modifiedBody))
			return res
		},
	})
	err := proxy.Start()
	if err != nil {
		log.Fatal(err)
	}

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)
	<-signalChannel

	// Clean up
	proxy.Close()
}
