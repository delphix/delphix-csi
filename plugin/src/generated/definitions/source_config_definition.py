#
# Copyright (c) 2019 by Delphix. All rights reserved.
#
from __future__ import absolute_import
from datetime import date, datetime

from generated.definitions.base_model_ import (
    Model, GeneratedClassesError, GeneratedClassesTypeError)
from generated import util

class SourceConfigDefinition(Model):
    """NOTE: This class is auto generated by the swagger code generator program.

    Do not edit the class manually.
    """

    def __init__(self, persistent_volume='', persistent_volume_claim='', name='', export_path='', validate=True):
        """SourceConfigDefinition - a model defined in Swagger. The type of some of these
        attributes can be defined as a List[ERRORUNKNOWN]. This just means they
        are a list of any type.

            :param persistent_volume: The persistent_volume of this SourceConfigDefinition.
            :type persistent_volume: str
            :param persistent_volume_claim: The persistent_volume_claim of this SourceConfigDefinition.
            :type persistent_volume_claim: str
            :param name: The name of this SourceConfigDefinition.
            :type name: str
            :param export_path: The export_path of this SourceConfigDefinition.
            :type export_path: str
            :param validate: If the validation should be done during init. This
            should only be called internally when calling from_dict.
            :type validate: bool
        """
        self.swagger_types = {
            'persistent_volume': str,
            'persistent_volume_claim': str,
            'name': str,
            'export_path': str
        }

        self.attribute_map = {
            'persistent_volume': 'persistent_volume',
            'persistent_volume_claim': 'persistent_volume_claim',
            'name': 'name',
            'export_path': 'export_path'
        }
        
        # Validating the attribute persistent_volume and then saving it.
        type_error = GeneratedClassesTypeError.type_error(SourceConfigDefinition,
                                                          'persistent_volume',
                                                          persistent_volume,
                                                          str,
                                                          False)
        if validate and type_error:
            raise type_error
        self._persistent_volume = persistent_volume

        # Validating the attribute persistent_volume_claim and then saving it.
        type_error = GeneratedClassesTypeError.type_error(SourceConfigDefinition,
                                                          'persistent_volume_claim',
                                                          persistent_volume_claim,
                                                          str,
                                                          False)
        if validate and type_error:
            raise type_error
        self._persistent_volume_claim = persistent_volume_claim

        # Validating the attribute name and then saving it.
        if validate and name is None:
            raise GeneratedClassesError(
                "The required parameter 'name' must not be 'None'.")
        type_error = GeneratedClassesTypeError.type_error(SourceConfigDefinition,
                                                          'name',
                                                          name,
                                                          str,
                                                          True)
        if validate and type_error:
            raise type_error
        self._name = name

        # Validating the attribute export_path and then saving it.
        type_error = GeneratedClassesTypeError.type_error(SourceConfigDefinition,
                                                          'export_path',
                                                          export_path,
                                                          str,
                                                          False)
        if validate and type_error:
            raise type_error
        self._export_path = export_path
    @classmethod
    def from_dict(cls, dikt):
        """Returns the dict as a model

        :param dikt: A dict.
        :type: dict
        :return: The sourceConfigDefinition of this SourceConfigDefinition.
        :rtype: SourceConfigDefinition
        """
        return util.deserialize_model(dikt, cls)

    @property
    def persistent_volume(self):
        """Gets the persistent_volume of this SourceConfigDefinition.

        Kubernetes Persistent Volume Name

        :return: The persistent_volume of this SourceConfigDefinition.
        :rtype: str
        """
        return self._persistent_volume

    @persistent_volume.setter
    def persistent_volume(self, persistent_volume):
        """Sets the persistent_volume of this SourceConfigDefinition.

        Kubernetes Persistent Volume Name

        :param persistent_volume: The persistent_volume of this SourceConfigDefinition.
        :type persistent_volume: str
        """
        # Validating the attribute persistent_volume and then saving it.
        type_error = GeneratedClassesTypeError.type_error(SourceConfigDefinition,
                                                          'persistent_volume',
                                                          persistent_volume,
                                                          str,
                                                          False)
        if type_error:
            raise type_error
        self._persistent_volume = persistent_volume

    @property
    def persistent_volume_claim(self):
        """Gets the persistent_volume_claim of this SourceConfigDefinition.

        Kubernetes Persistent Volume Claim Name

        :return: The persistent_volume_claim of this SourceConfigDefinition.
        :rtype: str
        """
        return self._persistent_volume_claim

    @persistent_volume_claim.setter
    def persistent_volume_claim(self, persistent_volume_claim):
        """Sets the persistent_volume_claim of this SourceConfigDefinition.

        Kubernetes Persistent Volume Claim Name

        :param persistent_volume_claim: The persistent_volume_claim of this SourceConfigDefinition.
        :type persistent_volume_claim: str
        """
        # Validating the attribute persistent_volume_claim and then saving it.
        type_error = GeneratedClassesTypeError.type_error(SourceConfigDefinition,
                                                          'persistent_volume_claim',
                                                          persistent_volume_claim,
                                                          str,
                                                          False)
        if type_error:
            raise type_error
        self._persistent_volume_claim = persistent_volume_claim

    @property
    def name(self):
        """Gets the name of this SourceConfigDefinition.

        Delphix Unique Identifier for the Source

        :return: The name of this SourceConfigDefinition.
        :rtype: str
        """
        return self._name

    @name.setter
    def name(self, name):
        """Sets the name of this SourceConfigDefinition.

        Delphix Unique Identifier for the Source

        :param name: The name of this SourceConfigDefinition.
        :type name: str
        """
        # Validating the attribute name and then saving it.
        if name is None:
            raise GeneratedClassesError(
                "The required parameter 'name' must not be 'None'.")
        type_error = GeneratedClassesTypeError.type_error(SourceConfigDefinition,
                                                          'name',
                                                          name,
                                                          str,
                                                          True)
        if type_error:
            raise type_error
        self._name = name

    @property
    def export_path(self):
        """Gets the export_path of this SourceConfigDefinition.

        Export Path on Delphix Engine (DO NOT FILL THIS)

        :return: The export_path of this SourceConfigDefinition.
        :rtype: str
        """
        return self._export_path

    @export_path.setter
    def export_path(self, export_path):
        """Sets the export_path of this SourceConfigDefinition.

        Export Path on Delphix Engine (DO NOT FILL THIS)

        :param export_path: The export_path of this SourceConfigDefinition.
        :type export_path: str
        """
        # Validating the attribute export_path and then saving it.
        type_error = GeneratedClassesTypeError.type_error(SourceConfigDefinition,
                                                          'export_path',
                                                          export_path,
                                                          str,
                                                          False)
        if type_error:
            raise type_error
        self._export_path = export_path