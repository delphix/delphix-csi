package delphix

// APISessionStruct - Describes a Delphix web service session and is the result of an
// initial handshake.
// extends TypedObject
type APISessionStruct struct {
	// Client software identification token.
	// required = false
	// maxLength = 64
	Client string `json:"client,omitempty"`
	// Locale as an IETF BCP 47 language tag, defaults to 'en-US'.
	// format = locale
	// required = false
	Locale string `json:"locale,omitempty"`
	// Object type.
	// required = true
	// format = type
	Type string `json:"type,omitempty"`
	// Version of the API to use.
	// required = true
	Version *APIVersionStruct `json:"version,omitempty"`
}

// APIErrorStruct - Description of an error encountered during an API call.
// extends TypedObject
type APIErrorStruct struct {
	// Action to be taken by the user, if any, to fix the underlying
	// problem.
	Action string `json:"action,omitempty"`
	// Extra output, often from a script or other external process, that
	// may give more insight into the cause of this error.
	CommandOutput string `json:"commandOutput,omitempty"`
	// For validation errors, a map of fields to APIError objects. For
	// all other errors, a string with further details of the error.
	Details string `json:"details,omitempty"`
	// Results of diagnostic checks run, if any, if the job failed.
	Diagnoses []*DiagnosisResultStruct `json:"diagnoses,omitempty"`
	// A stable identifier for the class of error encountered.
	ID string `json:"id,omitempty"`
	// Object type.
	// required = true
	// format = type
	Type string `json:"type,omitempty"`
}

// APIVersionStruct - Describes an API version.
// extends TypedObject
type APIVersionStruct struct {
	// Major API version number.
	// minimum = 0
	// required = true
	Major *int `json:"major,omitempty"`
	// Micro API version number.
	// minimum = 0
	// required = true
	Micro *int `json:"micro,omitempty"`
	// Minor API version number.
	// minimum = 0
	// required = true
	Minor *int `json:"minor,omitempty"`
	// Object type.
	// required = true
	// format = type
	Type string `json:"type,omitempty"`
}

// LoginRequestStruct - Represents a Delphix user authentication request.
// extends TypedObject
type LoginRequestStruct struct {
	// Whether to keep session alive for all requests or only via
	// 'KeepSessionAlive' request headers. Defaults to ALL_REQUESTS if
	// omitted.
	// enum = [ALL_REQUESTS KEEP_ALIVE_HEADER_ONLY]
	// default = ALL_REQUESTS
	KeepAliveMode string `json:"keepAliveMode,omitempty"`
	// The password of the user to authenticate.
	// format = password
	// required = true
	Password string `json:"password,omitempty"`
	// The authentication domain.
	// enum = [DOMAIN SYSTEM]
	Target string `json:"target,omitempty"`
	// Object type.
	// required = true
	// format = type
	Type string `json:"type,omitempty"`
	// The username of the user to authenticate.
	// required = true
	Username string `json:"username,omitempty"`
}

// DeleteParametersStruct - The parameters to use as input to delete requests for MSSQL,
// PostgreSQL, AppData, ASE or MySQL.
// extends TypedObject
type DeleteParametersStruct struct {
	// Flag indicating whether to continue the operation upon failures.
	Force *bool `json:"force,omitempty"`
	// Object type.
	// required = true
	// format = type
	Type string `json:"type,omitempty"`
}

// DiagnosisResultStruct - Details from a diagnosis check that was run due to a failed operation.
// extends TypedObject
type DiagnosisResultStruct struct {
	// True if this was a check that did not pass.
	Failure *bool `json:"failure,omitempty"`
	// Localized message.
	Message string `json:"message,omitempty"`
	// Message code associated with the event.
	MessageCode string `json:"messageCode,omitempty"`
	// Object type.
	// required = true
	// format = type
	Type string `json:"type,omitempty"`
}

// JSONStruct - A dummy schema that is used to represent JSON.
type JSONStruct struct {
}

// AppDataSourceConfigStruct - An AppData source config with empty vFiles.
// extends AppDataSourceConfig
type AppDataSourceConfigStruct struct {
	// Object name.
	// format = name
	// required = false
	Name string `json:"name,omitempty"`
	// Object repository.
	// format = repository
	// required = false
	Repository string `json:"repository,omitempty"`
	// Object type.
	// format = type
	// required = false
	Type string `json:"type,omitempty"`
	// Environment OS user
	// format = type
	// required = false
	EnvironmentUser string `json:"environmentUser,omitempty"`
	// Additional configuration parameters
	// format = type
	// required = false
	Parameters *JSONStruct `json:"parameters,omitempty"`
	// Mount Path for vFile
	// format = type
	// required = false
	Path string `json:"path,omitempty"`
}

