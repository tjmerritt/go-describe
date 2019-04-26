package describe

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
	"unsafe"
)

type Foo int

type Iface interface {
	Method(int) string
}

type Obj struct {
	Field int
}

var nested = struct {
	A struct {
		B int
	}
}{
	struct {
		B int
	}{10},
}

func (o *Obj) Method(x int) string {
	return fmt.Sprintf("%d", o.Field+x)
}

func TestDescribeType(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "bool",
			args: args{
				v: true,
			},
			want: "bool",
		},
		{
			name: "int",
			args: args{
				v: 1,
			},
			want: "int",
		},
		{
			name: "int8",
			args: args{
				v: int8(1),
			},
			want: "int8",
		},
		{
			name: "int16",
			args: args{
				v: int16(1),
			},
			want: "int16",
		},
		{
			name: "int32",
			args: args{
				v: int32(1),
			},
			want: "int32",
		},
		{
			name: "int64",
			args: args{
				v: int64(1),
			},
			want: "int64",
		},
		{
			name: "uint",
			args: args{
				v: uint(1),
			},
			want: "uint",
		},
		{
			name: "uint8",
			args: args{
				v: uint8(1),
			},
			want: "uint8",
		},
		{
			name: "uint16",
			args: args{
				v: uint16(1),
			},
			want: "uint16",
		},
		{
			name: "uint32",
			args: args{
				v: uint32(1),
			},
			want: "uint32",
		},
		{
			name: "uint64",
			args: args{
				v: uint64(1),
			},
			want: "uint64",
		},
		{
			name: "uintptr",
			args: args{
				v: uintptr(1),
			},
			want: "uintptr",
		},
		{
			name: "float32",
			args: args{
				v: float32(1.2),
			},
			want: "float32",
		},
		{
			name: "float64",
			args: args{
				v: float64(1.2),
			},
			want: "float64",
		},
		{
			name: "complex64",
			args: args{
				v: complex64(1.2),
			},
			want: "complex64",
		},
		{
			name: "complex128",
			args: args{
				v: complex128(1.2),
			},
			want: "complex128",
		},
		{
			name: "string",
			args: args{
				v: "abc",
			},
			want: "string",
		},
		{
			name: "pointer to int",
			args: args{
				v: new(int),
			},
			want: "*int",
		},
		{
			name: "slice of int",
			args: args{
				v: []int{1, 2},
			},
			want: "[]int",
		},
		{
			name: "array of int",
			args: args{
				v: [2]int{1, 2},
			},
			want: "[2]int",
		},
		{
			name: "chan of int",
			args: args{
				v: make(chan int),
			},
			want: "chan int",
		},
		{
			name: "send chan of int",
			args: args{
				v: func() chan<- int {
					return make(chan int)
				}(),
			},
			want: "chan<- int",
		},
		{
			name: "recv chan of int",
			args: args{
				v: func() <-chan int {
					return make(chan int)
				}(),
			},
			want: "<-chan int",
		},
		{
			name: "map",
			args: args{
				v: make(map[int]int),
			},
			want: "map[int]int",
		},
		{
			name: "func",
			args: args{
				v: func() {},
			},
			want: "func ()",
		},
		{
			name: "func with return",
			args: args{
				v: func() int { return 2 },
			},
			want: "func () int",
		},
		{
			name: "func with multiple return",
			args: args{
				v: func() (int, int) { return 1, 2 },
			},
			want: "func () (int, int)",
		},
		{
			name: "func with one parameter",
			args: args{
				v: func(int) {},
			},
			want: "func (int)",
		},
		{
			name: "func with two parameters",
			args: args{
				v: func(int, int) {},
			},
			want: "func (int, int)",
		},
		{
			name: "named int",
			args: args{
				v: Foo(1),
			},
			want: "int",
		},
		{
			name: "func with builtin interface return",
			args: args{
				v: func() error { return nil },
			},
			want: "func () error",
		},
		{
			name: "func with package interface return",
			args: args{
				v: func() reflect.Type { return nil },
			},
			want: "func () reflect.Type",
		},
		{
			name: "func with local interface return",
			args: args{
				v: func() Iface { return &Obj{} },
			},
			want: "func () Iface",
		},
		{
			name: "struct",
			args: args{
				v: struct{ a int }{a: 10},
			},
			want: "struct {\n\ta int\n}",
		},
		{
			name: "struct with tag",
			args: args{
				v: struct {
					a int `tag:""`
				}{a: 10},
			},
			want: "struct {\n\ta int `tag:\"\"`\n}",
		},
		{
			name: "struct with anonymous field",
			args: args{
				v: struct{ Obj }{},
			},
			want: "struct {\n\tObj\n}",
		},
		{
			name: "nested struct",
			args: args{
				v: struct{ A struct{ B int } }{},
			},
			want: "struct {\n\tA struct {\n\t\tB int\n\t}\n}",
		},
		{
			name: "func with interface param",
			args: args{
				v: func(a interface {
					Method()
				}) {
				},
			},
			want: "func (interface {\n\t\tMethod()\n\t})",
		},
		{
			name: "unsafe pointer",
			args: args{
				v: unsafe.Pointer(uintptr(0)),
			},
			want: "unsafe.Pointer",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DescribeType(tt.args.v)
			if got != tt.want {
				t.Errorf("DescribeType() = %v, want %v", got, tt.want)
			}
			//                        fmt.Printf("-> %s\n", got)
		})
	}
}

