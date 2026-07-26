# Project: an interpreter for a small language

**Modules it exercises:** 02 (tokenising), 03 (interfaces, sum types), 04 (errors), 05 (generics).
**Rough size:** 800–1200 lines. **Difficulty:** ★★★

Lexer, parser, evaluator. The single best exercise for interface design in Go,
and the natural sequel to `03-types-interfaces/expr`.

## The language

```
let x = 5;
let add = fn(a, b) { a + b };
let result = add(x, 10);

if (result > 10) { puts("big") } else { puts("small") }

let people = ["ana", "bogdan"];
let ages = {"ana": 31, "bogdan": 24};

let counter = fn() {
  let count = 0;
  fn() { count = count + 1; count }   // a closure over count
};
```

Types: integers, booleans, strings, arrays, hashes, functions, null.
Builtins: `len`, `puts`, `first`, `rest`, `push`.

## Requirements

1. **Lexer** producing a token stream with positions. Handle multi-character
   operators (`==`, `!=`, `<=`), string escapes, comments, and EOF cleanly.
   Every token carries a line and column, because error messages without
   positions are useless.
2. **Pratt parser** (top-down operator precedence). This is the part worth
   learning properly: prefix and infix parse functions in a map keyed by token
   type, and a precedence table. It handles `-a * b + c(d)[e]` correctly with
   about 200 lines and no grammar hacks.
3. **AST as a sealed interface** — `Node`, `Statement`, `Expression` with an
   unexported method, exactly like `03/expr`. Every node has a `String()` so you
   can round-trip source → AST → source in a test.
4. **Tree-walking evaluator** with an `Environment` that has a parent pointer.
   Closures capture their defining environment — that's how the counter above
   works, and getting it right is the moment the whole thing clicks.
5. **Errors, not panics.** A runtime error (`5 + "x"`, unknown identifier,
   wrong argument count) is a value that propagates and stops evaluation, with
   the position. Reuse the taxonomy from `04-errors/apierr`.
6. **A REPL** with multi-line input and readable error output.

## Milestones

1. Lexer + a test that tokenises the sample above. Table-driven, one row per
   token.
2. Parse and evaluate integer arithmetic only. `2 + 3 * 4` must be 14.
3. `let`, identifiers, the environment.
4. Booleans, comparisons, `if/else`, `return`.
5. Functions, calls, closures. Write the counter test.
6. Strings, arrays, hashes, builtins.
7. The REPL, then error handling everywhere.

## The interesting problem

Closures. `fn() { count = count + 1; count }` must see the `count` of the
environment it was *created* in, not the one it is *called* in. A function
object therefore stores its environment alongside its body, and calling it
creates a new environment whose parent is that stored one. Three lines of code,
and the entire concept of lexical scope.

## Testing it

Every stage is independently testable, which is why this project is so good:

- lexer: input string → expected token slice.
- parser: input → `ast.String()`, compared against a canonical form.
- evaluator: table of `{input, expectedValue}` — dozens of one-line cases.
- fuzz the parser: no input, however malformed, may panic. It must always be an
  error with a position.

## Stretch

- Compile the AST to bytecode and write a stack VM. Benchmark it against the
  tree-walker on a recursive fib(30) — expect 3–10×.
- Add a mark-and-sweep garbage collector for objects.
- Add `while`, `for`, and a module system with `import`.
