
# JournalZ-ro

JournalZ-ro is a command-line tool for creating, finding, and merging entries in your journal, designed for efficient and tag-based entry management.
Use it for journaling, note-taking or whatever you like

## Features

- **New Entry**: Create a new entry with predefined templates.
- **Find Entries by Tag**: Quickly search through entries using tags.
- **Merge Entries**: Merge entries with matching tags into a single volume.

## Dependencies

- **Neovim** 
   *MUST BE BOUND TO `nvim` COMMAND*
- **Go(golang)**
   *Not version specific to my knowledge, but built on 1.23.1*

## Installation

## Usage

### New Entry

Create a new journal entry:
```bash
journalz-ro new
```
This command generates a new entry based on the template defined in `entry_template` and opens it in a new neovim instance.

### Find Entries by Tag

Find entries associated with a specific tag:
```bash
journalz-ro find <tag>
```
Find entries then refine your search, start a new search, delete entries or add them to a merge list.

### Merge Entries
Merge entries that share a specific tag. Merge commands happen from within the find command. This requires a name for the merged entry:
```bash
m <name>
```
The merged entry will be saved as `<name>` in the MERGE_DIR directory.

## Configuration

```bash
nvim ~/.config/journalz-ro/config.json
```

### Defaults

```json
        "SAVE_DIR":        "~/Documents/",
		"VOLUME_DIR":      "~/Documents/Volumes/",
		"ENTRY_TEMPLATE":  "~/.config/journalz-ro/entry_template.md",
		"VOLUME_TEMPLATE": "~/.config/journalz-ro/volume_template.md",
		"INSERT_ON_NEW":   true,
		"START_POS":       5,
		"SEPARATE_WIN":    true,
		"TERMINAL_APP":    "alacritty",

```
## Planned Features
1. Templates for entries and volumes
    - Customize how entries and volumes are formatted

## Thanks

Contributions to JournalZ-ro are welcome! Please feel free to open issues or submit pull requests.

This is my first open-source tool and my first public build. I'm trying to learn more of everything, including git/github. Code reviews and opinions are welcome. 

Thanks again, and I hope this tool can help some people to learn smoothly. 

I made this tool because I wanted to take notes in my configured and comfy neovim without having to switch brains or think much. Just bind `journalz-ro new` to a key combo, type my note, tag it and move on. No title thinking, just quick notes that can later be merged by tags. Then those can be taken and refined later. I think this approach has a lot of potential. 
