package triago

type (
	SectionError string
	OptionError  string
)

func (e SectionError) Error() string {
	return "section not found: " + string(e)
}

func (e OptionError) Error() string {
	return "option not found: " + string(e)
}
