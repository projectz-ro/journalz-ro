
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

Use git to clone the repo

```bash 
cd journalz-ro
go build -o journalz-ro
```
or if you don't like typing

```bash
cd journalz-ro
go build -o jzro
```
then I recommend moving it into your $PATH somewhere

```bash
sudo mv journalz-ro /usr/local/bin/
```

What you name the build will be the name you use to run the program.

## Usage

### New Entry

Create a new journal entry:
```bash
journalz-ro new [tags]
```
This command generates a new entry with the tags provided and opens it in neovim. You do not name notes. You tag them by subject.

This, I think, will make finding older notes easier and more rewarding. It also takes away that "what do I call this...." problem and let's you get straight to putting your thoughts down. 

### Find Entries by Tag

Find entries associated with a specific tag:
```bash
journalz-ro find <tag>
```
Find entries then refine your search, start a new search, delete entries or add them to a merge list.

### Merge Entries
Merge entries that share a specific tag into a Volume. Merge commands happen from within the find command. This requires a name for the volume:
```bash
m <name>
```
The volume will be saved as `<name>` in the VOLUME_DIR directory.

Naming volumes makes sense to me as they are more curated. You're gathering your thoughts about one or more related topics, into one easy reference and maybe even for cleaning up into a finished work.

## Configuration

```bash
nvim ~/.config/journalz-ro/config.json
```

### Defaults

```json
        "SAVE_DIR":        "~/Documents/",
		"VOLUME_DIR":      "~/Documents/Volumes/",
		"INSERT_ON_NEW":   true,
		"START_POS":       8,

```
## Planned Features
1. Templates for entries and volumes, to customize how they are formatted
2. Greater search and filtering functionality, like date ranges
3. Color Themes

## Thanks

Contributions to JournalZ-ro are welcome! Please feel free to open issues or submit pull requests.

This is my first open-source tool and my first public build. I'm trying to learn more of everything, including git/github. Code reviews and opinions are welcome. 

Thanks again, and I hope this tool can help some people to learn smoothly. 

I made this tool because I wanted to take notes in my configured and comfy neovim without having to switch brains or think much. Just bind `journalz-ro new` to a key combo, type my note, tag it and move on. No title thinking, just quick notes that can later be merged by tags. Then those can be taken and refined later. I think this approach has a lot of potential. 
