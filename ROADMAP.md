# scummgo roadmap

This document lists what is missing to offer a complete point & click adventure engine, ordered by
priority. It comes from an audit of the code (file references are approximate and may drift) and
from a comparison with [inkwell](https://git.teletypegames.org/engines/inkwell), another SCUMM-style
framework (Go + Ebitengine, v0.4.1 at the time of writing; its facts come from its published
documentation, not from running it).

Sizes are a rough guide: **S** = hours, **M** = a few days, **L** = a week or more.

## 1. Where scummgo is today

| Area | What is there |
|---|---|
| Scripting | Lua (`decl` / `body` / `boot` pattern), modules with `export`/`import`, per-script interpreters, cutscenes, futures (`:wait()`), `sleep`, `every`, flags, callbacks per verb |
| World | Rooms with horizontal scrolling, objects with classes and animated states, actors, walkboxes (4-vertex convex quads, pathfinding, enable/disable), room entrances |
| Actors | Walking, talking, custom named animations (`actor:play`), perspective scale per walkbox (1/16 steps, blended across boxes, speed and hotspot follow, pixelation when far), inventory, ego switching |
| UI | Classic 12-verb SCUMM pane, sentence line, 6-slot inventory, dialog trees (`CONTROL:sentencechoice()`), EGA palette, debug overlay (walkboxes, hotspots) |
| Audio | One music stream, sound effects |
| Tooling | `scummgo pack` / `dump`, YAML manifests (image, spritesheet, costume, sound, music, script), `.idx` / `.dat` format with optional gzip, two examples, README with tutorial |

**What scummgo has that inkwell does not:** a real scripting language (inkwell scripts are Go
code), an asset packer and a binary package format with a dump tool, the authentic SCUMM verb pane
and EGA look, perspective scaling with pixelation, and costume animation data in YAML.

## 2. Comparison with inkwell

| Capability | scummgo | inkwell | Gap |
|---|---|---|---|
| Save / load, slots | No | Yes (`SaveSlot`/`LoadSlot`, file and browser stores, versioned saves) | **Big** |
| Settings and menus | No | Yes (menu widget, language, voice, theme, key bindings) | **Big** |
| Localization | No | Yes (translations kept with each text, runtime language switch) | **Big** |
| Layers / walk behind | Single background, Y-sorted items | Layers with Z (background, cast, foreground) | **Big** |
| Hotspot shapes | Rectangles | Rects and polygons | Medium |
| Walkboxes | Convex quads | Polygons sharing an edge | Medium |
| Triggers and exits | Only by script or objects | Triggers (condition edge) and generated exit hotspots | Medium |
| Custom verbs | Fixed 12 verbs | `Verb` entities, verb bar and radial wheel | Medium |
| Item combination | Only `APPLICABLE` + hand-written callbacks | `OnUseWith`, `OnUseWithAny`, `UseWithFail` | Medium |
| Dialog | Manual loop over `choice` | `Dialog` / `DialogNode` with `Show` and `Once`, escalating lines | Medium |
| Transitions | None | Fade / FadeTo, stages (splash, title, ending) | Medium |
| Validation | None | `Game.Validate` checks references | Medium |
| Placeholders for missing art | `log.Fatalf` | Deterministic placeholder graphics | Medium |
| Debug tools | `D` overlay | F1 hotspot overlay, debug menu, in-game explorer | Small |
| Themes | EGA only | Four theme presets | Small |
| Voice-over | No | Yes (per translation) | Small |
| Platforms | Desktop | Desktop, js/wasm store, Steam extension | Medium |
| CI, lint, API docs | Workflow not wired up, README only | Woodpecker, revive, godoc published per tag | Medium |
| Gamepad, touch, mobile | No | Not found in its docs | Same |
| Parallax, lighting, particles | No | Not found in its docs | Same |

The point is not to copy inkwell's Go DSL: scummgo keeps its data-driven, Lua-first design. The table
is a checklist of capabilities that players and authors expect from a complete engine.

## 3. Backlog by priority

### P0 — Stability, legal and foundations

The engine must not crash on ordinary authoring mistakes, and the project must be safe to share.

- **P0-1 License and attribution (S).** There is no `LICENSE` file. Add one and check that it is
  compatible with the license of rvrt-win, from which this framework starts (see README).
- **P0-2 Wire up CI (S).** `workflows/build.yml` is at the repo root, so GitHub never runs it. Move
  it to `.github/workflows/`, use the Go version of `go.mod` (it pins 1.22 while `go.mod` says
  1.25.1), add `go vet`, and fix the `-timeout 10s` that kills anything slow.
- **P0-3 Missing `defaults()` crashes (S).** `ObjectCall` and `ActorCall` dereference
  `app.defaults` (`command_object.go`, `command_actor.go`) and panic when a verb has no callback and
  the game never called `defaults()`. Provide a built-in default response.
- **P0-4 Rooms without walkboxes panic (S).** `ActorWalkToPosition` and
  `Room.GetScriptChildEntity` use `wbmatrix` unconditionally. Fall back to a straight line walk.
- **P0-5 Stop using `log.Fatalf` / `panic` for ordinary errors (M).** A missing resource, bad
  image or audio file, corrupt package, missing room background, duplicate object tag, unknown script
  or invalid costume `dir` end the process (`encoding.go`, `image.go`, `music.go`, `sound.go`,
  `room.go`, `pack/costume.go`). Return errors, show a clear message and, for missing art, draw a
  placeholder so a game can run before all assets exist.
- **P0-6 Surface script errors (M).** Errors inside `cutscene(fn)` are swallowed (`enterCutscene`
  ignores `ProtectedCall`), and callback errors only reach the log. Report script, line and message,
  and show them in the debug overlay.
- **P0-7 Timers and cutscene cancel (M).** `every()` goroutines are never stopped (they survive room
  changes and app close) and `sleep()` is a real blocking sleep that `ESC` cannot interrupt, so
  skipping a cutscene waits for the current sleep. Make timers cancellable (per room, per script) and
  sleeps interruptible.
- **P0-8 Inventory beyond 6 items (M).** `ControlInventory` has 6 fixed slots, no scrolling, and
  reserves room for arrows that are never drawn. Add scrolling and remove the possible out-of-range
  access.
- **P0-9 Small bugs found in the audit (S).** `ActorWalkToItem` misses a `return` after the
  "in the inventory" error; `ScriptRun` returns the command instead of the script; `ScriptImport`
  does not cache the loaded script; hit-testing and drawing iterate Go maps, so overlapping
  hotspots are resolved non-deterministically (use ordered slices); `ESC` interrupts every loaded
  script because the current cutscene is not tracked.
- **P0-10 Testability (M).** Tests that call `New()` need a window and deadlock (two Lua tests are
  skipped). Add a headless mode (no window, no audio) and cover the untested core: actors and
  actions, futures, commands, dialog, control pane sentence logic, manifest parsers and the CLI.

### P1 — Core features of a complete adventure engine

- **P1-1 Save / load game (L).** *(already in the previous roadmap)* Capture the current room,
  actor positions, facing and rooms, inventories, object states, owners and classes, room flags, music
  and the ego. Design question to settle first: scripts keep state in Lua globals, which cannot be
  saved, so provide a persistent `state` / variables API. Versioned format, several slots, refuse
  corrupt or newer saves without touching the running game, Lua API `save(slot)` / `load(slot)`.
- **P1-2 Walk-behind layers (L).** Rooms have a single flat background, so actors cannot pass behind
  scenery. Add foreground layers or masks with a Z order (SCUMM "z-planes"). Largest visual gap
  compared with the classic games.
- **P1-3 Dialog quality (M).** Word wrap (authors write `\n` by hand today), click / key to skip a
  line, speech time by text length with a setting, use `ActorSpeak.Delay`, and helpers for dialog
  trees (conditional choices, "once" choices, ordered loops) so scripts do not rewrite the same
  `while true` pattern.
- **P1-4 Verbs and items (M).** Push, Pull, Turn on and Turn off exist in the UI but have no
  behaviour. Provide default open / close toggling driven by object states, item combination
  (`use A with B`, in both orders, with a generic failure message), and make the verb set
  configurable per game.
- **P1-5 Lua standard libraries (S).** Only `lua.BaseOpen` is loaded, so `math`, `string` and
  `table` are unavailable. Open a safe subset (no `os`, `io`) so authors can use `string.format`,
  `math.random` and `table.insert`.
- **P1-6 Transitions and camera (M).** Fade in / out between rooms and for cutscenes, smooth `camto`
  pan (it is instant today), screen shake. Fixed 2 px/frame camera speed should be configurable.
- **P1-7 Audio (M).** Volume, fades, loop flag, per-room music, several simultaneous sounds, and no
  reload from the loader when a stopped track is played again. Prerequisite for the settings menu.
- **P1-8 Regions, triggers and exits (M).** Invisible areas that fire a callback when an actor
  enters, and room exits by walking off the edge of the screen, instead of building every transition
  from objects.
- **P1-9 Actor and object API gaps (M).** Expose `Elev` and walk speed to Lua, change object
  position and visibility from scripts, object state machines with conditions, actor following,
  and ego switching with its own inventory display (Day of the Tentacle style).
- **P1-10 Intro, title and ending support (M).** *(already in the previous roadmap)* A startup
  sequence: splash and title screens, a skippable intro cutscene with image sequences or animation
  resources, and credits / ending screens.

### P2 — Player-facing polish

- **P2-1 Settings menu (M).** *(already in the previous roadmap)* Language, music and effects
  volume, text speed, fullscreen / zoom and key bindings, persisted in a settings file.
- **P2-2 Main, pause and save menus (M).** Title menu, pause, and save / load screens built on the
  settings and save systems.
- **P2-3 Localization (L).** String tables per language in the resources, a lookup function for
  scripts, localized verb names and engine messages, fonts with accented characters, runtime language
  switch.
- **P2-4 Window and display (S).** The window is a fixed integer zoom chosen at start. Add
  fullscreen toggle, resizable window with integer scaling and letterboxing.
- **P2-5 Input (M).** Keyboard shortcuts for verbs, configurable bindings, and gamepad (and later
  touch) with a software cursor.
- **P2-6 Polygonal walkboxes and hotspots (M).** Convex polygons with any number of vertices
  instead of exactly four, and polygon hotspots. Optional scale gradients inside a box.
- **P2-7 Voice-over (M).** Optional recorded line per text, played instead of the timer.

### P3 — Authoring tools and distribution

- **P3-1 Dev loop (M).** Run from loose files without packing (`scummgo run`) and hot reload of
  scripts and assets, so the `make` step is not needed on every change. The loader also reopens the
  data file for every resource; cache it.
- **P3-2 Project check (M).** `scummgo check`: validate references between resources, rooms,
  costumes and scripts, report missing assets and Lua syntax errors before running.
- **P3-3 Scaffolding (S).** `scummgo new` to create a project with `main.go`, a `Makefile` and a first
  room instead of copying an example.
- **P3-4 Aseprite integration (M).** Generate costume and spritesheet manifests from Aseprite tags,
  and keep the PNG in sync with the `.aseprite` sources.
- **P3-5 Debug tools (M).** Beyond the overlay: console, flag and room inspector, teleport, names
  of walkboxes and objects, frame stats, interactive walkbox and hotspot editing.
- **P3-6 Documentation (M).** Lua API reference generated from the bindings, guides for animations
  and scale, a changelog and tagged releases.
- **P3-7 Packaging (M).** Single-binary distribution with embedded resources (`go:embed`), builds
  for macOS, Windows and Linux, and an evaluation of web (wasm) support.
- **P3-8 Performance (S).** Avoid per-frame allocations in `Room.Draw` and repeated walkbox
  lookups in `Actor.scale()`; profile the main loop.

### P4 — Nice to have

Parallax scrolling and vertical scrolling, lighting and palette effects, a CRT filter, tweens for
actors and camera, particles, achievements, UI themes (Sierra-style coin, MI2 verb layouts), Steam
integration, mobile.

## 4. Done

These items were in the previous roadmap and are implemented:

- **Custom animations for actors.** Named animations declared in the costume and played with
  `actor:play`, plus `every()` and `actor:isspeaking()` for background animations.
- **Costume scale.** Scale per walkbox in 1/16 steps, blended across shared edges, with walking
  speed, hotspot and dialog position following it, and larger pixels for distant actors.
- **Debug hotspots.** The debug overlay outlines object hotspots next to the walkboxes.
