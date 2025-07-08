package main

import (
	"fmt"
	"os"
	"os/exec"
)

var possibleVersions = []string{
	"php7.3",
	"php8.2",
	"php8.4",
}

const (
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Reset  = "\033[0m"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println(Yellow + "Usage: phpchange [version]" + Reset)
		return
	}
	var enableVersion = os.Args[1]

	if len(possibleVersions) == 0 {
		fmt.Println(Yellow + "version should be meaningful php version. For example: php7.3, php8.2, php8.4" + Reset)
		return
	}

	enableVersion = "php" + enableVersion

	// Disable all possible PHP versions
	for _, version := range possibleVersions {
		cmd := exec.Command("a2dismod", version)
		_, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println(Red+"PHP version", version, "does not exist."+Reset)
		} else {
			fmt.Println(Green + "PHP version was successfully disabled: " + version + Reset)
		}
	}

	// Enable the specified PHP version
	cmsEnable := exec.Command("a2enmod", enableVersion)

	_, err := cmsEnable.CombinedOutput()
	if err != nil {
		fmt.Println(Red+"PHP version", enableVersion, "does not exist."+Reset)
		return
	} else {
		fmt.Println(Green + "New version " + enableVersion + " was successfully enabled" + Reset)
	}

	// update alternatives for PHP CLI
	cmdAlt := exec.Command("update-alternatives", "--set", "php", "/usr/bin/"+enableVersion)
	_, err = cmdAlt.CombinedOutput()
	if err != nil {
		fmt.Println(Red+"Failed to update alternatives for PHP CLI:", err, Red)
		return
	} else {
		fmt.Println(Green+"PHP CLI updated to version:", enableVersion+Reset)
	}

	// Restart Apache to apply changes
	cmd := exec.Command("systemctl", "restart", "apache2")
	_, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println(Red+"Failed to restart Apache:", err, Red)
		return
	} else {
		fmt.Println(Green + "Apache restarted successfully." + Reset)
	}

	fmt.Println(Green+"PHP version changed to:", enableVersion+Reset)
}
