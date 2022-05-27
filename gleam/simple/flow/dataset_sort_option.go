package flow

import (
	"fmt"
	"github.com/chrislusf/gleam/instruction"
	"strings"
)

type SortOption struct {
	orderByList []instruction.OrderBy
}

func Field(indexes ...int) *SortOption {
	ret := &SortOption{}
	for _, index := range indexes {
		ret.orderByList = append(ret.orderByList, instruction.OrderBy{
			Index: index,
			Order: instruction.Ascending,
		})
	}
	return ret
}


func (o *SortOption) String() string {
	var buf strings.Builder
	for _, orderBy := range o.orderByList {
		buf.WriteString(fmt.Sprintf("%d ", orderBy.Index))
		if orderBy.Order == instruction.Ascending {
			buf.WriteString("asc")
		} else {
			buf.WriteString("desc")
		}
	}
	return buf.String()
}


// return a list of indexes
func (o *SortOption) Indexes() []int {
	var ret []int
	for _, x := range o.orderByList {
		ret = append(ret, x.Index)
	}
	return ret
}


func OrderBy(index int, ascending bool) *SortOption {
	ret := &SortOption{
		orderByList: []instruction.OrderBy{
			{
				Index: index,
				Order: instruction.Descending,
			},
		},
	}
	if ascending {
		ret.orderByList[0].Order = instruction.Ascending
	}
	return ret
}


// LocalLimit take the local first n rows and skip all other rows.
func (d *Dataset) LocalLimit(name string, n int, offset int) *Dataset {
	ret, step := add1ShardTo1Step(d)
	ret.IsLocalSorted = d.IsLocalSorted
	ret.IsPartitionedBy = d.IsPartitionedBy
	step.SetInstruction(name, instruction.NewLocalLimit(n, offset))
	step.Description = fmt.Sprintf("local limit %d", n)
	return ret
}