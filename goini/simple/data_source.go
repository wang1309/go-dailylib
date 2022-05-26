package simple

import (
	"fmt"
	"io"
	"os"
)

// dataSource is an interface that returns object which can be read and closed.
type dataSource interface {
	ReadCloser() (io.ReadCloser, error)
}


// sourceFile represents an object that contains content on the local file system.
type sourceFile struct {
	name string
}

func (s sourceFile) ReadCloser() (_ io.ReadCloser, err error) {
	return os.Open(s.name)
}

func parseDataSource(source interface{}) (dataSource, error) {
	switch s := source.(type) {
	case string:
		return sourceFile{s}, nil
	/*case []byte:
		return &sourceData{s}, nil
	case io.ReadCloser:
		return &sourceReadCloser{s}, nil
	case io.Reader:
		return &sourceReadCloser{ioutil.NopCloser(s)}, nil*/
	default:
		return nil, fmt.Errorf("error parsing data source: unknown type %q", s)
	}
}