// CSIProvisionParametersStruct - The parameters to use as input to provision AppData.
// extends ProvisionParameters
type CSIProvisionParametersStruct struct {
	// The new container for the provisioned database.
	// required = true
	Container *AppDataContainerStruct `json:"container,omitempty"`
	// Whether or not to mark this VDB as a masked VDB. It will be marked
	// as masked if this flag or the masking job are set.
	// create = optional
	// update = readonly
	Masked *bool `json:"masked,omitempty"`
	// The Masking Job to be run when this dataset is provisioned or
	// refreshed.
	// update = readonly
	// format = objectReference
	// referenceTo = /delphix-masking-job.json
	// create = optional
	MaskingJob string `json:"maskingJob,omitempty"`
	// The source that describes an external database instance.
	// required = true
	Source *AppDataVirtualSourceStruct `json:"source,omitempty"`
	// The source config including dynamically discovered attributes of
	// the source.
	// required = true
	SourceConfig CSISourceConfigStruct `json:"sourceConfig,omitempty"`
	// The TimeFlow point, bookmark, or semantic location to base
	// provisioning on.
	// properties = map[type:map[default:TimeflowPointSemantic]]
	// required = true
	TimeflowPointParameters TimeflowPointParameters `json:"timeflowPointParameters,omitempty"`
	// Object type.
	// required = true
	// format = type
	Type string `json:"type,omitempty"`
}

// AppDataDirectSourceConfigStruct - Source config for directly linked AppData sources.
// extends AppDataSourceConfig
type AppDataDirectSourceConfigStruct struct {
	// Whether this source was discovered.
	Discovered *bool `json:"discovered,omitempty"`
	// The user used to create and manage the configuration.
	// create = optional
	// update = optional
	// format = objectReference
	// referenceTo = /delphix-source-environment-user.json
	EnvironmentUser string `json:"environmentUser,omitempty"`
	// Whether this source should be used for linking.
	// default = true
	// create = optional
	// update = optional
	LinkingEnabled *bool `json:"linkingEnabled,omitempty"`
	// The name of the config.
	// create = required
	// update = optional
	// maxLength = 256
	Name string `json:"name,omitempty"`
	// Alternate namespace for this object, for replicated and restored
	// objects.
	// format = objectReference
	// referenceTo = /delphix-namespace.json
	Namespace string `json:"namespace,omitempty"`
	// The list of parameters specified by the source config schema in
	// the toolkit. If no schema is specified, this list is empty.
	// create = optional
	// update = optional
	Parameters *JSONStruct `json:"parameters,omitempty"`
	// The path to the data to be synced.
	// create = optional
	// update = optional
	// maxLength = 1024
	Path string `json:"path,omitempty"`
	// The object reference.
	// format = objectReference
	// referenceTo = /delphix-persistent-object.json
	Reference string `json:"reference,omitempty"`
	// The object reference of the source repository.
	// update = optional
	// format = objectReference
	// referenceTo = /delphix-appdata-source-repository.json
	// create = required
	Repository string `json:"repository,omitempty"`
	// The toolkit associated with this source config.
	// format = objectReference
	// referenceTo = /delphix-toolkit.json
	Toolkit string `json:"toolkit,omitempty"`
	// Object type.
	// format = type
	// required = true
	Type string `json:"type,omitempty"`
}

