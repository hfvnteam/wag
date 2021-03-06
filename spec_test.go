package wag

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/tsavola/wag/diswag"
	"github.com/tsavola/wag/internal/sexp"
	"github.com/tsavola/wag/runner"
	"github.com/tsavola/wag/traps"
)

const (
	specTestDir = "testdata/wabt/third_party/testsuite"
)

// for i in $(ls -1 *.wast); do echo 'func Test_'$(echo $i | sed 's/.wast$//' | tr - _ | tr . _)'(t *testing.T) { test(t, "'$(echo $i | sed 's/.wast$//')'") }'; done

// func Test_binary(t *testing.T)                { test(t, "binary") }
// func Test_call_indirect(t *testing.T)         { test(t, "call_indirect") }
// func Test_comments(t *testing.T)              { test(t, "comments") }
// func Test_conversions(t *testing.T)           { test(t, "conversions") }
// func Test_exports(t *testing.T)               { test(t, "exports") }
// func Test_f32(t *testing.T)                   { test(t, "f32") }
// func Test_f64(t *testing.T)                   { test(t, "f64") }
// func Test_fac(t *testing.T)                   { test(t, "fac") }
// func Test_float_exprs(t *testing.T)           { test(t, "float_exprs") }
// func Test_float_memory(t *testing.T)          { test(t, "float_memory") }
// func Test_float_misc(t *testing.T)            { test(t, "float_misc") }
// func Test_func(t *testing.T)                  { test(t, "func") }
// func Test_func_ptrs(t *testing.T)             { test(t, "func_ptrs") }
// func Test_imports(t *testing.T)               { test(t, "imports") }
// func Test_left_to_right(t *testing.T)         { test(t, "left-to-right") }
// func Test_linking(t *testing.T)               { test(t, "linking") }
// func Test_memory(t *testing.T)                { test(t, "memory") }
// func Test_names(t *testing.T)                 { test(t, "names") }
// func Test_set_local(t *testing.T)             { test(t, "set_local") }
// func Test_skip_stack_guard_page(t *testing.T) { test(t, "skip-stack-guard-page") }
// func Test_start(t *testing.T)                 { test(t, "start") }
// func Test_tee_local(t *testing.T)             { test(t, "tee_local") }
// func Test_traps(t *testing.T)                 { test(t, "traps") }

