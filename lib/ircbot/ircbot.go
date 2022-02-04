package ircbot

import (
	"crypto/tls"
	"fmt"
	ircevent "github.com/thoj/go-ircevent"
	"jaken/lib/ircstate"
	"jaken/lib/pluginmgr"
	"regexp"
	"strings"
)

const (
	PRIVMSG = "PRIVMSG"

	cmdWhoAmI     = "whoami"
	cmdMeet       = "meet"
	cmdForget     = "forget"
	cmdAddRole    = "add-role"
	cmdRemoveRole = "del-role"
	cmdListRoles  = "list-roles"
	cmdAddPerm    = "add-perm"
	cmdDelPerm    = "del-perm"
	cmdListPerms  = "list-perms"
)

type Params struct {
	Server        string
	UseTLS        bool
	VerifyTLS     bool
	Channel       string
	Nickname      string
	Realname      string
	CommandPrefix string
	Owner         string
}

type IrcBot struct {
	conn    *ircevent.Connection
	params  *Params
	state   *ircstate.State
	plugins *pluginmgr.PluginManager
}

var (
	reCommand     = regexp.MustCompile("^([a-zA-Z0-9\\+\\-]{1,64})")
	reValidParams = regexp.MustCompile("^[a-zA-Z0-9_\\-\\+\\ ]+$")
)

func New(params *Params, state *ircstate.State, plugins *pluginmgr.PluginManager) (*IrcBot, error) {
	ircBot := &IrcBot{
		conn:    ircevent.IRC(params.Nickname, params.Realname),
		params:  params,
		state:   state,
		plugins: plugins,
	}

	ircBot.conn.VerboseCallbackHandler = true
	ircBot.conn.Debug = true
	ircBot.conn.UseTLS = params.UseTLS
	ircBot.conn.TLSConfig = &tls.Config{InsecureSkipVerify: !params.VerifyTLS}

	err := ircBot.conn.Connect(params.Server)
	if err != nil {
		return nil, fmt.Errorf("irc.connect: %v", err)
	}

	ircBot.conn.AddCallback("001", func(e *ircevent.Event) { ircBot.conn.Join(params.Channel) })
	ircBot.conn.AddCallback(PRIVMSG, ircBot.PrivMsg)

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
		bot.RunPlugin(channel, source, command, params)
	}
}
