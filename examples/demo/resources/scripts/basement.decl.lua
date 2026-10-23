basement = room {

    background = ref("resources:backgrounds/Basement"),
    walkboxes = {
        box0 = walkbox {
            vertices = {
                pos {x=5, y=140},
                pos {x=300, y=140},
                pos {x=300, y=127},
                pos {x=23, y=127},
            },
            scale = 1,
        },
        box1 = walkbox {
            vertices = {
                pos {x=300, y=140},
                pos {x=450, y=140},
                pos {x=450, y=127},
                pos {x=300, y=127},
            },
            scale = 1,
        },
    },
    door_hall = object {
        name = "hall",
        hotspot = rect {x=10, y=30, w=10, h=150},
        usedir = LEFT,
        usepos = pos {x=15, y=135}
    },
    lights = object {
        name = "lights",
        sprites = ref("resources:sprites/basement-lights"),
        pos = pos {x=250, y=144},
        hotspot = rect {x=150, y=100, w=20, h=20},
        usedir = RIGHT,
        usepos = pos {x=140, y=120},
        default = state {
            anim = fixedanim { row = 0, col = 0},
        },
        state = "default"
    },
    switch_lights = object {
        name = "lights switch",
        pos = pos {x=41, y=85},
        usepos = pos {x=41, y=120},
        hotspot = rect {x=41, y=85, w=10, h=10},
        usedir = UP
    },

    darkness = object {
        name = "dark side of the room",
        pos = pos {x=300, y=0},
        hotspot = rect {x=300, y=0, w=300, h=144},
        usedir = RIGHT,
        usepos = pos {x=300, y=120},

    },


}

export { room = basement }
