package version

import (
	"encoding/json"
	"fmt"
	"github.com/logrusorgru/aurora"
	"github.com/nimaism/takeit/internal/model"
	"log"
	"net/http"
)

const gitHubAPIURL = "https://api.github.com/repos/nimaism/takeit/releases/latest"

const currentVersion = "v0.1.1"

const banner = `
████████╗ █████╗ ██╗  ██╗███████╗██╗████████╗
╚══██╔══╝██╔══██╗██║ ██╔╝██╔════╝██║╚══██╔══╝
   ██║   ███████║█████╔╝ █████╗  ██║   ██║   
   ██║   ██╔══██║██╔═██╗ ██╔══╝  ██║   ██║   
   ██║   ██║  ██║██║  ██╗███████╗██║   ██║   
   ╚═╝   ╚═╝  ╚═╝╚═╝  ╚═╝╚══════╝╚═╝   ╚═╝
     		 ` + currentVersion + `
`

const newVersionBanner = `
███╗   ██╗███████╗██╗    ██╗    ██╗   ██╗██████╗ ██████╗  █████╗ ████████╗███████╗     █████╗ ██╗   ██╗ █████╗ ██╗██╗      █████╗ ██████╗ ██╗     ███████╗
████╗  ██║██╔════╝██║    ██║    ██║   ██║██╔══██╗██╔══██╗██╔══██╗╚══██╔══╝██╔════╝    ██╔══██╗██║   ██║██╔══██╗██║██║     ██╔══██╗██╔══██╗██║     ██╔════╝
██╔██╗ ██║█████╗  ██║ █╗ ██║    ██║   ██║██████╔╝██║  ██║███████║   ██║   █████╗      ███████║██║   ██║███████║██║██║     ███████║██████╔╝██║     █████╗  
██║╚██╗██║██╔══╝  ██║███╗██║    ██║   ██║██╔═══╝ ██║  ██║██╔══██║   ██║   ██╔══╝      ██╔══██║╚██╗ ██╔╝██╔══██║██║██║     ██╔══██║██╔══██╗██║     ██╔══╝  
██║ ╚████║███████╗╚███╔███╔╝    ╚██████╔╝██║     ██████╔╝██║  ██║   ██║   ███████╗    ██║  ██║ ╚████╔╝ ██║  ██║██║███████╗██║  ██║██████╔╝███████╗███████╗
╚═╝  ╚═══╝╚══════╝ ╚══╝╚══╝      ╚═════╝ ╚═╝     ╚═════╝ ╚═╝  ╚═╝   ╚═╝   ╚══════╝    ╚═╝  ╚═╝  ╚═══╝  ╚═╝  ╚═╝╚═╝╚══════╝╚═╝  ╚═╝╚═════╝ ╚══════╝╚══════╝
 							go install github.com/nimaism/takeit@latest
`

func ShowVersion() {
	fmt.Println(aurora.Bold(banner))
}

func CheckLatestVersion() {
	resp, err := http.Get(gitHubAPIURL)
	if err != nil {
		log.Fatalf("Error fetching latest version: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Unexpected status code: %d", resp.StatusCode)
	}

	var release model.Release
	err = json.NewDecoder(resp.Body).Decode(&release)
	if err != nil {
		log.Fatalf("Error decoding JSON response: %v", err)
	}

	latestVersion := release.TagName

	if latestVersion != currentVersion {
		fmt.Println(aurora.Blue(newVersionBanner))
	}
}
