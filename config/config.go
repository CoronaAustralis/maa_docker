package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
)

type ProfilesStruct struct {
	Connection struct {
		AdbPath string `toml:"adb_path"`
		Device  string `toml:"device"`
		Config  string `toml:"config"`
	} `toml:"connection"`

	Resource struct {
		GlobalResource       string `toml:"global_resource"`
		PlatformDiffResource string `toml:"platform_diff_resource"`
		UserResource         bool   `toml:"user_resource"`
	} `toml:"resource"`

	StaticOptions struct {
		CpuOcr bool `toml:"cpu_ocr"`
	} `toml:"static_options"`

	InstanceOptions struct {
		TouchMode           string `toml:"touch_mode"`
		DeploymentWithPause bool   `toml:"deployment_with_pause"`
		AdbLiteEnabled      bool   `toml:"adb_lite_enabled"`
		KillAdbOnExit       bool   `toml:"kill_adb_on_exit"`
	} `toml:"instance_options"`
}

// TaskCluster represents a cluster of tasks
type TaskCluster struct {
	Hash     string    `json:"hash"`
	IsEnable bool      `json:"isEnable"`
	Type     string    `json:"type"`
	Alias    string    `json:"alias"`
	Time     time.Time `json:"time"`
	Tasks    []string  `json:"tasks"`
}

type DStruct struct {
	ExecuteDir  string
	HomeDir     string
	InfrastDir	string
	ProfilesDir string
	FightDir    string
}

// Root represents the root JSON structure
type Config struct {
	TaskCluster     map[string]TaskCluster `json:"task_cluster"`
	TemplateCluster TaskCluster            `json:"template_cluster"`
}

var Conf = &Config{TaskCluster: make(map[string]TaskCluster)}
var D = &DStruct{}
var Profiles = &ProfilesStruct{}

func init() {
	maa_dev := os.Getenv("MAA_DEV")
	if maa_dev != "" {
		D.ExecuteDir = maa_dev
	} else {
		path, err := os.Executable()
		if err != nil {
			log.Panicln("Error:", err)
		}
		D.ExecuteDir = filepath.Dir(path)
	}

	configPath := filepath.Join(D.ExecuteDir, "./config/config.json")
	configBytes, err := os.ReadFile(configPath)
	if err != nil {
		log.Panicln("load config error: ", err)
	}
	err = json.Unmarshal(configBytes, Conf)
	if err != nil {
		log.Panicln("unmarshal config error: ", err)
	}
	log.Println(Conf)

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Panicln("Error getting home directory:", err)
	}

	D.HomeDir = homeDir

	D.ProfilesDir = filepath.Join(homeDir, ".config", "maa", "profiles")
	D.FightDir = filepath.Join(homeDir, ".config", "maa", "tasks")
	D.InfrastDir = filepath.Join(homeDir, ".config", "maa", "infrast")

	if _, err := toml.DecodeFile(filepath.Join(D.ProfilesDir, "default.toml"), Profiles); err != nil {
		log.Fatalf("Error decoding TOML file: %v", err)
	}
}

func UpdateProfile() {
	profilePath := filepath.Join(D.ProfilesDir, "default.toml")

	if v, err := toml.Marshal(Profiles); err != nil {
		log.Println("serialize failed, err: ", err)
	} else {
		err := os.WriteFile(profilePath, v, 0777)
		if err != nil {
			log.Println("serialize failed, err: ", err)
		}
	}
}

func UpdateConfig() {
	configPath := filepath.Join(D.ExecuteDir, "./config/config.json")

	if v, err := json.MarshalIndent(Conf, "", "  "); err != nil {
		log.Println("serialize failed, err: ", err)
	} else {
		err := os.WriteFile(configPath, v, 0777)
		if err != nil {
			log.Println("serialize failed, err: ", err)
		}
	}
}
