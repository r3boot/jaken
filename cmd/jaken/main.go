package main

import (
	"jaken/lib/ircbot"
	"jaken/lib/ircstate"
	"jaken/lib/pluginmgr"
)

const (
	dbPath     = "./jaken.db"
	pluginPath = "./plugins"
)

func main() {
	plugins := pluginmgr.New(&pluginmgr.PluginParams{
		PluginPath: pluginPath,
	})

	state, err := ircstate.New(dbPath)
	defer state.Close()

	bot, err := ircbot.New(&ircbot.Params{
		Server:        "irc.oftc.net:6697",
		UseTLS:        true,
		VerifyTLS:     false,
		Channel:       "#nurdbottest",
		Nickname:      "jaken",
		Realname:      "Jaken",
		CommandPrefix: "@",
		Owner:         "r3boot!~r3boot@shell.as65342.net",
	}, state, plugins)
	if err != nil {
		panic(err)
	}

	bot.Run()
}
