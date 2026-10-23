# scummgo — Point & Click Adventure Toolkit

**scummgo** is a toolkit for creating classic *point & click* graphic adventures
(think SCUMM-era games like Monkey Island) in **Go**, using **raylib** for
rendering and **Lua** for all the game logic.

You build your game out of three kinds of things:

1. **A thin Go launcher** — a few lines that open the window and boot your
   first script.
2. **Lua scripts** — rooms, actors, objects, dialogs and cutscenes. This is
   where 95% of your game lives.
3. **Assets** — backgrounds, sprite sheets, costumes, music and sounds,
   described by tiny YAML manifests and packed into two binary files
   (`resources.idx` + `resources.dat`).

The screen is a fixed **320×200** retro layout (zoomable): a **320×144
viewport** for the current room on top, and a **320×56 control pane** below
with the classic 12 SCUMM verbs (*Open, Close, Push, Pull, Walk to, Pick up,
Talk to, Give, Use, Look at, Turn on, Turn off*), the current action sentence
and the player's inventory.

> ⚠️ This is still work in progress.

## Requirements

- Go 1.25+
- A C compiler (raylib-go needs one via cgo)

## Quick start

Run the bundled examples:

```bash
# Hello world: one room (Melee Island) with an intro cutscene,
# a bucket to pick up, and some pirates to talk to.
cd examples/hello-scummgo
make            # packs resources/ into resources.idx + resources.dat
make run

# A bigger multi-room adventure (an office with 8+ rooms, elevator,
# dialog trees, animations...).
cd examples/demo
make && make run
```

`make` (packing) is only needed again when your assets change. `make dump`
writes `resources.dump.json`, a canonical JSON dump of the package used for
cross-implementation validation.

## Anatomy of a game

Every game is a normal Go module plus a `resources/` directory:

```
mygame/
├── go.mod                  # requires github.com/jsolana/scummgo
├── main.go                 # ~10 lines of Go, that's all
├── Makefile
└── resources/
    ├── backgrounds/        # room images (+ .yml manifest each)
    ├── costumes/           # actor sprite sheets (+ .yml manifest each)
    ├── sprites/            # object sprite sheets (+ .yml manifest each)
    ├── audio/              # music and sound files (+ .yml manifest each)
    └── scripts/            # Lua scripts (packed as-is, no manifest needed)
```

### `main.go`

```go
package main

import "github.com/jsolana/scummgo"

func main() {
	loader := scummgo.NewResourceFileLoader("./") // reads ./resources.idx/.dat
	app := scummgo.New(loader,
		scummgo.WithScreenCaption("My Game"),
		scummgo.WithScreenZoom(4),
		scummgo.WithDebugMode(),
	)
	app.RunCommand(scummgo.ScriptRun{
		ScriptRef: scummgo.NewResourceRef("resources", "scripts/boot"),
	})
	app.Run()
}
```

The `ScriptRef` is `package:path/without/extension`. The package name is the
base name you pass to the packer (`resources` by default); the path is the
file inside `resources/` minus its extension. So `resources/scripts/boot`
means `resources/scripts/boot.lua`.

App options: `WithScreenCaption`, `WithScreenZoom` (default 4),
`WithDebugMode` (enables the `D` debug toggle in-game).

### `Makefile`

```makefile
.PHONY: all
all: resources

.PHONY: resources
resources: resources.idx resources.dat

resources.idx resources.dat:
	go run ../../cmd/scummgo pack resources

.PHONY: run
run: resources
	go run main.go

.PHONY: clean
clean:
	rm -f resources.idx resources.dat
```

### Assets and manifests

Every binary asset has a small YAML manifest next to it. The manifest
`Type` determines how the file is packed:

```yaml
# backgrounds/Melee.yml        → a plain image
type: image
compression: none
data:
  source: Melee.jpg
```

```yaml
# audio/OnTheHill.yml          → a music track (loops)
type: music
compression: none
data:
  source: "OnTheHill.ogg"
```

```yaml
# audio/Cricket.yml            → a one-shot sound effect
type: sound
compression: none
data:
  source: "Cricket.wav"
```

