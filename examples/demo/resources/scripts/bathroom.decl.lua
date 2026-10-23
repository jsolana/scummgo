bathroom = room {
    -- todo
    background = ref("resources:backgrounds/Bathroom"),
    walkboxes = {
        box0 = walkbox {
            vertices = {
                pos {x=5, y=140},
                pos {x=310, y=140},
                pos {x=310, y=120},
                pos {x=23, y=120},
            },
            scale = 1,
        },
    },
    door_breakroom = object {
        name = "breakroom",
        hotspot = rect {x=310, y=30, w=20, h=140},
        usedir = RIGHT,
        usepos = pos {x=310, y=120}
    },
    -- police_line = object {
    --     class = APPLICABLE,
    --     name = "police line",
    --     sprites = ref("resources:sprites/Police-Line"),
    --     hotspot = rect {x=88, y=76, w=33, h=22},
    --     pos = pos {x=104, y=98},
    --     usedir = LEFT,
    --     usepos = pos {x=135, y=98},
    --     default = state {
    --         anim = fixedanim { row = 0, col = 0},
    --     },
    --     pickup = state {},
    --     state = "default"
    -- },
}

export { room = bathroom }
