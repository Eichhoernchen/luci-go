// Copyright 2017 The LUCI Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package ensure

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/luci/luci-go/cipd/client/cipd/common"
	. "github.com/smartystreets/goconvey/convey"
)

func f(lines ...string) string {
	return strings.Join(lines, "\n")
}

func p(pkg, ver string) common.Pin {
	return common.Pin{PackageName: pkg, InstanceID: ver}
}

var goodEnsureFiles = []struct {
	name   string
	file   string
	expect *ResolvedFile
}{
	{
		"old_style",
		f(
			"# comment",
			"",
			"path/to/package deadbeefdeadbeefdeadbeefdeadbeefdeadbeef",
			"path/to/other_package some_tag:version",
			"path/to/yet_another a_ref",
		),
		&ResolvedFile{"", common.PinSliceBySubdir{
			"": {
				p("path/to/package", "deadbeefdeadbeefdeadbeefdeadbeefdeadbeef"),
				p("path/to/other_package", "some_tag:version"),
				p("path/to/yet_another", "a_ref"),
			},
		}},
	},

	{
		"templates",
		f(
			"path/to/package/${os}-${arch} latest",
			"path/to/other/${platform} latest",
		),
		&ResolvedFile{"", common.PinSliceBySubdir{
			"": {
				p("path/to/package/test_os-test_arch", "latest"),
				p("path/to/other/test_os-test_arch", "latest"),
			},
		}},
	},

	{
		"optional_templates",
		f(
			"path/to/package/${os}-${arch=neep,test_arch} latest",
			"path/to/other/${platform=test_os-test_arch} latest",
		),
		&ResolvedFile{"", common.PinSliceBySubdir{
			"": {
				p("path/to/package/test_os-test_arch", "latest"),
				p("path/to/other/test_os-test_arch", "latest"),
			},
		}},
	},

	{
		"optional_templates_no_match",
		f(
			"path/to/package/${os=spaz}-${arch=neep,test_arch} latest",
			"path/to/package/${platform=neep-foo} latest",
		),
		&ResolvedFile{"", common.PinSliceBySubdir{}},
	},

	{
		"Subdir directives",
		f(
			"some/package latest",
			"",
			"@Subdir a/subdir with spaces",
			"some/package canary",
			"some/other/package tag:value",
			"",
			"@Subdir", // reset back to empty
			"cool/package beef",
		),
		&ResolvedFile{"", common.PinSliceBySubdir{
			"": {
				p("some/package", "latest"),
				p("cool/package", "beef"),
			},
			"a/subdir with spaces": {
				p("some/package", "canary"),
				p("some/other/package", "tag:value"),
			},
		}},
	},

	{
		"ServiceURL setting",
		f(
			"$ServiceURL https://cipd.example.com/path/to/thing",
			"",
			"some/package version",
		),
		&ResolvedFile{"https://cipd.example.com/path/to/thing", common.PinSliceBySubdir{
			"": {
				p("some/package", "version"),
			},
		}},
	},

	{
		"empty",
		"",
		&ResolvedFile{},
	},

	{
		"wacky spaces",
		f(
			"path/to/package           latest",
			"tabs/to/package\t\t\t\tlatest",
			"\ttabs/and/spaces  \t  \t  \tlatest   \t",
		),
		&ResolvedFile{"", common.PinSliceBySubdir{
			"": {
				p("path/to/package", "latest"),
				p("tabs/to/package", "latest"),
				p("tabs/and/spaces", "latest"),
			},
		}},
	},
}

func testResolver(pkg, vers string) (common.Pin, error) {
	if strings.Contains(vers, "error") {
		return p("", ""), errors.New("testResolver returned error")
	}
	return p(pkg, vers), nil
}

func TestGoodEnsureFiles(t *testing.T) {
	t.Parallel()

	Convey("good ensure files", t, func() {
		for _, tc := range goodEnsureFiles {
			Convey(tc.name, func() {
				buf := bytes.NewBufferString(tc.file)
				f, err := ParseFile(buf)
				So(err, ShouldBeNil)
				rf, err := f.ResolveWith(testResolver, map[string]string{
					"os":       "test_os",
					"arch":     "test_arch",
					"platform": "test_os-test_arch",
				})
				So(err, ShouldBeNil)
				So(rf, ShouldResemble, tc.expect)
			})
		}
	})
}
