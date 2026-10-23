colors = import("resources:scripts/colors.decl")
ego = import("resources:scripts/ego.decl").actor
floor4 = import("resources:scripts/floor4.decl").room
canteen = import("resources:scripts/canteen.decl").room

music = music("resources:audio/OnTheHill")

mark = actor {
    name = "mark",
    costume = ref("resources:costumes/Mark"),
    talkcolor = colors.red,
    usepos = pos {x=400, y=140},
    usedir = RIGHT
}

function mark:talkto()
    while true do
        local choice = CONTROL:sentencechoice()
        choice:add("How is the marketing\n campaing going?")
        choice:add("Do you need me to\n give out some flyers?")
        choice:add("See you")
        local s = choice:waitsay()
        if s == 1 then
            mark:say("Well... We are missing some\n designs to put everything\n together"):wait()
            mark:say("It's weird, Steve [Design VP]\n does not usually \nmake us wait"):wait()
            mark:say("Could you go ask\n him about this? I think\n He is in the Callroom"):wait()
            callrooms.last_known_state.should_run_initial_state = false
        end
        if s == 2 then
            mark:say("I don't think it's neccessary,\n but thank you"):wait()
        end
        if s == 3 then
            CONTROL:paneon()
            return
        end
    end
end

function canteen:enter(entrance)
    ego:enter(entrance)
    music:play()
    mark:show {
        pos = pos {x=430, y=140},
        lookat = LEFT,
    }
    ego:select()
    CONTROL:paneon()
    CONTROL:cursoron()
    self:camfollow(ego)
end

function canteen.door_corridor:walkto()
    floor4:show(floor4.door_canteen)
end
