package zro_utils

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
)

// Check that a file or folder exists
func PathExists(path string) bool {
	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}

// Counts all files in a folder that match the given regex by filename
func CountFiles(folder string, regexTerm string) (int, error) {
	files, err := os.ReadDir(folder)
	if err != nil {
		return 0, err
	}

	regex := regexp.MustCompile(regexTerm)
	count := 0
	for _, file := range files {
		if !file.IsDir() && regex.MatchString(file.Name()) {
			count++
		}
	}
	return count, nil
}

// TODO double check compatibility with UTF-8 chars

// Write lines to a file(create it if necessary) from an array of strings
func WriteLines(filePath string, lines []string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("could not create file: %v", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return fmt.Errorf("could not write line: %v", err)
		}
	}
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("could not flush to file: %v", err)
	}

	return nil
}

func ClearTerminal() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error clearing terminal:", err)
	}
}

// Get all lines in a file as string array
func GetLines(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

// Check if a slice of strings contains a target string
func SliceStrContains(slice []string, target string) bool {
	for _, s := range slice {
		if s == target {
			return true
		}
	}
	return false
}

// Open up a file in Neovim.
// Starting line number of cursor.
// Immediately insert mode.
// Pass "" to termApp to open files in current terminal, else pass config TERMINAL_APP
func OpenInNvim(filePath string, startPos int, insertMode bool) {
	mainCmd := "nvim"
	cmdArgs := []string{"+" + strconv.Itoa(startPos), filePath}

	if insertMode {
		cmdArgs = append(cmdArgs, "-c", "startinsert")
	}

	cmd := exec.Command(mainCmd, cmdArgs...)

	// Connect standard I/O
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error opening neovim: %v\n", err)
		return
	}
}
