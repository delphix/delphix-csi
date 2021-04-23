#
# Copyright (c) 2019 by Delphix. All rights reserved.
#
import pkg_resources, pkgutil

__path__ = __import__('pkgutil').extend_path(__path__, __name__)

resource_package = __name__
__version__ = pkgutil.get_data(resource_package, 'VERSION')
