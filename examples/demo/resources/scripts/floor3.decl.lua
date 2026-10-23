floor3 = room {
    background = ref("resources:backgrounds/Floor3Elevator"),
    walkboxes = {
        box0 = walkbox {
            vertices = {
                pos {x=5, y=140},
                pos {x=315, y=140},
                pos {x=315, y=98},
                pos {x=23, y=98},
            },
            scale = 1,
        },
    },
    door_elevator = object {
        name = "elevator",
        hotspot = rect {x=135, y=34, w=33, h=61},
        usedir = UP,
        usepos = pos {x=140, y=105}
    },
    door_call = object {
        name = "call rooms",
        hotspot = rect {x=0, y=30, w=30, h=140},
        usedir = LEFT,
        usepos = pos {x=15, y=113}
    },
    door_breakroom = object {
        name = "breakroom",
        hotspot = rect {x=310, y=30, w=30, h=140},
        usedir = RIGHT,
        usepos = pos {x=310, y=140}
    },
    poster = object {
        name = "poster",
        hotspot = rect {x=266, y=33, w=30, h=30},
        usedir = UP,
        usepos = pos {x=276, y=130}
    },
}

export { room = floor3 }
