basement = import("resources:scripts/basement.decl").room
colors = import("resources:scripts/colors.decl")
ego = import("resources:scripts/ego.decl").actor
floorb1 = import("resources:scripts/floorb1.decl").room

music = music("resources:audio/Basement")
beeps = sound("resources:audio/sfx/BasementBeeps")

rob = actor {
    name = "R.O.B",
    costume = ref("resources:costumes/Rob"),
    talkcolor = colors.yellow,
    usepos = pos {x=435, y=80},
    state = "idle"
}

function rob:talkto()
    ego:say("I always suspected about you being a machine..."):wait()
    rob:say("I always thought you were taller"):wait()
    rob:say("Life dissapoints me every day..."):wait()
end

function basement:enter(entrance)
    ego:enter(entrance)
    music:play()

    rob:show {
        pos = pos {x=443, y=86},
        lookat = DOWN,
    }

    ego:walkto(pos {x=33, y=135}):wait()

    if not basement:flag("IntroCompleted") then
        CONTROL:paneoff()
        CONTROL:cursoroff()
        ego:say("Hello Rob"):wait()
        ego:say("Still plotting world domination\nor just running diagnostics today?"):wait()
        basement:setflag("IntroCompleted")
    end

    ego:select()
    CONTROL:paneon()
    CONTROL:cursoron()
    self:camfollow(ego)
end

function basement.door_hall:lookat()
    ego:say("This secret door?\nBest investment ever—if only\nit came with a remote control!"):wait()
end

function basement.door_hall:walkto()
    ego:walkto(pos {x=10, y=140}):wait()
    floorb1:show(floorb1.door_basement)
end

function basement.switch_lights:lookat()
    ego:say("The lights switch"):wait()
end

function basement.darkness:walkto()

    if not basement:flag("LightsOn") then
        ego:walkto(pos {x=300, y=140}):wait()
        ego:say("I can't see anything!"):wait()
        ego:say("I shall risk a little more light..."):wait()
    end
end

function basement.switch_lights:use()
    ego:say("what if... I try turning it on... ejem, ejem!"):wait()
end

function basement.switch_lights:turnoff()
    if not basement:flag("LightsOn") then
        ego:say("Lights are already off..."):wait()
    end
end

function basement.switch_lights:turnon()
    if not basement:flag("LightsOn") then
        turn_on_lights()
        beeps:play()
    elseif basement:flag("LightsOn") then
        ego:say("Lights are already on!"):wait()
    end
end

function turn_on_lights()
    ego:walkto(pos {x=40, y=120}):wait()
    basement:setflag("LightsOn")
    ego:say("LET THERE BE LIGHT!"):wait()
    basement.darkness:classon(UNTOUCHABLE)
    basement.lights:classon(UNTOUCHABLE)
end
