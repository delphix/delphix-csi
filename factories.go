package delphix

// AppDataProvisionParametersFactory is just a simple function to instantiate the AppDataProvisionParametersStruct
func AppDataProvisionParametersFactory(
	Container *AppDataContainerStruct,
	Masked *bool,
	MaskingJob string,
	Source *AppDataVirtualSourceStruct,
	SourceConfig AppDataSourceConfig,
	TimeflowPointParameters TimeflowPointParameters,
) AppDataProvisionParametersStruct {
	return AppDataProvisionParametersStruct{
		Type:                    "AppDataProvisionParameters",
		Container:               Container,
		Masked:                  Masked,
		MaskingJob:              MaskingJob,
		Source:                  Source,
		SourceConfig:            SourceConfig,
		TimeflowPointParameters: TimeflowPointParameters,
	}
}

// CSIProvisionParametersFactory is just a simple function to instantiate the AppDataProvisionParametersStruct
func CSIProvisionParametersFactory(
	Container *AppDataContainerStruct,
	Masked *bool,
	MaskingJob string,
	Source *CSIAppDataVirtualSourceStruct,
	SourceConfig *CSISourceConfigStruct,
	TimeflowPointParameters TimeflowPointParameters,
) CSIAppDataProvisionParametersStruct {
	return CSIAppDataProvisionParametersStruct{
		Type:                    "AppDataProvisionParameters",
		Container:               Container,
		Masked:                  Masked,
		MaskingJob:              MaskingJob,
		Source:                  Source,
		SourceConfig:            SourceConfig,
		TimeflowPointParameters: TimeflowPointParameters,
	}
}

// AppDataRollbackParametersFactory is just a simple function to instantiate the RollbackParametersStruct
func AppDataRollbackParametersFactory(
	TimeflowPointParameters *AppDataTimeflowPointStruct,
) AppDataRollbackParametersStruct {
	return AppDataRollbackParametersStruct{
		Type:                    "RollbackParameters",
		TimeflowPointParameters: TimeflowPointParameters,
	}
}

// CreateAppDataSyncParameters is just a simple function to instantiate the AppDataSyncParametersStruct
func CreateAppDataSyncParameters(
	Resync *bool,
) AppDataSyncParametersStruct {
	return AppDataSyncParametersFactory(
		Resync,
	)
}

// CreateAppDataRollbackParameters is just a simple function to instantiate the RollbackParametersStruct
func CreateAppDataRollbackParameters(
	Timestamp string,
	Timeflow string,
) AppDataRollbackParametersStruct {
	return AppDataRollbackParametersFactory(
		&AppDataTimeflowPointStruct{
			Timestamp: Timestamp,
			Timeflow:  Timeflow,
			Type:      "TimeflowPointTimestamp",
		},
	)
}

// AppDataSyncParametersFactory is just a simple function to instantiate the AppDataSyncParametersStruct
func AppDataSyncParametersFactory(
	Resync *bool,
) AppDataSyncParametersStruct {
	return AppDataSyncParametersStruct{
		Type:   "AppDataSyncParameters",
		Resync: Resync,
	}
}
