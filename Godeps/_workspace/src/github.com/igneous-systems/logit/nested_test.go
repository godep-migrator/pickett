// Copyright (C) 2010, Kyle Lemons <kyle@kylelemons.net>.  All rights reserved.

package logit

import (
	"testing"
)

func TestNestedLoggerNoFilterCreate(t *testing.T) {
	sl := NewDefaultLogger(WARNING)
	if sl == nil {
		t.Fatalf("NewDefaultLogger should never return nil")
	}
	defer sl.Close(-1)
	nest1 := NewNestedLogger("my/log/program", sl)
	if nest1 == nil {
		t.Fatalf("NewNestedLogger should never return nil")
	}
	if nest1.master == nil {
		t.Fatalf("nil master for nested logger")
	}
	defer nest1.Close(-1)

	list := nest1.getFilterRef(false)
	if list == nil {
		t.Fatalf("Nil filter list")
	}
	if len(list) != 1 {
		t.Fatalf("NewNestedLogger produced invalid logger (incorrect map count) %v", len(list))
	}
	if list[0].name != "stdout" {
		t.Fatalf("NewNestedLogger produced invalid default name %v", list[0].name)
	}
	if list[0].filter == nil {
		t.Fatalf("NewNestedLogger produced nil logger ")
	}
	if list[0].lvl != WARNING {
		t.Fatalf("NewNestedLogger produced invalid logger (incorrect level) %d", list[0].lvl)
	}
	if list[0].nd != nil {
		t.Fatalf("NewNestedLogger produced invalid logger (non nil nd)")
	}
	if nest1.version != nest1.master.version {
		t.Fatalf("Mismatched initial versions %d %d", nest1.version, nest1.master.version)
	}

	oldversion := nest1.version

	sl.ModifyFilterLvl("stdout", INFO, nil, nil)
	if nest1.master.version == oldversion {
		t.Fatalf("Modify filter didn't change version")
	}
	list = nest1.getFilterRef(false)
	if list == nil {
		t.Fatalf("Nil filter list")
	}
	if len(list) != 1 {
		t.Fatalf("NewNestedLogger produced invalid logger (incorrect map count) %v", len(list))
	}
	if list[0].name != "stdout" {
		t.Fatalf("NewNestedLogger produced invalid default name %v", list[0].name)
	}
	if list[0].filter == nil {
		t.Fatalf("NewNestedLogger produced nil logger ")
	}
	if list[0].lvl != INFO {
		t.Fatalf("NewNestedLogger produced invalid logger (incorrect level) %d", list[0].lvl)
	}
	if list[0].nd != nil {
		t.Fatalf("NewNestedLogger produced invalid logger (non nil nd)")
	}
	if nest1.version != nest1.master.version {
		t.Fatalf("Mismatched initial versions %d %d", nest1.version, nest1.master.version)
	}
}