func Test_address(t *testing.T)                         { test(t, "address") }
func Test_address_offset_range_fail(t *testing.T)       { test(t, "address-offset-range.fail") }
func Test_block(t *testing.T)                           { test(t, "block") }
func Test_br(t *testing.T)                              { test(t, "br") }
func Test_br_if(t *testing.T)                           { test(t, "br_if") }
func Test_br_table(t *testing.T)                        { test(t, "br_table") }
func Test_break_drop(t *testing.T)                      { test(t, "break-drop") }
func Test_call(t *testing.T)                            { test(t, "call") }
func Test_endianness(t *testing.T)                      { test(t, "endianness") }
func Test_f32_cmp(t *testing.T)                         { test(t, "f32_cmp") }
func Test_f32_load32_fail(t *testing.T)                 { test(t, "f32.load32.fail") }
func Test_f32_load64_fail(t *testing.T)                 { test(t, "f32.load64.fail") }
func Test_f32_store32_fail(t *testing.T)                { test(t, "f32.store32.fail") }
func Test_f32_store64_fail(t *testing.T)                { test(t, "f32.store64.fail") }
func Test_f64_cmp(t *testing.T)                         { test(t, "f64_cmp") }
func Test_f64_load32_fail(t *testing.T)                 { test(t, "f64.load32.fail") }
func Test_f64_load64_fail(t *testing.T)                 { test(t, "f64.load64.fail") }
func Test_f64_store32_fail(t *testing.T)                { test(t, "f64.store32.fail") }
func Test_f64_store64_fail(t *testing.T)                { test(t, "f64.store64.fail") }
func Test_float_literals(t *testing.T)                  { test(t, "float_literals") }
func Test_forward(t *testing.T)                         { test(t, "forward") }
func Test_func_local_after_body_fail(t *testing.T)      { test(t, "func-local-after-body.fail") }
func Test_func_local_before_param_fail(t *testing.T)    { test(t, "func-local-before-param.fail") }
func Test_func_local_before_result_fail(t *testing.T)   { test(t, "func-local-before-result.fail") }
func Test_func_param_after_body_fail(t *testing.T)      { test(t, "func-param-after-body.fail") }
func Test_func_result_after_body_fail(t *testing.T)     { test(t, "func-result-after-body.fail") }
func Test_func_result_before_param_fail(t *testing.T)   { test(t, "func-result-before-param.fail") }
func Test_get_local(t *testing.T)                       { test(t, "get_local") }
func Test_globals(t *testing.T)                         { test(t, "globals") }
func Test_i32(t *testing.T)                             { test(t, "i32") }
func Test_i32_load32_s_fail(t *testing.T)               { test(t, "i32.load32_s.fail") }
func Test_i32_load32_u_fail(t *testing.T)               { test(t, "i32.load32_u.fail") }
func Test_i32_load64_s_fail(t *testing.T)               { test(t, "i32.load64_s.fail") }
func Test_i32_load64_u_fail(t *testing.T)               { test(t, "i32.load64_u.fail") }
func Test_i32_store32_fail(t *testing.T)                { test(t, "i32.store32.fail") }
func Test_i32_store64_fail(t *testing.T)                { test(t, "i32.store64.fail") }
func Test_i64(t *testing.T)                             { test(t, "i64") }
func Test_i64_load64_s_fail(t *testing.T)               { test(t, "i64.load64_s.fail") }
func Test_i64_load64_u_fail(t *testing.T)               { test(t, "i64.load64_u.fail") }
func Test_i64_store64_fail(t *testing.T)                { test(t, "i64.store64.fail") }
func Test_import_after_func_fail(t *testing.T)          { test(t, "import-after-func.fail") }
func Test_import_after_global_fail(t *testing.T)        { test(t, "import-after-global.fail") }
func Test_import_after_memory_fail(t *testing.T)        { test(t, "import-after-memory.fail") }
func Test_import_after_table_fail(t *testing.T)         { test(t, "import-after-table.fail") }
func Test_int_exprs(t *testing.T)                       { test(t, "int_exprs") }
func Test_int_literals(t *testing.T)                    { test(t, "int_literals") }
func Test_labels(t *testing.T)                          { test(t, "labels") }
func Test_loop(t *testing.T)                            { test(t, "loop") }
func Test_memory_redundancy(t *testing.T)               { test(t, "memory_redundancy") }
func Test_memory_trap(t *testing.T)                     { test(t, "memory_trap") }
func Test_nop(t *testing.T)                             { test(t, "nop") }
func Test_of_string_overflow_hex_u32_fail(t *testing.T) { test(t, "of_string-overflow-hex-u32.fail") }
func Test_of_string_overflow_hex_u64_fail(t *testing.T) { test(t, "of_string-overflow-hex-u64.fail") }
func Test_of_string_overflow_s32_fail(t *testing.T)     { test(t, "of_string-overflow-s32.fail") }
func Test_of_string_overflow_s64_fail(t *testing.T)     { test(t, "of_string-overflow-s64.fail") }
func Test_of_string_overflow_u32_fail(t *testing.T)     { test(t, "of_string-overflow-u32.fail") }
func Test_of_string_overflow_u64_fail(t *testing.T)     { test(t, "of_string-overflow-u64.fail") }
func Test_resizing(t *testing.T)                        { test(t, "resizing") }
func Test_return(t *testing.T)                          { test(t, "return") }
func Test_select(t *testing.T)                          { test(t, "select") }
func Test_stack(t *testing.T)                           { test(t, "stack") }
func Test_store_retval(t *testing.T)                    { test(t, "store_retval") }
func Test_switch(t *testing.T)                          { test(t, "switch") }
func Test_typecheck(t *testing.T)                       { test(t, "typecheck") }
func Test_unreachable(t *testing.T)                     { test(t, "unreachable") }
func Test_unwind(t *testing.T)                          { test(t, "unwind") }

