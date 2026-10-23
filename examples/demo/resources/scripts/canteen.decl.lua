canteen = room {
    -- todo
    background = ref("resources:backgrounds/Canteen"),
    walkboxes = {
        box0 = walkbox {
            vertices = {
                pos {x=5, y=140},
                pos {x=510, y=140},
                pos {x=510, y=130},
                pos {x=23, y=130},
            },
            scale = 1,
        },
    },
    door_corridor = object {
        name = "corridor",
        hotspot = rect {x=0, y=30, w=20, h=140},
        usedir = LEFT,
        usepos = pos {x=10, y=140}
    },
}

export { room = canteen }
