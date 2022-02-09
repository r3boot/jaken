package ircbot

import (
	"jaken/lib/common"
	"strings"
)

func (bot *IrcBot) SubmitCommand(channel, hostmask, nickname, command, arguments string) {
	if !bot.IsAuthorized(hostmask, command) {
		bot.conn.Privmsgf(channel, "Not authorized to run command")
		return
	}

	msg := common.CommandMessage{
		Channel:   channel,
		Hostmask:  hostmask,
		Nickname:  nickname,
		Command:   command,
		Arguments: arguments,
	}

	// Submit message to command topic
	bot.commandChan <- msg
}

func (bot *IrcBot) WhoAmI(channel, caller, nickname string) {
	if bot.state.HasNickname(nickname) {
		bot.conn.Privmsgf(channel, "You are %s (%s)", nickname, caller)
	} else {
		bot.conn.Privmsgf(channel, "I dont know you")
	}
}

func (bot *IrcBot) Test(channel, caller, nickname string) {
	if !bot.IsAuthorized(caller, "allow") {
		bot.conn.Privmsgf(channel, "Not authorized to run command")
		return
	}

	bot.conn.Privmsgf(channel, "So thats, like, just your test %s", nickname)
}

func (bot *IrcBot) Help(channel, caller string) {
	if !bot.IsAuthorized(caller, "allow") {
		bot.conn.Privmsgf(channel, "Not authorized to run command")
		return
	}

	bot.conn.Privmsgf(channel, "See !commands ")
}

func (bot *IrcBot) Commands(channel, caller string) {
	if !bot.IsAuthorized(caller, "allow") {
		bot.conn.Privmsgf(channel, "Not authorized to run command")
		return
	}

	commands := bot.GetAuthorizedCommands(caller)
	formattedCommands := strings.Join(commands, ", ")
	bot.conn.Privmsgf(channel, "Available commands: %s", formattedCommands)
}

func (bot *IrcBot) Meet(channel, caller, nickname string) {
	if !bot.IsAuthorized(caller, "users") {
		bot.conn.Privmsgf(channel, "Not authorized to run command")
		return
	}

	if nickname == "" {
		bot.conn.Privmsgf(channel, "Need a nickname")
		return
	}

	if !bot.state.HasNickname(nickname) {
		bot.state.AddNickname(nickname)
	}

	bot.AddHostmaskFor(nickname)

	bot.state.AddPermission(nickname, memberRoleName)

	bot.conn.Privmsgf(channel, "Pleased to meet you %s", nickname)
}

func (bot *IrcBot) Forget(channel, caller, nickname string) {
	if !bot.IsAuthorized(caller, "users") {
		bot.conn.Privmsgf(channel, "Not authorized to run command")
		return
	}

	if nickname == "" {
		bot.conn.Privmsgf(channel, "Need a nickname")
		return
	}

	if !bot.state.HasNickname(nickname) {
		bot.conn.Privmsg(channel, "Unknown nickname")
		return
	}

	bot.state.RemovePermissions(nickname)
	bot.state.RemoveNickname(nickname)
	bot.conn.Privmsgf(channel, "Forgot about %s", nickname)
}

func (bot *IrcBot) AddRole(channel, caller, role string) {
	if !bot.IsAuthorized(caller, "rbac") {
		bot.conn.Privmsgf(channel, "Not authorized to run command")
		return
	}

	if role == "" {
		bot.conn.Privmsgf(channel, "Usage: %s <role>", cmdAddRole)
		return
	}

	if !bot.IsValidString(role) {
		bot.conn.Privmsgf(channel, "Role name is not valid")
		return
	}

	if bot.state.HasRole(role) {
		bot.conn.Privmsgf(channel, "Role already exists")
		return
	}

	bot.state.AddRole(role)
	bot.conn.Privmsgf(channel, "Added role %s", role)
}

