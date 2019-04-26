package describe

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"sort"
	"sync"
)

type packageType int

// Type returns a string that could be used to define a type
func Type(v interface{}) string {
	var buf bytes.Buffer
	describeType(&buf, reflect.TypeOf(v), 0, false)
	return buf.String()
}

func describeFuncParams(f io.Writer, t reflect.Type, level int) {
	fmt.Fprintf(f, "(")

	for i := 0; i < t.NumIn(); i++ {
		if i > 0 {
			fmt.Fprintf(f, ", ")
		}
		describeType(f, t.In(i), level+1, true)
	}

	fmt.Fprintf(f, ")")

	if t.NumOut() > 0 {
		fmt.Fprintf(f, " ")

		if t.NumOut() > 1 {
			fmt.Fprintf(f, "(")
		}

		for i := 0; i < t.NumOut(); i++ {
			if i > 0 {
				fmt.Fprintf(f, ", ")
			}
			describeType(f, t.Out(i), level+1, true)
		}

		if t.NumOut() > 1 {
			fmt.Fprintf(f, ")")
		}
	}
}

func typeName(t reflect.Type) string {
	name := t.Name()
	if name == "" {
		return ""
	}
	path := t.PkgPath()
	if path == "" || path == reflect.TypeOf(packageType(0)).PkgPath() {
		if name == "bool" || name == "int" || name == "string" {
			return ""
		}
		return name
	}
	return fmt.Sprintf("%s.%s", path, name)
}

func describeType(f io.Writer, t reflect.Type, level int, name bool) {
	if t == nil {
		fmt.Fprintf(f, "nil")
		return
	}

	k := t.Kind()
	//        fmt.Printf("kind %s name %s\n", k.String(), t.Name())

	if name {
		tn := typeName(t)
		if tn != "" {
			fmt.Fprintf(f, "%s", tn)
			return
		}
	}

	switch k {
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8,
		reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64,
		reflect.Complex64, reflect.Complex128, reflect.String:
		fmt.Fprintf(f, "%s", k.String())
	case reflect.Array:
		fmt.Fprintf(f, "[%d]", t.Len())
		describeType(f, t.Elem(), level+1, true)
	case reflect.Chan:
		fmt.Fprintf(f, "%s ", t.ChanDir().String())
		describeType(f, t.Elem(), level+1, true)
	case reflect.Func:
		fmt.Fprintf(f, "func ")
		describeFuncParams(f, t, level)
	case reflect.Interface:
		fmt.Fprintf(f, "interface")
		if t.NumMethod() == 0 {
			fmt.Fprintf(f, "{}")
		} else {
			fmt.Fprintf(f, " {\n")

			for i := 0; i < t.NumMethod(); i++ {
				m := t.Method(i)

				if m.Type.Kind() == reflect.Func {
					fmt.Fprintf(f, "%s%s", indent(level+1), m.Name)
					describeFuncParams(f, m.Type, level+1)

				} else {
					fmt.Fprintf(f, "%s%s ", indent(level+1), m.Name)
					describeType(f, m.Type, level+1, true)
				}

				fmt.Fprintf(f, "\n")
			}

			fmt.Fprintf(f, "%s}", indent(level))
		}
	case reflect.Map:
		fmt.Fprintf(f, "map[")
		describeType(f, t.Key(), level+1, true)
		fmt.Fprintf(f, "]")
		describeType(f, t.Elem(), level+1, true)
	case reflect.Ptr:
		fmt.Fprintf(f, "*")
		describeType(f, t.Elem(), level+1, true)
	case reflect.Slice:
		fmt.Fprintf(f, "[]")
		describeType(f, t.Elem(), level+1, true)
	case reflect.Struct:
		fmt.Fprintf(f, "struct")
		if t.NumField() == 0 {
			fmt.Fprintf(f, "{}")
		} else {
			fmt.Fprintf(f, " {\n")
			for i := 0; i < t.NumField(); i++ {
				sf := t.Field(i)

				if sf.Anonymous {
					fmt.Fprintf(f, "%s", indent(level+1))
				} else {
					fmt.Fprintf(f, "%s%s ", indent(level+1), sf.Name)
				}

				describeType(f, sf.Type, level+1, true)

				if sf.Tag != "" {
					fmt.Fprintf(f, " `%s`", sf.Tag)
				}

				fmt.Fprintf(f, "\n")
			}

			fmt.Fprintf(f, "%s}", indent(level))
		}
	case reflect.UnsafePointer:
		fmt.Fprintf(f, "unsafe.Pointer")
	default:
		fmt.Fprintf(f, "type of unknown kind %s", k.String())
	}
}

