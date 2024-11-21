package config

import (
	"encoding/json"
	"fmt"
	utils "github.com/projectz-ro/journalz-ro/zro_utils"
	"os"
)

type Config struct {
	ENTRY_DIR     string `json:"ENTRY_DIR"`
	VOLUME_DIR    string `json:"VOLUME_DIR"`
	INSERT_ON_NEW bool   `json:"INSERT_ON_NEW"`
	START_POS     int    `json:"START_POS"`
}

var (
	ConfigDir      string = os.Getenv("HOME") + "/.config/journalz-ro/"
	ConfigFile     string = os.Getenv("HOME") + "/.config/journalz-ro/config.json"
	DEFAULT_CONFIG        = Config{
		ENTRY_DIR:     os.Getenv("HOME") + "/Documents/JournalZ-ro/",
		VOLUME_DIR:    os.Getenv("HOME") + "/Documents/JournalZ-ro/Volumes/",
		INSERT_ON_NEW: true,
		START_POS:     8,
	}
	CONFIG Config = DEFAULT_CONFIG

	CommandsList []string = []string{"'new'", "'find'"}
)

func LoadConfig() {
	CONFIG = DEFAULT_CONFIG

	// Create save paths
	if !utils.PathExists(CONFIG.ENTRY_DIR) {
		err := os.MkdirAll(CONFIG.ENTRY_DIR, 0755)
		if err != nil {
			fmt.Println("Error creating entry directory", err)
			return
		}
	}
	if !utils.PathExists(CONFIG.VOLUME_DIR) {
		err := os.MkdirAll(CONFIG.VOLUME_DIR, 0755)
		if err != nil {
			fmt.Println("Error creating volume directory", err)
			return
		}
	}

	// Create config path and templates
	if !utils.PathExists(ConfigDir) {
		err := os.MkdirAll(ConfigDir, 0755)
		if err != nil {
			fmt.Println("Error creating entry directory", err)
			return
		}

		return
	}

	if !utils.PathExists(ConfigFile) {
		return
	}

	data, err := os.ReadFile(ConfigFile)
	if err != nil {
		fmt.Println("Error reading file: config.json. Continuing with default configuration")
		return
	} else {
		err = json.Unmarshal(data, &CONFIG)
		if err != nil {
			fmt.Println("Error loading config. Continuing with default configuration")
			CONFIG = DEFAULT_CONFIG
		}
	}
}
