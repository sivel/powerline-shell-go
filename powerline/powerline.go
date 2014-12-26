package powerline

import (
	"bytes"
	"fmt"
)

type Powerline struct {
	ShTemplate    string
	BashTemplate  string
	ColorTemplate string
	Reset         string
	Lock          string
	Network       string
	Separator     string
	SeparatorThin string
	Ellipsis      string
	Segments      [][]string
}

func (p *Powerline) Color(prefix string, code string) string {
	return fmt.Sprintf(
		p.ShTemplate,
		fmt.Sprintf(p.ColorTemplate, prefix, code),
	)
}

func (p *Powerline) ForegroundColor(code string) string {
	return p.Color("38", code)
}

func (p *Powerline) BackgroundColor(code string) string {
	return p.Color("48", code)
}

func (p *Powerline) AppendSegment(segment []string) {
	if segment != nil {
		p.Segments = append(p.Segments, segment)
	}
}

func (p *Powerline) AppendSegments(segments [][]string) {
	for _, segment := range segments {
		p.AppendSegment(segment)
	}
}

func (p *Powerline) PrintSegments() string {
	var nextBackground string
	var buffer bytes.Buffer
	for i, Segment := range p.Segments {
		if (i + 1) == len(p.Segments) {
			nextBackground = p.Reset
		} else {
			nextBackground = p.BackgroundColor(p.Segments[i+1][1])
		}
		if len(Segment) == 3 {
			buffer.WriteString(fmt.Sprintf("%s%s %s %s%s%s", p.ForegroundColor(Segment[0]), p.BackgroundColor(Segment[1]), Segment[2], nextBackground, p.ForegroundColor(Segment[1]), p.Separator))
		} else {
			buffer.WriteString(fmt.Sprintf("%s%s %s %s%s%s", p.ForegroundColor(Segment[0]), p.BackgroundColor(Segment[1]), Segment[2], nextBackground, p.ForegroundColor(Segment[4]), Segment[3]))
		}
	}

	buffer.WriteString(p.Reset)

	return buffer.String()
}

func NewPowerline(shell string) Powerline {
	p := Powerline{
		Lock:          "\uE0A2",
		Network:       "\uE0A2",
		Separator:     "\uE0B0",
		SeparatorThin: "\uE0B1",
		Ellipsis:      "\u2026",
	}

	switch shell {
	case "bash":
		p.ShTemplate = "\\[\\e%s\\]"
		p.ColorTemplate = "[%s;5;%sm"
		p.Reset = "\\[\\e[0m\\]"

	case "zsh":
		p.ShTemplate = "%s"
		p.ColorTemplate = "%%{[%s;5;%sm%%}"
		p.Reset = "%{$reset_color%}"
	}
	return p
}
