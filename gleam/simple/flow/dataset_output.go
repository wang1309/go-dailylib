package flow

import (
	"bufio"
	"fmt"
	"github.com/chrislusf/gleam/pb"
	"github.com/chrislusf/gleam/util"
	"io"
	"os"
)

// Printlnf prints to os.Stdout in the specified format,
// adding an "\n" at the end of each format
func (d *Dataset) Printlnf(format string) *Dataset {
	return d.Fprintf(os.Stdout, format+"\n")
}


// Fprintf formats using the format for each row and writes to writer.
func (d *Dataset) Fprintf(writer io.Writer, format string) *Dataset {
	fn := func(r io.Reader) error {
		w := bufio.NewWriter(writer)
		defer w.Flush()
		return util.Fprintf(w, r, format)
	}
	return d.Output(fn)
}

// Output concurrently collects outputs from previous step to the driver.
func (d *Dataset) Output(f func(io.Reader) error) *Dataset {
	step := d.Flow.AddAllToOneStep(d, nil)
	step.IsOnDriverSide = true
	step.Name = "Output"
	step.Function = func(readers []io.Reader, writers []io.Writer, stat *pb.InstructionStat) error {
		errChan := make(chan error, len(readers))
		for i, reader := range readers {
			go func(i int, reader io.Reader) {
				errChan <- f(reader)
			}(i, reader)
		}
		for range readers {
			err := <-errChan
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to process output: %v\n", err)
				return err
			}
		}
		return nil
	}
	return d
}