func Test_describeType(t *testing.T) {
	type args struct {
		t     reflect.Type
		level int
		name  bool
	}
	tests := []struct {
		name  string
		args  args
		wantF string
	}{
		{
			name: "interface",
			args: args{
				t:     reflect.TypeOf(Iface(nil)),
				level: 0,
				name:  false,
			},
			wantF: "nil",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &bytes.Buffer{}
			describeType(f, tt.args.t, tt.args.level, tt.args.name)
			if gotF := f.String(); gotF != tt.wantF {
				t.Errorf("describeType() = %v, want %v", gotF, tt.wantF)
			}
		})
	}
}

func TestDescribeValue(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "bool: true",
			args: args{
				v: true,
			},
			want: "true",
		},
		{
			name: "bool: false",
			args: args{
				v: false,
			},
			want: "false",
		},
		{
			name: "int",
			args: args{
				v: 1,
			},
			want: "1",
		},
		{
			name: "int8",
			args: args{
				v: int8(1),
			},
			want: "int8(1)",
		},
		{
			name: "int16",
			args: args{
				v: int16(1),
			},
			want: "int16(1)",
		},
		{
			name: "int32",
			args: args{
				v: int32(1),
			},
			want: "int32(1)",
		},
		{
			name: "int64",
			args: args{
				v: int64(1),
			},
			want: "int64(1)",
		},
		{
			name: "uint",
			args: args{
				v: uint(1),
			},
			want: "uint(1)",
		},
		{
			name: "uint8",
			args: args{
				v: uint8(1),
			},
			want: "uint8(1)",
		},
		{
			name: "uint16",
			args: args{
				v: uint16(1),
			},
			want: "uint16(1)",
		},
		{
			name: "uint32",
			args: args{
				v: uint32(1),
			},
			want: "uint32(1)",
		},
		{
			name: "uint64",
			args: args{
				v: uint64(1),
			},
			want: "uint64(1)",
		},
		{
			name: "uintptr",
			args: args{
				v: uintptr(1),
			},
			want: "uintptr(1)",
		},
		{
			name: "float32",
			args: args{
				v: float32(1.2),
			},
			want: "float32(1.2)",
		},
		{
			name: "float64",
			args: args{
				v: float64(1.2),
			},
			want: "float64(1.2)",
		},
		{
			name: "complex64",
			args: args{
				v: complex64(1.2),
			},
			want: "complex64(1.2)",
		},
		{
			name: "complex64 with j",
			args: args{
				v: complex64(1.2 + 1.3i),
			},
			want: "complex64(1.2+1.3i)",
		},
		{
			name: "complex128",
			args: args{
				v: complex128(1.2),
			},
			want: "complex128(1.2)",
		},
		{
			name: "complex128 with j",
			args: args{
				v: complex128(1.2 + 1.3i),
			},
			want: "complex128(1.2+1.3i)",
		},
		{
			name: "string",
			args: args{
				v: "abc",
			},
			want: "\"abc\"",
		},
		{
			name: "pointer to int",
			args: args{
				v: new(int),
			},
			want: "&0",
		},
		{
			name: "zero slice of int",
			args: args{
				v: []int{},
			},
			want: "[]int{}",
		},
		{
			name: "slice of int",
			args: args{
				v: []int{1, 2},
			},
			want: "[]int{\n\t1,\n\t2,\n}",
		},
		{
			name: "zero array of int",
			args: args{
				v: [0]int{},
			},
			want: "[0]int{}",
		},
		{
			name: "array of int",
			args: args{
				v: [2]int{1, 2},
			},
			want: "[2]int{\n\t1,\n\t2,\n}",
		},
		{
			name: "chan of int",
			args: args{
				v: make(chan int),
			},
			want: "make(chan int)",
		},
		{
			name: "buffered chan of int",
			args: args{
				v: make(chan int, 5),
			},
			want: "make(chan int, 5)",
		},
		{
			name: "send chan of int",
			args: args{
				v: func() chan<- int {
					return make(chan int)
				}(),
			},
			want: "make(chan<- int)",
		},
		{
			name: "recv chan of int",
			args: args{
				v: func() <-chan int {
					return make(chan int)
				}(),
			},
			want: "make(<-chan int)",
		},
		{
			name: "map",
			args: args{
				v: make(map[int]int),
			},
			want: "map[int]int{}",
		},
		{
			name: "map with entries",
			args: args{
				v: map[int]int{
					1: 2,
					3: 4,
				},
			},
			want: `map[int]int{
	1: 2,
	3: 4,
}`,
		},
		{
			name: "func",
			args: args{
				v: func() {},
			},
			want: "func () {func0}",
		},
		{
			name: "func with return",
			args: args{
				v: func() int { return 2 },
			},
			want: "func () int {func1}",
		},
		{
			name: "func with multiple return",
			args: args{
				v: func() (int, int) { return 1, 2 },
			},
			want: "func () (int, int) {func2}",
		},
		{
			name: "func with one parameter",
			args: args{
				v: func(int) {},
			},
			want: "func (int) {func3}",
		},
		{
			name: "func with two parameters",
			args: args{
				v: func(int, int) {},
			},
			want: "func (int, int) {func4}",
		},
		{
			name: "named int",
			args: args{
				v: Foo(1),
			},
			want: "Foo(1)",
		},
		{
			name: "func with builtin interface return",
			args: args{
				v: func() error { return nil },
			},
			want: "func () error {func5}",
		},
		{
			name: "func with package interface return",
			args: args{
				v: func() reflect.Type { return nil },
			},
			want: "func () reflect.Type {func6}",
		},
		{
			name: "func with local interface return",
			args: args{
				v: func() Iface { return &Obj{} },
			},
			want: "func () Iface {func7}",
		},
		{
			name: "struct",
			args: args{
				v: struct{ A int }{A: 10},
			},
			want: "struct {\n\tA int\n}{\n\tA: 10,\n}",
		},
		{
			name: "struct with private field",
			args: args{
				v: struct{ a int }{a: 10},
			},
			want: "struct {\n\ta int\n}{\n\ta: ...,\n}",
		},
		{
			name: "struct with tag",
			args: args{
				v: struct {
					A int `tag:""`
				}{A: 10},
			},
			want: "struct {\n\tA int `tag:\"\"`\n}{\n\tA: 10,\n}",
		},
		{
			name: "struct with anonymous field",
			args: args{
				v: struct{ Obj }{},
			},
			want: "struct {\n\tObj\n}{\n\tObj{\n\t\tField: 0,\n\t},\n}",
		},
		{
			name: "nested struct",
			args: args{
				v: struct{ A struct{ B int } }{},
			},
			want: "struct {\n\tA struct {\n\t\tB int\n\t}\n}{\n\tA: struct {\n\t\tB int\n\t}{\n\t\tB: 0,\n\t},\n}",
		},
		{
			name: "func with interface param",
			args: args{
				v: func(a interface {
					Method()
				}) {
				},
			},
			want: "func (interface {\n\t\tMethod()\n\t}) {func8}",
		},
		{
			name: "unsafe pointer",
			args: args{
				v: unsafe.Pointer(uintptr(0)),
			},
			want: "unsafe.Pointer(0)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DescribeValue(tt.args.v); got != tt.want {
				t.Errorf("DescribeValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_describeValue(t *testing.T) {
	type args struct {
		t     reflect.Type
		v     reflect.Value
		level int
	}
	tests := []struct {
		name  string
		args  args
		wantF string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &bytes.Buffer{}
			describeValue(f, tt.args.t, tt.args.v, tt.args.level)
			if gotF := f.String(); gotF != tt.wantF {
				t.Errorf("describeValue() = %v, want %v", gotF, tt.wantF)
			}
		})
	}
}
