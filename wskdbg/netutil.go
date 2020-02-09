package main

import "net"

func myIPAddress() string {
	ifaces, err := net.Interfaces()
	FatalIf(err)
	// handle err
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		FatalIf(err)
		// handle err
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					return ipnet.IP.String()
				}
			}
		}
	}
	return "127.0.0.1"
}

func doPost(url string, data []byte) string {

	log.Printf(">>> %s\n=== %s\n", url, data)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")

	transport := http.Transport{}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("!!! Connection Error: %s -> %s\n", url, err)
		return ""
	}
	if resp.StatusCode != 200 {
		log.Printf("!!! Http Error: %s -> %s\n", url, resp.Status)
	}
	body, _ := ioutil.ReadAll(resp.Body)

	resp.Body.Close()
	transport.CloseIdleConnections()
	return string(body)
}
