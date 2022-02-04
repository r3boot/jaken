package ircbot

import "fmt"

func (bot *IrcBot) IsOwner(hostmask string) bool {
	return bot.params.Owner == hostmask
}

func (bot *IrcBot) IsAuthorized(hostmask, command string) bool {
	// The owner is allowed to run all commands
	if bot.IsOwner(hostmask) {
		fmt.Printf("Allowed %s for %s: owner\n", command, hostmask)
		return true
	}

	role := bot.state.GetBindingRole(command)
	if role == "" {
		return false
	}

	nickname := bot.state.GetNicknameForHostmask(hostmask)
	if nickname == "" {
		fmt.Printf("Denied %s to %s: user unknown\n", command, hostmask)
		return false
	}

	for _, perm := range bot.state.ListPermissions(nickname) {
		if perm == role {
			fmt.Printf("Allowed %s for %s\n", command, hostmask)
			return true
		}
	}

	// Disallow by default
	fmt.Printf("Denied %s to %s: unknown error", command, hostmask)
	return false
}
