import ast
import logging

from models import Object_, Executable_


class MicroVisitor(ast.NodeVisitor):
    def __init__(self, file_path: str, service_name: str, module_name: str = ""):
        self.file_path = file_path
        self.service_name = service_name
        self.source_code = None
        self.module_name = module_name
        self.last_visited_class = None
        self.classes = list()
        self.methods = list()
        self.logger = logging.Logger(__name__)
        if self.module_name == "":
            fullname = "{}".format(self.service_name)
        else:
            fullname = "{}.{}".format(self.service_name, self.module_name)
        self.visit_stack = [fullname]
        # self.invocations = 0

    def visit(self, item):
        if isinstance(item, ast.ClassDef):
            self.visit_ClassDef(item)
        elif isinstance(item, ast.FunctionDef):
            self.visit_FunctionDef(item)
        else:
            self.generic_visit(item)

    def visit_ClassDef(self, node):
        assert len(self.visit_stack) > 0
        fullname = "{}.{}".format(self.visit_stack[-1], node.name)
        object_ = Object_(False, False, node.name, fullname, self.file_path, self.service_name,
                              ast.get_source_segment(self.source_code, node))
        self.visit_stack.append(fullname)
        self.last_visited_class = fullname
        self.classes.append(object_)
        self.logger.debug("entering class{}: {}".format(len(self.classes), node.name))
        self.generic_visit(node)
        self.last_visited_class = None
        name = self.visit_stack.pop()
        assert name == fullname
        self.logger.debug("exiting class{}: {}".format(len(self.classes), node.name))

    def visit_FunctionDef(self, node):
        assert len(self.visit_stack) > 0
        fullname = "{}.{}".format(self.visit_stack[-1], node.name)
        executable_ = Executable_(fullname, node.name, self.last_visited_class, self.service_name,
                              ast.get_source_segment(self.source_code, node))
        self.visit_stack.append(fullname)
        self.methods.append(executable_)
        self.logger.debug("entering method{}: {}".format(len(self.methods), node.name))
        self.generic_visit(node)
        name = self.visit_stack.pop()
        assert name == fullname
        self.logger.debug("exiting method{}: {}".format(len(self.methods), node.name))

    def analyze(self):
        with open(self.file_path, "r") as f:
            self.source_code = f.read()
        try:
            root_node = ast.parse(self.source_code)
            self.generic_visit(root_node)
        except SyntaxError as e:
            self.logger.info("Syntax error encountered when parsing {}".format(self.module_name))
            return False
        return True