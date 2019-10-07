package main

import (
	"fmt"
	"sort"
	"strings"
)

type IOStatWidget struct{}

type TableRegionStatusList []TableRegionStatus

func (c TableRegionStatusList) Len() int {
	return len(c)
}
func (c TableRegionStatusList) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
func (c TableRegionStatusList) Less(i, j int) bool {
	if c[i].wbytes < c[j].wbytes {
		return true
	} else if c[i].wbytes == c[j].wbytes {
		if c[i].rbytes < c[j].rbytes {
			return true
		} else if c[i].rbytes == c[j].rbytes {
			switch strings.Compare(c[i].String(), c[j].String()) {
			case 1:
				return true
			case -1:
				return false
			}
		}
	}
	return false
}

func newIOStatWidget() Widget {
	return &IOStatWidget{}
}

func (c *IOStatWidget) GetText() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Top hotspots\n")
	if list, ok := Stat().Load(TABLES_IO_STATUS); ok {
		sort.Sort(list.(TableRegionStatusList))
		for _, r := range list.(TableRegionStatusList) {
			fmt.Fprintf(&sb, "%s\n", r)
		}
	}
	sb.WriteString("\n")
	return sb.String()
}
