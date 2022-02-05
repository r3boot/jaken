package ircbot

import "strings"

func (bot *IrcBot) RunPlugin(channel, caller, command, params string) {
	role := bot.plugins.GetRole(command)
	if !bot.IsAuthorized(caller, role) {
		bot.conn.Privmsgf(channel, "Not authorized to run command")
		return
	}
	response := bot.plugins.Run(command, params)
	bot.conn.Privmsg(channel, response)
}

func (bot *IrcBot) WhoAmI(channel, caller, nickname string) {
	if bot.state.HasNickname(nickname) {
		bot.conn.Privmsgf(channel, "You are %s (%s)", nickname, caller)
	} else {
		bot.conn.Privmsgf(channel, "I dont know you")
	}
}

func (bot *IrcBot) Test(channel, caller, nickname string) {
	if !bot.IsAuthorized(caller, "test") {
		bot.conn.Privmsgf(channel, "Not authorized to run command")
		return
	}

	bot.conn.Privmsgf(channel, "So thats, like, just your test %s", nickname)
}

func (bot *IrcBot) Meet(channel, caller, nickname string) {
	if !bot.IsAuthorized(caller, "meet") {
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

	bot.conn.Privmsgf(channel, "Pleased to meet you %s", nickname)
}

func (bot *IrcBot) Forget(channel, caller, nickname string) {
	if !bot.IsAuthorized(caller, "forget") {
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
	if !bot.IsAuthorized(caller, "role") {
		bot.conn.Privmsgf(channel, "Not authorized to run command")
		return
	}

	if role == "" {
		bot.conn.Privmsgf(channel, "Need a role name")
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
	if !bot.IsAuthorized(caller, "role") {
		bot.conn.Privmsgf(channel, "Not authorized to run command")
		return
	}

	if role == "" {
		bot.conn.Privmsgf(channel, "Need a role name")
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
	if !bot.IsAuthorized(caller, "role") {
		bot.conn.Privmsgf(channel, "Not authorized to run command")
		return
	}

	formattedRoles := ""

	roles := bot.state.GetRoles()

	roles = append(roles, bot.plugins.GetRoles()...)

	if len(roles) > 0 {
		formattedRoles = strings.Join(roles, ", ")
		bot.conn.Privmsgf(channel, "%s", formattedRoles)
	} else {
		bot.conn.Privmsgf(channel, "No roles defined")
	}
}

func (bot *IrcBot) AddPerm(channel, caller, params string) {
	if !bot.IsAuthorized(caller, "perm") {
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
	if !bot.IsAuthorized(caller, "perm") {
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

	if !bot.state.HasPermission(nickname, role) {
		bot.conn.Privmsgf(channel, "User does not have that permission")
		return
	}

	// Errors can be safely ignored
	bot.state.RemovePermission(nickname, role)
	bot.conn.Privmsgf(channel, "Revoked %s from %s", role, nickname)
}

func (bot *IrcBot) ListPerms(channel, caller, params string) {
	if !bot.IsAuthorized(caller, "perm") {
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