```yaml
# sprites/objects.yml          → a sprite sheet for objects
type: spritesheet
compression: none
data:
  source: objects.png
  frames:
    width: 30                  # frame size in the grid
    height: 18
```

Costumes (actor sprite sheets) map grid cells to animations. Four directions
(`right`, `left`, `up`, `down`; `flip: true` mirrors `right`), three actions
(`idle`, `walk`, `speak`):

```yaml
# costumes/Guybrush.yml
type: costume
compression: none
data:
  sprites:
    sheet: Guybrush.png
    width: 32
    height: 48
  animations:
    - action: walk
      dir: right
      frames:
      - row: 0
        columns: [0, 1, 2, 1, 0, 3, 4, 5, 4, 3]
        duration: 100          # ms per frame
    - action: walk
      dir: left
      flip: true
      frames:
      - row: 0
        columns: [0, 1, 2, 1, 0, 3, 4, 5, 4, 3]
        duration: 100
    # ...idle and speak for each direction
```

Costume animations are encoded in sorted action order so package layouts are
deterministic.

## Lua scripting: the decl / body / boot pattern

The examples organize scripts in three kinds of files:

- **`<thing>.decl.lua`** — *declarations*: defines entities (rooms, objects)
  and `export{}`s them. No behavior here.
- **`<thing>.body.lua`** — *behavior*: imports the decl and defines what
  happens (callbacks, dialog trees, cutscenes).
- **`boot.lua`** — the entry point: imports all body scripts, then shows the
  first room.

```lua
-- boot.lua
begin = import("resources:scripts/hall.decl").room

-- Load body scripts
import("resources:scripts/hall.body")
import("resources:scripts/basement.body")

begin:show(nil)   -- start the game in the hall, no entrance door
```

`import("resources:scripts/x")` runs a script once and returns its `export{}`
table. `export{}` publishes entities so other scripts can use them.

### Rooms

```lua
-- hall.decl.lua
hall = room {
    background = ref("resources:backgrounds/Hall"),
    walkboxes = {
        box0 = walkbox {
            vertices = {
                pos {x=10,  y=140},
                pos {x=310, y=140},
                pos {x=305, y=117},
                pos {x=10,  y=117},
            },
            scale = 1,             -- depth scaling (0..1), 1 = full size
        },
    },
    door_elevator = object {       -- hotspot-only object (no sprites)
        name = "elevator",
        hotspot = rect {x=115, y=50, w=15, h=20},
        usedir = UP,               -- direction the actor faces to interact
        usepos  = pos {x=126, y=97},   -- where the actor stands
    },
}

export { room = hall }
```

- **Walkboxes** are convex quads that define walkable floor. Actors walk
  through the graph of boxes; `scale` shrinks the actor for perspective.
  Boxes can start disabled (`enabled = false`) and be toggled at runtime with
  `room.boxN:enable() / :disable()` — a common trick to open new areas.
- **Objects** are interactive things. `hotspot` is the clickable area,
  `usepos`/`usedir` where the actor stands/faces, `sprites` + `pos` for
  objects with graphics, `class` for behavior flags (see below).

### Actors

```lua
-- ego.decl.lua
colors = import("resources:scripts/colors.decl")

juan = actor {
    name = "Juan",
    costume = ref("resources:costumes/Juan"),
    talkcolor = colors.white,
}

export { actor = juan }
```

```lua
-- colors.decl.lua
export {
    white   = color { r = 0xFF, g = 0xFF, b = 0xFF },
    magenta = color { r = 0xAA, g = 0x00, b = 0xAA },
}
```

Actors are placed and moved in behavior code:

```lua
ego:show   { pos = pos {x=140, y=140}, lookat = RIGHT }
ego:walkto(pos {x=200, y=140}):wait()   -- :wait() blocks until arrival
ego:say("Another day in the office."):wait()
ego:lookdir(LEFT)
ego:select()              -- make him the player-controlled character
ego:toinventory(bucket)   -- put an object in his inventory
ego:hide()
```

### Room lifecycle and camera

