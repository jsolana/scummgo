office = room {
    background = ref("resources:backgrounds/OfficeFloor"),
    walkboxes = {
        box0 = walkbox {
            vertices = {
                pos {x=35, y=140},
                pos {x=410, y=140},
                pos {x=410, y=121},
                pos {x=60, y=121},
            },
            scale = 1,
        },
    },
    door_callrooms = object {
        name = "call rooms",
        hotspot = rect {x=410, y=30, w=20, h=140},
        usedir = RIGHT,
        usepos = pos {x=420, y=130}
    },
    lost_closet = object {
        class = APPLICABLE,
        name = "lost and found",
        sprites = ref("resources:sprites/lost-closet"),
        hotspot = rect {x=225, y=90, w=40, h=25},
        pos = pos {x=240, y=120},
        usedir = LEFT,
        usepos = pos {x=240, y=120},
        closed = state {
            anim = fixedanim { row = 0, col = 0},
        },
        opened = state {
            anim = fixedanim { row = 0, col = 1},
        },
        pickup = state {},
        state = "closed"
    },
}

export { room = office }
