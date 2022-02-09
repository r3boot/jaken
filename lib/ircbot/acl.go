package ircbot

import "fmt"

var (
	adminCommands = []string{
		cmdMeet,
		cmdForget,
		cmdAddRole,
		cmdRemoveRole,
		cmdListRoles,
		cmdAddPerm,
		cmdDelPerm,
		cmdListPerms,
		cmdAddBinding,
		cmdDelBinding,
		cmdListBindings,
	}
	alwaysCommands = []string{
		cmdHelp,
		cmdTest,
		cmdCommands,
		cmdWhoAmI,
	}
)

func (bot *IrcBot) IsOwner(hostmask string) bool {
	return bot.params.Owner == hostmask
}

func (bot *IrcBot) IsAuthorized(hostmask, command string) bool {

	/*
		// The owner is allowed to run all commands
		if bot.IsOwner(hostmask) {
			fmt.Printf("Allowed %s for %s: owner\n", command, hostmask)
			return true
		}

	*/

	role := bot.state.GetBindingRole(command)
	if role == "" {
		// If we cannot find a binding, set it to the default one
		fmt.Printf("No role found for %s, using %s\n", command, memberRoleName)
		role = memberRoleName
	}

	if role == adminRoleName {
		fmt.Printf("Allowed %s for %s: admin\n", command, hostmask)
		return true
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
	fmt.Printf("Denied %s to %s: unknown error\n", command, hostmask)
	return false
}

func (bot *IrcBot) GetAuthorizedCommands(hostmask string) []string {
	var commands []string

	commands = append(commands, alwaysCommands...)

	if bot.IsOwner(hostmask) {
		commands = append(commands, adminCommands...)
		// commands = append(commands, bot.plugins.ListPlugins()...)
	} else {
		nickname := bot.state.GetNicknameForHostmask(hostmask)
		if nickname == "" {
			return nil
		}

		/*
			roles := bot.plugins.GetRoles()

			for _, perm := range bot.state.ListPermissions(nickname) {
				if perm == "admin" {
					commands = append(commands, adminCommands...)
					continue
				}
				for _, role := range roles {
					if role == perm {
						commands = append(commands, bot.plugins.ListCommandsForRole(role)...)
					}
				}
			}

		*/
	}

	return commands
}
