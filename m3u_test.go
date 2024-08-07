package m3u_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	m3u "roob.re/m3u"
)

func TestDecode(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name      string
		src       string
		decoder   m3u.Decoder
		expect    m3u.Playlist
		expectErr error
	}{
		{
			name: "List of paths",
			src: `
/foo/bar.flac
/boo/far.flac
			`,
			expect: m3u.Playlist{
				Tracks: []m3u.Track{
					{Path: "/foo/bar.flac"},
					{Path: "/boo/far.flac"},
				},
			},
		},
		{
			name: "Real playlist",
			src: `
#EXTM3U
#EXTALB:The Sacrament of Sin
#EXTART:Powerwolf
#EXTINF:287,Sacred & Wild
../Powerwolf/2018 - The Sacrament of Sin/CD2/01. Epica - Sacred & Wild.flac
#EXTALB:The Sacrament of Sin
#EXTART:Powerwolf
#EXTINF:227,Resurrection by Erection
../Powerwolf/2018 - The Sacrament of Sin/CD2/04. Battle Beast - Resurrection By Erection.flac
			`,
			expect: m3u.Playlist{
				Tracks: []m3u.Track{
					{
						Path: "../Powerwolf/2018 - The Sacrament of Sin/CD2/01. Epica - Sacred & Wild.flac",
						Ext: []string{
							"#EXTALB:The Sacrament of Sin",
							"#EXTART:Powerwolf",
							"#EXTINF:287,Sacred & Wild",
						},
					},
					{
						Path: "../Powerwolf/2018 - The Sacrament of Sin/CD2/04. Battle Beast - Resurrection By Erection.flac",
						Ext: []string{
							"#EXTALB:The Sacrament of Sin",
							"#EXTART:Powerwolf",
							"#EXTINF:227,Resurrection by Erection",
						},
					},
				},
			},
		},
		{
			name: "Playlist with hiccups",
			src: `
#EXTALB:The Sacrament of Sin
#EXTART:Powerwolf


#EXTINF:287,Sacred & Wild
../Powerwolf/2018 - The Sacrament of Sin/CD2/01. Epica - Sacred & Wild.flac
/another/file.flac
#EXTALB:The Sacrament of Sin
#EXTART:Powerwolf
#EXTINF:227,Resurrection by Erection

../Powerwolf/2018 - The Sacrament of Sin/CD2/04. Battle Beast - Resurrection By Erection.flac
			`,
			expect: m3u.Playlist{
				Tracks: []m3u.Track{
					{
						Path: "../Powerwolf/2018 - The Sacrament of Sin/CD2/01. Epica - Sacred & Wild.flac",
						Ext: []string{
							"#EXTALB:The Sacrament of Sin",
							"#EXTART:Powerwolf",
							"#EXTINF:287,Sacred & Wild",
						},
					},
					{Path: "/another/file.flac"},
					{
						Path: "../Powerwolf/2018 - The Sacrament of Sin/CD2/04. Battle Beast - Resurrection By Erection.flac",
						Ext: []string{
							"#EXTALB:The Sacrament of Sin",
							"#EXTART:Powerwolf",
							"#EXTINF:227,Resurrection by Erection",
						},
					},
				},
			},
		},
		{
			name: "Playlist-level attributes",
			decoder: m3u.Decoder{IsGlobalDirective: func(s string) bool {
				return strings.Contains(s, "PLAYLIST")
			}},
			src: `
#EXTM3U
#PLAYLIST:Title
/foo/bar.flac
/boo/far.flac
			`,
			expect: m3u.Playlist{
				Ext: []string{
					"#PLAYLIST:Title",
				},
				Tracks: []m3u.Track{
					{Path: "/foo/bar.flac"},
					{Path: "/boo/far.flac"},
				},
			},
		},
	} {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			actual, err := tc.decoder.Decode(strings.NewReader(tc.src))
			if !errors.Is(err, tc.expectErr) {
				t.Fatalf("Expected error to be %v, got %v", tc.expectErr, err)
			}

			if diff := cmp.Diff(actual, tc.expect); diff != "" {
				t.Fatalf("Decoded playlist does not match expected:\n%s", diff)
			}
		})
	}
}

