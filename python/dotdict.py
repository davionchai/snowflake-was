from collections.abc import Iterable


class DotDict(dict):
    def __init__(self, *args, **kwargs):
        super(DotDict, self).__init__(*args, **kwargs)

        args_with_kwargs = []
        for arg in args:
            args_with_kwargs.append(arg)
        args_with_kwargs.append(kwargs)
        args = args_with_kwargs

        for arg in args:
            if isinstance(arg, dict):
                for k, v in arg.items():
                    self[k] = v
                    if isinstance(v, dict):
                        self[k] = DotDict(v)
                    elif isinstance(v, str) or isinstance(v, bytes):
                        self[k] = v
                    elif isinstance(v, Iterable):
                        klass = type(v)
                        map_value: list[any] = []
                        for e in v:
                            map_e = DotDict(e) if isinstance(e, dict) else e
                            map_value.append(map_e)
                        self[k] = klass(map_value)

    def __getattr__(self, attr):
        return self.get(attr)

    def __setattr__(self, key, value):
        self.__setitem__(key, value)

    def __setitem__(self, key, value):
        super(DotDict, self).__setitem__(key, value)
        self.__dict__.update({key: value})

    def __delattr__(self, item):
        self.__delitem__(item)

    def __delitem__(self, key):
        super(DotDict, self).__delitem__(key)
        del self.__dict__[key]

    def __getstate__(self):
        return self.__dict__

    def __setstate__(self, d):
        self.__dict__.update(d)
