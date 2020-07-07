#
# Copyright (c) 2019 by Delphix. All rights reserved.
#

import pprint

import six

from generated import util


class Model(object):
    # swaggerTypes: The key is attribute name and the
    # value is attribute type.
    swagger_types = {}

    # attributeMap: The key is attribute name and the
    # value is json key in definition.
    attribute_map = {}

    @classmethod
    def from_dict(cls, dikt):
        """Returns the dict as a model"""
        return util.deserialize_model(dikt, cls)

    def to_dict(self):
        """Returns the model properties as a dict

        :rtype: dict
        """
        result = {}

        for attr, _ in six.iteritems(self.swagger_types):
            value = getattr(self, attr)
            attr = self.attribute_map[attr]
            if isinstance(value, list):
                result[attr] = list(map(
                    lambda x: x.to_dict() if hasattr(x, "to_dict") else x,
                    value
                ))
            elif hasattr(value, "to_dict"):
                result[attr] = value.to_dict()
            elif isinstance(value, dict):
                result[attr] = dict(map(
                    lambda item: (item[0], item[1].to_dict())
                    if hasattr(item[1], "to_dict") else item,
                    value.items()
                ))
            else:
                result[attr] = value

        return result

    def to_str(self):
        """Returns the string representation of the model

        :rtype: str
        """
        return pprint.pformat(self.to_dict())

    def __repr__(self):
        """For `print` and `pprint`"""
        return self.to_str()

    def __eq__(self, other):
        """Returns true if both objects are equal"""
        return self.to_dict() == other.to_dict()

    def __ne__(self, other):
        """Returns true if both objects are not equal"""
        return not self == other


class GeneratedClassesError(Exception):
    """Generic Plugin exception with generated classes from schemas defined
    by the Plugin writer.

    This exception will be thrown whenever a a generic generated classe error
    gets thrown.

    Args:
    message (str): A user-readable message describing the exception.

    Attributes:
    message (str): A user-readable message describing the exception.
    """
    @property
    def message(self):
        return self.args[0]

    def __init__(self, message):
        super(GeneratedClassesError, self).__init__(message)


