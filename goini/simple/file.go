package simple

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
)

type File struct {
	options     LoadOptions
	dataSources []dataSource

	// Should make things safe, but sometimes doesn't matter.
	BlockMode bool
	lock      sync.RWMutex

	// To keep data in order.
	sectionList []string
	// To keep track of the index of a section with same name.
	// This meta list is only used with non-unique section names are allowed.
	sectionIndexes []int

	// Actual data is stored here.
	sections map[string][]*Section

	ValueMapper
}


// newFile initializes File object with given data sources.
func newFile(dataSources []dataSource, opts LoadOptions) *File {
	if len(opts.KeyValueDelimiters) == 0 {
		opts.KeyValueDelimiters = "=:"
	}
	if len(opts.KeyValueDelimiterOnWrite) == 0 {
		opts.KeyValueDelimiterOnWrite = "="
	}
	if len(opts.ChildSectionDelimiter) == 0 {
		opts.ChildSectionDelimiter = "."
	}

	return &File{
		BlockMode:   true,
		dataSources: dataSources,
		sections:    make(map[string][]*Section),
		options:     opts,
	}
}

// Reload reloads and parses all data sources.
func (f *File) Reload() (err error) {
	for _, s := range f.dataSources {
		if err = f.reload(s); err != nil {
			// In loose mode, we create an empty default section for nonexistent files.
			if os.IsNotExist(err) && f.options.Loose {
				_ = f.parse(bytes.NewBuffer(nil))
				continue
			}
			return err
		}

		if f.options.ShortCircuit {
			return nil
		}
	}

	return nil
}


func (f *File) reload(s dataSource) error {
	r, err := s.ReadCloser()
	if err != nil {
		return err
	}
	defer r.Close()

	return f.parse(r)
}

// NewSection creates a new section.
func (f *File) NewSection(name string) (*Section, error) {
	if len(name) == 0 {
		return nil, errors.New("empty section name")
	}

	if (f.options.Insensitive || f.options.InsensitiveSections) && name != DefaultSection {
		name = strings.ToLower(name)
	}

	if f.BlockMode {
		f.lock.Lock()
		defer f.lock.Unlock()
	}

	if !f.options.AllowNonUniqueSections && inSlice(name, f.sectionList) {
		return f.sections[name][0], nil
	}

	f.sectionList = append(f.sectionList, name)

	// NOTE: Append to indexes must happen before appending to sections,
	// otherwise index will have off-by-one problem.
	f.sectionIndexes = append(f.sectionIndexes, len(f.sections[name]))

	sec := newSection(f, name)
	f.sections[name] = append(f.sections[name], sec)

	return sec, nil
}

// GetSection returns section by given name.
func (f *File) GetSection(name string) (*Section, error) {
	secs, err := f.SectionsByName(name)
	if err != nil {
		return nil, err
	}

	return secs[0], err
}


// SectionsByName returns all sections with given name.
func (f *File) SectionsByName(name string) ([]*Section, error) {
	if len(name) == 0 {
		name = DefaultSection
	}
	if f.options.Insensitive || f.options.InsensitiveSections {
		name = strings.ToLower(name)
	}

	if f.BlockMode {
		f.lock.RLock()
		defer f.lock.RUnlock()
	}

	secs := f.sections[name]
	if len(secs) == 0 {
		return nil, fmt.Errorf("section %q does not exist", name)
	}

	return secs, nil
}
