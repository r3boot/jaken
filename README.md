# Design principles
This bot is based on the premise of a loosely coupling between the bot and its plugins, so that plugin development can
be made more easily. It does this by assuming that plugins are standalone binaries that can be executed. The second
premise is that there needs to be a fine-grained ACL mechanism. This is done with roles applied to plugins.

# Access control
This bot features fine-grained RBAC with a command-level granularity. Every command can be bound to a role, and users
can get permissions per role. By default, this works on the command level, but it is possible to bind multiple commands
to a single role. During intialization of the bot, one default role is created for admin users and and the user
management commands are bound to this role.

# Built-in commands
In order to make the bot work, it includes several commands by itself, a list of these can be found below

## whoami
Display how you are recognized by the bot. Syntax is `whoami`.

## test
Test if you can talk with the bot. Syntax is `test`

## help
Get some info to get you started. Syntax is `help`

## commands
List all available commands. Syntax is `help`

## meet
Introduce a new user to the bot. Syntax is `meet <nickname>`. This will cause the bot to perform a whois lookup for this
user, and store the corresponding hostmask into the database

## forget
Remove a user from the bot. Syntax is `forget <nickname>`. This will delete all hostmasks of users matching the nickname.

## add-role
Define a new role. Syntax is `add-role <role>`.

## del-role
Remove a role. Syntax is `del-role <role>`

## list-roles
List all available roles. Syntax is `list-roles`

## add-perm
Grants a user permission to a role. Syntax is `add-perm <nickname> <role>`

## del-perm
Revokes a permission from a user. Syntax is `del-perm <nickname> <role>`

## list-perms
List all roles for a user. Syntax is `list-perms [<nickname>]`. By default the permissions for the calling user will
be shown. By specifying `<nickname>` you can lookup the permissions for another user.

# Plugin development guidelines
Plugins should be designed such that they can both run standalone (ie, from the commandline) as well as being able to
run from inside the bot. The bot will take the filename of the script, minus the extension, to create the IRC command.
Because of this, the plugins need to adhere to the following rules:

* Filename must be a single word and can optionally have an extension
* Incoming arguments to the plugins must be set via argv
* Replies to the bot need to be sent to stdout

By default, the command is bound to a role with the same name as the command. This can be overriden by prepending the
filename of the command with `<role name>_`. This role identifier (including the underscore) will be removed from the
IRC command.

# Upcoming features
* Additional information passing to plugins
* karma/infoitems, in-bot or not
* Possibility to reply via notifications
* Port some plugins
* Alias support
* Endpoint for incoming irc messages
* Make bot configurable via cfg file and cli arguments
* Find something for triggers (eg, !remind, leaving messages for ppl)
* Possibility to create plugins which run continuously and communicate via stdin/stdout
* sqlite3 -> add transactions and concurrency limits
* Possibility to join multiple channels