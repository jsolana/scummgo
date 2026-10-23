DEFAULT = defaults()

function DEFAULT.pickup()
    bob:say("I can't pick that up.")
end

function DEFAULT.use()
    bob:say("I can't use that.")
end

function DEFAULT.open()
    bob:say("I can't open that.")
end

function DEFAULT.close()
    bob:say("I can't close that.")
end

function DEFAULT.pull()
    bob:say("I can't pull that.")
end

function DEFAULT.push()
    bob:say("I can't push that.")
end

function DEFAULT.talkto(what)
    if what.type == "actor" then
        bob:say("It's not time for a chat.")
    else
        bob:say("I can't talk to that.")
    end
end

function DEFAULT.lookat()
    bob:say("There is nothing special about that.")
end

function DEFAULT.turnon()
    bob:say("I can't turn that on.")
end

function DEFAULT.turnoff()
    bob:say("I can't turn that off.")
end

function DEFAULT.give()
    bob:say("I can't give that.")
end

function DEFAULT.walkto()
end
