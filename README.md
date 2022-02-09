# Design principles
This bot is based on a couple of different design criteriums:
* There needs to be a loose coupling between the bot and the plugins
* Fine-grained ACLs up to the command level need to be possible

## Loose coupling
Loose coupling of the bot will allow rapid plugin development without involving the bot. Furthermore, since the bots
do not depend on any IRC semantics, it is possible to run the bot outside of the context of the bot. This is implemented
by adding mqtt in between.

## Fine-grained ACLs
In order to implement some form of access-control, ACLs need to be in place. This need to be on the hostmask/command
granularity, using roles, principals and subjects.

# Configuration
The configuration settings for this bot can be set in three different ways: commandline arguments, environment variables
and yaml-based configuration. The precedence is environment > commandline > yaml. The following options are available:

|yaml|argument|envvar| description                         |default value|
|----|--------|------|-------------------------------------|-------------|
|server|-server|IRCBOT_SERVER| Which server to connect to |localhost:6667|
|use_tls|-usetls|IRCBOT_USETLS| Use TLS to connect to the server | false |
|verify_tls|-verifytls|IRCBOT_VERIFYTLS| Verify TLS server certificate |true|
|channel|-channel|IRCBOT_CHANNEL| Default channel to join |#example |
|nickname|-nickname|IRCBOT_NICKNAME| Nickname to use |ircbot|   
|realname|-realname|IRCBOT_REALNAME|Realname to use |ircbot|
|owner|-owner|IRCBOT_OWNER|Hostmask of the owner |unset|
|command_prefix|-commandprefix|IRCBOT_COMMANDPREFIX|Commandprefix to use |!|
|db_path|-dbpath|IRCBOT_DBPATH|Path to the database |./jaken.db|
|plugin_path|-pluginpath|IRCBOT_PLUGINPATH|Path to the plugins|./plugins|


# Access control
This bot features fine-grained RBAC with a command-level granularity. Every command can be bound to a role, and users
can get permissions per role. By default, this works on the command level, but it is possible to bind multiple commands
to a single role. During intialization of the bot, one default role is created for admin users and and the user
management commands are bound to this role.

# Built-in commands
In order to make the bot work, it includes several commands by itself, a list of these can be found below

|command|role|description|
|-------|----|-----------|
|whoami|member|Display how you are recognized by the bot|
|test|member|Test if you can talk with the bot|
|help|member|Get some info to get you started|
|commands|member|List all commands that are available to you|
|meet|admin|Introduce a new user to the bot. This will cause the bot to perform a whois lookup for this user, and store the corresponding hostmask into the database|
|forget|admin|Remove a user from the bot. This will delete all hostmasks of users matching the nickname|
|add-role|admin|Define a new role|
|del-role|admin|Remove a role|
|list-roles|admin|List all available roles|
|add-perm|admin|Grants a user permission to a role|
|del-perm|admin|Revokes a permission from a user|
|list-perms|admin|List all roles for a user. By default the permissions for the calling user will be shown. By specifying `<nickname>` you can lookup the permissions for another user.|

# Mqtt topics
Several topics are available for communication to/from the bot, as can be seen in the table below. Examples for how to
use this can be found underneath the `plugins` directory,

| topic                                   |direction| description                                    |
|-----------------------------------------|---------|------------------------------------------------|
| from/irc/(channel)/(nickname)/message   |towards plugin| Raw feed of messages in (channel)              |
| from/irc/(channel)/(nickname)/(command) |towards plugin| Listen in a channel for (ControlChar)(command) |
| to/irc/(channel)/privmsg                |towards bot| Send reply in (channel) via PRIVMSG            |
| to/irc/(channel)/notice                 |towards bot| Send reply in (channel) via NOTICE             |
| to/irc/(channel)/topic                  |towards bot| Sets TOPIC for (channel)                       |

# Plugin development guidelines
Plugins can be written in any language, as long as they communicate via the mqtt topics that are available.

# Upcoming features
* Port some plugins
* Alias support
* sqlite3 -> add transactions and concurrency limits
* Possibility to join multiple channels