func TestEncode(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name      string
		expect    string
		playlist  m3u.Playlist
		expectErr error
	}{
		{
			name: "List of paths",
			expect: strings.TrimSpace(`
#EXTM3U
/foo/bar.flac
/boo/far.flac
			`) + "\n",
			playlist: m3u.Playlist{
				Tracks: []m3u.Track{
					{Path: "/foo/bar.flac"},
					{Path: "/boo/far.flac"},
				},
			},
		},
		{
			name: "Playlist level directives",
			expect: strings.TrimSpace(`
#EXTM3U
#PLAYLIST:Foobar
/foo/bar.flac
/boo/far.flac
			`) + "\n",
			playlist: m3u.Playlist{
				Ext: []string{
					"#PLAYLIST:Foobar",
				},
				Tracks: []m3u.Track{
					{Path: "/foo/bar.flac"},
					{Path: "/boo/far.flac"},
				},
			},
		},
		{
			name: "Real playlist",
			expect: strings.TrimSpace(`
#EXTM3U
#EXTALB:The Sacrament of Sin
#EXTART:Powerwolf
#EXTINF:287,Sacred & Wild
../Powerwolf/2018 - The Sacrament of Sin/CD2/01. Epica - Sacred & Wild.flac
#EXTALB:The Sacrament of Sin
#EXTART:Powerwolf
#EXTINF:227,Resurrection by Erection
../Powerwolf/2018 - The Sacrament of Sin/CD2/04. Battle Beast - Resurrection By Erection.flac
			`) + "\n",
			playlist: m3u.Playlist{
				Tracks: []m3u.Track{
					{
						Path: "../Powerwolf/2018 - The Sacrament of Sin/CD2/01. Epica - Sacred & Wild.flac",
						Ext: []string{
							"#EXTALB:The Sacrament of Sin",
							"#EXTART:Powerwolf",
							"#EXTINF:287,Sacred & Wild",
						},
					},
					{
						Path: "../Powerwolf/2018 - The Sacrament of Sin/CD2/04. Battle Beast - Resurrection By Erection.flac",
						Ext: []string{
							"#EXTALB:The Sacrament of Sin",
							"#EXTART:Powerwolf",
							"#EXTINF:227,Resurrection by Erection",
						},
					},
				},
			},
		},
	} {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			buf := &bytes.Buffer{}
			err := tc.playlist.Encode(buf)
			if !errors.Is(err, tc.expectErr) {
				t.Fatalf("Expected error to be %v, got %v", tc.expectErr, err)
			}

			if diff := cmp.Diff(buf.String(), tc.expect); diff != "" {
				t.Fatalf("Decoded playlist does not match expected:\n%s", diff)
			}
		})
	}
}

func TestDecodeEncode(t *testing.T) {
	t.Parallel()

	golden := strings.TrimSpace(`
#EXTM3U
#PLAYLIST:Powerwolf
#EXTALB:The Sacrament of Sin
#EXTART:Powerwolf
#EXTINF:287,Sacred & Wild
../Powerwolf/2018 - The Sacrament of Sin/CD2/01. Epica - Sacred & Wild.flac
#EXTALB:The Sacrament of Sin
#EXTART:Powerwolf
#EXTINF:227,Resurrection by Erection
../Powerwolf/2018 - The Sacrament of Sin/CD2/04. Battle Beast - Resurrection By Erection.flac
	`) + "\n"

	decoder := m3u.Decoder{IsGlobalDirective: func(s string) bool {
		return strings.HasPrefix(s, "#PLAYLIST:")
	}}
	playlist, err := decoder.Decode(strings.NewReader(golden))
	if err != nil {
		t.Fatalf("decoding golden playlist: %v", err)
	}

	buf := &bytes.Buffer{}
	err = playlist.Encode(buf)
	if err != nil {
		t.Fatalf("encoding playlist: %v", err)
	}

	if diff := cmp.Diff(buf.String(), golden); diff != "" {
		t.Fatalf("decoded-encoded playlist does not match itself:\n%s", diff)
	}
}

