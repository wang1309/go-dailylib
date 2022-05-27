package flow

import (
	"github.com/chrislusf/gleam/instruction"
	"github.com/chrislusf/gleam/pb"
	"github.com/chrislusf/gleam/script"
	"github.com/chrislusf/gleam/util"
	"io"
	"sync"
	"time"
)

const (
	OneShardToOneShard NetworkType = iota
	OneShardToAllShard
	AllShardToOneShard
	OneShardToEveryNShard
	LinkedNShardToOneShard
	MergeTwoShardToOneShard
	AllShardTOAllShard
)

type DatasetShardStatus int
type NetworkType int
type ModeIO int

type StepMetadata struct {
	IsRestartable bool
	IsIdempotent  bool
}


type Flow struct {
	Name     string
	Steps    []*Step
	Datasets []*Dataset
	HashCode uint32
}

type DasetsetMetadata struct {
	TotalSize int64
	OnDisk    ModeIO
}

type Step struct {
	Id             int
	Flow           *Flow
	InputDatasets  []*Dataset
	OutputDataset  *Dataset
	Function       func([]io.Reader, []io.Writer, *pb.InstructionStat) error
	Instruction    instruction.Instruction
	Tasks          []*Task
	Name           string
	Description    string
	NetworkType    NetworkType
	IsOnDriverSide bool
	IsPipe         bool
	IsGoCode       bool
	Script         script.Script
	Command        *script.Command // used in Pipe()
	Meta           *StepMetadata
	Params         map[string]interface{}
	RunLocked
}

type Dataset struct {
	Flow            *Flow
	Id              int
	Shards          []*DatasetShard
	Step            *Step
	ReadingSteps    []*Step
	IsPartitionedBy []int
	IsLocalSorted   []instruction.OrderBy
	Meta            *DasetsetMetadata
	RunLocked
}

type DatasetShard struct {
	Id            int
	Dataset       *Dataset
	ReadingTasks  []*Task
	IncomingChan  *util.Piper
	OutgoingChans []*util.Piper
	Counter       int64
	ReadyTime     time.Time
	CloseTime     time.Time
	Meta          *DasetsetShardMetadata
}

type Task struct {
	Id           int
	Step         *Step
	InputShards  []*DatasetShard
	InputChans   []*util.Piper // task specific input chans. InputShard may have multiple reading tasks
	OutputShards []*DatasetShard
	Stat         *pb.InstructionStat
}

type DasetsetShardMetadata struct {
	TotalSize int64
	Timestamp time.Time
	URI       string
	Status    DatasetShardStatus
	Error     error
}

type RunLocked struct {
	sync.Mutex
	StartTime time.Time
}