package main

import (
	"os"
)

type model struct {
	mode       viewMode
	options    []string
	choices    []os.DirEntry
	cursor     int
	selected   map[int]struct{}
	musicQueue MusicQueue
}

func initialModel(musicFiles []os.DirEntry) model {
	return model{
		mode:       viewMenu,
		options:    []string{"add", "play", "list", "stop", "resume", "back", "next"},
		choices:    musicFiles,
		selected:   make(map[int]struct{}),
		musicQueue: NewMusicQueue(),
	}
}
