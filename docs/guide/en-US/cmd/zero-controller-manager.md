## zero-controller-manager



### Synopsis

The zero controller manager is a daemon that embeds
the core control loops. In applications of robotics and
automation, a control loop is a non-terminating loop that regulates the state of
the system. In Zero , a controller is a control loop that watches the shared
state of the miner through the zero-apiserver and makes changes attempting to move the
current state towards the desired state.

```
zero-controller-manager [flags]
```

### Options

```
      --bind-address ip                           The IP address for the proxy server to serve on (set to '0.0.0.0' for all IPv4 interfaces and '::' for all   IPv6 interfaces). This parameter is ignored if a config file is specified by --config. (default 0.0.0.0)
      --concurrent-gc-syncs int32                 The number of garbage collector workers that are allowed to sync concurrently.This parameter is ignored if a config file is specified by --config. (default 20)
      --config string                             The path to the configuration file.
      --enable-garbage-collector                  Enables the generic garbage collector. MUST be synced with the corresponding flag of the kube-apiserver. This parameter is ignored if a config file is specified by --config. (default true)
      --feature-gates mapStringBool               A set of key=value pairs that describe feature gates for alpha/experimental features. Options are:
                                                  APIListChunking=true|false (BETA - default=true)
                                                  APIPriorityAndFairness=true|false (BETA - default=true)
                                                  APIResponseCompression=true|false (BETA - default=true)
                                                  APISelfSubjectReview=true|false (ALPHA - default=false)
                                                  APIServerIdentity=true|false (BETA - default=true)
                                                  APIServerTracing=true|false (ALPHA - default=false)
                                                  AggregatedDiscoveryEndpoint=true|false (ALPHA - default=false)
                                                  AllAlpha=true|false (ALPHA - default=false)
                                                  AllBeta=true|false (BETA - default=false)
                                                  AnyVolumeDataSource=true|false (BETA - default=true)
                                                  AppArmor=true|false (BETA - default=true)
                                                  CPUManagerPolicyAlphaOptions=true|false (ALPHA - default=false)
                                                  CPUManagerPolicyBetaOptions=true|false (BETA - default=true)
                                                  CPUManagerPolicyOptions=true|false (BETA - default=true)
                                                  CSIMigrationPortworx=true|false (BETA - default=false)
                                                  CSIMigrationRBD=true|false (ALPHA - default=false)
                                                  CSINodeExpandSecret=true|false (ALPHA - default=false)
                                                  CSIVolumeHealth=true|false (ALPHA - default=false)
                                                  ComponentSLIs=true|false (ALPHA - default=false)
                                                  ContainerCheckpoint=true|false (ALPHA - default=false)
                                                  ContextualLogging=true|false (ALPHA - default=false)
                                                  CronJobTimeZone=true|false (BETA - default=true)
                                                  CrossNamespaceVolumeDataSource=true|false (ALPHA - default=false)
                                                  CustomCPUCFSQuotaPeriod=true|false (ALPHA - default=false)
                                                  CustomResourceValidationExpressions=true|false (BETA - default=true)
                                                  DisableCloudProviders=true|false (ALPHA - default=false)
                                                  DisableKubeletCloudCredentialProviders=true|false (ALPHA - default=false)
                                                  DownwardAPIHugePages=true|false (BETA - default=true)
                                                  DynamicResourceAllocation=true|false (ALPHA - default=false)
                                                  EventedPLEG=true|false (ALPHA - default=false)
                                                  ExpandedDNSConfig=true|false (BETA - default=true)
                                                  ExperimentalHostUserNamespaceDefaulting=true|false (BETA - default=false)
                                                  GRPCContainerProbe=true|false (BETA - default=true)
                                                  GracefulNodeShutdown=true|false (BETA - default=true)
                                                  GracefulNodeShutdownBasedOnPodPriority=true|false (BETA - default=true)
                                                  HPAContainerMetrics=true|false (ALPHA - default=false)
                                                  HPAScaleToZero=true|false (ALPHA - default=false)
                                                  HonorPVReclaimPolicy=true|false (ALPHA - default=false)
                                                  IPTablesOwnershipCleanup=true|false (ALPHA - default=false)
                                                  InTreePluginAWSUnregister=true|false (ALPHA - default=false)
                                                  InTreePluginAzureDiskUnregister=true|false (ALPHA - default=false)
                                                  InTreePluginAzureFileUnregister=true|false (ALPHA - default=false)
                                                  InTreePluginGCEUnregister=true|false (ALPHA - default=false)
                                                  InTreePluginOpenStackUnregister=true|false (ALPHA - default=false)
                                                  InTreePluginPortworxUnregister=true|false (ALPHA - default=false)
                                                  InTreePluginRBDUnregister=true|false (ALPHA - default=false)
                                                  InTreePluginvSphereUnregister=true|false (ALPHA - default=false)
                                                  JobMutableNodeSchedulingDirectives=true|false (BETA - default=true)
                                                  JobPodFailurePolicy=true|false (BETA - default=true)
                                                  JobReadyPods=true|false (BETA - default=true)
                                                  KMSv2=true|false (ALPHA - default=false)
                                                  KubeletInUserNamespace=true|false (ALPHA - default=false)
                                                  KubeletPodResources=true|false (BETA - default=true)
                                                  KubeletPodResourcesGetAllocatable=true|false (BETA - default=true)
                                                  KubeletTracing=true|false (ALPHA - default=false)
                                                  LegacyServiceAccountTokenTracking=true|false (ALPHA - default=false)
                                                  LocalStorageCapacityIsolationFSQuotaMonitoring=true|false (ALPHA - default=false)
                                                  LogarithmicScaleDown=true|false (BETA - default=true)
                                                  LoggingAlphaOptions=true|false (ALPHA - default=false)
                                                  LoggingBetaOptions=true|false (BETA - default=true)
                                                  MatchLabelKeysInPodTopologySpread=true|false (ALPHA - default=false)
                                                  MaxUnavailableStatefulSet=true|false (ALPHA - default=false)
                                                  MemoryManager=true|false (BETA - default=true)
                                                  MemoryQoS=true|false (ALPHA - default=false)
                                                  MinDomainsInPodTopologySpread=true|false (BETA - default=false)
                                                  MinimizeIPTablesRestore=true|false (ALPHA - default=false)
                                                  MultiCIDRRangeAllocator=true|false (ALPHA - default=false)
                                                  NetworkPolicyStatus=true|false (ALPHA - default=false)
                                                  NodeInclusionPolicyInPodTopologySpread=true|false (BETA - default=true)
                                                  NodeOutOfServiceVolumeDetach=true|false (BETA - default=true)
                                                  NodeSwap=true|false (ALPHA - default=false)
                                                  OpenAPIEnums=true|false (BETA - default=true)
                                                  OpenAPIV3=true|false (BETA - default=true)
                                                  PDBUnhealthyPodEvictionPolicy=true|false (ALPHA - default=false)
                                                  PodAndContainerStatsFromCRI=true|false (ALPHA - default=false)
                                                  PodDeletionCost=true|false (BETA - default=true)
                                                  PodDisruptionConditions=true|false (BETA - default=true)
                                                  PodHasNetworkCondition=true|false (ALPHA - default=false)
                                                  PodSchedulingReadiness=true|false (ALPHA - default=false)
                                                  ProbeTerminationGracePeriod=true|false (BETA - default=true)
                                                  ProcMountType=true|false (ALPHA - default=false)
                                                  ProxyTerminatingEndpoints=true|false (BETA - default=true)
                                                  QOSReserved=true|false (ALPHA - default=false)
                                                  ReadWriteOncePod=true|false (ALPHA - default=false)
                                                  RecoverVolumeExpansionFailure=true|false (ALPHA - default=false)
                                                  RemainingItemCount=true|false (BETA - default=true)
                                                  RetroactiveDefaultStorageClass=true|false (BETA - default=true)
                                                  RotateKubeletServerCertificate=true|false (BETA - default=true)
                                                  SELinuxMountReadWriteOncePod=true|false (ALPHA - default=false)
                                                  SeccompDefault=true|false (BETA - default=true)
                                                  ServerSideFieldValidation=true|false (BETA - default=true)
                                                  SizeMemoryBackedVolumes=true|false (BETA - default=true)
                                                  StatefulSetAutoDeletePVC=true|false (ALPHA - default=false)
                                                  StatefulSetStartOrdinal=true|false (ALPHA - default=false)
                                                  StorageVersionAPI=true|false (ALPHA - default=false)
                                                  StorageVersionHash=true|false (BETA - default=true)
                                                  TopologyAwareHints=true|false (BETA - default=true)
                                                  TopologyManager=true|false (BETA - default=true)
                                                  TopologyManagerPolicyAlphaOptions=true|false (ALPHA - default=false)
                                                  TopologyManagerPolicyBetaOptions=true|false (BETA - default=false)
                                                  TopologyManagerPolicyOptions=true|false (ALPHA - default=false)
                                                  UserNamespacesStatelessPodsSupport=true|false (ALPHA - default=false)
                                                  ValidatingAdmissionPolicy=true|false (ALPHA - default=false)
                                                  VolumeCapacityPriority=true|false (ALPHA - default=false)
                                                  WinDSR=true|false (ALPHA - default=false)
                                                  WinOverlay=true|false (BETA - default=true)
                                                  WindowsHostNetwork=true|false (ALPHA - default=true)
      --healthz-bind-address ipport               The IP address with port for the health check server to serve on (set to '0.0.0.0:10256'  for all IPv4 interfaces and '[::]:10256' for all IPv6 interfaces). Set empty to disable. This parameter is ignored if a config file is specified by --config. (default 0.0.0.0:10256)
  -h, --help                                      help for zero-controller-manager
      --kubeconfig string                         Path to kubeconfig file with authorization and master location information.
      --leader-elect                              Start a leader election client and gain leadership before executing the main loop. Enable this when running replicated components for high availability. (default true)
      --leader-elect-lease-duration duration      The duration that non-leader candidates will wait after observing a leadership renewal until attempting to acquire leadership of a led but unrenewed leader slot. This is effectively the maximum duration that a leader can be stopped before it is replaced by another candidate. This is only applicable if leader election is enabled. (default 15s)
      --leader-elect-renew-deadline duration      The interval between attempts by the acting master to renew a leadership slot before it stops leading. This must be less than the lease duration. This is only applicable if leader election is enabled. (default 10s)
      --leader-elect-resource-lock string         The type of resource object that is used for locking during leader election. Supported options are 'leases', 'endpointsleases' and 'configmapsleases'. (default "leases")
      --leader-elect-resource-name string         The name of resource object that is used for locking during leader election. (default "zero-controller-manager")
      --leader-elect-resource-namespace string    The namespace of resource object that is used for locking during leader election. (default "kube-system")
      --leader-elect-retry-period duration        The duration the clients should wait between attempting acquisition and renewal of a leadership. This is only applicable if leader election is enabled. (default 2s)
      --log-flush-frequency duration              Maximum number of seconds between log flushes (default 5s)
      --logging-format string                     Sets the log format. Permitted formats: "text". (default "text")
      --master string                             The address of the Kubernetes API server (overrides any value in kubeconfig).
      --metrics-bind-address ipport               The IP address with port for the metrics server to serve on (set to '0.0.0.0:10249' for all IPv4 interfaces and '[::]:10249' for all IPv6 interfaces). Set empty to disable. This parameter is ignored if a config file is specified by --config. (default 127.0.0.1:10249)
      --mysql-database string                     Database name for the server to use.
      --mysql-host string                         MySQL service host address. If left blank, the following related mysql options will be ignored. (default "127.0.0.1:3306")
      --mysql-max-connection-life-time duration   Maximum connection life time allowed to connect to mysql. (default 10s)
      --mysql-max-idle-connections int32          Maximum idle connections allowed to connect to mysql. (default 100)
      --mysql-max-open-connections int32          Maximum open connections allowed to connect to mysql. (default 100)
      --mysql-password string                     Password for access to mysql, should be used pair with password.
      --mysql-username string                     Username for access to mysql service.
      --namespace string                          Namespace that the controller watches to reconcile zero-apiserver objects. This parameter is ignored if a config file is specified by --config.
      --node-image string                         The blockchain node image used by default.This parameter is ignored if a config file is specified by --config. (default "ccr.ccs.tencentyun.com/superproj/zero-toyblc-amd64:v0.0.1")
      --parallelism int32                         The amount of parallelism to process. Must be greater than 0. Defaults to 16.This parameter is ignored if a config file is specified by --config. (default 16)
      --sync-period duration                      The minimum interval at which watched resources are reconciled.This parameter is ignored if a config file is specified by --config. (default 10h0m0s)
  -v, --v Level                                   number for the log level verbosity
      --version version[=true]                    Print version information and quit
      --vmodule pattern=N,...                     comma-separated list of pattern=N settings for file-filtered logging (only works for text log format)
      --watch-filter-value string                 The label value used to filter events prior to reconciliation.This parameter is ignored if a config file is specified by --config.
      --write-config-to string                    If set, write the default configuration values to this file and exit.
```

###### Auto generated by spf13/cobra on 21-Jul-2023
