package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/eclipse/paho.mqtt.golang"
	"crypto/md5"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

const (
	mqttServer = "tcp://clouddata.usr.cn:1883"
	userName = ""
	password = ""
	deviceId = ""
	topicPrefix = "$USR/DevTx/"
	clientIdPrefix = "APP:"
)

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %X\n", msg.Payload())
	temp := 0.0
	humd := 0.0
	temp += float64(msg.Payload()[0])
	temp += float64(msg.Payload()[1])/100
	humd += float64(msg.Payload()[2])
	humd += float64(msg.Payload()[3])/100
	fmt.Printf("temp: %v, humd: %v\n",temp,humd)
}

func main() {
	//mqtt.DEBUG = log.New(os.Stdout, "", 0)
	mqtt.ERROR = log.New(os.Stdout, "", 0)
	opts := mqtt.NewClientOptions().AddBroker(mqttServer).SetClientID(clientIdPrefix+userName)
	opts.SetUsername(userName)

	hashPassword := md5.Sum([]byte(password))
	md5Password := hexutil.Encode(hashPassword[:])
	fmt.Println(md5Password[2:])

	opts.SetPassword(md5Password[2:])
	opts.SetKeepAlive(2 * time.Second)
	opts.SetDefaultPublishHandler(f)
	opts.SetPingTimeout(1 * time.Second)

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	if token := c.Subscribe(topicPrefix+deviceId, 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	time.Sleep(600 * time.Second)

	if token := c.Unsubscribe(topicPrefix+deviceId); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	c.Disconnect(250)

	time.Sleep(1 * time.Second)
}