// Value returns a string that could be used to declare an initial value
func Value(v interface{}) string {
	var buf bytes.Buffer
	describeValue(&buf, reflect.TypeOf(v), reflect.ValueOf(v), 0)
	return buf.String()
}

func basicValue(t reflect.Type, v reflect.Value) string {
	if t == nil {
		return "nil"
	}

	k := t.Kind()
	//        fmt.Printf("kind %s name %s\n", k.String(), t.Name())
	i := v.Interface()

	switch k {
	case reflect.Bool:
		if i.(bool) {
			return "true"
		}
		return "false"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%d", i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return fmt.Sprintf("%d", i)
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%g", i)
	case reflect.Complex64:
		r := real(i.(complex64))
		j := imag(i.(complex64))
		if j != 0.0 {
			return fmt.Sprintf("%g+%gi", r, j)
		}
		return fmt.Sprintf("%g", r)
	case reflect.Complex128:
		r := real(i.(complex128))
		j := imag(i.(complex128))
		if j != 0.0 {
			return fmt.Sprintf("%g+%gi", r, j)
		}
		return fmt.Sprintf("%g", r)
	case reflect.String:
		// Should probably do some decoding of the string to make special characters visible, but this
		// is good enough for now.
		return fmt.Sprintf("\"%s\"", v)
	}

	return ""
}

