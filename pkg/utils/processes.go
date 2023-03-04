package utils

import (
	"bytes"
	"log"
	"os/exec"
)

func getStdoutString(out bytes.Buffer) string {
	// Try to read output up to a newline character
	bytesOut, err := out.ReadBytes(10)
	if err != nil {
		log.Println("Error: reading command exec output,", err)
		return ""
	}
	outString := string(bytesOut[0 : len(bytesOut)-1])
	return outString
}

func runCommand(out *bytes.Buffer, cmdname string, cmdargs []string) error {
	cmd := exec.Command(cmdname, cmdargs...)
	cmd.Stdout = out
	err := cmd.Run()
	return err
}

func CheckProcessIsActive(process string) bool {
	var out bytes.Buffer
	cmdargs := []string{"is-active", process}
	err := runCommand(&out, "systemctl", cmdargs)

	if err != nil {
		log.Println("Error: couldn't run systemctl command,", err)
		return false
	}

	if getStdoutString(out) == "active" {
		return true
	}

	return false
}

func HandleProcess(process string, action string) bool {
	log.Printf("Trying to %s %s...\n", action, process)
	if action != "restart" && action != "stop" && action != "start" {
		log.Printf("'%s' is not a valid action\n", action)
		return false
	}

	var out bytes.Buffer
	cmdargs := []string{action, process}
	err := runCommand(&out, "systemctl", cmdargs)

	if err != nil {
		log.Printf("Error: systemctl %s failed: %s\n", action, err)
		return false
	}
	return true
}
