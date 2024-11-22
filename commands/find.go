package commands

import (
  "bufio"
  "fmt"
  "os"
  "sort"
  "strconv"
  "strings"
  "time"

  "github.com/projectz-ro/journalz-ro/config"
  "github.com/projectz-ro/journalz-ro/db"
  "github.com/projectz-ro/journalz-ro/ui"
  utils "github.com/projectz-ro/journalz-ro/zro_utils"
  "github.com/spf13/cobra"
)

var (
  inclusive bool
  first bool
  ascending bool
  descending bool
  originalsOnly bool

  searchResults []db.Entry
  searchTags []string
  mergeList []db.Entry
)

var findCmd = &cobra.Command {
  Use: "find [tags]",
  Short: "Find entries by tags",
  Args: cobra.MinimumNArgs(1),
  RunE: func(cmd *cobra.Command, args []string) error {
    searchTags = args
    for _,
    x: = range searchTags {
      x = strings.TrimSpace(strings.ToLower(x))
    }
    err: = newSearch()
    if err != nil {
      return fmt.Errorf("Error initiating search: %v", err)
    }
    promptLoop()
    return nil
  },
}

func init() {
  rootCmd.AddCommand(findCmd)

  // Flags
  findCmd.PersistentFlags().BoolVarP(&inclusive, "inclusive", "i", false, "Show entries including ANY of the provided tags")
  findCmd.PersistentFlags().BoolVarP(&first, "first", "f", false, "Return only the first file to match the provided tags")
  findCmd.PersistentFlags().BoolVarP(&ascending, "ascending", "a", false, "Sort by date/time in ascending order")
  findCmd.PersistentFlags().BoolVarP(&descending, "descending", "d", false, "Sort by date/time in descending order")
  findCmd.PersistentFlags().BoolVarP(&originalsOnly, "originals-only", "o", false, "Show only original entries (exclude volumes)")
}

func newSearch() error {

  err: = initialSearch()
  if err != nil {
    return fmt.Errorf("failed to search entries: %v", err)
  }

  refErr: = refineSearch()
  if refErr != nil {
    return fmt.Errorf("failed to refine search: %v", refErr)
  }
  return nil
}
func initialSearch() error {
  err: = db.USERDB.DB.Preload("Tags").Where("EXISTS (?)", db.USERDB.DB.Table("entry_tags").
    Select("entry_id").
    Joins("JOIN tags ON entry_tags.tag_id = tags.id").
    Where("tags.tag_name IN (?)", searchTags)).
  Find(&searchResults).Error
  if err != nil {
    return fmt.Errorf("failed to search entries by tags: %v", err)
  }
  return nil
}
func refineSearch() error {
  var results []db.Entry

  if !inclusive {
    tagSet: = make(map[string]bool)
    for _,
    tag: = range searchTags {
      tagSet[tag] = true
    }

    for _,
    entry: = range searchResults {
      remaining: = len(searchTags)
      for _,
      tag: = range searchTags {
        for _,
        etag: = range entry.Tags {
          if tag == etag.TagName {
            remaining = remaining - 1
            break
          }
        }
      }
      if remaining < 1 {
        results = append(results, entry)
      }
    }
  } else {
    results = append(results, searchResults...)
  }

  if originalsOnly {
    var filtered []db.Entry
    for _,
    res: = range results {
      if len(res.Originals) == 0 {
        filtered = append(results, res)
      }
    }
    results = filtered
  }
  if len(results) == 0 {
    fmt.Println("No results found")

  }
  if ascending && descending {
    return fmt.Errorf("cannot sort by both ascending and descending")
  }
  if ascending {
    sort.Slice(results, func(i, j int) bool {
      return results[i].CreatedAt.Before(results[j].CreatedAt)
    })
  }
  if descending {
    sort.Slice(results, func(i, j int) bool {
      return results[i].CreatedAt.After(results[j].CreatedAt)
    })
  }
  if first {
    utils.OpenInNvim(results[0].FilePath, config.CONFIG.START_POS, false)
  }

  searchResults = results
  return nil
}

