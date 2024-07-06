# 🎵 m3ucat

Like `cat` but for m3u playlists. It also deduplicates them.

`m3ucat` receives a list of file paths as arguments, which are concatenated and output to `stdout`:

```console
$ m3ucat [--allow-duplicates] </path/to.m3u, ...>
#EXTM3U
/path/to/track.mp3
/path/to/track.mp3
```

The special value `-` causes `m3ucat` to read a playlist from `stdin`. It can be used with regular files.

`m3ucat` preserves directives found in the source playlists.

### Deduplication

By default, `m3ucat` removes duplicates by comparing the path to a track. If a path has been seen before, it does not add it again. `m3ucat` ignores directives when checking if a track is a duplicate, and the directives that are output are those of the first track it saw.

Deduplication can be turned off by specifying `--allow-duplicates`.

### Advanced usage

M3U is not a strongly defined format. Some applications define directives that are playlist-level, instead of track level. By default, m3ucat assumes all directives other than `#EXTM3U` belong to the next track to be found in the file. If your application uses playlist-level directives, you may be interested in treating them as so, which guarantees that will survive concatenation. `m3ucat` can be given a list of directive prefixes that, if a directive matches them, will cause that directive to be attached to the playlist itself instead of the track:

```console
$ M3UCAT_GLOBAL_DIRECTIVE_PREFIXES='#PLAYLIST,#GONIC-TITLE' m3ucat [--allow-duplicates] </path/to.m3u, ...>
#EXTM3U
#PLAYLIST:Foobar
#GONIC-TITLE:"Somethingsomething"
/path/to/track.mp3
/path/to/track.mp3
```

`m3ucat` will preserve playlist-level attributes from all ingested playlists and concatenate them to the output playlist without deduplication.
