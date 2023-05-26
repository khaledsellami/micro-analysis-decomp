import os
from typing import List


def get_sources(source_path: str) -> List[str]:
    paths = list()
    if os.path.isdir(source_path):
        for f in os.listdir(source_path):
            paths += get_sources(os.path.join(source_path, f))
    elif source_path.endswith(".py"):
        paths = [source_path]
    return paths


def walk(source_path, prefix=0):
    # "   " + "    ".join(["|" for i in range(prefix)]),
    string = "\n".join([" ".join([str(prefix), " |   "*max(prefix-1,0), "|---"*(prefix>0),  os.path.basename(source_path)])])
    if os.path.isdir(source_path):
        # print("   ", end="")
        # print(*["|" for i in range(prefix)], sep="    ")
        # print(prefix, " |   "*max(prefix-1,0), "|---"*(prefix>0),   os.path.basename(source_path)+"/", sep=" ")
        string += "/"
        strings=  list()
        has_py, has_init = False, False
        for f in os.listdir(source_path):
            if f=="__init__.py":
                has_init = True
            rstring, child_has_py = walk(os.path.join(source_path, f), prefix+1)
            strings.append(rstring)
            has_py = has_py or child_has_py
        if has_py and not has_init:
            string += "*"
        strings = [s for s in strings if s is not None]
        if len(strings)>0:
            return "\n".join([string] + strings), has_py
    else:
        if source_path.endswith(".py"):
            # print("   ", end="")
            # print(*["|" for i in range(prefix)], sep="    ")
            # print(prefix, " |   "*max(prefix-1,0), "|---"*(prefix>0),  os.path.basename(source_path), sep=" ")
            return string, True
    return None, False