func TestPlaylist_Deduplicate(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name     string
		playlist m3u.Playlist
		expected m3u.Playlist
	}{
		{
			name: "No dupes",
			playlist: m3u.Playlist{
				Tracks: []m3u.Track{
					{Path: "/foo/bar.flac"},
					{Path: "/boo/far.flac"},
				},
			},
			expected: m3u.Playlist{
				Tracks: []m3u.Track{
					{Path: "/foo/bar.flac"},
					{Path: "/boo/far.flac"},
				},
			},
		},
		{
			name: "Dupes",
			playlist: m3u.Playlist{
				Tracks: []m3u.Track{
					{Path: "/foo/bar.flac"},
					{Path: "/foo/bar.flac"},
					{Path: "/boo/far.flac"},
				},
			},
			expected: m3u.Playlist{
				Tracks: []m3u.Track{
					{Path: "/foo/bar.flac"},
					{Path: "/boo/far.flac"},
				},
			},
		},
		{
			name: "Metadata is ignored",
			playlist: m3u.Playlist{
				Tracks: []m3u.Track{
					{Path: "/foo/bar.flac", Ext: []string{"foo"}},
					{Path: "/foo/bar.flac", Ext: []string{"bar"}},
					{Path: "/boo/far.flac"},
				},
			},
			expected: m3u.Playlist{
				Tracks: []m3u.Track{
					{Path: "/foo/bar.flac", Ext: []string{"foo"}},
					{Path: "/boo/far.flac"},
				},
			},
		},
		{
			name: "Playlist directives are preserved",
			playlist: m3u.Playlist{
				Ext: []string{
					"#SOMETHING",
				},
				Tracks: []m3u.Track{
					{Path: "/foo/bar.flac", Ext: []string{"foo"}},
					{Path: "/foo/bar.flac", Ext: []string{"bar"}},
					{Path: "/boo/far.flac"},
				},
			},
			expected: m3u.Playlist{
				Ext: []string{
					"#SOMETHING",
				},
				Tracks: []m3u.Track{
					{Path: "/foo/bar.flac", Ext: []string{"foo"}},
					{Path: "/boo/far.flac"},
				},
			},
		},
	} {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if diff := cmp.Diff(tc.playlist.Deduplicate(), tc.expected); diff != "" {
				t.Fatalf("Deduped playlist does not match expected:\n%s", diff)
			}
		})
	}
}

func TestPlaylist_Join(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name     string
		first    m3u.Playlist
		second   m3u.Playlist
		expected m3u.Playlist
	}{
		{
			name: "Tracks only",
			first: m3u.Playlist{
				Tracks: []m3u.Track{
					{Path: "/1/bar.flac"},
					{Path: "/2/far.flac"},
				},
			},
			second: m3u.Playlist{
				Tracks: []m3u.Track{
					{Path: "/3/bar.flac"},
					{Path: "/4/far.flac"},
				},
			},
			expected: m3u.Playlist{
				Tracks: []m3u.Track{
					{Path: "/1/bar.flac"},
					{Path: "/2/far.flac"},
					{Path: "/3/bar.flac"},
					{Path: "/4/far.flac"},
				},
			},
		},
		{
			name: "Metadata",
			first: m3u.Playlist{
				Ext: []string{
					"#EXTFOO",
				},
			},
			second: m3u.Playlist{
				Ext: []string{
					"#EXTBAR",
				},
				Tracks: []m3u.Track{
					{Path: "/3/bar.flac"},
					{Path: "/4/far.flac"},
				},
			},
			expected: m3u.Playlist{
				Ext: []string{
					"#EXTFOO",
					"#EXTBAR",
				},
				Tracks: []m3u.Track{
					{Path: "/3/bar.flac"},
					{Path: "/4/far.flac"},
				},
			},
		},
	} {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if diff := cmp.Diff(tc.first.Join(tc.second), tc.expected); diff != "" {
				t.Fatalf("Deduped playlist does not match expected:\n%s", diff)
			}
		})
	}
}