func (bot *IrcBot) RemoveRole(channel, caller, role string) {
	if !bot.IsAuthorized(caller, "rbac") {
		bot.conn.Privmsgf(channel, "Not authorized to run command")
		return
	}

	if role == "" {
		bot.conn.Privmsgf(channel, "Usage: %s <role>", cmdRemoveRole)
		return
	}

	if !bot.IsValidString(role) {
		bot.conn.Privmsgf(channel, "Role name is not valid")
		return
	}

	if !bot.state.HasRole(role) {
		bot.conn.Privmsgf(channel, "No such role")
		return
	}

	bot.state.RemoveRole(role)
	bot.conn.Privmsgf(channel, "Removed role %s", role)
}

func (bot *IrcBot) ListRoles(channel, caller string) {
	if !bot.IsAuthorized(caller, "rbac") {
		bot.conn.Privmsgf(channel, "Not authorized to run command")
		return
	}

	formattedRoles := ""

	roles := bot.state.GetRoles()

	if len(roles) > 0 {
		formattedRoles = strings.Join(roles, ", ")
		bot.conn.Privmsgf(channel, "%s", formattedRoles)
	} else {
		bot.conn.Privmsgf(channel, "No roles defined")
	}
}

func (bot *IrcBot) AddPerm(channel, caller, params string) {
	if !bot.IsAuthorized(caller, "rbac") {
		bot.conn.Privmsgf(channel, "Not authorized to run command")
		return
	}

	if params == "" || strings.Count(params, " ") != 1 {
		bot.conn.Privmsgf(channel, "Usage: %s <nickname> <role>", cmdAddPerm)
		return
	}

	tokens := strings.Split(params, " ")

	if !bot.IsValidString(tokens[0]) {
		bot.conn.Privmsgf(channel, "Nickname is not valid")
		return
	}
	nickname := tokens[0]

	if !bot.IsValidString(tokens[1]) {
		bot.conn.Privmsgf(channel, "Role is not valid")
		return
	}
	role := tokens[1]

	if !bot.state.HasNickname(nickname) {
		bot.conn.Privmsgf(channel, "Nickname does not exist")
		return
	}

	if !bot.state.HasRole(role) {
		bot.conn.Privmsgf(channel, "Role does not exist")
		return
	}

	if bot.state.HasPermission(nickname, role) {
		bot.conn.Privmsgf(channel, "User already has that permission")
		return
	}

	// Errors can be safely ignored
	bot.state.AddPermission(nickname, role)
	bot.conn.Privmsgf(channel, "Granted %s to %s", role, nickname)
}

func (bot *IrcBot) DeletePerm(channel, caller, params string) {
	if !bot.IsAuthorized(caller, "rbac") {
		bot.conn.Privmsgf(channel, "Not authorized to run command")
		return
	}

	if params == "" || strings.Count(params, " ") != 1 {
		bot.conn.Privmsgf(channel, "Usage: %s <nickname> <role>", cmdDelPerm)
		return
	}

	tokens := strings.Split(params, " ")

	if !bot.IsValidString(tokens[0]) {
		bot.conn.Privmsgf(channel, "Nickname is not valid")
		return
	}
	nickname := tokens[0]

	if !bot.IsValidString(tokens[1]) {
		bot.conn.Privmsgf(channel, "Role is not valid")
		return
	}
	role := tokens[1]

	if !bot.state.HasNickname(nickname) {
		bot.conn.Privmsgf(channel, "Nickname does not exist")
		return
	}

	if !bot.state.HasRole(role) {
		bot.conn.Privmsgf(channel, "Role does not exist")
		return
	}

	if !bot.state.HasPermission(nickname, role) {
		bot.conn.Privmsgf(channel, "User does not have that permission")
		return
	}

	// Errors can be safely ignored
	bot.state.RemovePermission(nickname, role)
	bot.conn.Privmsgf(channel, "Revoked %s from %s", role, nickname)
}