func TestNestedLoggerFilterMatch(t *testing.T) {
	sl := NewDefaultLogger(WARNING)
	defer sl.Close(-1)

	nd := NewNestedData()
	nd.AddEntry("this/is", INFO, nil)
	sl.ModifyFilterLvl("stdout", WARNING, nil, nd)

	nest1 := NewNestedLogger("this/is/my", sl)
	if nest1 == nil {
		t.Fatalf("NewNestedLogger should never return nil")
	}
	defer nest1.Close(-1)

	list := nest1.getFilterRef(false)
	if list == nil {
		t.Fatalf("Nil filter list")
	}
	if len(list) != 1 {
		t.Fatalf("NewNestedLogger produced invalid logger (incorrect map count) %v", len(list))
	}
	if list[0].name != "stdout" {
		t.Fatalf("NewNestedLogger produced invalid default name %v", list[0].name)
	}
	if list[0].filter == nil {
		t.Fatalf("NewNestedLogger produced nil logger ")
	}
	if list[0].lvl != INFO {
		t.Fatalf("NewNestedLogger produced invalid logger (incorrect level) %v", list[0].lvl)
	}
	if list[0].nd != nil {
		t.Fatalf("NewNestedLogger produced invalid logger (non nil nd)")
	}
	if nest1.version != nest1.master.version {
		t.Fatalf("Mismatched initial versions %d %d", nest1.version, nest1.master.version)
	}

	oldversion := nest1.version
	nd = NewNestedData()
	nd.AddEntry("this/is", DEBUG, nil)
	sl.ModifyFilterLvl("stdout", ERROR, nil, nd)

	if nest1.master.version == oldversion {
		t.Fatalf("Modify filter didn't change version")
	}
	list = nest1.getFilterRef(false)
	if list == nil {
		t.Fatalf("Nil filter list")
	}
	if len(list) != 1 {
		t.Fatalf("NewNestedLogger produced invalid logger (incorrect map count) %v", len(list))
	}
	if list[0].name != "stdout" {
		t.Fatalf("NewNestedLogger produced invalid default name %v", list[0].name)
	}
	if list[0].filter == nil {
		t.Fatalf("NewNestedLogger produced nil logger ")
	}
	if list[0].lvl != DEBUG {
		t.Fatalf("NewNestedLogger produced invalid logger (incorrect level) %d", list[0].lvl)
	}
	if list[0].nd != nil {
		t.Fatalf("NewNestedLogger produced invalid logger (non nil nd)")
	}
	if nest1.version != nest1.master.version {
		t.Fatalf("Mismatched initial versions %d %d", nest1.version, nest1.master.version)
	}
}

func TestNestedLoggerFilterNoMatch(t *testing.T) {
	sl := NewDefaultLogger(WARNING)
	defer sl.Close(-1)

	nd := NewNestedData()
	nd.AddEntry("bad/is", INFO, nil)
	sl.ModifyFilterLvl("stdout", WARNING, nil, nd)

	nest1 := NewNestedLogger("this/is/my", sl)
	if nest1 == nil {
		t.Fatalf("NewNestedLogger should never return nil")
	}
	defer nest1.Close(-1)

	list := nest1.getFilterRef(false)
	if list == nil {
		t.Fatalf("Nil filter list")
	}
	if len(list) != 1 {
		t.Fatalf("NewNestedLogger produced invalid logger (incorrect map count) %v", len(list))
	}
	if list[0].name != "stdout" {
		t.Fatalf("NewNestedLogger produced invalid default name %v", list[0].name)
	}
	if list[0].filter == nil {
		t.Fatalf("NewNestedLogger produced nil logger ")
	}
	if list[0].lvl != WARNING {
		t.Fatalf("NewNestedLogger produced invalid logger (incorrect level) %v", list[0].lvl)
	}
	if list[0].nd != nil {
		t.Fatalf("NewNestedLogger produced invalid logger (non nil nd)")
	}
	if nest1.version != nest1.master.version {
		t.Fatalf("Mismatched initial versions %d %d", nest1.version, nest1.master.version)
	}

	nd = NewNestedData()
	nd.AddEntry("this/i", INFO, nil)
	sl.ModifyFilterLvl("stdout", WARNING, nil, nd)

	list = nest1.getFilterRef(false)
	if list == nil {
		t.Fatalf("Nil filter list")
	}
	if len(list) != 1 {
		t.Fatalf("NewNestedLogger produced invalid logger (incorrect map count) %v", len(list))
	}
	if list[0].name != "stdout" {
		t.Fatalf("NewNestedLogger produced invalid default name %v", list[0].name)
	}
	if list[0].filter == nil {
		t.Fatalf("NewNestedLogger produced nil logger ")
	}
	if list[0].lvl != WARNING {
		t.Fatalf("NewNestedLogger produced invalid logger (incorrect level) %v", list[0].lvl)
	}
	if list[0].nd != nil {
		t.Fatalf("NewNestedLogger produced invalid logger (non nil nd)")
	}
	if nest1.version != nest1.master.version {
		t.Fatalf("Mismatched initial versions %d %d", nest1.version, nest1.master.version)
	}
}
