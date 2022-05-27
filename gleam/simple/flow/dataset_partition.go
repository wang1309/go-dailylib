package flow

import "github.com/chrislusf/gleam/instruction"

func (d *Dataset) RoundRobin(name string, n int) *Dataset {
	if n <= 1 {
		return d
	}
	shard := n * len(d.Shards)
	ret := d.Flow.NewNextDataset(shard)
	step := d.Flow.AddAllToAllStep(d, ret)
	step.SetInstruction(name, instruction.NewRoundRobin())
	return ret
}
