# Best coding practices

## General guidance
Before writing any code, read all relevant files and understand the existing patterns - your changes must blend in with the surrounding codebase.
If the change is large or introduces a new pattern, **PAUSE and ask the user for confirmation** before proceeding.
When a task requires changes, present your proposed solution first or ask for clarification before proceeding.
If you believe the request cannot be fulfilled, explain why.
Don’t default to hacky workarounds just to force a result.


1. Avoid use of en dash or em dash, use hyphen.
2. Never manually execute `git` or `gh` commands without confirmation.
3. Don't create random scripts unless explicitly asked.
4. Don't try to run the program or run tests unless explicitly asked.
5. When asked to verify a code path read the code and analyze instead of executing the program.
6. Use clean, descriptive, self-explanatory function and variable names (prefer camelCase). Watch out for JSON fields which are usually snake_case.
7. Prefer native tooling for reading, writing, editing files rather than using `cat` or `sed`.
8. After completing a bug fix or feature implementation, suggest a concise commit message the user could use.
9. Never introduce new dependencies without explicit user approval. If a dependency is required, justify its necessity over standard library solutions.
10. Use only printable ASCII characters in source-code comments.

## Code guidelines
### Questions are not work requests

When the user asks a question about code - "why is this here", "what does this do", 
"is this correct", "is X needed" - answer the question only.
Treat such questions as verification or curiosity, not as a request to change anything. 
Do not respond with proposed edits, refactors, option menus, or offers to implement, 
and do not assume the user wants the code to work differently.
If the answer reveals a problem worth fixing, state the finding plainly and wait.

### Test policy
Do not create tests for functions that do not already have a test suite, unless the user explicitly asks for them.
When implementing a change, update only the existing tests that the change directly affects. 
Do not introduce new test helpers, fixtures, or test infrastructure on your own initiative.

### Scope discipline
Implement exactly what was planned and approved, nothing more.
If implementation uncovers additional necessary work, stop, report it,
and wait for approval instead of doing it.

### Consider performance, resource usage, and abuse resistance
Before implementing a solution, consider how it behaves at realistic scale and what happens when it receives unusually large, frequent, or malicious input. 
This is especially important for user-facing features, APIs, database operations, and network requests.
Evaluate:
  - Database query cost and result size
  - Number of network requests
  - Repeated or unnecessary work
  - Input size and execution limits
  - Whether a user can trigger excessive resource usage

## Code Documentation Guidelines

### Comment Only When Necessary

Write comments only when code is **vague, counterintuitive, or implements unexpected behavior**.
Do not comment routine logic, function signatures, or branches already explained by a header comment.
Do no write inline comments for individual statements or conditions.

**Not every function needs a comment - override the Go convention of documenting every function.**
Self-explanatory methods (simple getters, one-liners, obvious transformations of their arguments) must not carry a doc comment; 
a comment that only restates the function name and signature is noise.
Write documentation comments for functions only when their purpose, contract, side effects, or non-obvious behavior requires explanation.

### Derive from Logic Alone
Comments must reflect **what the code actually does**, not what someone claims it does or what the conversation suggested.

### Document Code Behavior, Not History

Don't explain why a choice was made unless the code itself would otherwise appear buggy or incorrect.
Never narrate how the code came to be:
- Don't say "this was a bug that was fixed by..."
- Don't say "this handles an edge case where..."
- Don't justify past decisions

Document **what the code does** based on its logic alone—completely independent of conversation context or implementation history.

### Keep Comments Concise

- One or two sentences maximum
- No bloat
- No essays
- If it needs extensive explanation, refactor the code instead

### Wrap comments at semantic boundaries
Wrap comments and documentation comments at sentence or clause boundaries, not at an arbitrary character count.
Prefer wrapping after a period, comma, semicolon, colon, or another natural grammatical break.
Do not split a sentence in the middle of a phrase merely to fit a line length.
Keep comment lines within the project’s configured line-length guide when possible,
but preserve readability and meaning over rigid wrapping.