func promptLoop() error {
  currentMode: = ui.SearchDisplay
  currentMsg: = ""
  currentList: = searchResults

  for {

    if currentMode == ui.SearchDisplay {
      currentList = searchResults
    } else {
      currentList = mergeList
    }

    ui.Render(currentMode, currentList, searchTags, currentMsg)
    fmt.Print("Your decision: ")

    var input string
    scanner: = bufio.NewScanner(os.Stdin)
    if scanner.Scan() {
      input = scanner.Text()
    }

    inputArr: = strings.Split(strings.ToLower(strings.Trim(input, " ")), " ")
    newCmd: = inputArr[0]
    newArgs: = inputArr[1:]

    // Secondary Commands
    if ui.CurrentDisplay == ui.SearchDisplay {

      switch strings.ToLower(newCmd) {
        case "r":
          if len(searchResults) > 4 {
            if len(newArgs) > 0 {
              first = false
              inclusive = false
              ascending = false
              descending = false
              originalsOnly = false

              var tempTags []string
              for _,
              arg: = range newArgs {
                split: = strings.Split(arg, "")
                if split[0] == "-" {
                  switch split[1] {
                  case "i":
                    inclusive = true
                  case "f":
                    first = true
                  case "a":
                    ascending = true
                  case "d":
                    descending = true
                  case "o":
                    originalsOnly = true
                  default:
                    currentMode = ui.SearchDisplay
                    currentMsg = "Not a recognized flag"
                    break
                  }
                } else {
                  tempTags = append(tempTags, arg)
                }
              }
              searchTags = tempTags
              refineSearch()
            } else {
              currentMode = ui.SearchDisplay
              currentMsg = "You must supply at least one tag to search for."
              break
            }
          } else {
            currentMode = ui.SearchDisplay
            currentMsg = "Refinement is only available when there are 5 or more results."
            break
          }
        case "n":
          if len(newArgs) > 0 {
            first = false
            inclusive = false
            ascending = false
            descending = false
            originalsOnly = false

            var tempTags []string
            for _,
            arg: = range newArgs {
              split: = strings.Split(arg, "")
              if split[0] == "-" {
                switch split[1] {
                case "i":
                  inclusive = true
                case "f":
                  first = true
                case "a":
                  ascending = true
                case "d":
                  descending = true
                case "o":
                  originalsOnly = true
                default:
                  currentMode = ui.SearchDisplay
                  currentMsg = "Not a recognized flag"
                  break
                }
              } else {
                tempTags = append(tempTags, arg)
              }
            }
            searchTags = tempTags
            newSearch()
          } else {
            currentMode = ui.SearchDisplay
            currentMsg = "You must supply at least one tag to search for."
            break
          }
        case "a":

          if len(searchResults) > 0 {
            var tempList []string
            for _,
            arg: = range newArgs {
              selectedNumber,
              err: = strconv.Atoi(arg)
              if err != nil || selectedNumber < 1 || selectedNumber > len(searchResults) {
                fmt.Println("Invalid selection:"+arg, err)
              }
              mergeList = append(mergeList, searchResults[selectedNumber-1])
              tempList = append(tempList, strconv.Itoa(selectedNumber))
            }
            currentMode = ui.SearchDisplay
            currentMsg = strings.Join(tempList, ", ") + " added to volume list"
            break
          } else {
            currentMode = ui.SearchDisplay
            currentMsg = "ADD WHAT?!"
            break
          }
        case "d":
          for _,
          arg: = range newArgs {
            selectedNumber,
            err: = strconv.Atoi(arg)
            if err != nil || selectedNumber < 1 || selectedNumber > len(searchResults) {
              fmt.Println("Invalid selection:"+arg, err)
              break
            }
            delErr: = db.USERDB.DeleteEntry(searchResults[selectedNumber-1].ID)
            if delErr != nil {
              return fmt.Errorf("Error deleting entry: %v", delErr)
            }
            searchResults = append(
              searchResults[: selectedNumber-1], searchResults[selectedNumber:]...)
            currentMode = ui.SearchDisplay
            currentMsg = "Entry deleted"
            break
          }
          break
        case "v":
          if len(mergeList) > 0 {
            currentMode = ui.MergeDisplay
            currentMsg = ""
            break
          } else {
            currentMode = ui.SearchDisplay
            currentMsg = "Add something to your volume list first..."
            break
          }

        case "q":
          os.Exit(0)
        default:
          selectedNumber,
          err: = strconv.Atoi(input)
          if err != nil || selectedNumber < 1 || selectedNumber > len(searchResults) {
            currentMode = ui.SearchDisplay
            currentMsg = "Invalid selection. Please enter a valid option."
            break
          }
          utils.OpenInNvim(searchResults[selectedNumber-1].FilePath, 0, false)
      }
    } else {
      switch strings.ToLower(newCmd) {
      case "m":
        if len(mergeList) > 1 && newArgs[0] != "" {
          newVolume,
          err: = createVolume(newArgs)
          if err != nil {
            currentMode = ui.SearchDisplay
            currentMsg = fmt.Sprintf("Error merging entries: %v", err)
            break
          } else {
            fmt.Println("Volume Created Successfully")
            utils.OpenInNvim(newVolume, config.CONFIG.START_POS, false)
            os.Exit(0)
          }
        } else {
          currentMode = ui.SearchDisplay
          currentMsg = "Add at least two entries to your volume list first..."
          break
        }
      case "b":
        ui.Render(ui.SearchDisplay, searchResults, searchTags, "")
        currentMode = ui.SearchDisplay
        currentMsg = ""
        break
      case "d":
        if len(newArgs) > 0 && len(mergeList) > 0 {
          for _,
          arg: = range newArgs {
            selectedNumber,
            err: = strconv.Atoi(arg)
            if err != nil || selectedNumber < 1 || selectedNumber > len(mergeList) {
              currentMode = ui.MergeDisplay
              currentMsg = "Invalid selection. Please enter a valid option."
              break
            }

            mergeList = append(
              mergeList[: selectedNumber-1],
              mergeList[selectedNumber:]...)
          }
          currentMode = ui.MergeDisplay
          currentMsg = "Merge list updated"
          break

        }
        currentMode = ui.MergeDisplay
        currentMsg = "DELETE WHAT?!"
        break

      case "q":
        os.Exit(0)
      default:
        currentMode = ui.SearchDisplay
        currentMsg = "Invalid command"
        break
      }
    }
  }
}

