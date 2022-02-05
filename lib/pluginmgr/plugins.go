package pluginmgr

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

type PluginParams struct {
	PluginPath string
}

type PluginManager struct {
	plugins map[string]string
	params  *PluginParams
}

func New(params *PluginParams) *PluginManager {
	fs, err := os.Stat(params.PluginPath)
	if err != nil {
		log.Fatalf("os.Stat %v", err)
	}
	if !fs.IsDir() {
		log.Fatalf("os.Stat not a directory")
	}

	pm := &PluginManager{
		params: params,
	}
	pm.LoadPlugins()

	return pm
}

func (pm *PluginManager) LoadPlugins() {
	newPlugins := make(map[string]string)

	files, err := ioutil.ReadDir(pm.params.PluginPath)
	if err != nil {
		panic(fmt.Sprintf("Failed to read pluginmgr from %s", pm.params.PluginPath))
	}

	for _, fs := range files {
		if fs.IsDir() {
			continue
		}

		name := strings.Split(path.Base(fs.Name()), ".")[0]
		if strings.Contains(name, "_") {
			name = strings.Join(strings.Split(name, "_")[1:], "_")
		}
		path := fmt.Sprintf("%s/%s", pm.params.PluginPath, fs.Name())

		newPlugins[name] = path
	}

	pm.plugins = newPlugins
}

func (pm *PluginManager) GetPlugin(command string) string {
	fname, found := pm.plugins[command]
	if !found {
		return ""
	}
	return fname
}

func (pm *PluginManager) Run(command, params string) string {
	pm.LoadPlugins()
	fname := pm.GetPlugin(command)
	if fname == "" {
		return "No such command"
	}

	out, err := exec.Command(fname, params).Output()
	if err != nil {
		return err.Error()
	}

	return string(out)
}

func (pm *PluginManager) GetRole(command string) string {
	pm.LoadPlugins()
	fname := pm.GetPlugin(command)
	if fname == "" {
		return ""
	}

	if strings.Contains(fname, "_") {
		return strings.Split(fname, "_")[0]
	} else {
		return command
	}
}

func (pm *PluginManager) GetRoles() []string {
	var roles []string

	files, err := ioutil.ReadDir(pm.params.PluginPath)
	if err != nil {
		panic(fmt.Sprintf("Failed to read pluginmgr from %s", pm.params.PluginPath))
	}

	for _, fs := range files {
		if fs.IsDir() {
			continue
		}

		role := strings.Split(path.Base(fs.Name()), ".")[0]
		if strings.Contains(role, "_") {
			role = strings.Join(strings.Split(role, "_")[1:], "_")
		}

		roles = append(roles, role)
	}

	return roles
}

func (pm *PluginManager) ListPlugins() []string {
	var plugins []string

	pm.LoadPlugins()

	for plugin, _ := range pm.plugins {
		plugins = append(plugins, plugin)
	}

	return plugins
}

func (pm *PluginManager) ListCommandsForRole(role string) []string {
	var commands []string
	files, err := ioutil.ReadDir(pm.params.PluginPath)
	if err != nil {
		panic(fmt.Sprintf("Failed to read pluginmgr from %s", pm.params.PluginPath))
	}

	for _, fs := range files {
		if fs.IsDir() {
			continue
		}

		name := strings.Split(path.Base(fs.Name()), ".")[0]

		if strings.HasPrefix(name, role) {
			if strings.Contains(name, "_") {
				commands = append(commands, strings.Join(strings.Split(name, "_")[1:], "_"))
			} else {
				commands = append(commands, name)
			}
		}
	}

	return commands
}
