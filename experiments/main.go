package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/nanopack/mist/clients"
	"github.com/nanopack/mist/core"
)

func GetPublicIP() (string, error) {
	// we are using a pulib IP API, we're using ipify here, below are some others
	// https://www.ipify.org
	// http://myexternalip.com
	// http://api.ident.me
	// http://whatismyipaddress.com/api
	url := "https://api.ipify.org?format=text"
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s", ip), nil
}

func Do(msg mist.Message) {
	switch msg.Command {
	case "publish":
		log.Println(msg.Data)
	}
}

func main() {

	myIp, err := GetPublicIP()
	if err != nil {
		panic(err)
	}

	// client, err := clients.New("127.0.0.1:1445", "")
	client, err := clients.New("207.154.225.107:1445", "")
	if err != nil {
		panic(err)
	}

	// example commands (not handling errors for brevity)
	client.Ping()
	client.Subscribe([]string{"login"})
	client.Publish([]string{"login"}, myIp)
	client.List()
	// client.Unsubscribe([]string{"hello"})

	// do stuff with messages
	for {
		select {
		case msg := <-client.Messages():
			Do(msg)
		case <-time.After(time.Second * 1):
			// do something if messages are taking too long
		}
	}

	// do stuff with messages (alternate)
	// for msg := range client.Messages() {
	// 	fmt.Printf("MSG: %#v\n", msg)
	// }
}
