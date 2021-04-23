#!/usr/local/bin/python2
#!/usr/bin/env python2
# -*- coding: utf-8 -*-
# @author: Daniel Stolf <daniel.stolf@delphix.com> - 11/2019
# Changelog:
#
import pkgutil, base64
from dlpx.virtualization.platform import Mount, MountSpecification, Plugin, Status
from dlpx.virtualization import libs
from random import randrange

from generated.definitions import (
    RepositoryDefinition,
    SourceConfigDefinition,
    SnapshotDefinition,
)

plugin = Plugin()


#
# Below is an example of the repository discovery operation.
#
# NOTE: The decorators are defined on the 'plugin' object created above.
#
# Mark the function below as the operation that does repository discovery.
@plugin.discovery.repository()
def repository_discovery(source_connection):
    script_content = pkgutil.get_data('resources', 'getRepo.sh')
    result = libs.run_bash(source_connection, script_content)
    
    if result.stdout:
        #0 - kubectl path | 1- context name | 2- cluster name | 3- username | 4- API endpoint | 5- Server Version
        RepositoryDefinition.kubectl = result.stdout.split("|")[0]
        RepositoryDefinition.context = result.stdout.split("|")[1]
        RepositoryDefinition.cluster = result.stdout.split("|")[2]
        RepositoryDefinition.username = result.stdout.split("|")[3]
        RepositoryDefinition.apiendpoint = result.stdout.split("|")[4]
        RepositoryDefinition.version = result.stdout.split("|")[5]
        RepositoryName = "Kubernetes - {} - {}".format( RepositoryDefinition.context, RepositoryDefinition.version)
    else:
        return []
    return [RepositoryDefinition(name=RepositoryName)]

@plugin.discovery.source_config()
def source_config_discovery(source_connection, repository):
    #
    # To have automatic discovery of source configs, return a list of
    # SourceConfigDefinitions similar to the list of
    # RepositoryDefinitions above.
    #
    return [SourceConfigDefinition(name = "Empty CSI Volume (clone this for empty volumes)")]


@plugin.linked.post_snapshot()
def linked_post_snapshot(direct_source, repository, source_config):
    return SnapshotDefinition()


@plugin.virtual.configure()
def configure(virtual_source, snapshot, repository):
    mount_path = "{}/{}".format(virtual_source.parameters.mount_location,virtual_source.guid)
    environmentvars = {
        "MOUNTPATH": mount_path
    }
    script_content = pkgutil.get_data('resources', 'getExportPath.sh')
    result = libs.run_bash(virtual_source.connection,script_content,environmentvars)
    
    environmentvars = {
        "GUID": virtual_source.guid,
        "PVC_NAME": virtual_source.parameters.persistent_volume_claim,
        "PV_NAME": virtual_source.parameters.persistent_volume,
        "NAMESPACE": virtual_source.parameters.namespace,
        "CONTEXT": repository.context,
        "CLUSTER": repository.cluster,
        "APIENDPOINT": repository.apiendpoint,
        "KUBECTLBIN": repository.kubectl

    }
    source_config = SourceConfigDefinition(name=virtual_source.guid)
    source_config.export_path = result.stdout
    source_config.persistent_volume_claim = virtual_source.parameters.persistent_volume_claim
    source_config.persistent_volume = virtual_source.parameters.persistent_volume
    source_config.namespace = virtual_source.parameters.namespace

    return source_config


@plugin.virtual.reconfigure()
def reconfigure(virtual_source, repository, source_config, snapshot):
    mount_path = "{}/{}".format(virtual_source.parameters.mount_location,virtual_source.guid)
    environmentvars = {
        "MOUNTPATH": mount_path
    }

    script_content = pkgutil.get_data('resources', 'getExportPath.sh')
    result = libs.run_bash(virtual_source.connection,script_content,environmentvars)
    
    source_config = SourceConfigDefinition(name=virtual_source.guid)
    source_config.export_path = result.stdout
    source_config.persistent_volume_claim = virtual_source.parameters.persistent_volume_claim
    source_config.persistent_volume = virtual_source.parameters.persistent_volume
    source_config.namespace = virtual_source.parameters.namespace

    return source_config

@plugin.virtual.post_snapshot()
def virtual_post_snapshot(virtual_source, repository, source_config):
    return SnapshotDefinition()


@plugin.virtual.mount_specification()
def virtual_mount_specification(virtual_source, repository):
    mount_location = "{}/{}".format(virtual_source.parameters.mount_location,virtual_source.guid)
    mounts = [Mount(virtual_source.connection.environment, mount_location)]

    return MountSpecification(mounts)