// CSISourceConfigStruct - Source config for directly linked AppData sources.
// extends AppDataDirectSourceConfigStruct
type CSISourceConfigStruct struct {
	// Whether this source was discovered.
	Discovered *bool `json:"discovered,omitempty"`
	// The user used to create and manage the configuration.
	// create = optional
	// update = optional
	// format = objectReference
	// referenceTo = /delphix-source-environment-user.json
	EnvironmentUser string `json:"environmentUser,omitempty"`
	// Whether this source should be used for linking.
	// default = true
	// create = optional
	// update = optional
	LinkingEnabled *bool `json:"linkingEnabled,omitempty"`
	// The name of the config.
	// create = required
	// update = optional
	// maxLength = 256
	Name string `json:"name,omitempty"`
	// Alternate namespace for this object, for replicated and restored
	// objects.
	// format = objectReference
	// referenceTo = /delphix-namespace.json
	Namespace string `json:"namespace,omitempty"`
	// The list of parameters specified by the source config schema in
	// the toolkit. If no schema is specified, this list is empty.
	// create = optional
	// update = optional
	Parameters *CSISourceConfigParameterStruct `json:"parameters,omitempty"`
	// The path to the data to be synced.
	// create = optional
	// update = optional
	// maxLength = 1024
	Path string `json:"path,omitempty"`
	// The object reference.
	// format = objectReference
	// referenceTo = /delphix-persistent-object.json
	Reference string `json:"reference,omitempty"`
	// The object reference of the source repository.
	// update = optional
	// format = objectReference
	// referenceTo = /delphix-appdata-source-repository.json
	// create = required
	Repository string `json:"repository,omitempty"`
	// The toolkit associated with this source config.
	// format = objectReference
	// referenceTo = /delphix-toolkit.json
	Toolkit string `json:"toolkit,omitempty"`
	// Object type.
	// format = type
	// required = true
	Type string `json:"type,omitempty"`
}

// CSISourceConfigParameterStruct - CSI Driver Source Config.
// extends AppDataSourceConfig
type CSISourceConfigParameterStruct struct {
	// Name
	// required = false
	Name string `json:"name,omitempty"`
	// Path Export By Delphix for NFS Mount
	// required = false
	ExportPath string `json:"export_path,omitempty"`
	// CSI PVC Name
	// required = false
	PersistentVolumeClaim string `json:"persistent_volume_claim,omitempty"`
	// CSI PV Name
	// required = false
	PersistentVolume string `json:"persistent_volume,omitempty"`
	// K8S namespace
	// required = false
	Namespace string `json:"namespace,omitempty"`
	// Mount Path for vFile
	// required = false
	MountLocation string `json:"mount_location,omitempty"`
}

// AppDataAdditionalMountPointStruct - Specifies an additional location on which to mount a subdirectory of
// an AppData container.
// extends TypedObject
type AppDataAdditionalMountPointStruct struct {
	// Reference to the environment on which the file system will be
	// mounted.
	// update = optional
	// format = objectReference
	// referenceTo = /delphix-host-environment.json
	// create = required
	Environment string `json:"environment,omitempty"`
	// Absolute path on the target environment were the filesystem should
	// be mounted.
	// update = optional
	// format = unixpath
	// create = required
	MountPath string `json:"mountPath,omitempty"`
	// Relative path within the container of the directory that should be
	// mounted.
	// create = optional
	// update = optional
	// format = unixpath
	SharedPath string `json:"sharedPath,omitempty"`
	// Object type.
	// required = true
	// format = type
	Type string `json:"type,omitempty"`
}

// AppDataContainerStruct - Data container for AppData.
// extends DatabaseContainer
type AppDataContainerStruct struct {
	// The date this container was created.
	// format = date
	CreationTime string `json:"creationTime,omitempty"`
	// A reference to the currently active TimeFlow for this container.
	// format = objectReference
	// referenceTo = /delphix-timeflow.json
	CurrentTimeflow string `json:"currentTimeflow,omitempty"`
	// Optional user-provided description for the container.
	// maxLength = 1024
	// create = optional
	// update = optional
	Description string `json:"description,omitempty"`
	// A reference to the group containing this container.
	// format = objectReference
	// referenceTo = /delphix-group.json
	// create = required
	Group string `json:"group,omitempty"`
	// A global identifier for this container, including across Delphix
	// Engines.
	GUID string `json:"guid,omitempty"`
	// True if this container is a masked container.
	Masked *bool `json:"masked,omitempty"`
	// Object name.
	// create = required
	// update = optional
	// maxLength = 256
	Name string `json:"name,omitempty"`
	// Alternate namespace for this object, for replicated and restored
	// objects.
	// format = objectReference
	// referenceTo = /delphix-namespace.json
	Namespace string `json:"namespace,omitempty"`
	// Native operating system of the original database source system.
	Os string `json:"os,omitempty"`
	// Whether to enable high performance mode.
	// enum = [TEMPORARILY_ENABLED ENABLED DISABLED]
	// default = DISABLED
	// create = readonly
	// update = readonly
	PerformanceMode string `json:"performanceMode,omitempty"`
	// A reference to the previous TimeFlow for this container.
	// format = objectReference
	// referenceTo = /delphix-timeflow.json
	PreviousTimeflow string `json:"previousTimeflow,omitempty"`
	// Native processor type of the original database source system.
	Processor string `json:"processor,omitempty"`
	// A reference to the container this container was provisioned from.
	// format = objectReference
	// referenceTo = /delphix-container.json
	ProvisionContainer string `json:"provisionContainer,omitempty"`
	// The object reference.
	// format = objectReference
	// referenceTo = /delphix-persistent-object.json
	Reference string `json:"reference,omitempty"`
	// Runtime properties of this container.
	Runtime *AppDataContainerRuntimeStruct `json:"runtime,omitempty"`
	// Policies for managing LogSync and SnapSync across sources.
	// create = optional
	// update = optional
	SourcingPolicy *SourcingPolicyStruct `json:"sourcingPolicy,omitempty"`
	// The toolkit managing the data in the container.
	// referenceTo = /delphix-toolkit.json
	Toolkit string `json:"toolkit,omitempty"`
	// True if this container is a transformation container.
	Transformation *bool `json:"transformation,omitempty"`
	// Object type.
	// required = true
	// format = type
	Type string `json:"type,omitempty"`
}

