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
