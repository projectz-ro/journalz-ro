package ui

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/projectz-ro/journalz-ro/config"
	"github.com/projectz-ro/journalz-ro/db"
	utils "github.com/projectz-ro/journalz-ro/zro_utils"
)

type DisplayMode int

const (
	SearchDisplay DisplayMode = iota
	MergeDisplay
)

var (
	CurrentDisplay DisplayMode = SearchDisplay
	CurrentEntries []db.Entry
)

// Colors
const (
	Reset         = "\033[0m"
	Black         = "\033[30m"
	Red           = "\033[31m"
	Green         = "\033[32m"
	Yellow        = "\033[33m"
	Blue          = "\033[34m"
	Magenta       = "\033[35m"
	Cyan          = "\033[36m"
	White         = "\033[37m"
	BrightBlack   = "\033[90m"
	BrightRed     = "\033[91m"
	BrightGreen   = "\033[92m"
	BrightYellow  = "\033[93m"
	BrightBlue    = "\033[94m"
	BrightMagenta = "\033[95m"
	BrightCyan    = "\033[96m"
	BrightWhite   = "\033[97m"

	// Background Colors
	BgBlack         = "\033[40m"
	BgRed           = "\033[41m"
	BgGreen         = "\033[42m"
	BgYellow        = "\033[43m"
	BgBlue          = "\033[44m"
	BgMagenta       = "\033[45m"
	BgCyan          = "\033[46m"
	BgWhite         = "\033[47m"
	BgBrightBlack   = "\033[100m"
	BgBrightRed     = "\033[101m"
	BgBrightGreen   = "\033[102m"
	BgBrightYellow  = "\033[103m"
	BgBrightBlue    = "\033[104m"
	BgBrightMagenta = "\033[105m"
	BgBrightCyan    = "\033[106m"
	BgBrightWhite   = "\033[107m"

	// Effects
	Bold          = "\033[1m"
	Dim           = "\033[2m"
	Italic        = "\033[3m"
	Underline     = "\033[4m"
	Blink         = "\033[5m"
	Invert        = "\033[7m"
	Hidden        = "\033[8m"
	StrikeThrough = "\033[9m"
)

func displayEntries(title string, entries []db.Entry) error {

	fmt.Println(title)

	for i, entry := range entries {
		date := entry.CreatedAt.Format("01-02-2006")
		fmt.Println(Bold, Blue, strconv.Itoa(i+1)+") ", Reset, entry.Name, " | Created: ", date)
		// Only preview first 10
		if i < 10 {
			tempLines, err := utils.GetLines(entry.FilePath)
			if err != nil {
				fmt.Println("Error reading body of entry at "+entry.FilePath, err)
				return err
			}
			preview := tempLines[config.CONFIG.START_POS-1 : config.CONFIG.START_POS+6]
			if len(preview) < 1 {
				fmt.Println("\t", "No text available for preview")
			} else {
				for _, line := range preview {

					fmt.Println("\t", Green, line, Reset)
				}
			}
		}
		//Separator
		if i < len(entries)-1 {

			fmt.Println(Yellow, "================================================================================", Reset)
		}
	}
	return nil
}

func sectionTitle(title string, symbol string) string {
	titleArr := strings.Split(title, "")
	center := 20
	start := center - int(math.Ceil(float64(len(titleArr))/float64(2)))
	end := start + len(titleArr)
	var newTitle []string = make([]string, 80)
	for i := range newTitle {
		if i >= start && i < end {
			newTitle[i] = titleArr[i-start]
		} else {
			newTitle[i] = symbol
		}
	}

	return strings.Join(newTitle, "")
}

func Render(desiredMode DisplayMode, entriesList []db.Entry, searchTags []string, infoMsg string) {
	utils.ClearTerminal()

	var title string
	CurrentDisplay = desiredMode
	CurrentEntries = entriesList

	switch CurrentDisplay {
	case MergeDisplay:
		title = string(Blue + sectionTitle("MERGE LIST", "=") + Reset)
	case SearchDisplay:
		title = string(Blue + sectionTitle("SEARCH RESULTS", "=") + Reset)
		fmt.Println(Green, "SEARCH TAGS = ", Reset, strings.Join(searchTags, ","))
	default:
		title = string(BgBrightGreen + BrightCyan + sectionTitle("ERROR", "x") + Reset)
	}
	fmt.Println("")

	// Entry Display
	displayEntries(title, entriesList)

	// Options
	fmt.Println(BrightMagenta, sectionTitle("OPTIONS", "="), Reset)
	if CurrentDisplay == SearchDisplay {

		// E.g. r -i finance
		fmt.Println(Magenta + "[R]efine current search: " + Reset + "r -[opts] [tag]...")
		// E.g. n -a health
		fmt.Println(Magenta + "[N]ew search: " + Reset + "n -[opts] [tag]...")
		// E.g. a 1 4 12
		fmt.Println(Magenta + "[A]dd entry to volume list: " + Reset + "a [number]...")
		// E.g. w
		//TODO uncomment and make this work
		// fmt.Println(Magenta + "[W]hole list to volume list: " + Reset + "w")
		// E.g. d 1 4 12
		fmt.Println(Magenta + "[D]elete entry permanently: " + Reset + "d [number]...")

		// E.g. v
		fmt.Println(Magenta + "[V]iew current volume list: " + Reset + "v")
		// E.g. q
		fmt.Println(Magenta + "[Q]uit: " + Reset + "q")
		// E.g. 31
		fmt.Println(Magenta + "[#] Number of the file to open: " + Reset + "[number]")

	} else {

		// E.g. m 2024
		fmt.Println(Magenta + "[M]erge entries from merge list to single volume: " + Reset + "m [name]...")
		// E.g. d 2 12 6
		fmt.Println(Magenta + "[D]elete entries from merge list:" + Reset + "d [number]...")
		// E.g. b
		fmt.Println(Magenta + "[B]ack to results: " + Reset + "b")
		// E.g. q
		fmt.Println(Magenta + "[Q]uit: " + Reset + "q")
	}

	// Info
	if infoMsg != "" {
		fmt.Println(Red, sectionTitle("INFO", "#"), Reset)
		fmt.Println(Red, infoMsg, Reset)
	}

}