class GeneratedClassesTypeError(GeneratedClassesError):
    """Plugin exception

    Some Plugin specific errors (type errors, etc.) need to be fixed via the
    plugin code. Potentially actionable by plugin code.

    This exception will be thrown whenever the plugin writer tries to create
    a generated object with the wrong type.

    Args:
        message (str): A user-readable message describing the exception.

    Attributes:
        message (str): A user-readable message describing the exception.
    """

    def __init__(
        self,
        object_type,
        parameter_name,
        actual_type,
        expected_type,
        required):
        actual, expected = self.get_actual_and_expected_type(
            actual_type, expected_type)

        message = ("{}'s parameter '{}' was {} but should be of {}{}.".format(
            object_type.__name__,
            parameter_name,
            actual,
            expected,
            (' if defined', '')[required]))
        super(GeneratedClassesTypeError, self).__init__(message)

    def get_actual_and_expected_type(self, actual_type, expected_type):
        """ Takes in the the actual and expected types and generates a tuple of
        two strings that are then used to generate the output message.

        Args:
            actual_type (Type, List[Type], Set[Type],
                         or Set[Tuple[Type, Type]]):
            type(s) that was actually passed in for the parameter. This will
            either take the type and make it a str or join the types as a
            string and put it in brackets.
            expected_type (Type or List[Type], Set[Type], Dict[Type, Type]):
            The type of the parameter that was expected. Or if this is a
            container then we assume there is one element in it and that type
            is the expected type of the container. (For dicts this is the key)
            ie: if expected_type = {str} then the returned expected string with
            be something like "type 'dict with key basestring'"

        Returns:
            tuple (str, str): the actual and expected strings used for the
            types.
        """

        def _remove_angle_brackets(type_string):
            return type_string.replace('<', '').replace('>', '')

        if isinstance(expected_type, list):
            if len(expected_type) != 1:
                raise ValueError('The thrown GeneratedClassesTypeError should'
                                 ' have had a list of size 1 as the'
                                 ' expected_type')
            single_type = expected_type[0]
            if single_type.__module__ != '__builtin__':
                type_name = '{}.{}'.format(
                    single_type.__module__, single_type.__name__)
            else:
                type_name = single_type.__name__
            expected = "type 'list of {}'".format(type_name)
        elif isinstance(expected_type, set):
            if len(expected_type) != 1:
                raise ValueError('The thrown GeneratedClassesTypeError should'
                                 ' have had a set of size 1 as the'
                                 ' expected_type')
            single_type = expected_type.pop()
            if single_type.__module__ != '__builtin__':
                type_name = '{}.{}'.format(
                single_type.__module__, single_type.__name__)
            else:
                type_name = single_type.__name__
                expected = "a dict with keys type '{}'".format(type_name)
        elif isinstance(expected_type, dict):
            if len(expected_type) != 1:
                raise ValueError('The thrown GeneratedClassesTypeError should'
                                 ' have had a dict of size 1 as the'
                                 ' expected_type')
            key_type = expected_type.keys()[0]
            value_type = expected_type.values()[0]
            if key_type.__module__ != '__builtin__':
                key_type_name = '{}.{}'.format(
                    key_type.__module__, key_type.__name__)
            else:
                key_type_name = key_type.__name__
            if value_type.__module__ != '__builtin__':
                value_type_name = '{}.{}'.format(
                    value_type.__module__, value_type.__name__)
            else:
                value_type_name = value_type.__name__
                expected = "type 'dict of {}:{}'".format(
                    key_type_name, value_type_name)
        else:
            expected = _remove_angle_brackets(str(expected_type))

        if isinstance(actual_type, list):
            actual = 'a list of [{}]'.format(
                ', '.join(_remove_angle_brackets(str(single_type))
                          for single_type in actual_type))
        elif isinstance(actual_type, set):
            #
            # If it's a set, check that it is either a set of tuples or set of
            # types. In the case of tuples, we couldn't just pass in a dict
            # because keys have to be unique and we don't need that for the
            # actual types.
            #
            if (not all(isinstance(type_tuple, tuple)
                        for type_tuple in actual_type)):
                actual = 'a dict with keys of {}{}{}'.format(
                    '{',
                    ', '.join(_remove_angle_brackets(str(single_type))
                              for single_type in actual_type),
                    '}')
            else:
                actual = 'a dict of {}{}{}'.format(
                    '{',
                    ', '.join(['{0}:{1}'.format(
                        _remove_angle_brackets(str(k)),
                        _remove_angle_brackets(str(v))) for k, v in actual_type]),
                    '}')

        else:
            actual = _remove_angle_brackets(str(actual_type))

        return actual, expected

    @staticmethod
    def type_error(object_type, parameter_name, parameter, expected_type, required, element_type=None):
        """Checks the parameter to see if it is the expected type. Depending on
        what swagger returns sometimes the type we want to check is not correct.
        If the type is incorrect then return the error that we want to raise.

        :param object_type: The object type that is currently being created.
        :param parameter_name: the name of the parameter passed into the model.
        :param parameter: The parameter passed into the model.
        :param expected_type: the expected datatype from swagger.
        :param required: Whether the parameter was required when creating the model.
        :param element_type: If this is a dict or list then this tells us what
            type it's value should be.
        :return: GeneratedClassesTypeError
        """
        # First just return None if the parameter was not required and is None
        if not required and parameter is None:
            return None
        # Now check if the types are incorrect.
        if expected_type == float:
            if not isinstance(parameter, (float, int, long, complex)):
                return GeneratedClassesTypeError(object_type,
                                                 parameter_name,
                                                 type(parameter),
                                                 float,
                                                 required)
        elif expected_type == str:
            if not isinstance(parameter, basestring):
                return GeneratedClassesTypeError(object_type,
                                                 parameter_name,
                                                 type(parameter),
                                                 basestring,
                                                 required)
        elif expected_type == list:
            if element_type:
                if not isinstance(parameter, list):
                    return GeneratedClassesTypeError(object_type,
                                                     parameter_name,
                                                     type(parameter),
                                                     [element_type],
                                                     required)

                if element_type == float:
                    check = all(isinstance(elem, (float, int, long, complex))
                                for elem in parameter)
                else:
                    check = all(isinstance(elem, element_type)
                                for elem in parameter)
                if not check:
                    return GeneratedClassesTypeError(
                        object_type,
                        parameter_name,
                        [type(elem) for elem in parameter],
                        [element_type],
                        required)
            else:
                if not isinstance(parameter, list):
                    return GeneratedClassesTypeError(object_type,
                                                     parameter_name,
                                                     type(parameter),
                                                     list,
                                                     required)
        elif expected_type == object or expected_type == dict:
            if not isinstance(parameter, dict):
                return GeneratedClassesTypeError(object_type,
                                                 parameter_name,
                                                 type(parameter),
                                                 {basestring},
                                                 required)

            #
            # If the value (element_type) is provided we want to check both
            # the key and value. If it isn't then just check the key. If the
            # element type is float we want to allow int, float, long, and
            # complex.
            #
            if element_type:
                if element_type == float:
                    value_check = all(isinstance(v,
                                                 (float, int, long, complex))
                                      for v in parameter.values())
                else:
                    value_check = all(isinstance(v, element_type)
                                      for v in parameter.values())
                if (not all(isinstance(k, basestring)
                            for k in parameter.keys()) or not value_check):
                    return GeneratedClassesTypeError(
                        object_type,
                        parameter_name,
                        {(type(k), type(v)) for k, v in parameter.items()},
                        {basestring: element_type},
                        required)
            else:
                if not all(isinstance(k, basestring) for k in parameter.keys()):
                    return GeneratedClassesTypeError(
                        object_type,
                        parameter_name,
                        {type(k) for k in parameter.keys()},
                        {basestring},
                        required)

        else:
            if not isinstance(parameter, expected_type):
                return GeneratedClassesTypeError(object_type,
                                                 parameter_name,
                                                 type(parameter),
                                                 expected_type,
                                                 required)
        return None
