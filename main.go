// Copyright 2019 Sonatype Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	"github.com/sonatype-nexus-community/cheque/config"
	"github.com/sonatype-nexus-community/cheque/linker"
	"github.com/sonatype-nexus-community/cheque/logger"
	"os"
	"os/exec"
	"fmt"
)

func main() {
	args := []string{}

	// Remove cheque custom arguments
	for _, arg := range os.Args[1:] {
		switch arg {
		case "-Werror=cheque":
			default: args = append(args, arg)
    }
	}

	count := linker.DoLink(args)
	if count > 0 {
		if config.ExitWithError() {
			fmt.Fprintf(os.Stderr, "Error: Vulnerable dependencies found: %v\n", count)
			os.Exit(count)
		} else {
			fmt.Fprintf(os.Stderr, "Warning: Vulnerable dependencies found: %v\n", count)
		}
	}

	switch config.GetCommand() {
	case "cheque":
		break;
	default:
		var cmdPath = fmt.Sprint("/usr/bin/", config.GetCommand());

		_, err := os.Stat(cmdPath)
		if err != nil {
			logger.Fatal("Cannot find official command: " + cmdPath)
		} else {
			externalCmd := exec.Command(cmdPath, args...)

			externalCmd.Stdout = os.Stdout
			externalCmd.Stderr = os.Stderr

			if err := externalCmd.Run(); err != nil {
				if exitError, ok := err.(*exec.ExitError); ok {
					logger.Fatal(fmt.Sprintf("There was an issue running the command %s, and the issue is %s", config.GetCommand(), os.Stderr))
					os.Exit(exitError.ExitCode())
				}
			}
		}
		break;
	}

	os.Exit(0)
}
