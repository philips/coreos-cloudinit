package cloudinit

import (
	"log"

	"launchpad.net/goyaml"
)

const DefaultSSHKeyName = "coreos-cloudinit"

type CloudConfig struct {
	SSH_Authorized_Keys []string
	Coreos              struct {
		Etcd  struct{ Discovery_URL string }
		Fleet struct{ Autostart bool }
	}
	Write_Files []WriteFile
}

func NewCloudConfig(contents []byte) (*CloudConfig, error) {
	var cfg CloudConfig
	err := goyaml.Unmarshal(contents, &cfg)
	return &cfg, err
}

func (cc CloudConfig) String() string {
	bytes, err := goyaml.Marshal(cc)
	if err == nil {
		return string(bytes)
	} else {
		return ""
	}
}

func ApplyCloudConfig(cfg CloudConfig, sshKeyName string) error {
	if len(cfg.SSH_Authorized_Keys) > 0 {
		err := AuthorizeSSHKeys(sshKeyName, cfg.SSH_Authorized_Keys)
		if err == nil {
			log.Printf("Authorized SSH keys for core user")
		} else {
			return err
		}
	}

	if len(cfg.Write_Files) > 0 {
		for _, file := range cfg.Write_Files {
			if err := ProcessWriteFile("/", &file); err != nil {
				return err
			}
			log.Printf("Wrote file %s to filesystem", file.Path)
		}
	}

	if cfg.Coreos.Etcd.Discovery_URL != "" {
		err := PersistEtcdDiscoveryURL(cfg.Coreos.Etcd.Discovery_URL)
		if err == nil {
			log.Printf("Consumed etcd discovery url")
		} else {
			log.Fatalf("Failed to persist etcd discovery url to filesystem: %v", err)
		}
	}

	if cfg.Coreos.Fleet.Autostart {
		err := StartUnit("fleet.service")
		if err == nil {
			log.Printf("Started fleet service.")
		} else {
			return err
		}
	}

	return nil
}
