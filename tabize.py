global Set
poems = []

import os

class Poem:
    def __init__(self, title, dynasty, author, content, tags):
        self.title = title
        self.author = author
        self.dynasty = dynasty
        self.content = content
        self.tags = tags

    def __str__(self):
        return f"{self.title}\t{self.author}\t{self.dynasty}\t{self.content}\t{self.tags}\n"


def chop(s, p):
    i = s.index(p)
    if i == -1:
        return s
    # delete line break in head part
    heads = s[:i].split("\n")
    head = ""
    for v in heads:
        head += v
    global Set
    Set = s[i+len(p):]
    return head

if __name__ == "__main__":
	path = "/Users/apple/2022-日常工作/大熊猫繁育基地/08静静/180首诗.tab"

	with open(path, "r") as handle:
	     Set = handle.read()
	        
	while len(Set) > 2 :
	    item = Poem(chop(Set, "["), chop(Set, "]"), chop(Set, "\n"), chop(Set, "分类标签: "), chop(Set, "\n"))
	    poems.append(item)

	out = open("result.txt", "w")
	[out.write(str(v)) for v in poems]
	out.close()