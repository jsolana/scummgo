bathroom = import("resources:scripts/bathroom.decl").room
callrooms = import("resources:scripts/callrooms.decl").room
colors = import("resources:scripts/colors.decl")
ego = import("resources:scripts/ego.decl").actor
office = import("resources:scripts/office.decl").room

music = music("resources:audio/OnTheHill")

grace = actor {
    name = "grace",
    costume = ref("resources:costumes/Grace"),
    talkcolor = colors.yellow,
    usepos = pos {x=300, y=130},
    usedir = LEFT
}

function office:enter(entrance)
    self:camright()
    ego:enter(entrance)
    self:camfollow(ego)
    music:play()


    grace:show {
        pos = pos {x=270, y=130},
        lookat = RIGHT,
    }

    if bathroom:flag("ToiletIsClogged") then
        self:camfollow(ego)
        sleep(10000)
        ego:say("Hey, Grace, something weird is going on in the bathroom..."):wait()
        grace:say("What do you mean by 'weird'?"):wait()
        ego:say("Let's just say the toilet decided to take a break."):wait()
        grace:say("So... it's clogged."):wait()
        ego:say("Yeah, it's like the Niagara Falls, but... backward."):wait()
        grace:say("Fine, I'll handle it. But you owe me."):wait()
        ego:say("Deal! You're a lifesaver!"):wait()
        grace:walkto(pos {x=400, y=130}):wait()
    end
    ego:select()
    CONTROL:paneon()
    CONTROL:cursoron()

end

function office.lost_closet:open()
    office.lost_closet.opened:set()
end

function office.lost_closet:close()
    office.lost_closet.closed:set()
end

function office.door_callrooms:walkto()
    callrooms:show(callrooms.door_office)
end

office.graceAskedAboutCloset = false

function grace:talkto()
    while true do
    local choice = CONTROL:sentencechoice()
        choice:add("I need to take a look at that closet")
        if office.graceAskedAboutCloset then
            choice:add("Do you know I AM the CEO?")
        end
        choice:add("Never mind")
        local s = choice:waitsay()
        if s == 1 then
            grace:say("You guys are always making\na mess with that closet, go away!"):wait()
            office.graceAskedAboutCloset = true
        elseif office.graceAskedAboutCloset and s == 2 then
            grace:say("Oh, really?\nI AM the Office Manager"):wait()
            grace:say("Do you want to play THAT game?"):wait()
        else
            CONTROL:paneon()
            return
        end
    end
end
