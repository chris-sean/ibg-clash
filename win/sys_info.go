package win

import (
	"github.com/Dreamacro/clash/log"
	"golang.org/x/sys/windows"

	"os/exec"
	"regexp"
	"strings"
)

var MachineId string
var uuidRegexp *regexp.Regexp

func init() {
	sysId, err := SystemUUID()
	if err != nil {
		log.Errorln("get sys id err: %v", err)
	}

	MachineId = strings.ToLower(sysId)
	log.Infoln("machine_id: %v", MachineId)
}

// SystemUUID returns uuid by running command "wmic csproduct get uuid".
func SystemUUID() (uuid string, err error) {
	if uuidRegexp == nil {
		uuidRegexp, err = regexp.Compile("[0-9A-Z]{8}-([0-9A-Z]{4}-){3}[0-9A-Z]{12}")
	}

	sysDir, err := windows.GetSystemDirectory()
	if err != nil {
		return
	}

	cmd := exec.Command(sysDir+`\wbem\wmic`, "csproduct", "get", "UUID")
	output, err := cmd.Output()
	if err != nil {
		return
	}

	str := string(output)
	lines := strings.Split(str, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if uuidRegexp.MatchString(line) {
			return line, nil
		}
	}
	return
}