```lua
function hall:enter(entrance)   -- entrance = the door object used (or nil)
    if not hall:flag("IntroCompleted") then
        cutscene(function()
            -- non-interactive cinematic; the player can skip with ESC
            ego:show { pos = pos {x=-10, y=140}, lookat = LEFT }
            ego:walkto(pos {x=140, y=140}):wait()
            ego:say("Morning, Bridget!"):wait()
            -- ...dialog, walks, sleep(1000)...
        end)
        hall:setflag("IntroCompleted")   -- run-once logic via room flags
    end
    if entrance ~= nil then
        ego:enter(entrance)   -- place ego at the door's usepos, facing in
    end
    ego:select()
    CONTROL:paneon()
    self:camfollow(ego)       -- camera follows the player
end

function hall:exit()          -- optional, called before leaving the room
end
```

Rooms also offer `camto(pos)`, `camleft()`, `camright()` for manual camera
moves, and `flag(name)` / `setflag(name)` / `resetflag(name)` for state.

### Object verbs

Every control-pane verb maps to a callback on the object, named as the
lowercased verb without spaces: `open`, `close`, `push`, `pull`, `walkto`,
`pickup`, `talkto`, `give`, `use`, `lookat`, `turnon`, `turnoff`.

```lua
function hall.accescontrol:lookat()
    ego:say("It's an access control."):wait()
end

function hall.door_elevator:walkto()
    -- doors usually lead somewhere...
end

function melee.bucket:use(on)          -- `on` = second item of the sentence
    if on == melee.clock then
        guybrush:say("Time flies, but I don't think\nI can gather it in the bucket.")
    elseif on == pirates then
        melee.bucket:give(pirates)     -- delegate to another verb
    else
        DEFAULT.use(self, on)          -- fall back to the default response
    end
end
```

Objects with graphics can have **states** (e.g. a microwave on/off) declared
in the decl file and switched at runtime:

```lua
-- decl: default = state { anim = fixedanim { row = 0, col = 0 } },
--       turnon  = state { anim = sequenceanim { duration = 100, row = 1, columns = {0,1,2} } }
function breakroom.microwave:turnon()
    ego:walkto(pos {x=135, y=80}):wait()
    microwave_on:play()
    breakroom.microwave.turnon:set()    -- switch sprite state
    sleep(10000)
    microwave_off:play()
    breakroom.microwave.default:set()
end
```

### Default responses

Register catch-all verb handlers once (e.g. in a `defaults.body.lua` imported
by `boot.lua`) so every untouched object answers politely:

```lua
-- resources/scripts/defaults.body.lua
ego = import("resources:scripts/ego.decl").actor

DEFAULT = defaults()

function DEFAULT.pickup()
    ego:say("I can't pick that up.")
end
function DEFAULT.use()
    ego:say("I can't use that.")
end
function DEFAULT.lookat()
    ego:say("There is nothing special about that.")
end
-- ...open, close, push, pull, talkto, give, turnon, turnoff, walkto
```

> **Gotcha:** every `.lua` file is its own Lua interpreter. Globals do **not**
> leak between files — a variable defined in `a.lua` is `nil` inside `b.lua`.
> Always `import()` the entities you use, even in `DEFAULT` handlers.

### Object classes

Classes tell the verb UI how objects behave. Declared with
`class = APPLICABLE` (and toggled at runtime with `obj:classon(CLASS)` /
`obj:classoff(CLASS)`):

| Class         | Meaning                                                        |
|---------------|----------------------------------------------------------------|
| `PERSON`      | It's a person (right-click talks instead of looking)           |
| `UNTOUCHABLE` | Cannot be interacted with                                      |
| `PICKABLE`    | Can be picked up                                               |
| `OPENABLE`    | Right-click opens it                                           |
| `CLOSEABLE`   | Right-click closes it                                          |
| `APPLICABLE`  | Can be used as the first item of a `Use`/`Give` sentence       |

### Dialog trees

