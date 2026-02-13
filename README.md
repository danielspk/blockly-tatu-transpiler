# Blockly to Tatu Transpiler

Transpiler that converts [Blockly](https://blockly.com/) visual blocks _(JSON format)_ into [Tatu](https://github.com/danielspk/tatu-lang) source code.

## Quick Start

```go
package main

import (
    "fmt"

    "github.com/danielspk/blockly-tatu-transpiler/pkg/blockly"
)

func main() {
    jsonData := `{
        "blocks": {
            "languageVersion": 0,
            "blocks": [{
                "type": "text_print",
                "id": "block1",
                "inputs": {
                    "TEXT": {"block": {"type": "text", "id": "block2", "fields": {"TEXT": "Hello, Tatu!"}}}
                }
            }]
        }
    }`

    tatuCode, err := blockly.Transpile([]byte(jsonData))
    if err != nil {
        panic(err)
    }

    fmt.Println(tatuCode)
}

// Output: (print "Hello, Tatu!")
```

## Web Interface

Run the visual Blockly editor:

```bash
PORT=8080 go run cmd/main.go
```

Then open `http://localhost:8080` in your browser.

## Supported Blocks

- **Variables** - Get/set variables
- **Logic** - Boolean, null, comparisons, and/or/not operations, ternary
- **Control** - If/elseif/else, repeat, while/until, for, forEach, break/continue
- **Math** - Numbers, arithmetic, trigonometry, rounding, random, constants
- **Text** - Strings, print, join, length, case conversion, substring operations
- **Lists** - Create, get/set items, length, sort, search operations
- **Procedures** - Function definitions and calls (with/without return values)
- **Tatu** - Custom blocks for:
  - create variable
  - including external files
  - map creation
  - expression statement
  - comment
  - function call
  - standard library functions call: `math:*`, `str:*`, `vec:*`, `map:*`, `time:*`, `fs:*`, `json:*`, `regex:*`, `type checking` and `casting`

## Transpilation

### Comments
```
comment: 'This is a comment
→ ; This is a comment
```

### Variables
```
variables_create_tatu: var x = 10
→ (var x 10)

variables_set: x = 10
→ (set x 10)

variables_get: x
→ x
```

### Logic
```
logic_boolean: true
→ true

logic_compare: 5 = 5
→ (= 5 5)

logic_operation: true AND false
→ (and true false)

logic_negate: NOT true
→ (not true)

logic_ternary: if true then "yes" else "no"
→ (if true "yes" "no")

logic_null: null
→ nil
```

### Control Flow
```
controls_if: if x > 5 then print("big")
→ (if (> x 5)
     (print "big"))

controls_repeat_ext: repeat 10 times
→ (begin
     (var __i 0)
     (while (< __i 10)
       (begin
         [body]
         (set __i (+ __i 1)))))

controls_whileUntil: while x < 10
→ (while (< x 10)
     [body])

controls_whileUntil: until x >= 10
→ (while (not (>= x 10))
     [body])

controls_for: for i from 0 to 10 by 1
→ (begin
     (var i 0)
     (while (<= i 10)
       (begin
         [body]
         (set i (+ i 1)))))

controls_forEach: for item in list
→ (begin
     (var item nil)
     (for (var __i_foreach 0)
          (< __i_foreach (vec:len list))
          (begin
            (set item (vec:get list __i_foreach))
            [body]
            (set __i_foreach (+ __i_foreach 1)))))

controls_flow_statements: break
→ (set __break true)

expression_tatu: expr (get x)
→ x
```

### Math
```
math_number: 42
→ 42

math_arithmetic: 5 + 3
→ (+ 5 3)

math_single: square root of 16
→ (math:sqrt 16)

math_trig: sin(90)
→ (math:sin 90)

math_constant: π
→ (math:pi)

math_round: round 3.7
→ (math:round 3.7)

math_modulo: 10 mod 3
→ (% 10 3)

math_random_int: random from 1 to 10
→ (math:rand 1 10)

math_constrain: constrain 15 low 0 high 10
→ (math:min (math:max 15 0) 10)

math_number_property: 10 is even
→ (= (% 10 2) 0)
```

### Text
```
text: "hello"
→ "hello"

text_print: print "hello"
→ (print "hello")

text_join: join "hello" " " "world"
→ (str:concat (to-string "hello") (to-string " ") (to-string "world"))

text_length: length of "hello"
→ (str:len "hello")

text_isEmpty: "hello" is empty
→ (= (str:len "hello") 0)

text_charAt: letter 0 of "hello"
→ (str:slice "hello" 0 1)

text_getSubstring: substring of "hello" from 1 to 3
→ (str:slice "hello" 1 3)

text_append: append "world" to x
→ (set x (str:concat x "world"))

text_indexOf: find "world" in "hello world"
→ (str:index "hello world" "world")

text_changeCase: uppercase "hello"
→ (str:upper "hello")

text_trim: trim both sides of "  hello  "
→ (str:trim "  hello  ")
```

### Lists (Vectors)
```
lists_create_with: [1, 2, 3]
→ (vector 1 2 3)

lists_length: length of list
→ (vec:len list)

lists_isEmpty: list is empty
→ (= (vec:len list) 0)

lists_getIndex: get item 0 from list
→ (vec:get list 0)

lists_setIndex: set item 0 of list to 42
→ (vec:set list 0 42)

lists_indexOf: find 42 in list
→ (vec:find list 42)

lists_getSublist: sublist of list from 1 to 3
→ (vec:slice list 1 4)

lists_split: split "a,b,c" by ","
→ (str:split "a,b,c" ",")

lists_sort: sort list ascending
→ (vec:sort list)

lists_reverse: reverse list
→ (vec:reverse list)
```

### Tatu Maps
```
map_create_with_tatu: create map with key-value pairs
→ (map "name" "John" "age" 30 "active" true)
```

### Procedures (Functions)
```
procedures_defnoreturn: define greet(name)
→ (def greet (name)
     [body])

procedures_defreturn: define double(x) return x * 2
→ (def double (x)
     (* x 2))

procedures_callnoreturn: greet("Alice")
→ (greet "Alice")

procedures_callreturn: double(5)
→ (double 5)

procedures_ifreturn: if x > 5 return x
→ (if (> x 5) x)

function_call_tatu: add(5, 5)
→ (add 5 5)
```

### Tatu Include
```
include_file_tatu: include "lib/utils.tatu"
→ (include "lib/utils.tatu")
```

### Tatu Standard Library
```
math_call_tatu: math:sqrt(16)
→ (math:sqrt 16)

string_call_tatu: str:upper("hello")
→ (str:upper "hello")

vector_call_tatu: vec:push([1,2], 3)
→ (vec:push (vector 1 2) 3)

map_call_tatu: map:get(mymap, "key")
→ (map:get mymap "key")

json_call_tatu: json:encode(data)
→ (json:encode data)

type_check_call_tatu: is-number(42)
→ (is-number 42)

casting_call_tatu: to-string(42)
→ (to-string 42)

time_call_tatu: time:now()
→ (time:now)

filesystem_call_tatu: fs:read("file.txt")
→ (fs:read "file.txt")

regex_call_tatu: regex:matches("\\d+", "123")
→ (regex:matches "\\d+" "123")
```

## License

This software is distributed under the _MIT_ license. See the [LICENSE](LICENSE) file.