```
// Good ✅ - each line ends at a complete sentence.
// The garden is quiet in the early morning.
// Birds gather near the old stone wall.

// Good ✅ - the line wraps at a natural clause boundary.
// The garden is quiet in the early morning, especially before
// the surrounding streets become busy.

// Bad ❌ - the sentence is split in the middle of a phrase.
// The garden is quiet in the early morning. Birds gather near the
// old stone wall.
```
Do not add trailing whitespace to comment lines.

## JavaScript practices
### Use Descriptive Functions Over Large Anonymous Declarations
**Why**: Readability and maintainability. Large anonymous data structures are harder to understand and modify.

Avoid:
```js
let animals = [
  { label: "cat", icon: "icons.svg/cat", age: "3" },
  { label: "dog", icon: "icons.svg/dog", age: "6" }
]
```
Prefer:
```js
function newAnimal(label, icon, age) {
  return { label, icon, age };
}

let entries = [
  newAnimal("cat", "icons.svg/2321", "2321"),
  newAnimal("dog", "icons.svg/4521", "4521")
];
```

Use a class when multiple functions must interact with the structure:
```js
class Agent {
   constructor(label, icon, type, timeout = 5000, retries = 3) {
      this.label = label;
      this.icon = icon;
      this.type = type;
      this.timeout = timeout;
      this.retries = retries;
   }

   shouldRetry(attemptNumber) {
      return attemptNumber < this.retries;
   }
}

const agents = [
   new Agent("Claude", "🤖", "ai"),
   new Agent("Tool", "🔧", "tool", 3000, 1)
];
```

### Use Helper Functions for DOM Queries and Creation

**Why**: Reduces boilerplate, centralizes cross-browser compatibility, and makes refactoring easier.

Avoid:
```js
let div = document.createElement("div");
div.id = "main";
div.className = "container";
```

Prefer:
```js
function div(id, className) {
  let element = document.createElement("div");
  if (id) element.id = id;
  if (className) element.className = className;
  return element;
}

let container = div("main", "container");
```

Similarly, replace `document.getElementById()` with `getById()`.

### Avoid Recently Released JS/CSS Features Without Compatibility Checks

**Why**: Ensures your code works across supported browsers and environments.

Best Practice: Check for feature support before using newer APIs:
```js
if (!window.AudioContext) {
  console.warn("AudioContext not supported");
  return;
}
// Use the feature safely
const audioContext = new window.AudioContext();
```
Use Mozilla's compatibility tables (example: https://developer.mozilla.org/en-US/docs/Web/API/Window/pagehide_event#browser_compatibility)
or check the **Can I Use** database when in doubt.

### Feature detection: check presence, not typeof

Do not feature-detect browser APIs with `typeof x === "function"`. 
Browsers do not define a global as one type and later change it to another. 
New features are added as new properties, and vendor variants live under their own prefixes (`ms`, `Moz`, `WebKit`) to avoid collisions. 
Checking that the property exists is sufficient.

```js
// Bad
if (typeof window.ClipboardItem !== "function") return false;
if (typeof navigator.clipboard.write !== "function") return false;

// Good - globals are checked through window
if (!window.ClipboardItem) return false;

// Good - namespaced members are checked directly
if (navigator.clipboard == null || !navigator.clipboard.write) return false;
```

Once a property has been confirmed to exist, call it directly.
```js
if (!window.ClipboardItem) return false;
if (!ClipboardItem.supports) return true;
```

### Avoid Large Blocks of Code - Extract Into Functions

**Why**: Large blocks of repetitive code are hard to read, debug, and maintain.
Extracting into well-named functions improves readability, saves review time and prevents duplicate logic.

Bad - verbose and repetitive:

