package typictx

import "gopkg.in/urfave/cli.v1"

// Context of typical application
type Context struct {
	TypiApp

	Name        string
	Version     string
	Description string

	ReadmeTemplate string
	ReadmeFile     string

	Commands []cli.Command

	Modules []*Module
}

// ReadmeTemplateOrDefault return readme template field or the default template
func (c Context) ReadmeTemplateOrDefault() string {
	if c.ReadmeTemplate == "" {
		return defaultReadmeTemplate
	}
	return c.ReadmeTemplate
}

// ReadmeFileOrDefault return readme file field or default template
func (c Context) ReadmeFileOrDefault() string {
	if c.ReadmeFile == "" {
		return defaultReadmeFile
	}

	return c.ReadmeFile
}
