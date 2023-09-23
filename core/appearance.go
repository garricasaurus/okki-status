package core

import (
	"regexp"
	"strings"
	"text/template"
)

// Appearance encapsulates the attributes that control how a module
// on is rendered on the Bar.
// Appearance is typically defined through the configuration file.
type Appearance struct {
	Color           *Color     `yaml:"color"`
	Border          *Border    `yaml:"border"`
	Separator       *Separator `yaml:"separator"`
	Format          string     `yaml:"format"`
	formatTmpl      *template.Template
	FormatShort     string `yaml:"format_short"`
	formatShortTmpl *template.Template
	Align           string `yaml:"align"`
	MinWidth        int    `yaml:"min_width"`
	Urgent          bool   `yaml:"urgent"`
}

// CompileTemplates parses the format and short-format templates
// contained within the Appearance.
func (a *Appearance) CompileTemplates() (err error) {
	a.formatTmpl, err = template.New("format").Parse(a.Format)
	if err != nil {
		return
	}
	a.formatShortTmpl, err = template.New("format_short").Parse(a.FormatShort)
	return err
}

// ExecuteFormat processes the format template
func (a *Appearance) ExecuteFormat(ctx Update) string {
	if a.Format == "" {
		return ""
	}
	return executeTemplate(a.formatTmpl, ctx)
}

// ExecuteFormatShort processes the short-format template
func (a *Appearance) ExecuteFormatShort(ctx Update) string {
	if a.FormatShort == "" {
		return ""
	}
	return executeTemplate(a.formatShortTmpl, ctx)
}

func executeTemplate(tmpl *template.Template, ctx Update) string {
	sb := &strings.Builder{}
	err := tmpl.Execute(sb, ctx)
	if err != nil {
		return "?"
	}
	return sb.String()
}

// Color represents a color on the Bar, and must be
// specified using the #RRGGBBAA or #RRGGBB notation.
type Color struct {
	Foreground string `yaml:"foreground"`
	Background string `yaml:"background"`
}

// Border encapsulates the border settings for a module.
type Border struct {
	Color  string `yaml:"color"`
	Top    int    `yaml:"top"`
	Bottom int    `yaml:"bottom"`
	Left   int    `yaml:"left"`
	Right  int    `yaml:"right"`
}

// Separator contains the swaybar seprarator related settings.
type Separator struct {
	Enabled    bool `yaml:"enabled"`
	BlockWidth int  `yaml:"block_width"`
}

// Variant contains an appearance variant for the module.
// The appearance Variant is only applied, when the Regex expression
// matches the actual status string of the module.
// This can be used to conditionally control the appearance of
// the modules, for example setting a different background color,
// or different format text based on the status string.
type Variant struct {
	Appearance *Appearance `yaml:"appearance"`
	re         *regexp.Regexp
	Pattern    string `yaml:"pattern"`
}

// Compile the regex pattern contained in the Variant.
func (v *Variant) Compile() error {
	if v.Pattern != "" {
		re, err := regexp.Compile(v.Pattern)
		if err != nil {
			return err
		}
		v.re = re
	}
	return nil
}

// Match indicates if s matches using the normal regex pattern.
func (v *Variant) Match(s string) bool {
	if v.re == nil {
		return false
	}
	return v.re.MatchString(s)
}
