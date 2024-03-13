package main

import (
	"context"
	log "github.com/sirupsen/logrus"
	"net"
	"net/http"
	"time"
)

func noRedirect(req *http.Request, via []*http.Request) error {
	return http.ErrUseLastResponse
}

func perform_request(src_ip string) bool {
	lIP := net.ParseIP(src_ip)
	lPort := 0
	localAddr := &net.TCPAddr{IP: lIP, Port: lPort}

	keepAlive := 10 * time.Second

	dialer := &net.Dialer{
		LocalAddr: localAddr,
		Timeout:   10 * time.Second,
		KeepAlive: keepAlive,
		DualStack: true,
	}

	conn, err := dialer.Dial("tcp", net.JoinHostPort("example.com", "80"))

	if err != nil {
		log.Error(err)
		return false
	}

	time.Sleep(100 * time.Millisecond)
	defer conn.Close()

	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return conn, nil
		},
		DisableKeepAlives:  false,
		DisableCompression: true, // Disable gzip encoding
	}
	// Create an HTTP client with the custom transport
	client := &http.Client{
		Transport:     transport,
		CheckRedirect: noRedirect,
	}

	host_and_port := net.JoinHostPort("example.com", "80")
	httpUrl := "http://" + host_and_port + "/"

	// Create an HTTP GET request with custom headers
	req, err := http.NewRequest("GET", httpUrl, nil)
	if err != nil {
		log.Error("Error creating request:", err)
		return false
	}

	// Set custom headers
	req.Host = "example.com"
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko)")
	req.Close = false

	// Do the request
	// Create a context with timeout (for the request send)
	timeout_for_troute := 60 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout_for_troute)
	defer cancel()

	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		log.Error("Error ", err)
		return false
	}
	log.Error(resp)
	return true
}

func main() {
	//perform_request("38.110.46.22")
	perform_request("2001:550:9005::11")
}
