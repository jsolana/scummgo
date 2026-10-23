bathroom = import("resources:scripts/bathroom.decl").room
breakroom = import("resources:scripts/breakroom.decl").room
floor3 = import("resources:scripts/floor3.decl").room
ego = import("resources:scripts/ego.decl").actor

music = music("resources:audio/OnTheHill")
microwave_on = sound("resources:audio/sfx/MicrowaveOn")
microwave_off = sound("resources:audio/sfx/MicrowaveOff")

function breakroom:enter(entrance)
    music:play()
    ego:enter(entrance)
    ego:select()
    CONTROL:paneon()
    CONTROL:cursoron()
    self:camfollow(ego)
end


function breakroom.door_corridor:walkto()
    floor3:show(floor3.door_breakroom)
end

function breakroom.door_bathroom:walkto()
    bathroom:show(bathroom.door_breakroom)
end

function breakroom.microwave:turnon()
    ego:walkto(pos {x=135, y=80}):wait()
    microwave_on:play()
    sleep(1500)
    breakroom.microwave.turnon:set()
    sleep(10000)
    microwave_off:play()
    breakroom.microwave.default:set()
end