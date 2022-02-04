# Design principles
This bot is based on the premise of a loosely coupling between the bot and its plugins, so that plugin development can
be made more easily. It does this by assuming that plugins are standalone binaries that can be executed. The second
premise is that there needs to be a fine-grained ACL mechanism. This is done with roles applied to plugins.

# Built-in commands
In order to make the bot work, it includes several commands by itself, a list of these can be found below

## meet
Introduce a new user to the bot. Syntax is `meet nickname`. This will cause the bot to perform a whois lookup for this
user, and store the corresponding hostmask into the database

## forget
Remove a user from the bot. Syntax is `forget nickname`. This will delete all hostmasks of users matching the nickname.

## whoami
Display how you are recognized by the bot. Syntax is `whoami`.

# Plugin development guidelines
Plugins should be designed such that they can both run standalone (ie, from the commandline) as well as being able to
run from inside the bot. Because of this, the plugins need to adhere to the following rules:
* Incoming arguments to the plugins must be set via argv
* Replies to the bot need to be sent to stdout

# Upcoming features
* Fine-grained ACLs
* Additional information passing to plugins
* karma/infoitems, in-bot or not
* Possibility to reply via notifications
* Port some plugins
* Alias support
* Endpoint for incoming irc messages