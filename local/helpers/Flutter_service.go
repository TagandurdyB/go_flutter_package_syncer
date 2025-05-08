package helpers

import "os/exec"

type FlutterService struct{}

func (service FlutterService) GetFlutterVersion() (output []byte, status bool) {
	cmd := exec.Command("flutter", "--version")
	output, err := cmd.Output()
	if err != nil {
		status = false
	} else {
		status = true
	}
	return
}

func (service FlutterService) FlutterPubGet() (output []byte, status bool) {
	cmd := exec.Command("flutter", "pub", "get")
	output, err := cmd.Output()
	if err != nil {
		status = false
	} else {
		status = true
	}
	return
}
