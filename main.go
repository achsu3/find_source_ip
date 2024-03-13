package main

import (
	"context"
	log "github.com/sirupsen/logrus"
	"net"
	"net/http"
	"time"
	"flag"
	"fmt"
	"strings"
)

var (
	flagIP = flag.String("ip", "6", "run over IPv4 or IPv6: 4 or 6, default=6")
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
	// perform_request("2001:550:9005::11")
	// loop through each IP on each interface
	flag.Parse()
	log.Println("IP Version: ", *flagIP)

    interfaces, err := net.Interfaces()
    if err != nil {
        fmt.Println("Error:", err)
    }

	local_addr_ip := ""
    for _, iface := range interfaces {
        // Get addresses associated with the interface
        addrs, err := iface.Addrs()
        if err != nil {
            log.Println("Error:", err)
            continue
        }

        // Print the interface name
        fmt.Printf("Interface: %s\n", iface.Name)
	
		// fmt.Printf("Found Interface: %s\n", iface.Name)
		// Print each address associated with the interface
		for _, a := range addrs {
			if ipnet, ok := a.(*net.IPNet); ok {
				
				if *flagIP == "4" && strings.Contains(ipnet.IP.String(), ".") {
					local_addr_ip = ipnet.IP.String()

					if perform_request(local_addr_ip){
						fmt.Println(local_addr_ip)
						return
					}

				
				}
				if *flagIP == "6"  && strings.Contains(ipnet.IP.String(), ":") {
					local_addr_ip = ipnet.IP.String()
					if perform_request(local_addr_ip){
						fmt.Println(local_addr_ip)
						return
					}
				}
			}
		}
	
    }

	fmt.Println("0")
	return 
}
