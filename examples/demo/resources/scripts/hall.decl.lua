hall = room {
    background = ref("resources:backgrounds/Hall"),
    walkboxes = {
        box0 = walkbox {
            vertices = {
                pos {x=10, y=140},
                pos {x=310, y=140},
                pos {x=305, y=117},
                pos {x=10, y=117},
            },
            scale = 1,
        },
        box1 = walkbox {
            vertices = {
                pos {x=182, y=108},
                pos {x=200, y=108},
                pos {x=200, y=117},
                pos {x=182, y=117},
            },
            scale = 1,
        },
        box2 = walkbox {
            vertices = {
                pos {x=184, y=108},
                pos {x=110, y=108},
                pos {x=110, y=100},
                pos {x=184, y=100},
            },
            scale = 1,
        },
    },
    accesscontrol = object {
        class = APPLICABLE,
        name = "access control",
        sprites = ref("resources:sprites/Access-Control"),
        pos = pos {x=165, y=110},
        hotspot = rect {x=156, y=84, w=25, h=35},
        usedir = LEFT,
        usepos = pos {x=200, y=90},
        default = state {
            anim = fixedanim { row = 0, col = 0},
        },
        pickup = state {},
        state = "default"
    },
    door_hub = object {
        name = "hub",
        hotspot = rect {x=310, y=30, w=10, h=150},
        usedir = RIGHT,
        usepos = pos {x=290, y=140}
    },
    door_elevator = object {
        name = "elevator",
        hotspot = rect {x=115, y=50, w=15, h=20},
        usedir = UP,
        usepos = pos {x=126, y=97}
    },
}

export { room = hall }
