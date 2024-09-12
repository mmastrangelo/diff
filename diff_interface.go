/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package diff

import "reflect"

func (d *Differ) diffInterface(path []string, a, b reflect.Value, parent interface{}) error {
	if a.Kind() == reflect.Invalid {
		d.cl.Add(CREATE, path, nil, exportInterface(b))
		return nil
	}

	if b.Kind() == reflect.Invalid {
		d.cl.Add(DELETE, path, exportInterface(a), nil)
		return nil
	}

	if a.Kind() != b.Kind() {
		return ErrTypeMismatch
	}

	if a.IsNil() && b.IsNil() {
		return nil
	}

	if a.IsNil() {
		// if b is a slice of map[string]interface{}, set a to an empty slice of the same type
		// This allows us to identify a change from a nil slice to a populated slice as create operations
		// rather than an update of the entire slice
		aUpdated := false
		if b.Elem().Kind() == reflect.Slice {
			if b.Elem().Len() == 0 {
				return nil
			} else if b.Elem().Index(0).Elem().Kind() == reflect.Map {
				var aVal interface{}
				a = reflect.ValueOf(&aVal).Elem()
				a.Set(reflect.MakeSlice(b.Elem().Type(), 0, 0))

				aUpdated = true
			}
		}

		if !aUpdated {
			d.cl.Add(UPDATE, path, nil, exportInterface(b), parent)
			return nil
		}
	}

	if b.IsNil() {
		// if a is a slice of map[string]interface{}, set b to an empty slice of the same type
		// This allows us to identify a change from a populated slice to a nil slice as delete operations
		// rather than an update of the entire slice
		bUpdated := false
		if a.Elem().Kind() == reflect.Slice {
			if a.Elem().Len() == 0 {
				return nil
			} else if a.Elem().Index(0).Elem().Kind() == reflect.Map {
				var bVal interface{}
				b = reflect.ValueOf(&bVal).Elem()
				b.Set(reflect.MakeSlice(a.Elem().Type(), 0, 0))

				bUpdated = true
			}
		}

		if !bUpdated {
			d.cl.Add(UPDATE, path, exportInterface(a), nil, parent)
			return nil
		}
	}

	return d.diff(path, a.Elem(), b.Elem(), parent)
}
