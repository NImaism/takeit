package runner

import (
	"fmt"
	"github.com/logrusorgru/aurora"
)

func (r *Runner) printScanResult(target, discussion, documentation string, isVulnerable bool) {
	if isVulnerable {
		if r.Option.Silent {
			fmt.Println(target)
		}
		r.printMessage(fmt.Sprintf("[+] %s is vulnerable | %s | %s", target, discussion, documentation), "", true)
	} else {
		if !r.Option.Silent {
			r.printMessage(fmt.Sprintf("[-] %s is not vulnerable", target), "", false)
		}
	}
}

func (r *Runner) printMessage(message, target string, isPositive bool) {
	if r.Option.Silent {
		return
	}
	if r.Option.NoColors {
		if target != "" {
			fmt.Printf("[+] %s: %s\n", message, target)
		} else {
			fmt.Println(message)
		}
	} else {
		colorFunc := aurora.Red
		if isPositive {
			colorFunc = aurora.Green
		}

		if target != "" {
			fmt.Println(colorFunc(fmt.Sprintf("[+] %s: %s", message, target)))
		} else {
			fmt.Println(colorFunc(message))
		}
	}
}
