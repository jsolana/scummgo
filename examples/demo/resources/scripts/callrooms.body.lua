
callrooms = import("resources:scripts/callrooms.decl").room
ego = import("resources:scripts/ego.decl").actor
floor3 = import("resources:scripts/floor3.decl").room
office = import("resources:scripts/office.decl").room

music = music("resources:audio/OnTheHill")

tyler = actor {
    name = "tyler",
    costume = ref("resources:costumes/Tyler"),
    usepos = pos {x=290, y=140},
    usedir = LEFT
}

function callrooms:enter(entrance)
    music:play()
    ego:enter(entrance)

    tyler:show {
        pos = pos {x=240, y=120},
        lookat = RIGHT,
    }
    ego:select()
    CONTROL:paneon()
    CONTROL:cursoron()
    self:camfollow(ego)
end

function callrooms.door_floor3:walkto()
    floor3:show(floor3.door_call)
end

function callrooms.door_office:walkto()
    office:show(office.door_callrooms)
end

function tyler:talkto()
    while true do
        local choice = CONTROL:sentencechoice()
        if not callrooms.last_known_state.should_run_initial_state then
            choice:add("I need the designs for the marketing campaign")
            choice:add("How are you?")
            choice:add("Bye!")
            local s = choice:waitsay()
            if s == 1 then
                tyler:say("We lost the design book!\nWe can't finish our work without it"):wait()
                tyler:say("Find me that book and that missing work is done"):wait()
            end
            if s == 2 then
                tyler:say("Can't complain"):wait()
            end
            if s == 3 then
                CONTROL:paneon()
                return
            end
        elseif callrooms.last_known_state.should_run_initial_state then
            choice:add("How are you?")
            choice:add("Bye!")
            local s = choice:waitsay()
            if s == 1 then
                tyler:say("Can't complain"):wait()
            end
            if s == 2 then
                CONTROL:paneon()
                return
            end
        end
    end
end
