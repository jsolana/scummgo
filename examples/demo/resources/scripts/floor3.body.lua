basement = import("resources:scripts/basement.decl").room
breakroom = import("resources:scripts/breakroom.decl").room
callrooms = import("resources:scripts/callrooms.decl").room
floor3 = import("resources:scripts/floor3.decl").room
floor4 = import("resources:scripts/floor4.decl").room
hall = import("resources:scripts/hall.decl").room
ego = import("resources:scripts/ego.decl").actor

music = music("resources:audio/OnTheHill")

function floor3:enter(entrance)
    music:play()
    ego:enter(entrance)

    ego:select()
    CONTROL:paneon()
    CONTROL:cursoron()
    self:camfollow(ego)
end


function floor3.door_elevator:walkto()
    while true do
        local choice = CONTROL:sentencechoice()
        choice:add("4rd floor")
        choice:add("Hall")
        choice:add("B1 floor")
        choice:add("Exit")
        local s = choice:wait()
        if s == 1 then
            floor4:show(floor4.door_elevator)
            return
        end
        if s == 2 then
            hall:show(hall.door_elevator)
            return
        end
        if s == 3 then
            basement:show(basement.door_hall)
            return
        end
        if s == 4 then
            CONTROL:paneon()
            return
        end
    end
end

function floor3.door_call:walkto()
    callrooms:show(callrooms.door_floor3)
end

function floor3.door_breakroom:walkto()
    breakroom:show(breakroom.door_corridor)
end

function floor3.poster:lookat()
    ego:say("Be Shitified"):wait()
    ego:say("Our ideas, words, actions,\nand goals are aligned."):wait()
    ego:say("Our purpose is to change the world,\nand we are committed to clicks, banners,\nthe metrics, and the buzzword of the week"):wait()
    ego:lookdir(DOWN)
    ego:say("Cool!"):wait()
end