func (bot *IrcBot) ListPerms(channel, caller, params string) {
	if !bot.IsAuthorized(caller, "rbac") {
		bot.conn.Privmsgf(channel, "Not authorized to run command")
		return
	}

	nickname := bot.state.GetNicknameForHostmask(caller)
	if len(params) > 0 && bot.IsValidString(params) {
		nickname = bot.state.GetNicknameForHostmask(caller)
	}
	if nickname == "" {
		bot.conn.Privmsgf(channel, "Nickname unknown")
		return
	}

	perms := bot.state.ListPermissions(nickname)
	if len(perms) > 0 {
		bot.conn.Privmsgf(channel, strings.Join(perms, ", "))
	} else {
		bot.conn.Privmsgf(channel, "No permissions defined for nickname")
	}
}

func (bot *IrcBot) AddBinding(channel, caller, params string) {
	if !bot.IsAuthorized(caller, "rbac") {
		bot.conn.Privmsgf(channel, "Not authorized to run command")
		return
	}

	if params == "" || strings.Count(params, " ") != 1 {
		bot.conn.Privmsgf(channel, "Usage: %s <command> <role>", cmdAddBinding)
		return
	}

	tokens := strings.Split(params, " ")

	if !bot.IsValidString(tokens[0]) {
		bot.conn.Privmsgf(channel, "command is not valid")
		return
	}
	command := tokens[0]

	if !bot.IsValidString(tokens[1]) {
		bot.conn.Privmsgf(channel, "Role is not valid")
		return
	}
	role := tokens[1]

	if !bot.state.HasRole(role) {
		bot.conn.Privmsgf(channel, "Role does not exist")
		return
	}

	if bot.state.HasBinding(command, role) {
		bot.conn.Privmsgf(channel, "Binding does not exist")
		return
	}

	// Errors can be safely ignored
	bot.state.AddBinding(command, role)
	bot.conn.Privmsgf(channel, "Bound %s to %s", command, role)
}

func (bot *IrcBot) DeleteBinding(channel, caller, params string) {
	if !bot.IsAuthorized(caller, "rbac") {
		bot.conn.Privmsgf(channel, "Not authorized to run command")
		return
	}

	if params == "" || strings.Count(params, " ") != 1 {
		bot.conn.Privmsgf(channel, "Usage: %s <command> <role>", cmdDelPerm)
		return
	}

	tokens := strings.Split(params, " ")

	if !bot.IsValidString(tokens[0]) {
		bot.conn.Privmsgf(channel, "Command is not valid")
		return
	}
	command := tokens[0]

	if !bot.IsValidString(tokens[1]) {
		bot.conn.Privmsgf(channel, "Role is not valid")
		return
	}
	role := tokens[1]

	if !bot.state.HasRole(role) {
		bot.conn.Privmsgf(channel, "Role does not exist")
		return
	}

	if !bot.state.HasBinding(command, role) {
		bot.conn.Privmsgf(channel, "Binding does not exist")
		return
	}

	// Errors can be safely ignored
	bot.state.RemoveBinding(command, role)
	bot.conn.Privmsgf(channel, "Removed %s from %s", role, command)
}

func (bot *IrcBot) ListBindings(channel, caller, role string) {
	if !bot.IsAuthorized(caller, "rbac") {
		bot.conn.Privmsgf(channel, "Not authorized to run command")
		return
	}

	if role == "" {
		bot.conn.Privmsgf(channel, "Usage: %s <role>", cmdListBindings)
		return
	}

	if !bot.state.HasRole(role) {
		bot.conn.Privmsgf(channel, "Role does not exist")
		return
	}

	bindings := bot.state.GetBindingsForRole(role)
	if len(bindings) > 0 {
		bot.conn.Privmsgf(channel, strings.Join(bindings, ", "))
	} else {
		bot.conn.Privmsgf(channel, "No bindings defined for %s", role)
	}
}
