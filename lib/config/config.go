package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	yaml "gopkg.in/yaml.v2"
)

const FileName = ".lazy-awslogs"

type Stream struct {
	Name string
}
type Group struct {
	Name    string
	Streams []Stream
}
type Environment struct {
	Name          string
	Profile       string
	Region        string
	DefaultGroup  string
	DefaultStream string
	Groups        []Group
}
type Configuration struct {
	Current      string
	Environments []Environment
}
type EnvironmentKey int

const (
	Name EnvironmentKey = iota
	Profile
	Region
	DefaultGroup
	DefaultStream
)

// Load returns a Configuration object.
func Load() (c Configuration) {
	err := yaml.Unmarshal(configurationOpen(), &c)
	if err != nil {
		log.Fatalf("error: %v", err)
		os.Exit(1)
	}
	return
}

// save writes self(Configuration object).
func (this *Configuration) save() {
	encodedBytes, _ := yaml.Marshal(&this)
	ioutil.WriteFile(configurationFilePath(), encodedBytes, 0644)
}

// configurationOpen returns a Configuration bytes.
func configurationOpen() []byte {
	filename := configurationFilePath()
	if _, err := os.Stat(filename); err != nil {
		return nil
	}

	file, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("error: %v", err)
		os.Exit(1)
	}
	return file
}

// configurationFilePath returns a Configuration path.
func configurationFilePath() string {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return home + "/" + FileName + ".yaml"
}

// CurrentEnvironment returns Environment current.
func (this *Configuration) CurrentEnvironment() (environment Environment) {
	for i, env := range this.Environments {
		if env.Name == this.Current {
			environment = this.Environments[i]
			break
		}
	}
	return
}

// CurrentEnvironment returns Environment current pointer.
func (this *Configuration) CurrentEnvironmentP() (environment *Environment) {
	for i, env := range this.Environments {
		if env.Name == this.Current {
			environment = &this.Environments[i]
			break
		}
	}
	return
}

// UpdateCurrent updates CurrentEnvironment.
func (this *Configuration) UpdateCurrent(current string) (err error) {
	exist := false
	for _, env := range this.Environments {
		if env.Name == current {
			this.Current = current
			exist = true
			this.save()
			break
		}
	}
	if !exist {
		err = errors.New(current + " is not exist.")
	}
	return
}

// SetCurrentProfile sets profile name in current Environment.
func (this *Configuration) SetCurrentEnvironmentParam(key EnvironmentKey, value string) {
	env := this.CurrentEnvironmentP()
	switch key {
	case Profile:
		env.Profile = value
	case Region:
		env.Region = value
	case DefaultGroup:
		env.DefaultGroup = value
	case DefaultStream:
		env.DefaultStream = value
	}
	this.save()
	return
}

// AddEnvironment adds a Environment.
func (this *Configuration) AddEnvironment(newEnvironment Environment) (err error) {
	exist := false
	for _, env := range this.Environments {
		if env.Name == newEnvironment.Name {
			exist = true
			break
		}
	}
	if exist {
		err = errors.New(newEnvironment.Name + " already exist.")
	}
	this.Environments = append(this.Environments, newEnvironment)

	if len(this.Environments) == 1 {
		this.Current = newEnvironment.Name
	}

	this.save()
	return
}

// RemoveEnvironment remove a Environment.
func (this *Configuration) RemoveEnvironment(targetName string) (err error) {
	// This func create a new []Environment without targetName's Environment
	removedEnvironments := func() (newEnvironments []Environment) {
		for _, env := range this.Environments {
			if env.Name != targetName {
				newEnvironments = append(newEnvironments, env)
			}
		}
		return
	}

	exist := false
	for _, env := range this.Environments {
		if env.Name == targetName {
			exist = true
			if this.Current == env.Name {
				err = errors.New(targetName + " is current Environment.")
				break
			}
			this.Environments = removedEnvironments()
			this.save()
			break
		}
	}
	if !exist {
		err = errors.New(targetName + " is not exist.")
	}
	return
}

// ReadGroups returns groups in Config.
func (this *Configuration) ReadGroups() (groups []Group) {
	for i, env := range this.Environments {
		if env.Name == this.Current {
			groups = this.Environments[i].Groups
			break
		}
	}
	return
}

// UpdateGroups updates groups in Config.
func (this *Configuration) UpdateGroups(groups []Group) {
	for i, env := range this.Environments {
		if env.Name == this.Current {
			this.Environments[i].Groups = groups
			break
		}
	}
	this.save()
}

// ReadStreams returns streams in Config.
func (this *Configuration) ReadStreams(groupName string) (streams []Stream) {
	// update groups of current environment
	for i, env := range this.Environments {
		if env.Name == this.Current {
			for j, group := range this.Environments[i].Groups {
				if group.Name == groupName {
					streams = this.Environments[i].Groups[j].Streams
					break
				}
			}
		}
	}
	return
}

// UpdateStreams update streams in Config.
func (this *Configuration) UpdateStreams(groupName string, streams []Stream) {
	// update groups of current environment
	for i, env := range this.Environments {
		if env.Name == this.Current {
			for j, group := range this.Environments[i].Groups {
				if group.Name == groupName {
					this.Environments[i].Groups[j].Streams = streams
					break
				}
			}
		}
	}
	this.save()
}
