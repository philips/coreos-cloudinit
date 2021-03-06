package cloudinit

import (
	"fmt"
	"log"
	"path"

	"github.com/coreos/go-systemd/dbus"
)

type Script []byte

func StartUnit(name string) error {
	conn, err := dbus.New()
	if err != nil {
		return err
	}

	_, err = conn.StartUnit(name, "replace")
	return err
}

func ExecuteScript(scriptPath string) (string, error) {
	props := []dbus.Property{
		dbus.PropDescription("Unit generated and executed by coreos-cloudinit on behalf of user"),
		dbus.PropExecStart([]string{"/bin/bash", scriptPath}, false),
	}

	base := path.Base(scriptPath)
	name := fmt.Sprintf("coreos-cloudinit-%s.service", base)

	log.Printf("Creating transient systemd unit '%s'", name)

	conn, err := dbus.New()
	if err != nil {
		return "", err
	}

	_, err = conn.StartTransientUnit(name, "replace", props...)
	return name, err
}
