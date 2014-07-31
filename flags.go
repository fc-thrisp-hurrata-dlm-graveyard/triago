package triago

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

// Registers a flag set to be parsed. Register all flag sets
// before calling this function. flag.CommandLine is automatically
// registered.
func (c *Config) RegisterFlagSet(flagSetName string, set *flag.FlagSet) {
	c.flags[flagSetName] = set
}

// Sets a flag's value and persists the changes to the disk.
func (c *Config) Set(flagSetName string, f *flag.Flag) error {
	name := f.Name
	value := f.Value.String()

	c.AddOption(flagSetName, name, value)

	if c.FileName != "" {
		return c.WriteFile(c.FileName,
			0644,
			fmt.Sprintf("action flagset: flag: value = %s:%s:%s",
				flagSetName, name, value))
	}
	return nil
}

// Deletes a flag from config file and persists the changes to the disk.
func (c *Config) Delete(flagSetName, flagName string) error {
	c.RemoveOption(flagSetName, flagName)
	if c.FileName != "" {
		return c.WriteFile(c.FileName,
			0644,
			fmt.Sprintf("action delete flagset: flag = %s:%s",
				flagSetName, flagName))
	}
	return nil
}

// Parses the config file for the provided flag set.
// If the flags are already set, values are overwritten
// by the values in the config file. Defaults are not set
// if the flag is not in the file.
func (c *Config) ParseSet(flagSetName string, set *flag.FlagSet) {
	set.VisitAll(func(f *flag.Flag) {
		val := getEnv(c.EnvPrefix, flagSetName, f.Name)
		if val != "" {
			set.Set(f.Name, val)
			return
		}

		val, err := c.String(flagSetName, f.Name)
		if err == nil {
			set.Set(f.Name, val)
		}
	})
}

// Parses all the registered flag sets, including the command
// line set and sets values from the config file if they are
// not already set.
func (c *Config) Parse() {
	for name, set := range c.flags {
		alreadySet := make(map[string]bool)
		set.Visit(func(f *flag.Flag) {
			alreadySet[f.Name] = true
		})
		set.VisitAll(func(f *flag.Flag) {
			// if not already set, set it from dict if exists
			if alreadySet[f.Name] {
				return
			}

			val := getEnv(c.EnvPrefix, name, f.Name)
			if val != "" {
				set.Set(f.Name, val)
				return
			}

			val, err := c.String(name, f.Name)
			if err == nil {
				set.Set(f.Name, val)
			}
		})
	}
}

// Parses command line flags and then, all of the registered
// flag sets with the values provided in the config file.
func (c *Config) ParseAll() {
	if !flag.Parsed() {
		flag.Parse()
	}
	c.Parse()
}

// Looks up variable in environment
func getEnv(envPrefix, flagSetName, flagName string) string {
	// If we haven't set an EnvPrefix, don't lookup vals in the ENV
	if envPrefix == "" {
		return ""
	}
	// Append a _ to flagSetName if it exists.
	if flagSetName != "" {
		flagSetName += "_"
	}
	flagName = strings.Replace(flagName, ".", "_", -1)
	envKey := strings.ToUpper(envPrefix + flagSetName + flagName)
	return os.Getenv(envKey)
}
