package ircbot

import (
	"crypto/tls"
	"fmt"
	ircevent "github.com/thoj/go-ircevent"
	"jaken/lib/broker"
	"jaken/lib/common"
	"jaken/lib/ircstate"
	"regexp"
	"strings"
)

const (
	PRIVMSG = "PRIVMSG"

	cmdWhoAmI     = "whoami"
	cmdMeet       = "meet"
	cmdForget     = "forget"
	cmdTest       = "test"
	cmdHelp       = "help"
	cmdCommands   = "commands"
	cmdAddRole    = "add-role"
	cmdRemoveRole = "del-role"
	cmdListRoles  = "list-roles"
	cmdAddPerm    = "add-perm"
	cmdDelPerm    = "del-perm"
	cmdListPerms  = "list-perms"
)

type Params struct {
	Server         string
	UseTLS         bool
	VerifyTLS      bool
	Channel        string
	Nickname       string
	Realname       string
	CommandPrefix  string
	Owner          string
	UnfilteredChan chan common.RawMessage
	CommandChan    chan common.CommandMessage
	PrivmsgChan    chan common.FromMessage
	NoticeChan     chan common.FromMessage
	TopicChan      chan common.TopicMessage
}

type IrcBot struct {
	conn           *ircevent.Connection
	params         *Params
	state          *ircstate.State
	mqtt           *broker.Mqtt
	builtIn        []string
	unfilteredChan chan common.RawMessage
	commandChan    chan common.CommandMessage
	privmsgChan    chan common.FromMessage
	noticeChan     chan common.FromMessage
	topicChan      chan common.TopicMessage
}

var (
	reCommand     = regexp.MustCompile("^([a-zA-Z0-9\\+\\-]{1,64})")
	reValidParams = regexp.MustCompile("^[a-zA-Z0-9_\\-\\+\\ ]+$")
)

func New(params *Params, state *ircstate.State, mqtt *broker.Mqtt) (*IrcBot, error) {
	ircBot := &IrcBot{
		conn:           ircevent.IRC(params.Nickname, params.Realname),
		params:         params,
		state:          state,
		mqtt:           mqtt,
		unfilteredChan: params.UnfilteredChan,
		commandChan:    params.CommandChan,
		privmsgChan:    params.PrivmsgChan,
		noticeChan:     params.NoticeChan,
		topicChan:      params.TopicChan,
	}

	ircBot.conn.VerboseCallbackHandler = false
	ircBot.conn.Debug = false
	ircBot.conn.UseTLS = params.UseTLS
	ircBot.conn.TLSConfig = &tls.Config{
		ServerName:         strings.Split(params.Server, ":")[0],
		InsecureSkipVerify: !params.VerifyTLS,
	}

	err := ircBot.conn.Connect(params.Server)
	if err != nil {
		return nil, fmt.Errorf("irc.connect: %v", err)
	}

	ircBot.conn.AddCallback("001", func(e *ircevent.Event) { ircBot.conn.Join(params.Channel) })
	ircBot.conn.AddCallback(PRIVMSG, ircBot.PrivMsg)

	// Listen for replies
	go ircBot.replyWorker()

	return ircBot, nil
}

func (bot *IrcBot) Run() {
	bot.conn.Loop()
}

func (bot *IrcBot) PrivMsg(e *ircevent.Event) {
	if len(e.Arguments) != 2 {
		return
	}

	nickname := e.Nick
	source := e.Source
	channel := e.Arguments[0]
	line := e.Arguments[1]

	// Submit line into feed topic
	bot.unfilteredChan <- common.RawMessage{
		Channel:  channel,
		Hostmask: source,
		Nickname: nickname,
		Message:  line,
	}

	// Check if we are dealing with a command
	if !strings.HasPrefix(line, bot.params.CommandPrefix) {
		return
	}
	commandString := line[1:]

	commandResult := reCommand.FindAllStringSubmatch(commandString, -1)
	if len(commandResult) != 1 {
		fmt.Printf("Invalid command received\n")
		return
	}
	command := commandResult[0][1]

	params := ""
	rawParams := strings.Join(strings.Split(line, " ")[1:], " ")
	if len(rawParams) > 0 {
		paramsResult := reValidParams.FindAllStringSubmatch(rawParams, -1)
		if len(paramsResult) != 1 {
			fmt.Printf("Invalid parameters received\n")
			return
		}
		params = rawParams
	}

	switch command {
	case cmdWhoAmI:
		bot.WhoAmI(channel, source, nickname)
	case cmdTest:
		bot.Test(channel, source, nickname)
	case cmdHelp:
		bot.Help(channel, source)
	case cmdCommands:
		bot.Commands(channel, source)
	case cmdMeet:
		bot.Meet(channel, source, params)
	case cmdForget:
		bot.Forget(channel, source, params)
	case cmdAddRole:
		bot.AddRole(channel, source, params)
	case cmdRemoveRole:
		bot.RemoveRole(channel, source, params)
	case cmdListRoles:
		bot.ListRoles(channel, source)
	case cmdAddPerm:
		bot.AddPerm(channel, source, params)
	case cmdDelPerm:
		bot.DeletePerm(channel, source, params)
	case cmdListPerms:
		bot.ListPerms(channel, source, params)
	default:
		bot.SubmitCommand(channel, source, nickname, command, params)
	}
}