```lua
function pirates:talkto()
    while true do
        local choice = CONTROL:sentencechoice()
        choice:add("Hello, guys! Are you really\nsure you don't have the keys?")
        choice:add("I'm sure you have them! I know you!")
        choice:add("I think I will go to the Scumm bar. See you!")
        local s = choice:waitsay()   -- waits for pick; ego says the line aloud
        if s == 1 then
            pirates:say("Yes, we are sure!"):wait()
        end
        if s == 2 then
            pirates:say("No, we don't have them!"):wait()
        end
        if s == 3 then
            CONTROL:paneon()
            return
        end
    end
end
```

Use `choice:wait()` instead of `choice:waitsay()` when the player's pick
should not be spoken (e.g. elevator floor menus). `\n` breaks dialog lines.
Actor speech (`actor:say(text)`) stays on screen for a duration computed from
its length (at least 2s); append `:wait()` to block the script until it
finishes.

### Music and sound

```lua
music = music("resources:audio/Hall")     -- looped soundtrack
beeps = sound("resources:audio/sfx/Beeps") -- one-shot effect

music:play()   -- or :stop()
beeps:play()
```

### Control pane & inventory

`CONTROL` is a global object:

- `CONTROL:paneon() / :paneoff()` — show/hide the verb pane (use `paneoff`
  during cutscenes and intros).
- `CONTROL:cursoron() / :cursoroff()` — show/hide the mouse cursor.
- `CONTROL:sentencechoice()` — start a dialog tree (see above).

The inventory lives in the control pane: `ego:toinventory(obj)` adds an item;
clicking an inventory item and then a room item builds a `Use item with X` /
`Give item to X` sentence (the target must be `APPLICABLE`).

Right-click is a shortcut action: *Talk to* people, *Open/Close* openable
things, *Look at* everything else. Left-click on empty floor walks there.
`ESC` interrupts the current cutscene; `D` toggles the debug overlay
(walkboxes and hotspots) when the app was created `WithDebugMode()`.

### Animations

Two animation builders for object states:

- `fixedanim { row = 0, col = 0 }` — a single static frame of the sheet.
- `sequenceanim { duration = 100, row = 6, columns = {5, 6, 7} }` — a looping
  sequence of cells with a per-frame duration in ms.

## Tutorial: a tiny two-room adventure

Let's build *The Secret of the Boring Office* from scratch: a hall with a
mysterious light switch and a door to a basement. You'll explore both rooms,
get default responses for verbs you haven't scripted, and use room flags for
run-once events.

### 1. Scaffold

```bash
mkdir mygame && cd mygame
go mod init example.com/mygame
go get github.com/jsolana/scummgo
mkdir -p resources/{backgrounds,costumes,sprites,audio,scripts}
```

Create `main.go` and `Makefile` exactly as shown above (adjust the module
path; keep `go run ../../cmd/scummgo` only if you live inside this repo —
otherwise install the CLI once with `go install github.com/jsolana/scummgo/cmd/scummgo@latest`
and use `scummgo pack resources`).

### 2. Art

Drop two 320×144 (or larger; the camera pans) backgrounds:

- `resources/backgrounds/hall.png`
- `resources/backgrounds/basement.png`

And a character sheet `resources/costumes/hero.png` (a grid of 32×48 frames).
Write the manifests:

```yaml
# resources/backgrounds/hall.yml
type: image
compression: none
data:
  source: hall.png
```

```yaml
# resources/backgrounds/basement.yml
type: image
compression: none
data:
  source: basement.png
```

```yaml
# resources/costumes/hero.yml   — minimal: idle + walk to the right
type: costume
compression: none
data:
  sprites:
    sheet: hero.png
    width: 32
    height: 48
  animations:
    - action: idle
      dir: right
      frames:
      - row: 1
        columns: [0]
        duration: 1000
    - action: walk
      dir: right
      frames:
      - row: 0
        columns: [0, 1, 2, 1]
        duration: 100
```

(No audio? Skip it — music/sound are optional.)

### 3. Scripts

```lua
-- resources/scripts/colors.decl.lua
export {
    white = color { r = 0xFF, g = 0xFF, b = 0xFF },
}

-- resources/scripts/ego.decl.lua
colors = import("resources:scripts/colors.decl")

hero = actor {
    name = "Alex",
    costume = ref("resources:costumes/hero"),
    talkcolor = colors.white,
}

export { actor = hero }
```

