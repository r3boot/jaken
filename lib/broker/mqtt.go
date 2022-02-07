package broker

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"jaken/lib/common"
	"log"
	"time"
)

/*
 * Topic(s) for sending data
 * from/ircbot/unfiltered: receives an unfiltered steam of PRIVMSGs
 * from/ircbot/commands: receives lines which are marked as a command
 *
 * Topics for receiving data
 * to/ircbot/privmsg: send a reply via PRIVMSG
 * to/ircbot/notice: send a reply via NOTICE
 */

const (
	MaxInFlight = 100
)

type Params struct {
	Server         string
	ClientId       string
	Username       string
	Password       string
	UnfilteredChan chan common.ToMessage
	CommandChan    chan common.ToMessage
}

type Mqtt struct {
	client         mqtt.Client
	config         *Params
	unfilteredChan chan common.ToMessage
	commandChan    chan common.ToMessage
	privmsgChan    chan common.FromMessage
	notifyChan     chan common.FromMessage
}

func New(params *Params) *Mqtt {
	broker := &Mqtt{
		config:         params,
		unfilteredChan: params.UnfilteredChan,
		commandChan:    params.CommandChan,
		privmsgChan:    make(chan common.FromMessage, MaxInFlight),
		notifyChan:     make(chan common.FromMessage, MaxInFlight),
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s", params.Server))
	opts.SetClientID(params.ClientId)

	if params.Username != "" {
		opts.SetUsername(params.Username)
	}

	if params.Password != "" {
		opts.SetPassword(params.Password)
	}

	broker.client = mqtt.NewClient(opts)
	token := broker.client.Connect()
	for !token.WaitTimeout(3 * time.Second) {
	}
	if err := token.Error(); err != nil {
		log.Fatal(err)
	}

	go broker.channelReader()

	return broker
}

func (mqtt *Mqtt) channelReader() {
	for {
		select {
		case msg := <-mqtt.unfilteredChan:
			fmt.Printf("unfiltered: %v\n", msg)
		case msg := <-mqtt.commandChan:
			fmt.Printf("command: %v\n", msg)
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
}