// AppDataContainerRuntimeStruct - Runtime properties of an AppData container.
// extends DBContainerRuntime
type AppDataContainerRuntimeStruct struct {
	// True if the LogSync is enabled and running for this container.
	LogSyncActive *bool `json:"logSyncActive,omitempty"`
	// The pre-provisioning runtime for the container.
	PreProvisioningStatus *PreProvisioningRuntimeStruct `json:"preProvisioningStatus,omitempty"`
	// Object type.
	// required = true
	// format = type
	Type string `json:"type,omitempty"`
}

// AppDataProvisionParametersStruct - The parameters to use as input to provision AppData.
// extends ProvisionParameters
type AppDataProvisionParametersStruct struct {
	// The new container for the provisioned database.
	// required = true
	Container *AppDataContainerStruct `json:"container,omitempty"`
	// Whether or not to mark this VDB as a masked VDB. It will be marked
	// as masked if this flag or the masking job are set.
	// create = optional
	// update = readonly
	Masked *bool `json:"masked,omitempty"`
	// The Masking Job to be run when this dataset is provisioned or
	// refreshed.
	// update = readonly
	// format = objectReference
	// referenceTo = /delphix-masking-job.json
	// create = optional
	MaskingJob string `json:"maskingJob,omitempty"`
	// The source that describes an external database instance.
	// required = true
	Source *AppDataVirtualSourceStruct `json:"source,omitempty"`
	// The source config including dynamically discovered attributes of
	// the source.
	// required = true
	SourceConfig AppDataSourceConfig `json:"sourceConfig,omitempty"`
	// The TimeFlow point, bookmark, or semantic location to base
	// provisioning on.
	// properties = map[type:map[default:TimeflowPointSemantic]]
	// required = true
	TimeflowPointParameters TimeflowPointParameters `json:"timeflowPointParameters,omitempty"`
	// Object type.
	// required = true
	// format = type
	Type string `json:"type,omitempty"`
}

// CSIAppDataProvisionParametersStruct - The parameters to use as input to provision AppData.
// extends ProvisionParameters
type CSIAppDataProvisionParametersStruct struct {
	// The new container for the provisioned database.
	// required = true
	Container *AppDataContainerStruct `json:"container,omitempty"`
	// Whether or not to mark this VDB as a masked VDB. It will be marked
	// as masked if this flag or the masking job are set.
	// create = optional
	// update = readonly
	Masked *bool `json:"masked,omitempty"`
	// The Masking Job to be run when this dataset is provisioned or
	// refreshed.
	// update = readonly
	// format = objectReference
	// referenceTo = /delphix-masking-job.json
	// create = optional
	MaskingJob string `json:"maskingJob,omitempty"`
	// The source that describes an external database instance.
	// required = true
	Source *CSIAppDataVirtualSourceStruct `json:"source,omitempty"`
	// The source config including dynamically discovered attributes of
	// the source.
	// required = true
	SourceConfig AppDataSourceConfig `json:"sourceConfig,omitempty"`
	// The TimeFlow point, bookmark, or semantic location to base
	// provisioning on.
	// properties = map[type:map[default:TimeflowPointSemantic]]
	// required = true
	TimeflowPointParameters TimeflowPointParameters `json:"timeflowPointParameters,omitempty"`
	// Object type.
	// required = true
	// format = type
	Type string `json:"type,omitempty"`
}