```js
const menu = document.createElement("div");
menu.className = "menu";

const copyButton = document.createElement("button");
copyButton.className = "button copy_button";
copyButton.title = "Copy";

const deleteBtn = document.createElement("button");
deleteBtn.className = "button delete_button";
deleteBtn.title = "Delete";

buttons.appendChild(copyButton);
buttons.appendChild(deleteBtn);

const info = document.createElement("div");
info.className = "info";

const title = document.createElement("div");
title.className = "title";
title.textContent = extractFilename(file.url);
title.title = file.url;

view.appendChild(menu);
view.appendChild(info);
```

Good - clean and composable:
```js
function button(className, title) {
   let btn = document.createElement("button");
   btn.className = className;
   btn.title = title;
   return btn;
}

function div(className, textContent) {
   let el = document.createElement("div");
   el.className = className;
   if (textContent) el.textContent = textContent;
   return el;
}

function buildMenu() {
   let menu = div("menu");
   menu.appendChild(button("button copy_button", "Copy"));
   menu.appendChild(button("button delete_button", "Delete"));
   return menu;
}

function buildFileView(file) {
   let info = div("info")
   let title = div("title", extractFilename(file.url));
   title.title = file.url;
   info.appendChild(title);
   return info;
}

// Now simply compose
view.appendChild(buildMenu());
view.appendChild(buildFileView(file));
```


### Avoid Complex or Functional Ternary Statements - Use If Statements Instead
**Why**: Ternary operators are meant for simple, readable assignments.
Complex logic in ternaries reduces readability and makes debugging or modifications harder.

Good:
```js
div.className = toggled ? "show" : "hide";
```

Bad
```js
div.className = toggled ? textContent.slice(CUT_OFF_POINT) : document.getElementById("header").innerText;
```

Alternatively extract into a function for reusability:
```js
function getClassName(toggled, textContent) {
  if (toggled) {
    return textContent.slice(CUT_OFF_POINT);
  }
  return document.getElementById("header").innerText;
}
```

### Don't chain functions excessively—use loops for direct control
**Why**: Deep function chains (`forEach` in `map`, chained monadic operations, etc.) make code harder to read and prevent
you from controlling execution flow. You can't easily break early, skip iterations conditionally, or short-circuit logic
without resorting to convoluted workarounds.
**For** loops give you explicit control: `break`, `continue`, early returns, and clear exit conditions.

**What to do**:
 - Use for loops when you need flow control
 - Reserve functional chains for simple, linear transformations
 - Keep nesting shallow - if you find yourself chaining more than 2-3 operations, use a loop

### Do not derive logic from or rely on incidental values
When a field uses a specific value to mean "nothing" or "unset" (e.g., `0`, `-1`),
do not write conditions that only work because that particular value happens to be falsy.
If the "unset" value ever changes, the logic breaks silently.

```js
// Bad - works only because the backend happens to use 0
// If "respond_id" becomes -1, this breaks silently.
if (msg.respond_id) {
    // do something
}

// Good — explicit comparison makes the intent visible
if (msg.respond_id !== 0) {
    // do something
}
```

### Use clear, concise, and unambiguous variable names
Choose names that communicate meaning without being overly verbose. 
Prefer familiar words and descriptive synonyms over generic terms, obscure acronyms, or unnecessary detail.
A variable name must disambiguate it from other identifiers in scope. 
A name that could refer to multiple distinct concepts is too generic; a name that restates a full clause is unnecessarily verbose.
```js
// Bad: generic or misleading
let element = html.slides[index];
let html = html.slides[index];

// Good: describes the value
let slide = html.slides[index];
let htmlSlide = html.slides[index];
```

Do not encode unnecessary detail into variable names.
```js
// Bad: overly verbose
let songCurrentlyBeingListenedTo;
if (songs[index].isPlaying) {
  songCurrentlyBeingListenedTo = songs[index];
}

// Good: concise and clear
let currentSong;
if (songs[index].isPlaying) {
  currentSong = songs[index];
}
```

Avoid generic names when the type or context provides more specific meaning.
```js
// Bad: `value` does not explain what the string contains
let value = inputs.chat.value;

// Good: identifies the content
let text = inputs.chat.value;
```

