// Code generated by capnpc-go. DO NOT EDIT.

package test

import (
	capnp "capnproto.org/go/capnp/v3"
	text "capnproto.org/go/capnp/v3/encoding/text"
	schemas "capnproto.org/go/capnp/v3/schemas"
)

type Test struct{ capnp.Struct }

// Test_TypeID is the unique identifier for the type Test.
const Test_TypeID = 0xde17d5ca34295e24

func NewTest(s *capnp.Segment) (Test, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	return Test{st}, err
}

func NewRootTest(s *capnp.Segment) (Test, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	return Test{st}, err
}

func ReadRootTest(msg *capnp.Message) (Test, error) {
	root, err := msg.Root()
	return Test{root.Struct()}, err
}

func (s Test) String() string {
	str, _ := text.Marshal(0xde17d5ca34295e24, s.Struct)
	return str
}

func (s Test) Id() (string, error) {
	p, err := s.Struct.Ptr(0)
	return p.Text(), err
}

func (s Test) HasId() bool {
	return s.Struct.HasPtr(0)
}

func (s Test) IdBytes() ([]byte, error) {
	p, err := s.Struct.Ptr(0)
	return p.TextBytes(), err
}

func (s Test) SetId(v string) error {
	return s.Struct.SetText(0, v)
}

// Test_List is a list of Test.
type Test_List struct{ capnp.List }

// NewTest creates a new list of Test.
func NewTest_List(s *capnp.Segment, sz int32) (Test_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1}, sz)
	return Test_List{l}, err
}

func (s Test_List) At(i int) Test { return Test{s.List.Struct(i)} }

func (s Test_List) Set(i int, v Test) error { return s.List.SetStruct(i, v.Struct) }

func (s Test_List) String() string {
	str, _ := text.MarshalList(0xde17d5ca34295e24, s.List)
	return str
}

// Test_Future is a wrapper for a Test promised by a client call.
type Test_Future struct{ *capnp.Future }

func (p Test_Future) Struct() (Test, error) {
	s, err := p.Future.Struct()
	return Test{s}, err
}

const schema_d856645e9d045ec3 = "x\xda\x12\xf0s`2d\xdd\xcf\xc8\xc0\x10(\xc2\xca" +
	"\xf6_%N\xd3\xe4\xd4U\xf1{\x0c\x82\xdc\x8c\xff\x0f" +
	"\xc7\xb1\xcc\x8dK\x09\xbb\xc1\xc0\xca\xc8\xce\xc0 x\xb4" +
	"I\xf0$\x98\xb6g\xd0\xfd_\x92Z\\\xa2\x97\x9cX" +
	"\xc0\x98W`\x15\x92Z\\\xc2\x10\xc0\xc8\x18\xc8\xc2\xcc" +
	"\xc2\xc0\xc0\xc2\xc8\xc0 \xc8+\xc5\xc0\x10\xc8\xc1\xcc\x18" +
	"(\xc2\xc4\xc8\x9c\x99\xc2\xc8\xc3\xc0\xc4\xc8\xc3\xc0\x08\x08" +
	"\x00\x00\xff\xffC\x93\x18F"

func init() {
	schemas.Register(schema_d856645e9d045ec3,
		0xde17d5ca34295e24)
}
