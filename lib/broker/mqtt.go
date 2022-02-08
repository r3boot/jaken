package broker

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"jaken/lib/common"
	"log"
	"strings"
	"time"
)

/*
 * Topic(s) for sending data
 * from/irc/<channel>/<nickname>/feed: receives an unfiltered steam of PRIVMSGs
 * from/irc/<channel>/<nickname>/<command>: receives lines which are marked as a command
 *
 * Topics for receiving data
 * to/irc/<recipient>/privmsg: send a reply via PRIVMSG
 * to/irc/<recipient>/notice: send a reply via NOTICE
 */

const (
	MaxInFlight = 100

	SubmitTopicPrefix  = "from/irc"
	ReceiveTopicPrefix = "to/irc"

	privmsgTopic = ReceiveTopicPrefix + "/+/privmsg"
	noticeTopic  = ReceiveTopicPrefix + "/+/notice"
	topicTopic   = ReceiveTopicPrefix + "/+/topic"
)

type Params struct {
	Server         string
	ClientId       string
	Username       string
	Password       string
	UnfilteredChan chan common.RawMessage
	CommandChan    chan common.CommandMessage
	PrivmsgChan    chan common.FromMessage
	NoticeChan     chan common.FromMessage
	TopicChan      chan common.TopicMessage
}

type Mqtt struct {
	client         mqtt.Client
	config         *Params
	unfilteredChan chan common.RawMessage
	commandChan    chan common.CommandMessage
	privmsgChan    chan common.FromMessage
	privmsgToken   mqtt.Token
	noticeChan     chan common.FromMessage
	noticeToken    mqtt.Token
	topicChan      chan common.TopicMessage
	topicToken     mqtt.Token
}

func New(params *Params) *Mqtt {
	broker := &Mqtt{
		config:         params,
		unfilteredChan: params.UnfilteredChan,
		commandChan:    params.CommandChan,
		privmsgChan:    params.PrivmsgChan,
		noticeChan:     params.NoticeChan,
		topicChan:      params.TopicChan,
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

	broker.privmsgToken = broker.client.Subscribe(privmsgTopic, 0, broker.privmsgReceivedCallback)
	broker.privmsgToken.Wait()

	broker.noticeToken = broker.client.Subscribe(noticeTopic, 0, broker.noticeReceivedCallback)
	broker.noticeToken.Wait()

	broker.topicToken = broker.client.Subscribe(topicTopic, 0, broker.topicReceivedCallback)
	broker.topicToken.Wait()

	go broker.channelReader()

	return broker
}

func (mqtt *Mqtt) channelReader() {
	for {
		select {
		case msg := <-mqtt.unfilteredChan:
			{
				channel := strings.Replace(msg.Channel, "#", "", -1)
				topic := fmt.Sprintf("%s/%s/%s/message", SubmitTopicPrefix, channel, msg.Nickname)
				token := mqtt.client.Publish(topic, 0, false, msg.Message)
				token.Wait()
			}
		case msg := <-mqtt.commandChan:
			channel := strings.Replace(msg.Channel, "#", "", -1)
			topic := fmt.Sprintf("%s/%s/%s/%s", SubmitTopicPrefix, channel, msg.Nickname, msg.Command)
			token := mqtt.client.Publish(topic, 0, false, msg.Arguments)
			token.Wait()
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
}

func (mqtt *Mqtt) privmsgReceivedCallback(client mqtt.Client, message mqtt.Message) {
	channel := fmt.Sprintf("#%s", strings.Split(message.Topic(), "/")[2])
	mqtt.privmsgChan <- common.FromMessage{
		Recipient: channel,
		Message:   string(message.Payload()),
	}
}

func (mqtt *Mqtt) noticeReceivedCallback(client mqtt.Client, message mqtt.Message) {
	channel := fmt.Sprintf("#%s", strings.Split(message.Topic(), "/")[2])
	mqtt.noticeChan <- common.FromMessage{
		Recipient: channel,
		Message:   string(message.Payload()),
	}
}

func (mqtt *Mqtt) topicReceivedCallback(client mqtt.Client, message mqtt.Message) {
	channel := fmt.Sprintf("#%s", strings.Split(message.Topic(), "/")[2])
	mqtt.topicChan <- common.TopicMessage{
		Channel: channel,
		Topic:   string(message.Payload()),
	}
}
