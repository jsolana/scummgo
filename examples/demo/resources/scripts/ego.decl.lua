colors = import("resources:scripts/colors.decl")

bob = actor {
    name = "Bob",
    costume = ref("resources:costumes/Bob"),
    talkcolor = colors.white,
}

export { actor = bob }
