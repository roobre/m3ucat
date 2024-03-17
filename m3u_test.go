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
				{Path: "/foo/bar.flac"},
				{Path: "/boo/far.flac"},
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
	} {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			actual, err := m3u.Decode(strings.NewReader(tc.src))
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
				{Path: "/foo/bar.flac"},
				{Path: "/boo/far.flac"},
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
#EXTALB:The Sacrament of Sin
#EXTART:Powerwolf
#EXTINF:287,Sacred & Wild
../Powerwolf/2018 - The Sacrament of Sin/CD2/01. Epica - Sacred & Wild.flac
#EXTALB:The Sacrament of Sin
#EXTART:Powerwolf
#EXTINF:227,Resurrection by Erection
../Powerwolf/2018 - The Sacrament of Sin/CD2/04. Battle Beast - Resurrection By Erection.flac
	`) + "\n"

	playlist, err := m3u.Decode(strings.NewReader(golden))
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
				{Path: "/foo/bar.flac"},
				{Path: "/boo/far.flac"},
			},
			expected: m3u.Playlist{
				{Path: "/foo/bar.flac"},
				{Path: "/boo/far.flac"},
			},
		},
		{
			name: "Dupes",
			playlist: m3u.Playlist{
				{Path: "/foo/bar.flac"},
				{Path: "/foo/bar.flac"},
				{Path: "/boo/far.flac"},
			},
			expected: m3u.Playlist{
				{Path: "/foo/bar.flac"},
				{Path: "/boo/far.flac"},
			},
		},
		{
			name: "Metadata is ignored",
			playlist: m3u.Playlist{
				{Path: "/foo/bar.flac", Ext: []string{"foo"}},
				{Path: "/foo/bar.flac", Ext: []string{"bar"}},
				{Path: "/boo/far.flac"},
			},
			expected: m3u.Playlist{
				{Path: "/foo/bar.flac", Ext: []string{"foo"}},
				{Path: "/boo/far.flac"},
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
