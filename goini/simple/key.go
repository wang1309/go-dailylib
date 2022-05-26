package simple

import (
	"errors"
	"strings"
)

// ValueMapper represents a mapping function for values, e.g. os.ExpandEnv
type ValueMapper func(string) string


// Key represents a key under a section.
type Key struct {
	s               *Section
	Comment         string
	name            string
	value           string
	isAutoIncrement bool
	isBooleanType   bool

	isShadow bool
	shadows  []*Key

	nestedValues []string
}

// newKey simply return a key object with given values.
func newKey(s *Section, name, val string) *Key {
	return &Key{
		s:     s,
		name:  name,
		value: val,
	}
}


func (k *Key) addNestedValue(val string) error {
	if k.isAutoIncrement || k.isBooleanType {
		return errors.New("cannot add nested value to auto-increment or boolean key")
	}

	k.nestedValues = append(k.nestedValues, val)
	return nil
}


func (k *Key) addShadow(val string) error {
	if k.isShadow {
		return errors.New("cannot add shadow to another shadow key")
	} else if k.isAutoIncrement || k.isBooleanType {
		return errors.New("cannot add shadow to auto-increment or boolean key")
	}

	if !k.s.f.options.AllowDuplicateShadowValues {
		// Deduplicate shadows based on their values.
		if k.value == val {
			return nil
		}
		for i := range k.shadows {
			if k.shadows[i].value == val {
				return nil
			}
		}
	}

	shadow := newKey(k.s, k.name, val)
	shadow.isShadow = true
	k.shadows = append(k.shadows, shadow)
	return nil
}

// transformValue takes a raw value and transforms to its final string.
func (k *Key) transformValue(val string) string {
	if k.s.f.ValueMapper != nil {
		val = k.s.f.ValueMapper(val)
	}

	// Fail-fast if no indicate char found for recursive value
	if !strings.Contains(val, "%") {
		return val
	}
	for i := 0; i < depthValues; i++ {
		vr := varPattern.FindString(val)
		if len(vr) == 0 {
			break
		}

		// Take off leading '%(' and trailing ')s'.
		noption := vr[2 : len(vr)-2]

		// Search in the same section.
		// If not found or found the key itself, then search again in default section.
		nk, err := k.s.GetKey(noption)
		if err != nil || k == nk {
			nk, _ = k.s.f.Section("").GetKey(noption)
			if nk == nil {
				// Stop when no results found in the default section,
				// and returns the value as-is.
				break
			}
		}

		// Substitute by new value and take off leading '%(' and trailing ')s'.
		val = strings.Replace(val, vr, nk.value, -1)
	}
	return val
}

// String returns string representation of value.
func (k *Key) String() string {
	return k.transformValue(k.value)
}