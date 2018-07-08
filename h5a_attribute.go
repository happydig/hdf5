// Copyright ©2017 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hdf5

// #include "hdf5.h"
// #include <stdlib.h>
// #include <string.h>
import "C"

import (
	"reflect"
	"unsafe"
)

type Attribute struct {
	Identifier
}

func newAttribute(id C.hid_t) *Attribute {
	return &Attribute{Identifier{id}}
}

func createAttribute(id C.hid_t, name string, dType *Datatype, dSpace *Dataspace, acpl *PropList) (*Attribute, error) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	hid := C.H5Acreate2(id, cName, dType.id, dSpace.id, acpl.id, P_DEFAULT.id)
	if err := checkId(hid); err != nil {
		return nil, err
	}
	return newAttribute(hid), nil
}

func openAttribute(id C.hid_t, name string) (*Attribute, error) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	hid := C.H5Aopen(id, cName, P_DEFAULT.id)
	if err := checkId(hid); err != nil {
		return nil, err
	}
	return newAttribute(hid), nil
}

// Access the type of an attribute
func (s *Attribute) GetType() Identifier {
	ftype := C.H5Aget_type(s.id)
	return Identifier{ftype}
}

// Close releases and terminates access to an attribute.
func (s *Attribute) Close() error {
	return s.closeWith(h5aclose)
}

func h5aclose(id C.hid_t) C.herr_t {
	return C.H5Aclose(id)
}

// Space returns an identifier for a copy of the dataspace for a attribute.
func (s *Attribute) Space() *Dataspace {
	hid := C.H5Aget_space(s.id)
	if int(hid) > 0 {
		return newDataspace(hid)
	}
	return nil
}

// Read reads raw data from a attribute into a buffer.
func (s *Attribute) Read(data interface{}, dType *Datatype) error {
	var addr unsafe.Pointer
	v := reflect.ValueOf(data)

	switch v.Kind() {

	case reflect.Array:
		addr = unsafe.Pointer(v.UnsafeAddr())

	case reflect.String:
		str := (*reflect.StringHeader)(unsafe.Pointer(v.UnsafeAddr()))
		addr = unsafe.Pointer(str.Data)

	case reflect.Ptr:
		addr = unsafe.Pointer(v.Pointer())

	default:
		addr = unsafe.Pointer(v.UnsafeAddr())
	}

	rc := C.H5Aread(s.id, dType.id, addr)
	err := h5err(rc)
	return err
}

// Write writes raw data from a buffer to an attribute.
func (s *Attribute) Write(data interface{}, dType *Datatype) error {
	var addr unsafe.Pointer
	v := reflect.Indirect(reflect.ValueOf(data))
	switch v.Kind() {

	case reflect.Array:
		addr = unsafe.Pointer(v.UnsafeAddr())

	case reflect.String:
		str := C.CString(v.Interface().(string))
		defer C.free(unsafe.Pointer(str))
		addr = unsafe.Pointer(&str)

	case reflect.Ptr:
		addr = unsafe.Pointer(v.Pointer())

	default:
		addr = unsafe.Pointer(v.UnsafeAddr())
	}

	rc := C.H5Awrite(s.id, dType.id, addr)
	err := h5err(rc)
	return err
}
