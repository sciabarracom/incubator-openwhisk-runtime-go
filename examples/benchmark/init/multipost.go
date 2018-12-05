package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Init data
var Init = flag.String("init", "", "init data")

// Run data
var Run = flag.String("run", "", "run data")

// Repeat flag
var Repeat = flag.Int("repeat", 1, "repeat count")

// Debug flag
var Debug = flag.Bool("debug", false, "debugging")

func doPost(url string, data []byte) {
	if *Debug {
		log.Printf(">>> %s\n", url)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Connection Error: %s -> %s\n", url, err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Printf("Http Error: %s -> %s\n", url, resp.Status)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	if *Debug {
		fmt.Printf("<<< %s", string(body))
	}
}

func main() {
	flag.Parse()
	var initData []byte
	var runData []byte

	if Init != nil {
		initData, _ = ioutil.ReadFile(*Init)
	}

	if Run != nil {
		runData, _ = ioutil.ReadFile(*Run)
	}

	for _, port := range flag.Args() {
		//fmt.Println(port)
		if len(initData) > 0 {
			u := fmt.Sprintf("http://localhost:%s/init", port)
			doPost(u, initData)
		}
		if len(runData) > 0 {
			for i := 0; i < *Repeat; i++ {
				u := fmt.Sprintf("http://localhost:%s/run", port)
				doPost(u, runData)
			}
		}
	}
}