// AppDataSnapshotStruct - Snapshot of an AppData TimeFlow.
// extends TimeflowSnapshot
type AppDataSnapshotStruct struct {
	// A value in the set {CONSISTENT, INCONSISTENT, CRASH_CONSISTENT}
	// indicating what type of recovery strategies must be invoked when
	// provisioning from this snapshot.
	Consistency string `json:"consistency,omitempty"`
	// Reference to the database of which this TimeFlow is a part.
	// format = objectReference
	// referenceTo = /delphix-container.json
	Container string `json:"container,omitempty"`
	// Point in time at which this snapshot was created. This may be
	// different from the time corresponding to the TimeFlow.
	// format = date
	CreationTime string `json:"creationTime,omitempty"`
	// The location within the parent TimeFlow at which this snapshot was
	// initiated.
	FirstChangePoint *AppDataTimeflowPointStruct `json:"firstChangePoint,omitempty"`
	// The location of the snapshot within the parent TimeFlow
	// represented by this snapshot.
	LatestChangePoint *AppDataTimeflowPointStruct `json:"latestChangePoint,omitempty"`
	// The JSON payload conforming to the DraftV4 schema based on the
	// type of application data being manipulated.
	Metadata *JSONStruct `json:"metadata,omitempty"`
	// Boolean value indicating if a virtual database provisioned from
	// this snapshot will be missing nologging changes.
	MissingNonLoggedData *bool `json:"missingNonLoggedData,omitempty"`
	// Object name.
	// create = readonly
	// update = readonly
	Name string `json:"name,omitempty"`
	// Alternate namespace for this object, for replicated and restored
	// objects.
	// referenceTo = /delphix-namespace.json
	// format = objectReference
	Namespace string `json:"namespace,omitempty"`
	// The object reference.
	// format = objectReference
	// referenceTo = /delphix-persistent-object.json
	Reference string `json:"reference,omitempty"`
	// Retention policy, in days. A value of -1 indicates the snapshot
	// should be kept forever.
	// update = optional
	Retention *int `json:"retention,omitempty"`
	// Runtime properties of the snapshot.
	Runtime *AppDataSnapshotRuntimeStruct `json:"runtime,omitempty"`
	// Boolean value indicating that this snapshot is in a transient
	// state and should not be user visible.
	Temporary *bool `json:"temporary,omitempty"`
	// TimeFlow of which this snapshot is a part.
	// referenceTo = /delphix-timeflow.json
	// format = objectReference
	Timeflow string `json:"timeflow,omitempty"`
	// Time zone of the source database at the time the snapshot was
	// taken.
	Timezone string `json:"timezone,omitempty"`
	// The toolkit associated with this snapshot.
	// format = objectReference
	// referenceTo = /delphix-toolkit.json
	Toolkit string `json:"toolkit,omitempty"`
	// Object type.
	// required = true
	// format = type
	Type string `json:"type,omitempty"`
	// Version of database source repository at the time the snapshot was
	// taken.
	Version string `json:"version,omitempty"`
}

// AppDataRollbackParametersStruct - The parameters to use as input to rollback requests for AppData
type AppDataRollbackParametersStruct struct {
	// The TimeFlow point, bookmark, or semantic location to roll the
	// database back to.
	// required = true
	// properties = map[type:map[default:TimeflowPointSemantic]]
	TimeflowPointParameters *AppDataTimeflowPointStruct `json:"timeflowPointParameters,omitempty"`
	// Object type.
	// format = type
	// required = true
	Type string `json:"type,omitempty"`
}

// AppDataSnapshotRuntimeStruct - Runtime (non-persistent) properties of AppData TimeFlow snapshots.
type AppDataSnapshotRuntimeStruct struct {
	// True if this snapshot can be used as the basis for provisioning a
	// new TimeFlow.
	Provisionable *bool `json:"provisionable,omitempty"`
	// Object type.
	// required = true
	// format = type
	Type string `json:"type,omitempty"`
}

