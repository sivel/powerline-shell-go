// Copyright 2014 Matt Martz <matt@sivel.net>
// All Rights Reserved.
//
//    Licensed under the Apache License, Version 2.0 (the "License"); you may
//    not use this file except in compliance with the License. You may obtain
//    a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//    WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//    License for the specific language governing permissions and limitations
//    under the License.

package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

func getCurrentWorkingDir() (string, []string) {
	dir, err := filepath.Abs(".")
	if err != nil {
		log.Fatal(err)
	}
	userDir := strings.Replace(dir, os.Getenv("HOME"), "~", 1)
	parts := strings.Split(userDir, "/")
	return dir, parts
}

func getVirtualEnv() (string, []string, string) {
	var parts []string
	virtualEnv := os.Getenv("VIRTUAL_ENV")
	if virtualEnv == "" {
		return "", parts, ""
	}

	parts = strings.Split(virtualEnv, "/")

	virtualEnvName := path.Base(virtualEnv)
	return virtualEnv, parts, virtualEnvName
}

func isWritableDir(dir string) bool {
	tmpPath := path.Join(dir, ".powerline-write-test")
	_, err := os.Create(tmpPath)
	if err != nil {
		return false
	}
	os.Remove(tmpPath)
	return true
}

func getGitInformation() (string, bool) {
	var status string
	var staged bool
	stdout, _ := exec.Command("git", "status", "--ignore-submodules").Output()
	reBranch := regexp.MustCompile(`^(HEAD|On branch) (\S+)`)
	matchBranch := reBranch.FindStringSubmatch(string(stdout))
	if len(matchBranch) > 0 {
		if matchBranch[2] == "detached" {
			status = "(Detached)"
		} else {
			status = matchBranch[2]
		}
	}

	reStatus := regexp.MustCompile(`Your branch is (ahead|behind).*?([0-9]+) comm`)
	matchStatus := reStatus.FindStringSubmatch(string(stdout))
	if len(matchStatus) > 0 {
		status = fmt.Sprintf("%s %s", status, matchStatus[2])
		if matchStatus[1] == "behind" {
			status = fmt.Sprintf("%s\u21E3", status)
		} else if matchStatus[1] == "ahead" {
			status = fmt.Sprintf("%s\u21E1", status)
		}
	}

	staged = !strings.Contains(string(stdout), "nothing to commit")
	if strings.Contains(string(stdout), "Untracked files") {
		status = fmt.Sprintf("%s +", status)
	}

	return status, staged
}

type Powerline struct {
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

func (p *Powerline) Segment(content string, fg string, bg string) string {
	foreground := fmt.Sprintf(p.BashTemplate, fmt.Sprintf(p.ColorTemplate, "38", fg))
	background := fmt.Sprintf(p.BashTemplate, fmt.Sprintf(p.ColorTemplate, "48", bg))
	return fmt.Sprintf("%s%s %s", foreground, background, content)
}

func (p *Powerline) Color(prefix string, code string) string {
	return fmt.Sprintf(p.BashTemplate, fmt.Sprintf(p.ColorTemplate, prefix, code))
}

func (p *Powerline) ForegroundColor(code string) string {
	return p.Color("38", code)
}

func (p *Powerline) BackgroundColor(code string) string {
	return p.Color("48", code)
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

func main() {
	home := false
	p := Powerline{
		BashTemplate:  "\\[\\e%s\\]",
		ColorTemplate: "[%s;5;%sm",
		Reset:         "\\[\\e[0m\\]",
		Lock:          "\uE0A2",
		Network:       "\uE0A2",
		Separator:     "\uE0B0",
		SeparatorThin: "\uE0B1",
		Ellipsis:      "\u2026",
	}
	cwd, cwdParts := getCurrentWorkingDir()
	_, _, virtualEnvName := getVirtualEnv()
	if virtualEnvName != "" {
		p.Segments = append(p.Segments, []string{"00", "35", virtualEnvName})
	}
	if cwdParts[0] == "~" {
		cwdParts = cwdParts[1:len(cwdParts)]
		p.Segments = append(p.Segments, []string{"15", "31", "~"})
		home = true
	}
	if len(cwdParts) >= 4 {
		p.Segments = append(p.Segments, []string{"250", "237", cwdParts[1], p.SeparatorThin, "244"})
		p.Segments = append(p.Segments, []string{"250", "237", p.Ellipsis, p.SeparatorThin, "244"})
		p.Segments = append(p.Segments, []string{"254", "237", cwdParts[len(cwdParts)-1]})
	} else if len(cwdParts) == 3 {
		if home {
			p.Segments = append(p.Segments, []string{"250", "237", cwdParts[0], p.SeparatorThin, "244"})
		} else {
			p.Segments = append(p.Segments, []string{"250", "237", cwdParts[1], p.SeparatorThin, "244"})
		}
		p.Segments = append(p.Segments, []string{"250", "237", p.Ellipsis, p.SeparatorThin, "244"})
		p.Segments = append(p.Segments, []string{"254", "237", cwdParts[len(cwdParts)-1]})
	} else if len(cwdParts) != 0 {
		p.Segments = append(p.Segments, []string{"254", "237", cwdParts[len(cwdParts)-1]})
	}

	if !isWritableDir(cwd) {
		p.Segments = append(p.Segments, []string{"254", "124", p.Lock})
	}

	gitStatus, gitStaged := getGitInformation()
	if gitStatus != "" {
		if gitStaged {
			p.Segments = append(p.Segments, []string{"15", "161", gitStatus})
		} else {
			p.Segments = append(p.Segments, []string{"0", "148", gitStatus})
		}
	}

	p.Segments = append(p.Segments, []string{"15", "236", "\\$"})

	fmt.Print(p.PrintSegments())
}
