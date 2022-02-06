package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

const (
	DefaultCfgFile       = "ircbot.yml"
	DefaultServer        = "localhost:6667"
	DefaultUseTls        = false
	DefaultVerifyTls     = true
	DefaultChannel       = "#example"
	DefaultNickname      = "ircbot"
	DefaultRealname      = "ircbot"
	DefaultCommandPrefix = "!"
	DefaultOwner         = ""
	DefaultDbPath        = "./jaken.db"
	DefaultPluginPath    = "./plugins"

	HelpCfgFile       = "Yaml file containing configuration (" + DefaultCfgFile + ")"
	HelpServer        = "IRC server to connect to (" + DefaultServer + ")"
	HelpUseTls        = "Use TLS to connect to the server (false)" // TODO
	HelpVerifyTls     = "Verify TLS server certificate (true)"     // TODO
	HelpChannel       = "Default channel to join (" + DefaultChannel + ")"
	HelpNickname      = "Nickname to use (" + DefaultNickname + ")"
	HelpRealName      = "Realname to use (" + DefaultRealname + ")"
	HelpCommandPrefix = "Commandprefix to use (" + DefaultCommandPrefix + ")"
	HelpOwner         = "Hostmask of the owner of this bot (unset)"
	HelpDbPath        = "Path to the database (" + DefaultDbPath + ")"
	HelpPluginPath    = "Path to the plugins (" + DefaultPluginPath + ")"
)

type Settings struct {
	Server        string `yaml:"server"`
	UseTls        bool   `yaml:"use_tls"`
	VerifyTls     bool   `yaml:"verify_tls"`
	Channel       string `yaml:"channel"`
	Nickname      string `yaml:"nickname"`
	Realname      string `yaml:"realname"`
	Owner         string `yaml:"owner"`
	CommandPrefix string `yaml:"command_prefix"`
	DbPath        string `yaml:"dbpath"`
	PluginPath    string `yaml:"pluginPath"`
}

func defaultSettings() *Settings {
	return &Settings{
		Server:        DefaultServer,
		UseTls:        DefaultUseTls,
		VerifyTls:     DefaultVerifyTls,
		Channel:       DefaultChannel,
		Nickname:      DefaultNickname,
		Realname:      DefaultRealname,
		Owner:         DefaultOwner,
		CommandPrefix: DefaultCommandPrefix,
		DbPath:        DefaultDbPath,
		PluginPath:    DefaultPluginPath,
	}
}

func loadConfiguration(cfgfile string) (*Settings, error) {
	DefaultSettings := defaultSettings()
	content, err := ioutil.ReadFile(cfgfile)
	if err != nil {
		return DefaultSettings, fmt.Errorf("loadConfiguration ioutil.ReadFile: %v", err)
	}

	settings := &Settings{}
	err = yaml.Unmarshal(content, settings)
	if err != nil {
		return DefaultSettings, fmt.Errorf("loadConfiguration yaml.Unmarshal: %v", err)
	}

	return settings, nil
}

func stringYamlOrFlagOrEnvVar(d, y, f, e string) string {
	result := y

	if len(f) > 0 && f != d {
		result = f
	}

	envVar, ok := os.LookupEnv(e)
	if ok && envVar != d {
		result = envVar
	}

	return result
}

func boolYamlOrFlagOrEnvVar(d, y, f bool, e string) bool {
	result := y

	if f != d {
		result = f
	}

	envVar, ok := os.LookupEnv(e)
	parsedEnvVar := false
	if envVar == "1" {
		parsedEnvVar = true
	}
	if ok && parsedEnvVar != d {
		result = parsedEnvVar
	}

	return result
}

func Load(cfgfile string, args *Settings) *Settings {
	yml, _ := loadConfiguration(cfgfile)

	return &Settings{
		Server:        stringYamlOrFlagOrEnvVar(DefaultServer, yml.Server, args.Server, "IRCBOT_SERVER"),
		UseTls:        boolYamlOrFlagOrEnvVar(DefaultUseTls, yml.UseTls, args.UseTls, "IRCBOT_USETLS"),
		VerifyTls:     boolYamlOrFlagOrEnvVar(DefaultVerifyTls, yml.VerifyTls, args.VerifyTls, "IRCBOT_VERIFYTLS"),
		Channel:       stringYamlOrFlagOrEnvVar(DefaultChannel, yml.Channel, args.Channel, "IRCBOT_CHANNEL"),
		Nickname:      stringYamlOrFlagOrEnvVar(DefaultNickname, yml.Nickname, args.Nickname, "IRCBOT_NICKNAME"),
		Realname:      stringYamlOrFlagOrEnvVar(DefaultRealname, yml.Realname, args.Realname, "IRCBOT_REALNAME"),
		Owner:         stringYamlOrFlagOrEnvVar(DefaultOwner, yml.Owner, args.Owner, "IRCBOT_OWNER"),
		CommandPrefix: stringYamlOrFlagOrEnvVar(DefaultCommandPrefix, yml.CommandPrefix, args.CommandPrefix, "IRCBOT_COMMANDPREFIX"),
		DbPath:        stringYamlOrFlagOrEnvVar(DefaultDbPath, yml.DbPath, args.DbPath, "IRCBOT_DBPATH"),
		PluginPath:    stringYamlOrFlagOrEnvVar(DefaultPluginPath, yml.PluginPath, args.PluginPath, "IRCBOT_PLUGINPATH"),
	}
}