// AppDataSourceRuntimeStruct - Runtime (non-persistent) properties of an AppData source.
// extends SourceRuntime
type AppDataSourceRuntimeStruct struct {
	// True if the source is JDBC accessible. If false then no properties
	// can be retrieved.
	Accessible *bool `json:"accessible,omitempty"`
	// The time that the 'accessible' propery was last checked.
	// format = date
	AccessibleTimestamp string `json:"accessibleTimestamp,omitempty"`
	// Size of the database in bytes.
	// base = 1024
	// units = B
	DatabaseSize float64 `json:"databaseSize,omitempty"`
	// Status indicating whether the source is enabled. A source has a
	// 'PARTIAL' status if its sub-sources are not all enabled.
	// enum = [ENABLED PARTIAL DISABLED]
	Enabled string `json:"enabled,omitempty"`
	// The reason why the source is not JDBC accessible.
	NotAccessibleReason string `json:"notAccessibleReason,omitempty"`
	// Status of the source. 'Unknown' if all attempts to connect to the
	// source failed.
	// enum = [RUNNING INACTIVE PENDING CANCELED FAILED CHECKING UNKNOWN]
	Status string `json:"status,omitempty"`
	// Object type.
	// required = true
	// format = type
	Type string `json:"type,omitempty"`
}

// AppDataSyncParametersStruct - The parameters to use as input to sync an AppData source.
// extends SyncParameters
type AppDataSyncParametersStruct struct {
	// Whether or not to force a non-incremental load of data prior to
	// taking a snapshot.
	// default = false
	Resync *bool `json:"resync,omitempty"`
	// Object type.
	// required = true
	// format = type
	Type string `json:"type,omitempty"`
}

// AppDataTimeflowPointParametersStruct - An AppData source config with empty vFiles.
// extends AppDataSourceConfig
type AppDataTimeflowPointParametersStruct struct {
	// Container Reference.
	// format = repository
	// required = false
	Container string `json:"container,omitempty"`
	// Object repository.
	// format = repository
	// required = false
	Repository string `json:"repository,omitempty"`
	// A snapshot location.
	// format = repository
	// required = false
	Location string `json:"location,omitempty"`
	// Timestamp of the snapshot.
	// format = repository
	// required = false
	Timestamp string `json:"timestamp,omitempty"`
	// Object type.
	// format = type
	// required = false
	Type string `json:"type,omitempty"`
}

// AppDataTimeflowPointStruct - A unique point within an AppData TimeFlow.
// extends TimeflowPoint
type AppDataTimeflowPointStruct struct {
	// The TimeFlow location.
	// create = optional
	Location string `json:"location,omitempty"`
	// Reference to TimeFlow containing this point.
	// required = true
	// format = objectReference
	// referenceTo = /delphix-timeflow.json
	Timeflow string `json:"timeflow,omitempty"`
	// The logical time corresponding to the TimeFlow location.
	// format = date
	// create = optional
	Timestamp string `json:"timestamp,omitempty"`
	// Object type.
	// required = true
	// format = type
	Type string `json:"type,omitempty"`
}

