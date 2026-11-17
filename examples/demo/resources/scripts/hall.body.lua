colors = import("resources:scripts/colors.decl")
floor3 = import("resources:scripts/floor3.decl").room
floor4 = import("resources:scripts/floor4.decl").room
floorb1 = import("resources:scripts/floorb1.decl").room
hall = import("resources:scripts/hall.decl").room
ego = import("resources:scripts/ego.decl").actor

music = music("resources:audio/Hall")

alice = actor {
    name = "alice",
    costume = ref("resources:costumes/Alice"),
    talkcolor = colors.magenta,
    usepos = pos {x=115, y=140},
    usedir = RIGHT
}


function hall:enter(entrance)
    if not hall:flag("IntroCompleted") then
        ego:select()
        cutscene(function()
            music:play()

            alice:show {
                pos = pos {x=350, y=140},
                lookat = RIGHT,
            }

            ego:show {
                pos = pos {x=-10, y=140},
                lookat = LEFT,
            }

            ego:walkto(pos {x=140, y=140})
            alice:walkto(pos {x=180, y=140}):wait()
            alice:say("Good morning, Bob!\nThank goodness you've arrived!"):wait()
            ego:say("Morning, Alice! What's up?"):wait()
            alice:say("We have a problem with\nthe Shitify campaign."):wait()
            ego:say("Do you mean our great and awesome\nmarketing campaign?"):wait()
            alice:say("Yes, that one."):wait()
            ego:say("What's the problem?"):wait()
            alice:say("Aparently, the Marketing VP\nhas a severe blocker."):wait()
            alice:say("You should go talk to Mark Miller."):wait()
            alice:say("You will find him at the third floor."):wait()
            ego:say("Ok, I'll do. Thanks!"):wait()
            alice:say("Good luck!"):wait()
            alice:walkto(pos {x=220, y=140}):wait()
            sleep(1000)
            alice:lookdir(LEFT)
            alice:play("checkwatch"):wait()
            alice:say("I don't know what's going on today\nbut the office seems a madhouse."):wait()
            alice:say("I mean more than usual."):wait()
            alice:say("See you later!"):wait()
            alice:walkto(pos {x=370, y=140}):wait()
            ego:lookdir(DOWN)
            sleep(1000)
            ego:say("Another day in the office."):wait()
            ego:say("Go for it!"):wait()

        end)
        hall:setflag("IntroCompleted")
        ego:show { pos = pos {x=140, y=140}, lookat = RIGHT }
        alice:hide()
    end
    if entrance ~= nil then
        ego:enter(entrance)
    else
        ego:show {
            pos = pos {x=140, y=140},
            lookat = RIGHT,
        }
    end
    self:camfollow(ego)
    CONTROL:paneon()
end

function hall.accesscontrol:lookat()
    ego:say("It's an access control."):wait()
end

function hall.accesscontrol:turnon()
    ego:say("Should I use the pinky?"):wait()
end

function hall.accesscontrol:open()
    ego:say("This fingerprint scanner isn't working again..."):wait()
    ego:say("I've tried my thumb, my index finger,\neven my pinky. Nothing works!"):wait()
    ego:say("As if using my pinky was going to magically fix it."):wait()
end

function hall.door_hub:lookat()
    ego:say("There are still people at the party"):wait()
    ego:say("Wait..."):wait()
    ego:say("Is that a goat?"):wait()
end

function hall.door_hub:walkto()
    ego:say("Maybe later..."):wait()
end

function hall.door_elevator:walkto()
    while true do
        local choice = CONTROL:sentencechoice()
        choice:add("4th floor")
        choice:add("3rd floor")
        choice:add("B1 floor")
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
            floorb1:show(floorb1.door_elevator)
            return
        end
        if s == 4 then
            CONTROL:paneon()
            return
        end
    end
end
