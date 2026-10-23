basement = import("resources:scripts/basement.decl").room
ego = import("resources:scripts/ego.decl").actor
floor3 = import("resources:scripts/floor3.decl").room
floor4 = import("resources:scripts/floor4.decl").room
floorb1 = import("resources:scripts/floorb1.decl").room
hall = import("resources:scripts/hall.decl").room

music = music("resources:audio/OnTheHill")

function floorb1:enter(entrance)
    ego:enter(entrance)
    music:play()
    ego:select()
    CONTROL:paneon()
    CONTROL:cursoron()
    self:camfollow(ego)
end


function floorb1.door_elevator:walkto()
    while true do
        local choice = CONTROL:sentencechoice()
        choice:add("4th floor")
        choice:add("3rd floor")
        choice:add("Hall")
        choice:add("Exit")
        local s = choice:wait()
        if s == 1 then
            floor4:show(floor4.door_elevator)
            return
        end
        if s == 2 then
            floor3:show(floor3.door_elevator)
            return
        end
        if s == 3 then
            hall:show(hall.door_elevator)
            return
        end
        if s == 3 then
            CONTROL:paneon()
            return
        end
    end
end

function floorb1.door_basement:walkto()
    basement:show(basement.door_hall)
end
