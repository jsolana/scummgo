callrooms = room {
    background = ref("resources:backgrounds/CallRooms"),
    walkboxes = {
        box0 = walkbox {
            vertices = {
                pos {x=5, y=120},
                pos {x=310, y=120},
                pos {x=310, y=98},
                pos {x=23, y=98},
            },
            scale = 1,
        },
    },
    door_floor3 = object {
        name = "corridor",
        hotspot = rect {x=310, y=30, w=20, h=140},
        usedir = RIGHT,
        usepos = pos {x=310, y=105}
    },
    door_office = object {
        name = "office",
        hotspot = rect {x=0, y=30, w=30, h=140},
        usedir = LEFT,
        usepos = pos {x=4, y=105}
    },
}

export { room = callrooms }
