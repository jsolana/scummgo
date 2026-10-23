basement = import("resources:scripts/basement.decl").room
canteen = import("resources:scripts/canteen.decl").room
ego = import("resources:scripts/ego.decl").actor
floor3 = import("resources:scripts/floor4.decl").room
floor4 = import("resources:scripts/floor4.decl").room
hall = import("resources:scripts/hall.decl").room

music = music("resources:audio/OnTheHill")

function floor4:enter(entrance)
    music:play()
    ego:enter(entrance)
    ego:select()
    CONTROL:paneon()
    CONTROL:cursoron()
    self:camfollow(ego)
end

function floor4.door_elevator:walkto()
    ego:walkto(pos {x=160, y=105}):wait()
    while true do
        local choice = CONTROL:sentencechoice()
        choice:add("3rd floor")
        choice:add("Hall")
        choice:add("B1 floor")
        choice:add("Exit")
        local s = choice:wait()
        if s == 1 then
            floor3:show(floor3.door_elevator)
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
        if s == 3 then
            CONTROL:paneon()
            return
        end
    end
end

function floor4.door_canteen:walkto()
    canteen:show(canteen.door_corridor)
end
