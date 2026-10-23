bathroom = import("resources:scripts/bathroom.decl").room
breakroom = import("resources:scripts/breakroom.decl").room
ego = import("resources:scripts/ego.decl").actor

music = music("resources:audio/OnTheHill")

function bathroom:enter(entrance)
    music:play()
    ego:enter(entrance)
    ego:select()
    CONTROL:paneon()
    CONTROL:cursoron()
    self:camfollow(ego)
end

function bathroom.door_breakroom:walkto()
    breakroom:show(breakroom.door_bathroom)
end