```lua
-- resources/scripts/hall.decl.lua
hall = room {
    background = ref("resources:backgrounds/hall"),
    walkboxes = {
        floor = walkbox {
            vertices = {
                pos {x=5,   y=140},
                pos {x=315, y=140},
                pos {x=315, y=117},
                pos {x=5,   y=117},
            },
            scale = 1,
        },
    },
    door_basement = object {
        name = "basement door",
        hotspot = rect {x=280, y=60, w=30, h=80},
        usedir = RIGHT,
        usepos  = pos {x=270, y=138},
    },
    lightswitch = object {
        name = "light switch",
        hotspot = rect {x=20, y=70, w=10, h=14},
        usedir = UP,
        usepos  = pos {x=25, y=138},
    },
}

export { room = hall }

-- resources/scripts/basement.decl.lua
basement = room {
    background = ref("resources:backgrounds/basement"),
    walkboxes = {
        floor = walkbox {
            vertices = {
                pos {x=5,   y=140},
                pos {x=315, y=140},
                pos {x=315, y=117},
                pos {x=5,   y=117},
            },
            scale = 1,
        },
    },
    door_hall = object {
        name = "stairs",
        hotspot = rect {x=0, y=40, w=20, h=100},
        usedir = LEFT,
        usepos  = pos {x=15, y=138},
    },
}

export { room = basement }
```

```lua
-- resources/scripts/hall.body.lua
ego = import("resources:scripts/ego.decl").actor
basement = import("resources:scripts/basement.decl").room
hall = import("resources:scripts/hall.decl").room

function hall:enter(entrance)
    if entrance ~= nil then
        ego:enter(entrance)
    else
        ego:show { pos = pos {x=150, y=138}, lookat = RIGHT }
    end
    ego:select()
    CONTROL:paneon()
    self:camfollow(ego)
end

function hall.lightswitch:lookat()
    ego:say("A light switch. Mysteries await\nbehind every flick."):wait()
end

function hall.lightswitch:turnon()
    ego:say("It was already on... or off.\nHard to tell these days."):wait()
end

function hall.door_basement:walkto()
    basement:show(basement.door_hall)   -- travel through this door
end
```

```lua
-- resources/scripts/basement.body.lua
ego = import("resources:scripts/ego.decl").actor
hall = import("resources:scripts/hall.decl").room
basement = import("resources:scripts/basement.decl").room

function basement:enter(entrance)
    ego:enter(entrance)
    if not basement:flag("Visited") then
        basement:setflag("Visited")
        cutscene(function()
            ego:say("So dark... and so boring.\nDouble disappointment."):wait()
        end)
    end
    ego:select()
    CONTROL:paneon()
    self:camfollow(ego)
end

function basement.door_hall:walkto()
    hall:show(hall.door_basement)
end
```

```lua
-- resources/scripts/defaults.body.lua
ego = import("resources:scripts/ego.decl").actor

DEFAULT = defaults()

function DEFAULT.lookat()
    ego:say("Nothing special about that.")
end
function DEFAULT.use()
    ego:say("That doesn't work.")
end
function DEFAULT.pickup()
    ego:say("I can't pick that up.")
end
-- ...add the rest as you need them

-- resources/scripts/boot.lua
import("resources:scripts/defaults.body")
import("resources:scripts/hall.body")
import("resources:scripts/basement.body")

hall = import("resources:scripts/hall.decl").room

hall:show(nil)
```

### 4. Pack and run

```bash
make     # writes resources.idx + resources.dat
make run # or: go run main.go
```

Click the door: you'll walk to it and be transported to the basement. Click
the stairs to come back. Right-click anything for the quick action. Press
`D` for the debug overlay. That's the full loop — from here, grow your game
by adding rooms, actors, states and puzzles exactly like `examples/demo`
does.

## CLI reference

```bash
# Pack YAML/Lua assets into resources.idx + resources.dat
go run ./cmd/scummgo pack -o resources ./path/to/resources

# Dump a package as canonical JSON (for cross-impl validation with scummrs)
go run ./cmd/scummgo dump -o resources.dump.json resources
```
