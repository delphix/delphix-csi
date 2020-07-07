# Table of Contents
1.  [Introduction](#introduction)
2.  [Design Proposal](#design)
3.  [Minimum Requirements](#minimum)
4.  [Recommended Mecahnism](#recommended)
5.  [Kubernetes Reference Architecture](#reference-arch)
6.  [Useful Links](#links)

## <a id="introduction"></a>Introduction

The Container Storage Interface (CSI) is a standard for exposing arbitrary block and file storage storage systems to containerized workloads on Container Orchestration Systems (COs) like Kubernetes. Using CSI third-party storage providers can write and deploy plugins exposing new storage systems in Kubernetes without ever having to touch the core Kubernetes code.

We want to build a a CSI driver to enable Delphix to act as a volume provider to COs like Kubernetes and Docker Swarm.

Now, since we're talking about Container **Orchestration** Systems, that implies some changes in the way Delphix has historically worked. Delphix would "just" be a volume provider and most of the orchestration we usually do would be handed to the CO and its objects, resources etc. Any future plugin that might evolve from this would talk to the CO API to create resources and objects that references Delphix volumes,

Here's a short story about a new volume/vFile creation:

* user interacts with `kube-apiserver`
* our controller is triggered by the creation of a `PersistentVolumeClaim` (PVC) object
* the controller talks to Delphix, asks for a volume. If everything is ok a `PersistentVolume` object is created 
* The `PersistentVolumeClaim` is updated with the volume object information
* Once the container that uses this volume is schedule on a `node`, the volume is mounted on the same `node` and the mount is propagated to the container

## <a id="design"></a>Design Proposal

A CSI driver is deployed on Kubernetes as two components: a controller and a per-node component.

The controller component can be deployed on any node of the cluster and consists of the CSI Driver that implements the CSI Controller service and a group of sidecar containers. The sidecar containers are the ones that will intereact with Kubernetes objects and APIs and then make calls to the driver's CSI Controller Service.

The CSI Driver Controller Service, in turn, we'll send requests to Delphix and relay Delphix's responses to the sidecar containers, that will update the Kubernetes API with Delphix infos.

The Kubernetes kubelet runs on every node and is responsible for making the CSI Node service calls. These calls mount and unmount the storage volume from the storage system, making it available for the Pod to use. Kubelet makes calls to the CSI driver through a UNIX domain socket shared on the host via a `HostPath` volume. There is also a second UNIX domain socket that the node-driver-registrar uses to register the CSI driver to kubelet.

### Create Volume
```
    +------+
    | User |
    +------+
      |(API Object Creation - PersistentVolumeClaim)
      |(a PersistentVolumeClaim can have a DataSource, be it another PVC or a PVC's Snapshot)
      V
+-----------------------+
| Kubernetes API        |
|                       |<-------------+
+-----------------------+              |
  |(API Object Creation Trigger)       |
  |                                    | (API Object Creation - PersistentVolume)
  V                                    | (with all the Delphix metadata we need )
+----------------------+               | (API Object Update - PersistentVolumeClaim - success, its PV etc)
| external-provisioner |---------------+
| (sidecar)            |<------+
+----------------------+       |
  | (Unix Domain Socket)       | response with vFile/volume information
  | (CreateVolume)             | 
  V                            |
+-------------------------+    |
|  CSI Controller Service |----+
| (using delphix-go-sdk)  |
+-------------------------+
  | (API call)    A
  |               |
  V               | response (OK/NOK, Volume metadata, mount instructions etc)
+---------+       |
| Delphix |-------+
+---------+
```
**DataSources**
The external-provisioner provides the (alpha) ability to request a volume be pre-populated from a data source during provisioning... and so does Delphix.

Since Kubernetes will only know about volumes provided BY Delphix TO Kubernetes, and nothing about vFiles provided by Delphix to other systems, we can assume that the we can only clone DataSources that are themselves Kubernetes volumes.

Using vFiles, VDBs and other Delphix DataSources as Kubernetes PVC DataSources is an intriguing concept, that we'll maybe implement as a Delphix Plugin at a later time.

For the record, a high level approach would be:
* Delphix create volumes from dSources/VDBs/vFiles that are outside of Kubernetes
* Delphix connects to Kubernetes APIs and creates a `PersistentVolume` (not a Claim, a Volume), passing along all the metadata necessary to link Kubernetes `PersistenVolume` and Delphix vFiles mount points
* This `PersistentVolume` objects can, in turn, be used in a `StatefulSet` Kubernetes object declaration, and the volume (Delphix vFile) will be mounted to a container

**Snapshot**
The external-provisioner also supports Snapshot `DataSources`... and so does Delphix! if a Snapshot CRD is specified as a data source on a PVC object, the sidecar container fetches the information about the snapshot by fetching `SnapshotContent` object inside Kubernetes API. It thenpopulates a `CreateVolume` to Delphix with the Snapshot information, to specify to Delphix that the new volume should be populated using the specified snapshot.

Again, there'll be nothing inside Kubernetes API about Delphix vFiles/VDBs/dSources and their snapshots that are not related to that specific Kubernetes Cluster and its objects.  See `DataSources` for a proposed (long term) approach to bring data from outside Kubernetes.

**PersistentVolumeClaim (clone)**

Cloning is also implemented by specifying a kind: of type `PersistentVolumeClaim` in the DataSource field of a Provision request. It's the responsbility of the external-provisioner to verify that the claim specified in the `DataSource` object exists, is in the same storage class as the volume being provisioned and that the claim is currently Bound.

### Delete Volume
```
    +------+
    | User |
    +------+
      |(API Object Deletion - PersistentVolumeClaim)
      |(If this volume is a DataSource for another volume, does Kubernetes relay the request, or does it do some validation?)
      V
+-----------------------+
| Kubernetes API        |
|                       |<-------------+
+-----------------------+              |
  |(API Object Deletetion Trigger)     |
  |                                    | (API Object Deletion - PersistentVolume)
  V                                    | (API Object Deletion - PersistentVolumeClaim)
+----------------------+               | 
| external-provisioner |---------------+
| (sidecar)            |<------+
+----------------------+       |
  | (Unix Domain Socket)       | OK/NOK response
  | (DeleteVolume)             | 
  V                            |
+-------------------------+    |
|  CSI Controller Service |----+
| (using delphix-go-sdk)  |
+-------------------------+
  | (API call)    A
  |               |
  V               | response (If it has dependencies - NOK - else - OK)
+---------+       |
| Delphix |-------+
+---------+
```

### Create Snapshot
```
    +------+
    | User |
    +------+
      |(API Object Creation - VolumeSnapshot)
      |(A VolumeSnapshot has to reference a snapshot from an existing volume)
      V
+-----------------------+
| Kubernetes API        |
|                       |<-------------+
+-----------------------+              |
  |(API Object Creation Trigger)       |
  |                                    | (API Object Creation - VolumeSnapshotContent)
  V                                    | (with all the Delphix metadata we need )
+----------------------+               | (API Object Update - VolumeSnapshot - success, its VolumeSnapshotContent etc)
| external-snapshotter |---------------+
| (sidecar)            |<------+
+----------------------+       |
  | (Unix Domain Socket)       | response with snapshot information
  | (CreateSnapshot)           | 
  V                            |
+-------------------------+    |
|  CSI Controller Service |----+
| (using delphix-go-sdk)  |
+-------------------------+
  | (API call)    A
  |               |
  V               | response (OK/NOK, snapshot metadata etc)
+---------+       |
| Delphix |-------+
+---------+
```

### Delphix Retention Policy - Delete Snapshot (long term implementation)
```
    +---------+
    |Retention|
    |Policy   |
    +---------+
       |
       | (Snapshot Deletion)
       V
    +---------+
    | Delphix |
    +---------+
      |(API Object Deletion - VolumeSnapshot)
      |
      V
+-----------------------+
| Kubernetes API        |
|                       |<-------------+
+-----------------------+              |
  |(API Object Deletion Trigger)       |
  |                                    | (API Object Deletion - VolumeSnapshotContent)
  V                                    | (API Object Deletion - VolumeSnapshot)
+----------------------+               | 
| external-snapshotter |---------------+
| (sidecar)            |<------+
+----------------------+       |
  | (Unix Domain Socket)       | OK
  | (DeleteSnapshot)           | 
  V                            |
+-------------------------+    |
|  CSI Controller Service |----+
| (using delphix-go-sdk)  |
+-------------------------+
  | (API call)    A
  | (Delete       |
  V  Snapshot)    | response (Probably OK, Delphix policy has already deleted it, right?)
+---------+       |
| Delphix |-------+
+---------+
```

## <a id="minimum"></a>Minimum Requirements

The only requirements are around how Kubernetes (master and node) components find and communicate with a CSI driver.

Specifically, the following is dictated by Kubernetes regarding CSI:

* Kubelet to CSI Driver Communication
    * Kubelet directly issues CSI calls (like `NodeStageVolume`, `NodePublishVolume`, etc.) to CSI drivers via a Unix Domain Socket to mount and unmount volumes.
    * Kubelet discovers CSI drivers (and the Unix Domain Socket to use to interact with a CSI driver) via the kubelet plugin registration mechanism.
    * Therefore, all CSI drivers deployed on Kubernetes MUST register themselves using the kubelet plugin registration mechanism on each supported node.
* Master to CSI Driver Communication
    * Kubernetes master components do not communicate directly (via a Unix Domain Socket or otherwise) with CSI drivers.
    * Kubernetes master components interact only with the Kubernetes API.
    * Therefore, CSI drivers that require operations that depend on the Kubernetes API (like volume create, volume attach, volume snapshot, etc.) MUST watch the Kubernetes API and trigger the appropriate CSI operations against it.

## <a id="recommended"></a>Recommended Mechanism

The Kubernetes development team has established a "Recommended Mechanism" for developing, deploying, and testing CSI Drivers on Kubernetes. It aims to reduce boilerplate code and simplify the overall process for CSI Driver developers.

This "Recommended Mechanism" makes use of the following components:

* Kubernetes CSI Sidecar Containers
* Kubernetes CSI objects
* CSI Driver Testing tools

To implement a CSI driver using this mechanism, a CSI driver developer should:

1. Create a containerized application implementing the Identity, Node, and optionally the Controller services described in the [CSI specification](https://github.com/container-storage-interface/spec/blob/master/spec.md#rpc-interface) (the CSI driver container).
    * See [Developing CSI Driver](https://kubernetes-csi.github.io/docs/developing.html) for more information.
2. Unit test it using csi-sanity.
    * See [Driver - Unit Testing](https://kubernetes-csi.github.io/docs/unit-testing.html) for more information.
3. Define Kubernetes API YAML files that deploy the CSI driver container along with appropriate sidecar containers.
    * See [Deploying in Kubernetes](https://kubernetes-csi.github.io/docs/deploying.html) for more information.
4. Deploy the driver on a Kubernetes cluster and run end-to-end functional tests on it.
    * See [Driver - Functional Testing](https://kubernetes-csi.github.io/docs/functional-testing.html)

At a minimum, CSI drivers must implement the following CSI services:

* CSI Identity service
    * Enables callers (Kubernetes components and CSI sidecar containers) to identify the driver and what optional functionality it supports.
* CSI Node service
    * Only `NodePublishVolume`, `NodeUnpublishVolume`, and `NodeGetCapabilitie`s are required.
    * Required methods enable callers to make a volume available at a specified path and discover what optional functionality the driver supports.

See more at [the official docs](https://kubernetes-csi.github.io/docs/developing.html)

Understanding the CSI Sidecars, what they do, how they interact with Kubernetes API and which services they call on the CSI Driver Endpoint... it actually helps  lot in understanding the scope of what needs to be developed. So we're going more in depth on them.

### CSI Sidecar Containers ###

The CSI Sidecar containers contain common logic to watch the Kubernetes API, trigger appropriate operations agains the Volume Driver container and update Kubernetes API as appropriate.

They should be bundled together with our CSI driver containers and be deployed together in a pod.

This should reduce boilerplate code and the amount of Kubernetes specific code for us to worry about. 

Also, separations of concerns: code that interacts with Kubernetes API and code that interacts with Delphix API are isolated and will live in different containers.

#### external-provisioner 

A sidecar container that whatches the Kubernetes API for for `PersistentVolumeClaim` objects

**TL/DR**
Since we support dynamic volume provisioning, we should use this sidecar and :
* advertise the CSI `CREATE_DELETE_VOLUME` controller capability
* implement `CREATE_VOLUME` and `DELETE_VOLUME` functions
* deployed as controller

Calls `CreateVolume` againts the CSI endpoint to provision new volume

`CreateVolume` is triggered when a new Kubernetes `PersistentVolumeClaim` object, if it references a StorageClass and the name in the provisioner field of the `StorageClass` matches the name returned by the `GetPluginInfo` call in the specified CSI endpoint.

Once the new volume is succesffully provisioned, the sidecar container creates a Kubernetes `PersistentVolume` object to represent the volume.

The deletion of a PVC object bound to a PV corresponding to the driver with a delete reclaim policy causes the sidecar container to trigger a `DeleteVolume` operation against the Delphix endpoint.

More on `DataSource` above, on the `CreateVolume` story.

#### external-attacher
A sidecar container that watches Kubernetes API server for `VolumeAttachment` objects and triggers Controller[Publish|Unpublish]Volume operations againts CSI Endpoint

**TL/DR**
To integrate with Kubernetes attach/detach hooks, we should 
* advertise CSI `PUBLISH_UNPUBLISH_VOLUME` controller cabapability
* Calls `ControllerPublishVolume` and `ControllerUnpublishVolume` (services we should implement)
* Deployed as a controller

#### external-snapshotter
Sidecar container that watches the Kubernetes API for VolumneSnapshot and VolumeSnapshotContent CRD objects.

**TL/DR**
* advertise the CSI `CREATE_DELETE_SNAPSHOT` controller capability
* Calls `CreateSnapshot` and `DeleteSnapshot` on the controller  (services we should implement)
* deployed as controller

The creation of a VolumeSnapshot object referencind a SnapshotClass CRD object corresponding to our driver causes the sidecar container to trigger a `CreateSnapshot` operation against the CSI endpoint. When a new snapshot is successfully created provisioned, the sidecar container will create a Kubernetes `VolumeSnapshotContent` object representing the snashot.



#### external-resizer

Sidecar container that watches the Kubernetes API for `PersistentVolumeClaim` object edits and triggers ControllerExpandVolume operations against a CSI endpoint if user requested more storage on `PersistentVolumeClaim` object.

**TL/DR**
* I don't think Delphix has any kind of restraint on the vFile sizes, so we'll probably not use this one at first
* DON'T advertise `VolumeExpansion` capability

#### node-driver-registrar
A sidecar container that fetches driver information (using NodeGetInfo) from a CSI endpoint and registers it with the kubelet on that node using the kubelet plugin registration mechanism.

**TL/DR**
* Calls `NodeGetInfo`, `NodeStageVolume` and `NodePublishVolume` (services we should implement)
* Deployer per node (`DaemonSet`)

Kubelet directly issues CSI NodeGetInfo, NodeStageVolume, and NodePublishVolume calls against CSI drivers. It uses the kubelet plugin registration mechanism to discover the unix domain socket to talk to the CSI driver. Therefore, all CSI drivers should use this sidecar container to register themselves with kubelet.

#### cluster-driver-registrar
A sidecar container that registers a CSI Driver with a Kubernetes cluster by creating a `CSIDriver` Object which enables the driver to customize how Kubernetes interacts with it.

**TL/DR**
* Since we'll support `ControllerPublishVolume` and probably won't need pod metadata at first, we can probably skip this
* We should probably publish the `CSIDriver` Object, although it's not needed

#### livenessprobe

This enables Kubernetes to automatically detect issues with the driver and restart the pod to try and fix the issue. We should use this to improve the availability of our driver.

Deployed as part of both, controller and node deployments.

## <a id="reference-arch"></a> Kubernetes Reference Architecture
[Kubernetes basic architecture](https://kubernetes.io/docs/concepts/overview/components/) 

![Kubernetes Architecture](./images/kubernetes-architecture.png "Kubernetes Architecture")

## <a id="links"></a>Useful Links 
* [https://kubernetes-csi.github.io/docs/introduction.html](CSI Docs)
* [https://medium.com/google-cloud/understanding-the-container-storage-interface-csi-ddbeb966a3b](Understanding the CSI)
* [https://arslan.io/2018/06/21/how-to-write-a-container-storage-interface-csi-plugin/](How toi develop a CSI plugin)
* There's a [Gluster CSI Driver](https://github.com/gluster/gluster-csi-driver) which might actually help us a lot on building our own driver