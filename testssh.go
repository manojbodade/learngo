package main

import (
	"strings"

	tk "github.com/eaciit/toolkit"
	. "github.com/ekobudy/sshclient"
)

func main() {
	var SshClient SshSetting

	SshClient.SSHAuthType = SSHAuthType_NoPassword
	//	SshClient.SSHHost = "52.74.15.96:22"
	//	SshClient.SSHHost = "52.76.247.30:22"
	SshClient.SSHHost = "10.20.172.190:22"
	//	SshClient.SSHUser = "developer"
	//	SshClient.SSHUser = "ec2-user"
	SshClient.SSHUser = "eaciit"

	//	SshClient.SSHKeyLocation = "/Users/mazte/pubkey/devgoeaciit.pem"
	//	SshClient.SSHKeyLocation = "/Users/mazte/pubkey/clv-key.pem"
	//	SshClient.SSHPassword = "Bismillah"
	//	shc := []string{"cat /proc/version", "cat /etc/system-release"}
	shc := []string{"df -P"}
	resp, err := SshClient.RunCommandSshAsMap(shc...)
	if err != nil {
		tk.Println("error ", err.Error())
	}
	//	tk.Printfn("RESULT %v ", resp)
	if len(resp) > 0 {
		str := resp[0].Output
		strline := strings.Split(str, "\r\n")
		for cx0, strls := range strline {
			strix := strings.Join(strings.Fields(strls), "`")
			if cx0 > 0 && !strings.Contains(strix, "tmpfs") {
				tk.Println("line ", cx0, "string ", strix, ";BY \t ", len(strings.Split(strix, "`")))
				//      tk.Println()
			}
		}
	}
}
