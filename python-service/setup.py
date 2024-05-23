# setup.py
from setuptools import setup, find_packages
import os
import re


MODULE_NAME = "microanalyzer"

with open(os.path.join(os.path.dirname(__file__), MODULE_NAME, "_version.py")) as f:
    version_file = f.read()
    version_match = re.search(r"^__version__ *= *['\"]([^'\"]*)['\"]", version_file, re.M)
    if version_match:
        __version__ = version_match.group(1)
    else:
        raise RuntimeError("Unable to find version string.")

if __name__ == "__main__":
    setup(
        name=MODULE_NAME,
        version=__version__,
        packages=find_packages(exclude=['tests', 'experimenting', 'analysis']),
        python_requires=">=3.6",
    )