// AppDataVirtualSourceStruct - A virtual AppData source.
// extends AppDataManagedSource
type AppDataVirtualSourceStruct struct {
	// Locations to mount subdirectories of the AppData in addition to
	// the normal target mount point. These paths will be mounted and
	// unmounted as part of enabling and disabling this source.
	// create = optional
	// update = optional
	AdditionalMountPoints []*AppDataAdditionalMountPointStruct `json:"additionalMountPoints,omitempty"`
	// Indicates whether Delphix should automatically restart this
	// virtual source when target host reboot is detected.
	// update = optional
	// create = required
	AllowAutoVDBRestartOnHostReboot *bool `json:"allowAutoVDBRestartOnHostReboot,omitempty"`
	// Reference to the configuration for the source.
	// format = objectReference
	// referenceTo = /delphix-source-config.json
	// create = optional
	Config string `json:"config,omitempty"`
	// Reference to the container being fed by this source, if any.
	// format = objectReference
	// referenceTo = /delphix-container.json
	Container string `json:"container,omitempty"`
	// A user-provided description of the source.
	Description string `json:"description,omitempty"`
	// Hosts that might affect operations on this source. Property will
	// be null unless the includeHosts parameter is set when listing
	// sources.
	Hosts []string `json:"hosts,omitempty"`
	// Flag indicating whether the source is a linked source in the
	// Delphix system.
	Linked *bool `json:"linked,omitempty"`
	// Flag indicating whether it is allowed to collect logs, potentially
	// containing sensitive information, related to this source.
	// default = false
	// create = optional
	// update = optional
	LogCollectionEnabled *bool `json:"logCollectionEnabled,omitempty"`
	// Object name.
	// format = objectName
	// create = optional
	// update = optional
	// maxLength = 256
	Name string `json:"name,omitempty"`
	// Alternate namespace for this object, for replicated and restored
	// objects.
	// format = objectReference
	// referenceTo = /delphix-namespace.json
	Namespace string `json:"namespace,omitempty"`
	// User-specified operation hooks for this source.
	// update = optional
	// create = optional
	Operations *VirtualSourceOperationsStruct `json:"operations,omitempty"`
	// The JSON payload conforming to the DraftV4 schema based on the
	// type of application data being manipulated.
	// create = required
	// update = optional
	Parameters *JSONStruct `json:"parameters,omitempty"`
	// The object reference.
	// format = objectReference
	// referenceTo = /delphix-persistent-object.json
	Reference string `json:"reference,omitempty"`
	// Runtime properties of this source.
	Runtime *AppDataSourceRuntimeStruct `json:"runtime,omitempty"`
	// Flag indicating whether the source is used as a staging source for
	// pre-provisioning. Staging sources are managed by the Delphix
	// system.
	Staging *bool `json:"staging,omitempty"`
	// Status of this source.
	// enum = [DEFAULT PENDING_UPGRADE]
	Status string `json:"status,omitempty"`
	// The toolkit associated with this source.
	// format = objectReference
	// referenceTo = /delphix-toolkit.json
	Toolkit string `json:"toolkit,omitempty"`
	// Object type.
	// required = true
	// format = type
	Type string `json:"type,omitempty"`
	// Flag indicating whether the source is a virtual source in the
	// Delphix system.
	Virtual *bool `json:"virtual,omitempty"`
}

// CSIAppDataVirtualSourceStruct - A virtual AppData source.
// extends AppDataManagedSource
type CSIAppDataVirtualSourceStruct struct {
	// Locations to mount subdirectories of the AppData in addition to
	// the normal target mount point. These paths will be mounted and
	// unmounted as part of enabling and disabling this source.
	// create = optional
	// update = optional
	AdditionalMountPoints []*AppDataAdditionalMountPointStruct `json:"additionalMountPoints,omitempty"`
	// Indicates whether Delphix should automatically restart this
	// virtual source when target host reboot is detected.
	// update = optional
	// create = required
	AllowAutoVDBRestartOnHostReboot *bool `json:"allowAutoVDBRestartOnHostReboot,omitempty"`
	// Reference to the configuration for the source.
	// format = objectReference
	// referenceTo = /delphix-source-config.json
	// create = optional
	Config string `json:"config,omitempty"`
	// Reference to the container being fed by this source, if any.
	// format = objectReference
	// referenceTo = /delphix-container.json
	Container string `json:"container,omitempty"`
	// A user-provided description of the source.
	Description string `json:"description,omitempty"`
	// Hosts that might affect operations on this source. Property will
	// be null unless the includeHosts parameter is set when listing
	// sources.
	Hosts []string `json:"hosts,omitempty"`
	// Flag indicating whether the source is a linked source in the
	// Delphix system.
	Linked *bool `json:"linked,omitempty"`
	// Flag indicating whether it is allowed to collect logs, potentially
	// containing sensitive information, related to this source.
	// default = false
	// create = optional
	// update = optional
	LogCollectionEnabled *bool `json:"logCollectionEnabled,omitempty"`
	// Object name.
	// format = objectName
	// create = optional
	// update = optional
	// maxLength = 256
	Name string `json:"name,omitempty"`
	// Alternate namespace for this object, for replicated and restored
	// objects.
	// format = objectReference
	// referenceTo = /delphix-namespace.json
	Namespace string `json:"namespace,omitempty"`
	// User-specified operation hooks for this source.
	// update = optional
	// create = optional
	Operations *VirtualSourceOperationsStruct `json:"operations,omitempty"`
	// The JSON payload conforming to the DraftV4 schema based on the
	// type of application data being manipulated.
	// create = required
	// update = optional
	Parameters *CSISourceConfigParameterStruct `json:"parameters,omitempty"`
	// The object reference.
	// format = objectReference
	// referenceTo = /delphix-persistent-object.json
	Reference string `json:"reference,omitempty"`
	// Runtime properties of this source.
	Runtime *AppDataSourceRuntimeStruct `json:"runtime,omitempty"`
	// Flag indicating whether the source is used as a staging source for
	// pre-provisioning. Staging sources are managed by the Delphix
	// system.
	Staging *bool `json:"staging,omitempty"`
	// Status of this source.
	// enum = [DEFAULT PENDING_UPGRADE]
	Status string `json:"status,omitempty"`
	// The toolkit associated with this source.
	// format = objectReference
	// referenceTo = /delphix-toolkit.json
	Toolkit string `json:"toolkit,omitempty"`
	// Object type.
	// required = true
	// format = type
	Type string `json:"type,omitempty"`
	// Flag indicating whether the source is a virtual source in the
	// Delphix system.
	Virtual *bool `json:"virtual,omitempty"`
}

