breakroom = room {
    -- todo
    background = ref("resources:backgrounds/Breakroom"),
    walkboxes = {
        box0 = walkbox {
            vertices = {
                pos {x=5, y=140},
                pos {x=315, y=140},
                pos {x=315, y=105},
                pos {x=5, y=105},
            },
            scale = 1,
        },
        box1 = walkbox {
            vertices = {
                pos {x=110, y=105},
                pos {x=205, y=105},
                pos {x=205, y=90},
                pos {x=110, y=90},
            },
            scale = 1,
        },
    },
    door_corridor = object {
        name = "corridor",
        hotspot = rect {x=0, y=30, w=20, h=150},
        usedir = LEFT,
        usepos = pos {x=15, y=140}
    },
    door_bathroom = object {
        name = "bathroom",
        hotspot = rect {x=24, y=46, w=20, h=50},
        usedir = UP,
        usepos = pos {x=30, y=110}
    },
    microwave = object {
        class = APPLICABLE,
        name = "microwave",
        sprites = ref("resources:sprites/microwave"),
        hotspot = rect {x=100, y=50, w=18, h=11},
        pos = pos {x=115, y=64},
        usedir = UP,
        usepos = pos {x=135, y=80},
        default = state {
            anim = fixedanim { row = 0, col = 0},
        },
        turnon = state {
            anim = sequenceanim { duration = 100, row = 1, columns = {0,1,2}},
        },
        state = "default"
    }
}

export { room = breakroom }
