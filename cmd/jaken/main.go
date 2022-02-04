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
		Server:        "irc.mononoke.nl:6697",
		UseTLS:        true,
		VerifyTLS:     false,
		Channel:       "#bottest",
		Nickname:      "jaken",
		Realname:      "Jaken",
		CommandPrefix: "1",
		Owner:         "r3boot!~r3boot@cloaked",
	}, state, plugins)
	if err != nil {
		panic(err)
	}

	bot.Run()
}
