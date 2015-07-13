// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package context

import (
	"fmt"
	"sort"
)

// ListStatus indicates the status of the import.
type ListStatus byte

func (ls ListStatus) String() string {
	switch ls {
	case StatusUnknown:
		return "?"
	case StatusMissing:
		return "m"
	case StatusStd:
		return "s"
	case StatusLocal:
		return "l"
	case StatusExternal:
		return "e"
	case StatusInternal:
		return "i"
	case StatusUnused:
		return "u"
	case StatusProgram:
		return "p"
	case StatusVendor:
		return "v"
	}
	return ""
}

const (
	// StatusUnknown indicates the status was unable to be obtained.
	StatusUnknown ListStatus = iota
	// StatusMissing indicates import not found in GOROOT or GOPATH.
	StatusMissing
	// StatusStd indicates import found in GOROOT.
	StatusStd
	// StatusLocal indicates import is part of the local project.
	StatusLocal
	// StatusExternal indicates import is found in GOPATH and not copied.
	StatusExternal
	// StatusInternal indicates import has been copied locally under internal.
	StatusInternal
	// StatusUnused indicates import has been copied, but is no longer used.
	StatusUnused
	// StatusProgram indicates the import is a main package but internal or vendor.
	StatusProgram
	// StatusVendor indicates theimport is in the vendor folder.
	StatusVendor
)

// ListItem represents a package in the current project.
type ListItem struct {
	Status     ListStatus
	Path       string
	VendorPath string
}

func (li ListItem) String() string {
	if len(li.VendorPath) == 0 || li.VendorPath == li.Path {
		return fmt.Sprintf("%s %s", li.Status, li.Path)
	}
	return fmt.Sprintf("%s %s [%s]", li.Status, li.Path, li.VendorPath)
}

type listItemSort []ListItem

func (li listItemSort) Len() int      { return len(li) }
func (li listItemSort) Swap(i, j int) { li[i], li[j] = li[j], li[i] }
func (li listItemSort) Less(i, j int) bool {
	if li[i].Status == li[j].Status {
		return li[i].Path < li[j].Path
	}
	return li[i].Status > li[j].Status
}

// ListStatus obtains the current package status list.
func (ctx *Context) ListStatus() ([]ListItem, error) {
	var err error
	if !ctx.loaded {
		err = ctx.loadPackage()
		if err != nil {
			return nil, err
		}
	}
	list := make([]ListItem, 0, len(ctx.Package))
	for _, pkg := range ctx.Package {
		li := ListItem{
			Status:     pkg.Status,
			Path:       pkg.CanonicalPath,
			VendorPath: pkg.LocalPath,
		}
		list = append(list, li)
	}
	// Sort li by Status, then Path.
	sort.Sort(listItemSort(list))

	return list, nil
}