func test(t *testing.T, name string) {
	const (
		parallel = true
	)

	if parallel {
		t.Parallel()
	}

	filename := path.Join(specTestDir, name) + ".wast"

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}

	quiet := false

	if strings.HasSuffix(name, ".fail") {
		quiet = true

		defer func() {
			x := recover()
			if x == nil {
				t.Error()
			} else {
				t.Logf("expected panic: %s", x)
			}
		}()
	}

	for i := 0; len(data) > 0; i++ {
		data = testModule(t, data, fmt.Sprintf("%s-%d", name, i), quiet)
	}
}

func testModule(t *testing.T, data []byte, filename string, quiet bool) []byte {
	const (
		maxTextSize   = 0x100000
		maxRODataSize = 0x100000
		maxMemorySize = 0x100000
		stackSize     = 4096 // limit stacktrace length

		timeout     = time.Second * 3
		dumpExps    = false
		dumpText    = false
		dumpROData  = false
		dumpGlobals = false
		dumpMemory  = false
	)

	module, data := sexp.ParsePanic(data)
	if module == nil {
		return nil
	}

	if name := module[0].(string); name != "module" {
		t.Logf("%s not supported", name)
		return data
	}

	module = append([]interface{}{
		module[0],
		[]interface{}{
			"import",
			sexp.Quote("wag"),
			sexp.Quote("get_arg"),
			[]interface{}{
				"func",
				"$get_arg",
				[]interface{}{"result", "i64"},
			},
		},
		[]interface{}{
			"import",
			sexp.Quote("wag"),
			sexp.Quote("set_result"),
			[]interface{}{
				"func",
				"$set_result",
				[]interface{}{"param", "i32"},
			},
		},
	}, module[1:]...)

	var realStartName string
	exports := make(map[string]string)

	for i := 1; i < len(module); {
		item, ok := module[i].([]interface{})
		if !ok {
			i++
			continue
		}

		var itemName string
		switch x := item[0].(type) {
		case string:
			itemName = x

		case sexp.Quoted:
			itemName = x.String()
		}

		switch itemName {
		case "start":
			realStartName = item[1].(string)
			module = append(module[:i], module[i+1:]...)

		case "export":
			exports[item[1].(string)] = item[2].(string)
			module = append(module[:i], module[i+1:]...)

		case "func":
			if len(item) > 1 {
				if expo, ok := item[1].([]interface{}); ok && expo[0].(string) == "export" {
					item[1] = "$" + expo[1].(sexp.Quoted).String()
				}

				if s, ok := item[1].(string); ok && len(item) > 2 {
					if expo, ok := item[2].([]interface{}); ok && expo[0].(string) == "export" {
						exports[expo[1].(sexp.Quoted).String()] = s[1:]
					}
				}
			}
			i++

		default:
			i++
		}
	}

	testTypes := make(map[int]string)

	testFunc := []interface{}{
		"func",
		"$test",
		[]interface{}{
			"call",
			"$set_result",
			[]interface{}{"i32.const", "0xbadc0de"},
		},
	}

	if realStartName != "" {
		testFunc = append(testFunc, []interface{}{
			"if",
			[]interface{}{
				"i64.eq",
				[]interface{}{"call", "$get_arg"},
				[]interface{}{"i64.const", "0"},
			},
			[]interface{}{
				"block",
				[]interface{}{"call", realStartName},
				[]interface{}{
					"set_global",
					"$test_result",
					[]interface{}{"i64.const", "777"},
				},
				[]interface{}{"return"},
			},
		})
	}

	var idCount int

	for {
		id := idCount

		assert, tail := sexp.ParsePanic(data)
		if assert == nil {
			data = tail
			break
		}

		testType := assert[0].(string)
		if testType == "module" {
			break
		}

		idCount++
		data = tail

		var argCount int
		var exprType string

		for _, x := range assert[1:] {
			if expr, ok := x.([]interface{}); ok {
				argCount++

				exprName := expr[0].(string)
				if strings.Contains(exprName, ".") {
					exprType = strings.SplitN(exprName, ".", 2)[0]
					break
				}
			}
		}

		if argCount > 1 && exprType == "" {
			t.Fatalf("can't figure out type of %s", sexp.Stringify(assert, true))
		}

		invoke2call(exports, assert[1:])

		var test []interface{}

		switch testType {
		case "assert_return":
			if argCount > 1 {
				var check interface{}

				switch exprType {
				case "f32", "f64":
					bitsType := strings.Replace(exprType, "f", "i", 1)

					check = []interface{}{
						bitsType + ".eq",
						[]interface{}{
							bitsType + ".reinterpret/" + exprType,
							assert[1],
						},
						[]interface{}{
							bitsType + ".reinterpret/" + exprType,
							assert[2],
						},
					}

				default:
					check = []interface{}{exprType + ".eq", assert[1], assert[2]}
				}

				test = []interface{}{
					"block",
					[]interface{}{
						"call",
						"$set_result",
						check,
					},
					[]interface{}{"return"},
				}
			} else {
				test = append([]interface{}{"block"}, assert[1:]...)
				test = append(test, []interface{}{
					"block",
					[]interface{}{
						"call",
						"$set_result",
						[]interface{}{"i32.const", "1"},
					},
					[]interface{}{"return"},
				})
			}

		case "assert_trap":
			test = []interface{}{
				"block",
				assert[1],
				[]interface{}{"return"},
			}

		case "invoke":
			n := assert[1].(sexp.Quoted).String()
			name, found := exports[n]
			if !found {
				name = n
			}

			test = []interface{}{
				"block",
				append([]interface{}{"call", "$" + name}, assert[2:]...),
				[]interface{}{
					"call",
					"$set_result",
					[]interface{}{"i32.const", "-1"},
				},
				[]interface{}{"return"},
			}

		default:
			testType = ""
		}

		testTypes[id] = testType

		if test != nil {
			testFunc = append(testFunc, []interface{}{
				"if",
				[]interface{}{
					"i64.eq",
					[]interface{}{"call", "$get_arg"},
					[]interface{}{"i64.const", strconv.Itoa(id)},
				},
				test,
			})
		}
	}

	module = append(module, testFunc)
	module = append(module, []interface{}{"start", "$test"})

	if dumpExps {
		fmt.Println(sexp.Stringify(module, true))
	}

	{
		wasmReadCloser := wast2wasm(sexp.Unparse(module), quiet)
		defer wasmReadCloser.Close()
		wasm := bufio.NewReader(wasmReadCloser)

		var timedout bool

		p, err := runner.NewProgram(maxTextSize, maxRODataSize)
		if err != nil {
			t.Fatal(err)
		}
		defer func() {
			if !timedout {
				p.Close()
			}
		}()

		var m Module
		m.load(wasm, runner.Env, p.Text, p.ROData, p.RODataAddr(), nil)
		p.Seal()
		p.SetData(m.Data())
		p.SetFunctionMap(m.FunctionMap())
		p.SetCallMap(m.CallMap())
		minMemorySize, maxMemorySize := m.MemoryLimits()

		if dumpText && testing.Verbose() {
			diswag.PrintTo(os.Stdout, m.Text())
		}

		if dumpROData {
			buf := m.ROData()
			for i := 0; len(buf) > 0; i++ {
				if len(buf) > 4 {
					t.Logf("read-only data #%d*8: 0x%08x 0x%08x", i, binary.LittleEndian.Uint32(buf[:4]), binary.LittleEndian.Uint32(buf[4:8]))
					buf = buf[8:]
				} else {
					t.Logf("read-only data #%d*8: 0x%08x", i, binary.LittleEndian.Uint32(buf[:4]))
					buf = buf[4:]
				}
			}
		}

		if dumpGlobals {
			data, memoryOffset := m.Data()
			buf := data[:memoryOffset]

			if len(buf) == 0 {
				t.Log("no globals")
			}

			for i := 0; len(buf) > 0; i++ {
				t.Logf("global #%d: 0x%016x", i, binary.LittleEndian.Uint64(buf))
				buf = buf[8:]
			}
		}

		if dumpMemory {
			data, memoryOffset := m.Data()
			t.Logf("memory: %#v", data[memoryOffset:])
		}

		memGrowSize := maxMemorySize
		if maxMemorySize > 0 && memGrowSize > maxMemorySize {
			memGrowSize = maxMemorySize
		}

		r, err := p.NewRunner(minMemorySize, memGrowSize, stackSize)
		if err != nil {
			t.Fatal(err)
		}
		defer func() {
			if !timedout {
				r.Close()
			}
		}()

		if realStartName != "" {
			var printBuf bytes.Buffer
			result, err := r.Run(0, m.Signatures(), &printBuf)
			if printBuf.Len() > 0 {
				t.Logf("run: module %s: print:\n%s", filename, string(printBuf.Bytes()))
			}
			if err != nil {
				t.Fatal(err)
			}
			if result != 777 {
				t.Fatalf("0x%x", result)
			}
			t.Logf("run: module %s: start", filename)
		}

		for id := 0; id < idCount; id++ {
			testType := testTypes[id]
			if testType == "" {
				t.Logf("run: module %s: test #%d: not supported", filename, id)
				continue
			}

			var printBuf bytes.Buffer
			var result int32
			var panicked interface{}
			done := make(chan struct{})

			go func() {
				defer close(done)
				defer func() {
					panicked = recover()
				}()
				result, err = r.Run(int64(id), m.Signatures(), &printBuf)
			}()

			timer := time.NewTimer(timeout)

			select {
			case <-done:
				timer.Stop()

			case <-timer.C:
				timedout = true
				t.Fatalf("run: module %s: test #%d: timeout", filename, id)
			}

			if printBuf.Len() > 0 {
				t.Logf("run: module %s: test #%d: print output:\n%s", filename, id, string(printBuf.Bytes()))
			}

			if panicked != nil {
				t.Fatalf("run: module %s: test #%d: panic: %v", filename, id, panicked)
			}

			var stackBuf bytes.Buffer
			if err := r.WriteStacktraceTo(&stackBuf); err == nil {
				if stackBuf.Len() > 0 {
					t.Logf("run: module %s: test #%d: stacktrace:\n%s", filename, id, string(stackBuf.Bytes()))
				}
			} else {
				t.Errorf("run: module %s: test #%d: stacktrace error: %v", filename, id, err)
			}

			if err != nil {
				if _, ok := err.(traps.Id); ok {
					if testType == "assert_trap" {
						t.Logf("run: module %s: test #%d: pass", filename, id)
					} else {
						t.Errorf("run: module %s: test #%d: FAIL due to unexpected trap", filename, id)
					}
				} else {
					t.Fatal(err)
				}
			} else {
				if testType == "assert_return" {
					switch result {
					case 1:
						t.Logf("run: module %s: test #%d: pass", filename, id)

					case 0:
						t.Errorf("run: module %s: test #%d: FAIL", filename, id)

					default:
						t.Fatalf("run: module %s: test #%d: bad result: 0x%x", filename, id, result)
					}
				} else if testType == "invoke" {
					switch result {
					case -1:
						t.Logf("run: module %s: test #%d: invoke", filename, id)

					default:
						t.Fatalf("run: module %s: test #%d: bad result: 0x%x", filename, id, result)
					}
				} else {
					t.Errorf("run: module %s: test #%d: FAIL due to unexpected return (result: 0x%x)", filename, id, result)
				}
			}
		}
	}

	return data
}

func invoke2call(exports map[string]string, x interface{}) {
	if item, ok := x.([]interface{}); ok {
		if s, ok := item[0].(string); ok && s == "invoke" {
			item[0] = "call"

			s := item[1].(sexp.Quoted).String()
			if name, found := exports[s]; found {
				item[1] = "$" + name
			} else {
				item[1] = "$" + s
			}
		}

		for _, e := range item {
			invoke2call(exports, e)
		}
	}
}