func createVolume(newArgs []string) (string, error) {
  allTags: = getVolTags()
  allOriginals: = getVolOg()

  day: = time.Now().Format("Monday")
  date: = time.Now().Format("01/02/2006")
  time: = time.Now().Format("15:04")

  maxWidth: = 80

  name: = strings.Join(newArgs, " ")
  filepath: = config.CONFIG.VOLUME_DIR + name + ".md"

  _,
  err: = db.USERDB.InsertEntry(
    name,
    allTags,
    allOriginals,
    filepath,
  )
  if err != nil {
    return "",
    fmt.Errorf("Error adding new entry: %v", err)
  }

  lines: = []string {
    "",
    fmt.Sprintf("%*s", maxWidth, day),
    fmt.Sprintf("%*s", maxWidth, date),
    fmt.Sprintf("%*s", maxWidth, time),
    fmt.Sprintf("# %s", name),
    "",
    "---",
    "",
    "",
  }
  // TODO warn for long running process if they add a crazy number of entries
  // perhaps limit the number of entries later

  for _,
  file: = range mergeList {
    tempLines,
    err: = utils.GetLines(file.FilePath)
    if err != nil {
      return "",
      fmt.Errorf("Error reading original entries: %v", err)
    }
    lines = append(lines, tempLines[config.CONFIG.START_POS-1:]...)
    lines = append(lines, "")
    lines = append(lines, "---")
    lines = append(lines, "")
  }
  writeErr: = utils.WriteLines(filepath, lines)
  if writeErr != nil {
    return "",
    fmt.Errorf("Error writing new file: %v", writeErr)
  }

  return filepath,
  nil
}
func getVolTags() []string {
  tagSet: = make(map[string]struct {})

  for _,
  entry: = range mergeList {
    for _,
    tag: = range entry.Tags {
      tagSet[tag.TagName] = struct {} {}
    }
  }

  var uniqueTags []string
  for tag: = range tagSet {
    uniqueTags = append(uniqueTags, tag)
  }

  return uniqueTags
}

func getVolOg() []db.Entry {
  var temp []db.Entry
  for _,
  entry: = range mergeList {
    temp = append(temp, entry)
    if len(entry.Originals) > 0 {
      temp = append(temp, entry.Originals...)
    }
  }
  return temp

}