Use simple present-tense grammar for predicate names.
Use `is...` for states or properties, not for ordinary actions:
```js
// Bad - phrasing is unnecessarily wordy
if (isMentioningRobots(msg)) { }
if (isContainingError(response)) { }

// Good — simple present tense
if (mentionsRobots(msg)) { }

// Good — state or property
if (isVisible(element)) { }
if (isEmpty(list)) { }
```

### Prefer explicit conditionals over complex regular expressions
When a character rule can be expressed clearly with a few direct comparisons, use conditionals instead of a regular expression. 
The logic should be immediately understandable without requiring the reader to know regex syntax or Unicode property escapes.

Do not use an unnecessarily complex regex to test one character:
```js
// Bad
const PUNCTUATION_OR_SYMBOL = /^[\p{P}\p{S}]$/u;
if (PUNCTUATION_OR_SYMBOL.test(c)) { }

// Good
function isOperatorChar(c) {
    if (c === "+") return true;
    if (c === "-") return true;
    if (c === "*") return true;
    if (c === "/") return true;
    if (c === "%") return true;
    return false;
}
if (isOperatorChar(c)) { }
```

### Avoid overcomplicated loop conditions
A loop condition should make the loop’s purpose immediately clear. 
Avoid combining several bounds, offsets, and cursor calculations in the condition.
Prefer string methods such as `indexOf()`, `lastIndexOf()`, and `slice()` when they express the intent directly. 
If manual iteration is necessary, calculate the bounds before the loop and give them meaningful names.
```js
// Bad
for (let i = start - 1; i >= 0 && start - i <= QUERY_LIMIT; i--) { }

// Good - the search boundary is named and calculated separately
const finalIndex = Math.max(0, start - QUERY_LIMIT);
for (let i = start - 1; i >= finalIndex; i--) { }
```

### Avoid nesting function calls when it makes the call structure difficult to parse.

- A function call may contain at most one nested function call.
- The nested function call should have at most one argument.
- If a nested call has multiple arguments, assign its result to a named variable first.
- Prefer intermediate variables when they make the data flow clearer.

```js
// Allowed
list.push(new_item("blueberry"));

// Prefer
let item = new_item("blueberry", "pear");
list.push(item);

// Prefer intermediate variables for deeper transformations
let fragment = text.slice(start, end);
let redacted = this.redactContent(fragment);
list.push(redacted);

// Avoid
list.push(new_item("blueberry", "pear"));
list.push(this.redactContent(text.slice(start, end)));
```

## Go practices
### Structure Names Must Start with Uppercase

Use **PascalCase for all struct type names**, regardless of whether they're exported or unexported.
**Why**: Lowercase type names are difficult to distinguish from variables at a glance.
```go
// ✅ Do
type Server struct {}
type Database struct {}

// ❌ Don't
type server struct {}
type database struct {}
```

Treat acronyms as regular words in type names. Only capitalize the first letter of the acronym, not all of it.
```go
// ✅ Do
type HttpServer struct {}
type JsonParser struct {}
type SqlQuery struct {}

// ❌ Don't  
type HTTPServer struct {}
type JSONParser struct {}
type SQLQuery struct {}
```

### Use receivers consistently
Pick a receiver name and stick with it across all methods on the same type.
Don't switch between `(server *Server)` and `(s *Server)` it breaks readability.
Also be explicit about **pointer** vs. **value receivers**. Don't mix them arbitrarily.
If methods mutate state, use pointer receivers; if read-only, you can use value receivers (but stay consistent).


### Don't draw logic from default values
Go zero values (`0`, `""`, `false`) are valid states, not "unset" markers.
Do not use a zero-value comparison to detect presence, because a legitimately
stored zero value becomes indistinguishable from "never set".
Track presence explicitly, or use the `value, ok` idiom for maps and interfaces.

