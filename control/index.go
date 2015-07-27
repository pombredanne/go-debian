/* {{{ Copyright (c) Paul R. Tagliamonte <paultag@debian.org>, 2015
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE. }}} */

package control

import (
	"bufio"
	"strings"

	"pault.ag/go/debian/dependency"
	"pault.ag/go/debian/version"
)

type PackageIndex struct {
	Paragraph

	Package  string
	Binaries []string

	Version    version.Version
	Maintainer string

	BuildDepends dependency.Dependency
	Architecture []dependency.Arch

	StandardsVersion string
	Format           string
	Files            []string
	VcsBrowser       string
	VcsGit           string
	Homepage         string
	Directory        string
	Priority         string
	Section          string

	/*
		TODO:
			Checksums-Sha1:
			Checksums-Sha256:
			Package-List:
	*/
}

func ParsePackagesIndex(reader *bufio.Reader) (ret []PackageIndex, err error) {
	ret = []PackageIndex{}

	for {
		block, err := ParsePackagesIndexParagraph(reader)
		if err != nil {
			return ret, err
		}
		if block != nil {
			ret = append(ret, *block)
		} else {
			break
		}
	}

	return
}

// Given a bufio.Reader, produce a PackageIndex struct to encapsulate the
// data contained within.
func ParsePackagesIndexParagraph(reader *bufio.Reader) (ret *PackageIndex, err error) {

	/* a PackageIndex is a Paragraph, with some stuff. So, let's first take
	 * the bufio.Reader and produce a stock Paragraph. */
	src, err := ParseParagraph(reader)
	if err != nil {
		return nil, err
	}

	if src == nil {
		return nil, nil
	}

	version, err := version.Parse(src.Values["Version"])
	if err != nil {
		return nil, err
	}

	arch, err := dependency.ParseArchitectures(src.Values["Architecture"])
	if err != nil {
		return nil, err
	}

	ret = &PackageIndex{
		Paragraph: *src,

		Package:  src.Values["Package"],
		Binaries: strings.Split(src.Values["Binary"], ", "),

		Version:    version,
		Maintainer: src.Values["Maintainer"],

		BuildDepends: src.getOptionalDependencyField("Build-Depends"),
		Architecture: arch,

		// Directory        string
		// Priority         string
		// Section          string

		VcsBrowser: src.Values["VcsBrowser"],
		VcsGit:     src.Values["VcsGit"],

		Directory: src.Values["Directory"],
		Priority:  src.Values["Priority"],
		Section:   src.Values["Section"],

		Format:           src.Values["Format"],
		StandardsVersion: src.Values["Standards-Version"],
		Homepage:         src.Values["Homepage"],

		Files: strings.Split(src.Values["Files"], "\n"),
	}

	return
}

// vim: foldmethod=marker
