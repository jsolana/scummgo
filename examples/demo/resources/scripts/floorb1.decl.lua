floor_minus1 = room {
    background = ref("resources:backgrounds/FloorBasementElevator"),
    walkboxes = {
        box0 = walkbox {
            vertices = {
                pos {x=5, y=140},
                pos {x=310, y=120},
                pos {x=310, y=98},
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
    door_basement = object {
        name = "basement",
        hotspot = rect {x=300, y=30, w=30, h=140},
        usedir = RIGHT,
        usepos = pos {x=300, y=140}
    },
}

export { room = floor_minus1 }
