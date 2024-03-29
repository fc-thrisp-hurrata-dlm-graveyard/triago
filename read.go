package triago

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strings"
)

// _read is the base to read a file and get the configuration representation.
// That representation can be queried with GetString, etc.
func _read(fname string, c *Config) (*Config, error) {
	file, err := os.Open(fname)
	if err != nil {
		return nil, err
	}

	if err = c.read(bufio.NewReader(file)); err != nil {
		return nil, err
	}

	if err = file.Close(); err != nil {
		return nil, err
	}

	return c, nil
}

// Read reads a configuration file and returns its representation.
// All arguments, except `fname`, are related to `New()`
func Read(fname string, comment, separator string, preSpace, postSpace bool, envprefix string, hasflags bool) (*Config, error) {
	return _read(fname, New(comment, separator, preSpace, postSpace, envprefix, fname, hasflags))
}

// ReadDefault reads a configuration file and returns its representation.
// It uses values by default.
func ReadDefault(fname string) (*Config, error) {
	return _read(fname, NewDefault())
}

func (c *Config) read(buf *bufio.Reader) (err error) {
	var section, option string

	for {
		l, err := buf.ReadString('\n') // parse line-by-line
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		l = strings.TrimSpace(l)

		// Switch written for readability (not performance)
		switch {
		// Empty line and comments
		case len(l) == 0, l[0] == '#', l[0] == ';':
			continue

		// New section
		case l[0] == '[' && l[len(l)-1] == ']':
			option = "" // reset multi-line value
			section = strings.TrimSpace(l[1 : len(l)-1])
			c.AddSection(section)

		// No new section and no section defined so
		//case section == "":
		//return os.NewError("no section defined")

		// Other alternatives
		default:
			i := strings.IndexAny(l, "=:")

			switch {
			// Option and value
			case i > 0:
				i := strings.IndexAny(l, "=:")
				option = strings.TrimSpace(l[0:i])
				value := strings.TrimSpace(stripComments(l[i+1:]))
				c.AddOption(section, option, value)
			// Continuation of multi-line value
			case section != "" && option != "":
				prev, _ := c.RawString(section, option)
				value := strings.TrimSpace(stripComments(l))
				c.AddOption(section, option, prev+"\n"+value)
			default:
				return errors.New("could not parse line: " + l)
			}
		}
	}
	return nil
}
