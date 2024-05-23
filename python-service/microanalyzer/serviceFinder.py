import os


class ServiceNode:
    def __init__(self, path, parent):
        self.path = path
        self.name = os.path.basename(path)
        self.parent = parent
        self.children = list()
        self.has_py = False
        self.has_init = False
        self.is_service = False
        self.can_be_service = False
        self.npy = 0

    def update(self):
        for c in self.children:
            if c.name=="__init__.py":
                self.has_init = True
            elif c.has_py:
                self.has_py = True
        self.can_be_service = self.has_py and not self.has_init


class ServiceFinder:

    def __init__(self, source_path):
        self.source_path = source_path
        self.nodes = list()
        self.services = list()
        self.root = None

    def create_node(self, path, parent=None):
        assert os.path.isdir(path)
        node = ServiceNode(path, parent)
        self.nodes.append(node)
        for f in os.listdir(path):
            if f=="__init__.py":
                node.has_init = True
                node.npy += 1
            elif f.endswith(".py"):
                node.has_py = True
                node.npy += 1
            elif os.path.isdir(os.path.join(path, f)):
                node.children.append(self.create_node(os.path.join(path, f), node))
        node.update()
        return node

    def create_root(self):
        self.root = self.create_node(self.source_path, None)

    def find_services_root(self, current: ServiceNode):
        assert current in self.nodes
        potential_services = [n for n in current.children if n.can_be_service]
        if len(potential_services)==0:
            if current.can_be_service:
                return current.path, [current.path]
            else:
                raise Exception("service root can't be found!")
        if len(potential_services)>1:
            if current.npy==0:
                return current.path, [n.path for n in  potential_services]
            else:
                print(current.path)
                return current.path, [n.path for n in  potential_services]
        elif len(potential_services)==1:
            return self.find_services_root(potential_services[0])

    def get_services(self):
        assert self.root is not None
        root_path, services = self.find_services_root(self.root)
        return root_path, services
