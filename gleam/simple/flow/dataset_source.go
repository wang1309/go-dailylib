package flow

import (
	"github.com/chrislusf/gleam/pb"
	"io"
)

type Sourcer interface {
	Generate(*Flow) *Dataset
}

// Read accepts a function to read data into the flow, creating a new dataset.
// This allows custom complicated pre-built logic for new data sources.
func (fc *Flow) Read(s Sourcer) (ret *Dataset) {
	return s.Generate(fc)
}

// Source produces data feeding into the flow.
// Function f writes to this writer.
// The written bytes should be MsgPack encoded []byte.
// Use util.EncodeRow(...) to encode the data before sending to this channel
func (fc *Flow) Source(name string, f func(io.Writer, *pb.InstructionStat) error) (ret *Dataset) {
	ret = fc.NewNextDataset(1)
	step := fc.AddOneToOneStep(nil, ret)
	step.IsOnDriverSide = true
	step.Name = name
	step.Function = func(readers []io.Reader, writers []io.Writer, stats *pb.InstructionStat) error {
		errChan := make(chan error, len(writers))
		for _, writer := range writers {
			go func(writer io.Writer) {
				errChan <- f(writer, stats)
			}(writer)
		}
		for range writers {
			err := <-errChan
			if err != nil {
				return err
			}
		}
		return nil
	}
	return
}
