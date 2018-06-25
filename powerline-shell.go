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
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/sivel/powerline-shell-go/powerline"
)

type Colors struct {
	HomeFG           string
	HomeBG           string
	CwdFG            string
	CwdSeparatorThin string
	CwdBG            string
	VirtualEnvFG     string
	VirtualEnvBG     string
	LockFG           string
	LockBG           string
	GitFG            string
	GitBG            string
	GitStagedFG      string
	GitStagedBG      string
	PromptFG         string
	PromptBG         string
}

func NewColors(pallet string) Colors {
	if pallet == "light" {
		return Colors{
			HomeFG:           "000",
			HomeBG:           "166",
			CwdFG:            "237",
			CwdSeparatorThin: "244",
			CwdBG:            "250",
			VirtualEnvFG:     "015",
			VirtualEnvBG:     "161",
			LockFG:           "254",
			LockBG:           "080",
			GitFG:            "015",
			GitBG:            "055",
			GitStagedFG:      "000",
			GitStagedBG:      "035",
			PromptFG:         "000",
			PromptBG:         "250",
		}
	}

	return Colors{
		HomeFG:           "015",
		HomeBG:           "031",
		CwdFG:            "250",
		CwdSeparatorThin: "244",
		CwdBG:            "237",
		VirtualEnvFG:     "000",
		VirtualEnvBG:     "035",
		LockFG:           "254",
		LockBG:           "124",
		GitFG:            "000",
		GitBG:            "148",
		GitStagedFG:      "015",
		GitStagedBG:      "161",
		PromptFG:         "015",
		PromptBG:         "236",
	}
}

func getCurrentWorkingDir() (string, []string) {
	dir, err := filepath.Abs(".")
	if err != nil {
		log.Fatal(err)
	}
	userDir := strings.Replace(dir, os.Getenv("HOME"), "~", 1)
	userDir = strings.TrimSuffix(userDir, "/")
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
	reBranch := regexp.MustCompile(`^(HEAD detached at|HEAD detached from|On branch) (\S+)`)
	matchBranch := reBranch.FindStringSubmatch(string(stdout))
	if len(matchBranch) > 0 {
		if matchBranch[2] == "detached" {
			status = matchBranch[2]
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

func addCwd(colors Colors, cwdParts []string, ellipsis string, separator string) [][]string {
	segments := [][]string{}
	home := false
	if cwdParts[0] == "~" {
		cwdParts = cwdParts[1:len(cwdParts)]
		home = true
	}

	if home {
		segments = append(segments, []string{colors.HomeFG, colors.HomeBG, "~"})

		if len(cwdParts) > 2 {
			segments = append(segments, []string{colors.CwdFG, colors.CwdBG, cwdParts[0], separator, colors.CwdSeparatorThin})
			segments = append(segments, []string{colors.CwdFG, colors.CwdBG, ellipsis, separator, colors.CwdSeparatorThin})
		} else if len(cwdParts) == 2 {
			segments = append(segments, []string{colors.CwdFG, colors.CwdBG, cwdParts[0], separator, colors.CwdSeparatorThin})
		}
	} else {
		if len(cwdParts[len(cwdParts)-1]) == 0 {
			segments = append(segments, []string{colors.CwdFG, colors.CwdBG, "/"})
		}

		if len(cwdParts) > 3 {
			segments = append(segments, []string{colors.CwdFG, colors.CwdBG, cwdParts[1], separator, colors.CwdSeparatorThin})
			segments = append(segments, []string{colors.CwdFG, colors.CwdBG, ellipsis, separator, colors.CwdSeparatorThin})
		} else if len(cwdParts) > 2 {
			segments = append(segments, []string{colors.CwdFG, colors.CwdBG, cwdParts[1], separator, colors.CwdSeparatorThin})
		}
	}

	if len(cwdParts) != 0 && len(cwdParts[len(cwdParts)-1]) > 0 {
		segments = append(segments, []string{colors.CwdFG, colors.CwdBG, cwdParts[len(cwdParts)-1]})
	}

	return segments
}

func addVirtulEnvName(colors Colors) []string {
	_, _, virtualEnvName := getVirtualEnv()
	if virtualEnvName != "" {
		return []string{colors.VirtualEnvFG, colors.VirtualEnvBG, virtualEnvName}
	}

	return nil
}

func addLock(colors Colors, cwd string, lock string) []string {
	if !isWritableDir(cwd) {
		return []string{colors.LockFG, colors.LockBG, lock}
	}

	return nil
}

func addGitInfo(colors Colors) []string {
	gitStatus, gitStaged := getGitInformation()
	if gitStatus != "" {
		if gitStaged {
			return []string{colors.GitStagedFG, colors.GitStagedBG, gitStatus}
		} else {
			return []string{colors.GitFG, colors.GitBG, gitStatus}
		}
	} else {
		return nil
	}
}

func addDollarPrompt(colors Colors) []string {
	return []string{colors.PromptFG, colors.PromptBG, "\\$"}
}

func main() {
	shell := "bash"
	pallet := "dark"

	if len(os.Args) > 1 {
		shell = os.Args[1]
		if len(os.Args) == 3 {
			pallet = os.Args[2]
		}
	}

	p := powerline.NewPowerline(shell)
	cwd, cwdParts := getCurrentWorkingDir()

	colors := NewColors(pallet)

	p.AppendSegment(addVirtulEnvName(colors))
	p.AppendSegments(addCwd(colors, cwdParts, p.Ellipsis, p.SeparatorThin))
	p.AppendSegment(addLock(colors, cwd, p.Lock))
	p.AppendSegment(addGitInfo(colors))
	p.AppendSegment(addDollarPrompt(colors))

	fmt.Print(p.PrintSegments())
}