// PreProvisioningRuntimeStruct - Runtime properties for pre-provisioning of a MSSQL database container.
// extends TypedObject
type PreProvisioningRuntimeStruct struct {
	// Timestamp of the last update to the status.
	LastUpdateTimestamp string `json:"lastUpdateTimestamp,omitempty"`
	// User action required to resolve any error that the
	// pre-provisioning run encountered.
	PendingAction string `json:"pendingAction,omitempty"`
	// Indicates the current state of pre-provisioning for the database.
	// enum = [ACTIVE INACTIVE FAULTED UNKNOWN]
	PreProvisioningState string `json:"preProvisioningState,omitempty"`
	// Response taken based on the status of the pre-provisioning run.
	Response string `json:"response,omitempty"`
	// The status of the pre-provisioning run.
	Status string `json:"status,omitempty"`
	// Object type.
	// format = type
	// required = true
	Type string `json:"type,omitempty"`
}

// SourcingPolicyStruct - Database policies for managing SnapSync and LogSync across sources for
// a MSSQL container.
// extends TypedObject
type SourcingPolicyStruct struct {
	// True if LogSync should run for this database.
	// create = optional
	// update = optional
	// default = false
	LogsyncEnabled *bool `json:"logsyncEnabled,omitempty"`
	// Object type.
	// required = true
	// format = type
	Type string `json:"type,omitempty"`
}

// VirtualSourceOperationsStruct - Describes operations which are performed on virtual sources at various
// times.
// extends TypedObject
type VirtualSourceOperationsStruct struct {
	// Operations to perform when initially creating the virtual source
	// and every time it is refreshed.
	// create = optional
	// update = optional
	ConfigureClone []SourceOperation `json:"configureClone,omitempty"`
	// Operations to perform after refreshing a virtual source. These
	// operations can be used to restore any data or configuration backed
	// up in the preRefresh operations.
	// create = optional
	// update = optional
	PostRefresh []SourceOperation `json:"postRefresh,omitempty"`
	// Operations to perform after rewinding a virtual source. These
	// operations can be used to automate processes once the rewind is
	// complete.
	// create = optional
	// update = optional
	PostRollback []SourceOperation `json:"postRollback,omitempty"`
	// Operations to perform after snapshotting a virtual source.
	// create = optional
	// update = optional
	PostSnapshot []SourceOperation `json:"postSnapshot,omitempty"`
	// Operations to perform after starting a virtual source.
	// create = optional
	// update = optional
	PostStart []SourceOperation `json:"postStart,omitempty"`
	// Operations to perform after stopping a virtual source.
	// update = optional
	// create = optional
	PostStop []SourceOperation `json:"postStop,omitempty"`
	// Operations to perform before refreshing a virtual source. These
	// operations can backup any data or configuration from the running
	// source before doing the refresh.
	// create = optional
	// update = optional
	PreRefresh []SourceOperation `json:"preRefresh,omitempty"`
	// Operations to perform before rewinding a virtual source. These
	// operations can backup any data or configuration from the running
	// source prior to rewinding.
	// create = optional
	// update = optional
	PreRollback []SourceOperation `json:"preRollback,omitempty"`
	// Operations to perform before snapshotting a virtual source. These
	// operations can quiesce any data prior to snapshotting.
	// create = optional
	// update = optional
	PreSnapshot []SourceOperation `json:"preSnapshot,omitempty"`
	// Operations to perform before starting a virtual source.
	// create = optional
	// update = optional
	PreStart []SourceOperation `json:"preStart,omitempty"`
	// Operations to perform before stopping a virtual source.
	// update = optional
	// create = optional
	PreStop []SourceOperation `json:"preStop,omitempty"`
	// Object type.
	// required = true
	// format = type
	Type string `json:"type,omitempty"`
}
