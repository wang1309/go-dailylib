package simple

import (
	"errors"
	"fmt"
	"strings"
)

// Section represents a config section.
type Section struct {
	f        *File
	Comment  string
	name     string
	keys     map[string]*Key
	keyList  []string
	keysHash map[string]string

	isRawSection bool
	rawBody      string
}


func newSection(f *File, name string) *Section {
	return &Section{
		f:        f,
		name:     name,
		keys:     make(map[string]*Key),
		keyList:  make([]string, 0, 10),
		keysHash: make(map[string]string),
	}
}

// NewBooleanKey creates a new boolean type key to given section.
func (s *Section) NewBooleanKey(name string) (*Key, error) {
	key, err := s.NewKey(name, "true")
	if err != nil {
		return nil, err
	}

	key.isBooleanType = true
	return key, nil
}

// NewKey creates a new key to given section.
func (s *Section) NewKey(name, val string) (*Key, error) {
	if len(name) == 0 {
		return nil, errors.New("error creating new key: empty key name")
	} else if s.f.options.Insensitive || s.f.options.InsensitiveKeys {
		name = strings.ToLower(name)
	}

	if s.f.BlockMode {
		s.f.lock.Lock()
		defer s.f.lock.Unlock()
	}

	if inSlice(name, s.keyList) {
		if s.f.options.AllowShadows {
			if err := s.keys[name].addShadow(val); err != nil {
				return nil, err
			}
		} else {
			s.keys[name].value = val
			s.keysHash[name] = val
		}
		return s.keys[name], nil
	}

	s.keyList = append(s.keyList, name)
	s.keys[name] = newKey(s, name, val)
	s.keysHash[name] = val
	return s.keys[name], nil
}


// Section assumes named section exists and returns a zero-value when not.
func (f *File) Section(name string) *Section {
	sec, err := f.GetSection(name)
	if err != nil {
		if name == "" {
			name = DefaultSection
		}
		sec, _ = f.NewSection(name)
		return sec
	}
	return sec
}

// Key assumes named Key exists in section and returns a zero-value when not.
func (s *Section) Key(name string) *Key {
	key, err := s.GetKey(name)
	if err != nil {
		// It's OK here because the only possible error is empty key name,
		// but if it's empty, this piece of code won't be executed.
		key, _ = s.NewKey(name, "")
		return key
	}
	return key
}

// GetKey returns key in section by given name.
func (s *Section) GetKey(name string) (*Key, error) {
	if s.f.BlockMode {
		s.f.lock.RLock()
	}
	if s.f.options.Insensitive || s.f.options.InsensitiveKeys {
		name = strings.ToLower(name)
	}
	key := s.keys[name]
	if s.f.BlockMode {
		s.f.lock.RUnlock()
	}

	if key == nil {
		// Check if it is a child-section.
		sname := s.name
		for {
			if i := strings.LastIndex(sname, s.f.options.ChildSectionDelimiter); i > -1 {
				sname = sname[:i]
				sec, err := s.f.GetSection(sname)
				if err != nil {
					continue
				}
				return sec.GetKey(name)
			}
			break
		}
		return nil, fmt.Errorf("error when getting key of section %q: key %q not exists", s.name, name)
	}
	return key, nil
}