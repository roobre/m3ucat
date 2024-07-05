package m3u

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type Playlist struct {
	Ext    []string
	Tracks []Track
}

type Track struct {
	Ext  []string
	Path string
}

func Decode(src io.Reader) (Playlist, error) {
	scn := bufio.NewScanner(src)

	playlist := Playlist{}
	track := Track{}
	for scn.Scan() {
		line := strings.TrimSpace(scn.Text())
		if line == "" {
			continue
		}

		if line == "#EXTM3U" {
			continue
		}

		if strings.HasPrefix(line, "#") {
			track.Ext = append(track.Ext, line)
			continue
		}

		track.Path = line
		playlist.Tracks = append(playlist.Tracks, track)
		track = Track{}
	}

	return playlist, scn.Err()
}

func (p Playlist) Encode(w io.Writer) error {
	_, err := fmt.Fprintln(w, "#EXTM3U")
	if err != nil {
		return err
	}

	for _, track := range p.Tracks {
		for _, ext := range track.Ext {
			_, err := fmt.Fprintln(w, ext)
			if err != nil {
				return err
			}
		}

		_, err := fmt.Fprintln(w, track.Path)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p Playlist) Deduplicate() Playlist {
	added := map[string]bool{}
	newP := Playlist{}

	for _, track := range p.Tracks {
		if added[track.Path] {
			continue
		}

		added[track.Path] = true
		newP.Tracks = append(newP.Tracks, track)
	}

	return newP
}
