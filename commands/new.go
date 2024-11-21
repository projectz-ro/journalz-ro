package commands

import (
	"fmt"
	"time"

	"github.com/projectz-ro/journalz-ro/config"
	"github.com/projectz-ro/journalz-ro/db"
	utils "github.com/projectz-ro/journalz-ro/zro_utils"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new [tags]",
	Short: "Create a new entry",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := createEntry(args)
		if err != nil {
			return fmt.Errorf("Error creating entry: %v", err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}

func createEntry(args []string) error {

	currentDate := time.Now().Format("2006-01-02_15-04-05")
	day := time.Now().Format("Monday")
	date := time.Now().Format("01/02/2006")
	time := time.Now().Format("15:04")

	maxWidth := 80

	title := fmt.Sprintf("Entry_%s.md", currentDate)
	filepath := config.CONFIG.ENTRY_DIR + title

	_, err := db.USERDB.InsertEntry(
		title,
		args,
		nil,
		filepath,
	)
	if err != nil {
		return fmt.Errorf("Error adding new entry: %v", err)
	}

	lines := []string{
		"",
		fmt.Sprintf("%*s", maxWidth, day),
		fmt.Sprintf("%*s", maxWidth, date),
		fmt.Sprintf("%*s", maxWidth, time),
		"",
		"---",
		"",
		"",
	}

	writeErr := utils.WriteLines(filepath, lines)
	if writeErr != nil {
		return fmt.Errorf("Error writing new file: %v", writeErr)
	}

	fmt.Println("Created new entry:", filepath)
	utils.OpenInNvim(filepath, config.CONFIG.START_POS, true)

	return nil
}
