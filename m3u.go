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

type Decoder struct {
	// IsGlobalDirective, is used to check whether a directive is assumed to belong to the playlist itself, instead of
	// the next track in the file.
	// Examples of this include #PLAYLIST:, but other implementations may use others.
	// IsGlobalDirectve receives the whole directive line, including the leading #.
	IsGlobalDirective func(string) bool
}

var defaultDecoder = Decoder{}

func Decode(src io.Reader) (Playlist, error) {
	return defaultDecoder.Decode(src)
}

func (d Decoder) Decode(src io.Reader) (Playlist, error) {
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
			if d.IsGlobalDirective != nil && d.IsGlobalDirective(line) {
				playlist.Ext = append(playlist.Ext, line)
			} else {
				track.Ext = append(track.Ext, line)
			}

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

	for _, ext := range p.Ext {
		_, err := fmt.Fprintln(w, ext)
		if err != nil {
			return err
		}
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
	newP := Playlist{}
	newP.Ext = append(newP.Ext, p.Ext...)

	added := map[string]bool{}
	for _, track := range p.Tracks {
		if added[track.Path] {
			continue
		}

		added[track.Path] = true
		newP.Tracks = append(newP.Tracks, track)
	}

	return newP
}

func (p Playlist) Join(other Playlist) Playlist {
	newP := Playlist{}
	newP.Tracks = append(newP.Tracks, p.Tracks...)
	newP.Tracks = append(newP.Tracks, other.Tracks...)

	newP.Ext = append(newP.Ext, p.Ext...)
	newP.Ext = append(newP.Ext, other.Ext...)

	return newP
}
