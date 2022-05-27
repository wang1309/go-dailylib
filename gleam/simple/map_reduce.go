package simple

import (
	"flag"
	"fmt"
	"github.com/chrislusf/gleam/pb"
	"os"
	"reflect"
	"runtime"
	"sync"
)


type MapperId string
type ReducerId string
type Mapper func([]interface{}) error
type Reducer func(x, y interface{}) (interface{}, error)


type MapperObject struct {
	Mapper Mapper
	Name string
}

type ReducerObject struct {
	Reducer Reducer
	Name    string
}


var (
	HasInitalized bool

	taskOption gleamTaskOption

	stat = &pb.ExecutionStat{} // TsEmit() needs this global value

	mappers      map[MapperId]MapperObject
	reducers     map[ReducerId]ReducerObject

	mappersLock  sync.Mutex
	reducersLock sync.Mutex
)

func init() {
	mappers = make(map[MapperId]MapperObject)
	reducers = make(map[ReducerId]ReducerObject)
}

type gleamTaskOption struct {
	Mapper          string
	Reducer         string
	KeyFields       string
	ExecutorAddress string
	HashCode        uint
	StepId          int
	TaskId          int
	IsProfiling     bool
}

type gleamRunner struct {
	Option *gleamTaskOption
}

func init() {
	flag.StringVar(&taskOption.Mapper, "gleam.mapper", "", "the generated mapper name")
	flag.StringVar(&taskOption.Reducer, "gleam.reducer", "", "the generated reducer name")
	flag.StringVar(&taskOption.KeyFields, "gleam.keyFields", "", "the 1-based key fields")
	flag.StringVar(&taskOption.ExecutorAddress, "gleam.executor", "", "executor address")
	flag.UintVar(&taskOption.HashCode, "flow.hashcode", 0, "flow hashcode")
	flag.IntVar(&taskOption.StepId, "flow.stepId", -1, "flow step id")
	flag.IntVar(&taskOption.TaskId, "flow.taskId", -1, "flow task id")
	flag.BoolVar(&taskOption.IsProfiling, "gleam.profiling", false, "profiling all steps")
}

// Init determines whether the driver program will execute the mapper/reducer or not.
// If the command line invokes the mapper or reducer, execute it and exit.
// This function will invoke flag.Parse() first.
func Init() {
	HasInitalized = true

	flag.Parse()

	if taskOption.Mapper != "" || taskOption.Reducer != "" {
		runner := &gleamRunner{Option: &taskOption}
		runner.runMapperReducer()
		os.Exit(0)
	}
}

func GetMapper(mapperId MapperId) (mapper MapperObject, found bool) {
	mappersLock.Lock()
	defer mappersLock.Unlock()

	mapper, found = mappers[mapperId]

	return
}

// RegisterMapper register a mapper function to process a command
func RegisterMapper(fn Mapper) MapperId {
	mappersLock.Lock()
	defer mappersLock.Unlock()

	mapperId := MapperId(fmt.Sprintf("m%d", len(mappers)+1))
	mappers[mapperId] = MapperObject{fn, runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()}

	return mapperId
}

func GetReducer(reducerId ReducerId) (reducer ReducerObject, found bool) {
	reducersLock.Lock()
	defer reducersLock.Unlock()

	reducer, found = reducers[reducerId]

	return
}

func RegisterReducer(fn Reducer) ReducerId {
	reducersLock.Lock()
	defer reducersLock.Unlock()

	reducerId := ReducerId(fmt.Sprintf("r%d", len(reducers)+1))
	reducers[reducerId] = ReducerObject{fn, runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()}

	return reducerId
}
