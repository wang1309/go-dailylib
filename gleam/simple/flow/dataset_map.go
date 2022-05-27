package flow

import (
	"dailylib/gleam/simple"
	"github.com/chrislusf/gleam/script"
	"os"
)

// Mapper runs the mapper registered to the mapperId.
// This is used to execute pure Go code.
func (d *Dataset) Map(name string, mapperId simple.MapperId) *Dataset {
	ret, step := add1ShardTo1Step(d)
	step.Name = name + ".Map"
	step.IsPipe = false
	step.IsGoCode = true

	ex, _ := os.Executable()

	mapper, _ := simple.GetMapper(mapperId)
	step.Description = mapper.Name

	var args []string
	args = append(args, os.Args[1:]...)
	args = append(args, "-gleam.mapper", string(mapperId))
	step.Command = &script.Command{
		Path: ex,
		Args: args,
	}
	return ret
}


func add1ShardTo1Step(d *Dataset) (ret *Dataset, step *Step) {
	ret = d.Flow.NewNextDataset(len(d.Shards))
	step = d.Flow.AddOneToOneStep(d, ret)
	return
}