func describeValue(f io.Writer, t reflect.Type, v reflect.Value, level int) {
	if t == nil {
		fmt.Fprintf(f, "nil")
		return
	}

	k := t.Kind()
	//        fmt.Printf("kind %s name %s\n", k.String(), t.Name())
	tn := typeName(t)

	switch k {
	case reflect.Bool, reflect.Int, reflect.String:
		bv := basicValue(t, v)
		if tn != "" {
			fmt.Fprintf(f, "%s(%s)", tn, bv)
		} else {
			fmt.Fprintf(f, "%s", bv)
		}
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint,
		reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32,
		reflect.Float64, reflect.Complex64, reflect.Complex128:
		bv := basicValue(t, v)
		if tn == "" {
			tn = k.String()
		}
		fmt.Fprintf(f, "%s(%s)", tn, bv)
	case reflect.Array:
		describeType(f, t, level, true)
		if v.Len() == 0 {
			fmt.Fprintf(f, "{}")
		} else {
			fmt.Fprintf(f, "{\n")

			for j := 0; j < v.Len(); j++ {
				fmt.Fprintf(f, "%s", indent(level+1))
				describeValue(f, t.Elem(), v.Index(j), level+1)
				fmt.Fprintf(f, ",\n")
			}

			fmt.Fprintf(f, "%s", indent(level))
			fmt.Fprintf(f, "}")
		}
	case reflect.Chan:
		fmt.Fprintf(f, "make(")
		describeType(f, t, level, true)
		c := v.Cap()
		if c > 0 {
			fmt.Fprintf(f, ", %d)", c)
		} else {
			fmt.Fprintf(f, ")")
		}
	case reflect.Func:
		fmt.Fprintf(f, "func ")
		describeFuncParams(f, t, level)
		fmt.Fprintf(f, " {func%d}", objectNumber(v))
	case reflect.Interface:
		describeType(f, t, level, true)
		fmt.Fprintf(f, "{}")
	case reflect.Map:
		describeType(f, t, level, true)
		if v.Len() == 0 {
			fmt.Fprintf(f, "{}")
		} else {
			fmt.Fprintf(f, "{\n")

			kt := t.Key()
			keys := v.MapKeys()
			sort.Slice(keys, func(i, j int) bool { return less(kt, keys[i], keys[j]) })
			for _, k := range keys {
				fmt.Fprintf(f, "%s", indent(level+1))
				describeValue(f, t.Key(), k, level+1)
				fmt.Fprintf(f, ": ")
				describeValue(f, t.Elem(), v.MapIndex(k), level+1)
				fmt.Fprintf(f, ",\n")
			}

			fmt.Fprintf(f, "%s}", indent(level))
		}
	case reflect.Ptr:
		fmt.Fprintf(f, "&")
		describeValue(f, t.Elem(), v.Elem(), level+1)
	case reflect.Slice:
		describeType(f, t, level, true)
		if v.Len() == 0 {
			fmt.Fprintf(f, "{}")
		} else {
			fmt.Fprintf(f, "{\n")

			for j := 0; j < v.Len(); j++ {
				fmt.Fprintf(f, "%s", indent(level+1))
				describeValue(f, t.Elem(), v.Index(j), level+1)
				fmt.Fprintf(f, ",\n")
			}

			fmt.Fprintf(f, "%s}", indent(level))
		}
	case reflect.Struct:
		describeType(f, t, level, true)
		fmt.Fprintf(f, "{\n")

		for i := 0; i < t.NumField(); i++ {
			sf := t.Field(i)
			fv := v.Field(i)

			fmt.Fprintf(f, "%s", indent(level+1))
			if !sf.Anonymous {
				if sf.PkgPath != "" && sf.PkgPath != reflect.TypeOf(packageType(0)).PkgPath() {
					fmt.Fprintf(f, "%s.%s: ", sf.PkgPath, sf.Name)
				} else {
					fmt.Fprintf(f, "%s: ", sf.Name)
				}
			}
			if sf.Name == "" || ('A' <= sf.Name[0] && sf.Name[0] <= 'Z') {
				describeValue(f, sf.Type, fv, level+1)
			} else {
				fmt.Fprintf(f, "...")
			}
			fmt.Fprintf(f, ",\n")
		}

		fmt.Fprintf(f, "%s}", indent(level))
	case reflect.UnsafePointer:
		fmt.Fprintf(f, "unsafe.Pointer(%x)", v.Pointer())
	default:
		fmt.Fprintf(f, "type of unknown kind %s", k.String())
	}
}

func indent(level int) string {
	tabs := "\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t"
	t := tabs
	for len(t) < level {
		t = t + t
	}
	return t[0:level]
}

var objNums struct {
	sync.RWMutex
	Table map[uintptr]int
}

func objectNumber(v reflect.Value) int {
	ptr := v.Pointer()
	objNums.RLock()
	if n, ok := objNums.Table[ptr]; ok {
		objNums.RUnlock()
		return n
	}
	objNums.RUnlock()
	objNums.Lock()
	defer objNums.Unlock()
	if n, ok := objNums.Table[ptr]; ok {
		return n
	}
	if objNums.Table == nil {
		objNums.Table = make(map[uintptr]int)
	}
	n := len(objNums.Table)
	objNums.Table[ptr] = n
	return n
}

func less(t reflect.Type, a, b reflect.Value) bool {
	k := t.Kind()
	if k == reflect.String {
		return a.String() < b.String()
	}
	if k == reflect.Int || k == reflect.Int8 || k == reflect.Int16 || k == reflect.Int32 || k == reflect.Int64 {
		return a.Int() < b.Int()
	}
	if k == reflect.Uint || k == reflect.Uint8 || k == reflect.Uint16 || k == reflect.Uint32 || k == reflect.Uint64 {
		return a.Uint() < b.Uint()
	}
	if k == reflect.Float32 || k == reflect.Float64 {
		return a.Float() < b.Float()
	}
	return false
}
