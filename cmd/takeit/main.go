package main

import (
	"fmt"
	"github.com/logrusorgru/aurora"
	"github.com/nimaism/takeit/internal/options"
	"github.com/nimaism/takeit/internal/pattern"
	"github.com/nimaism/takeit/internal/runner"
	"github.com/nimaism/takeit/internal/verison"
)

func main() {
	option := options.Options{}
	if err := option.ParseFlags(); err != nil {
		fmt.Println(aurora.Red(fmt.Sprintf("error parsing flags: %v", err)))
		return
	}

	if !option.Silent {
		version.ShowVersion()
	}

	if len(option.Targets) < 1 {
		fmt.Println(aurora.Red("Input is empty!"))
		return
	}

	if !option.DisableUpdateCheck {
		version.CheckLatestVersion()
		if err := pattern.UpdatePatterns(option.ExcludePatterns); err != nil {
			fmt.Println(aurora.Red(fmt.Sprintf("error updating patterns: %v", err)))
			return
		}
	}

	engine, err := runner.NewRunner(&option)
	if err != nil {
		fmt.Println(aurora.Red(fmt.Sprintf("error creating runner: %v", err)))
		return
	}

	engine.Run()
}
