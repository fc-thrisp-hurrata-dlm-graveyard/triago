package triago

import (
	"flag"
	"strings"
)

const (
	// Default section name.
	DEFAULT_SECTION = "DEFAULT"
	// Maximum allowed depth when recursively substituting variable names.
	_DEPTH_VALUES = 200

	DEFAULT_COMMENT       = "# "
	ALTERNATIVE_COMMENT   = "; "
	DEFAULT_SEPARATOR     = ":"
	ALTERNATIVE_SEPARATOR = "="
	DEFAULT_FILENAME      = "/tmp/config.ini"
)

// Config is the representation of configuration settings.
type (
	Config struct {
		comment   string
		separator string

		// Sections order
		lastIdSection int            // Last section identifier
		idSection     map[string]int // Section : position

		// The last option identifier used for each section.
		lastIdOption map[string]int // Section : last identifier

		// Section -> option : value
		data map[string]map[string]*tValue

		// Useful Options
		FileName  string
		EnvPrefix string
		flags     map[string]*flag.FlagSet
	}

	// tValue holds the input position for a value.
	tValue struct {
		position int    // Option order
		v        string // value
	}
)

// New creates an empty configuration representation.
// This representation can be filled with AddSection and AddOption and then
// saved to a file using WriteFile.
func New(comment, separator string, preSpace, postSpace bool, envprefix string, filename string, hasflags bool) *Config {
	if comment != DEFAULT_COMMENT && comment != ALTERNATIVE_COMMENT {
		panic("comment character not valid")
	}

	if separator != DEFAULT_SEPARATOR && separator != ALTERNATIVE_SEPARATOR {
		panic("separator character not valid")
	}

	// == Get spaces around separator
	if preSpace {
		separator = " " + separator
	}

	if postSpace {
		separator += " "
	}
	//==

	c := new(Config)

	c.comment = comment
	c.separator = separator
	c.idSection = make(map[string]int)
	c.lastIdOption = make(map[string]int)

	c.data = make(map[string]map[string]*tValue)

	c.AddSection(DEFAULT_SECTION) // Default section always exists.

	c.EnvPrefix = envprefix

	if filename != "" {
		c.FileName = filename
	} else {
		c.FileName = DEFAULT_FILENAME
	}

	if hasflags {
		c.flags = make(map[string]*flag.FlagSet)
		c.RegisterFlagSet("", flag.CommandLine)
	}

	return c
}

// NewDefault creates a configuration representation with values by default.
func NewDefault() *Config {
	return New(DEFAULT_COMMENT, DEFAULT_SEPARATOR, false, true, "", "", false)
}

// Merge merges the given configuration "source" with this one ("target").
//
// Any option (under any section) from source that is not in
// target will be copied into target. When the target already has an option with
// the same name and section then it is overwritten.
func (target *Config) Merge(source *Config) {
	if source == nil || source.data == nil || len(source.data) == 0 {
		return
	}

	for section, option := range source.data {
		for optionName, optionValue := range option {
			target.AddOption(section, optionName, optionValue.v)
		}
	}
}

func stripComments(l string) string {
	// Comments are preceded by space or TAB
	for _, c := range []string{" ;", "\t;", " #", "\t#"} {
		if i := strings.Index(l, c); i != -1 {
			l = l[0:i]
		}
	}
	return l
}