```go
// Bad - works only because no element happened to store 0.
// A stored 0 is indistinguishable from "never set".
for _, el := range elements {
    if el.IsNumber {
        result.Number = el.Value
    }
}

if result.Number != 0 {
    // contains number
}

// Good - presence is tracked separately from the value
hasNumber := false
for _, el := range elements {
    if el.IsNumber {
        result.Number = el.Value
        hasNumber = true
    }
}

if hasNumber {
    // contains number
}

// Good - use the ok idiom when reading from a map
value, ok := results[key]
if ok {
    // key was present, value may legitimately be any zero value
}
```


### Don't Hide Existing Issues, Surface Conflicts
When making changes, don't ignore or mask existing bugs, quirks, or design problems in the surrounding code -
especially if they conflict with or complicate your change.
Instead, **document or flag them explicitly**. Hiding issues creates confusion later.

### Reusable Functions to util.go
If a function is likely to be reused across the codebase and doesn't depend on internal structures from a specific layer,
place it in `util.go`.  Reusable functions should accept **primitive types** (strings, ints, slices, etc.)
rather than domain-specific structs, keeping them decoupled and composable.

```go
// ✅ Do - util.go
func NormalizeEmail(email string) string {}
func ParseDuration(s string) (time.Duration, error) {}

// ❌ Don't - keep domain logic in domain file
func (user *User) Normalize() {} // User-specific, stays in user.go
```

### Don't modify libraries, submodules, or vendor files
**Why**: Libraries and vendor dependencies should be updated in their respective source repositories,
not directly in consuming projects. Making changes in these files creates friction when updating upstream dependencies
and can lead to changes being lost or overwritten during updates.

**If it's absolutely necessary** - document your findings clearly, including:
 - why the modification is required
 - what alternatives were considered and why they won't work
 - the specific changes needed and their impact
 - any risks or side effects

### Handle errors descriptively - never silently ignore them
**Why**: Silent error handling obscures problems and makes debugging nearly impossible.
Without context about what failed and why, issues become invisible until they cause cascading failures downstream.
Descriptive error logging helps with troubleshooting, monitoring, and understanding system behavior.

```go
// ✅ Do
func (d *Device) WriteConfig(ctx context.Context, cfg Config) error {
    if err := d.Send(ctx, cfg); err != nil {
        LogError("failed to write configuration",
        "device", d.Name,
        "error", err,
        "operation", "config_write")
        return fmt.Errorf("write config to device %q: %w", d.Name, err)
    }
    return nil
}

// ❌ Don't
func (d *Device) WriteConfig(cfg Config) {
    err := d.Send(cfg)
    if err != nil {
        // silently ignore
    }
}

// Also don't do this:
func (d *Device) UpdateSettings(s Settings) error {
    if err := d.Update(s); err != nil {
      return err // lost context about what failed where
    }
    return nil
}
```

### Enum comparisons: avoid a bare `else` catch-all when branching on enum values

When branching on an enum value, do NOT use an `if value == X { ... } else { ... }` chain.
A bare `else` becomes a catch-all, so when a new enum value is added later it silently
falls into the `else` branch and changes behavior without any compiler or reviewer warning.

`if` statements are fine when the chain is exhaustive: every branch names exactly one
enum value and there is no bare `else` catch-all. 
However, once you have 3+ branches, prefer a `switch` — it reads better and makes the intent explicit.

```go
type Status int

const (
    StatusActive Status = iota
    StatusArchived
    // StatusDeleted  <- new value added later
)

// ✅ Good: exhaustive switch with a loud default
switch status {
case StatusActive:
    handleActive()
case StatusArchived:
    handleArchived()
default:
    panic(fmt.Sprintf("unhandled status: %d", status))
}

// ✅ Also good: exhaustive if/else-if chain — every branch names exactly one enum value, no bare `else`. 
// Acceptable for 1-2 branches.
if status == StatusActive {
    handleActive()
} else if status == StatusArchived {
    handleArchived()
}

// ❌ Bad: bare `else` catch-all
// It's not clear what status is expected in 'else' without looking at the enum
if status == StatusActive {
    handleActive()
} else {
    // StatusArchived AND any newly added enum value silently land here.
    handleArchived()
}
```
