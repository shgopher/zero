## zero-demo

Launch a zero demo server

### Synopsis

The demo server is a standard, specification-compliant demo 
example of the zero service.

Find more zero-demo information at:
    https://github.com/superproj/zero/blob/master/docs/guide/en-US/cmd/zero-demo.md

```
zero-demo [flags]
```

### Options

```
  -c, --config FILE                                      Read configuration from specified FILE, support JSON, TOML, YAML, HCL, or Java properties formats.
      --db.database string                               Database name for the server to use.
      --db.host string                                   MySQL service host address. If left blank, the following related mysql options will be ignored. (default "127.0.0.1:3306")
      --db.log-mode int                                  Specify gorm log level. (default 1)
      --db.max-connection-life-time duration             Maximum connection life time allowed to connect to mysql. (default 10s)
      --db.max-idle-connections int                      Maximum idle connections allowed to connect to mysql. (default 100)
      --db.max-open-connections int                      Maximum open connections allowed to connect to mysql. (default 100)
      --db.password string                               Password for access to mysql, should be used pair with password.
      --db.username string                               Username for access to mysql service.
      --enable-tls                                       Enable TLS for grpc server.
      --feature-gates mapStringBool                      A set of key=value pairs that describe feature gates for alpha/experimental features. Options are:
                                                         AllAlpha=true|false (ALPHA - default=false)
                                                         AllBeta=true|false (BETA - default=false)
                                                         ContextualLogging=true|false (ALPHA - default=false)
                                                         LoggingAlphaOptions=true|false (ALPHA - default=false)
                                                         LoggingBetaOptions=true|false (BETA - default=true)
                                                         MachinePool=true|false (ALPHA - default=false)
      --grpc.addr string                                 Specify the gRPC server bind address and port. (default "0.0.0.0:8081")
      --grpc.network string                              Specify the network for the gRPC server. (default "tcp")
      --grpc.timeout duration                            Timeout for server connections. (default 30s)
  -h, --help                                             help for zero-demo
      --http.addr string                                 Specify the HTTP server bind address and port. (default "0.0.0.0:8443")
      --http.network string                              Specify the network for the HTTP server. (default "tcp")
      --http.timeout duration                            Timeout for server connections. (default 30s)
      --jaeger.env string                                Specify the deployment environment(dev/test/staging/prod). (default "dev")
      --jaeger.server string                             Server is the url of the Jaeger server. (default "http://127.0.0.1:14268/api/traces")
      --jaeger.service-name string                       Specify the service name for jaeger resource.
      --jwt.expired duration                             JWT token expiration time. (default 2h0m0s)
      --jwt.key string                                   Private key used to sign jwt token. (default "zero666")
      --jwt.max-refresh duration                         This field allows clients to refresh their token until MaxRefresh has passed. (default 2h0m0s)
      --jwt.signing-method string                        JWT token signature method. (default "HS512")
      --kubeconfig string                                Path to kubeconfig file with authorization and master location information.
      --log.disable-caller                               Disable output of caller information in the log.
      --log.disable-stacktrace                           Disable the log to record a stack trace for all messages at or above panic level.
      --log.enable-color                                 Enable output ansi colors in plain format logs.
      --log.format FORMAT                                Log output FORMAT, support plain or json format. (default "console")
      --log.level LEVEL                                  Minimum log output LEVEL. (default "info")
      --log.output-paths strings                         Output paths of log. (default [stdout])
      --metrics.allow-metric-labels stringToString       The map from metric-label to value allow-list of this label. The key's format is <MetricName>,<LabelName>. The value's format is <allowed_value>,<allowed_value>...e.g. metric1,label1='v1,v2,v3', metric1,label2='v1,v2,v3' metric2,label1='v1,v2,v3'. (default [])
      --metrics.disabled-metrics strings                 This flag provides an escape hatch for misbehaving metrics. You must provide the fully qualified metric name in order to disable it. Disclaimer: disabling metrics is higher in precedence than showing hidden metrics.
      --metrics.show-hidden-metrics-for-version string   The previous version for which you want to show hidden metrics. Only the previous minor version is meaningful, other values will not be allowed. The format is <major>.<minor>, e.g.: '1.16'. The purpose of this format is make sure you have the opportunity to notice if the next release hides additional metrics, rather than being surprised when they are permanently removed in the release after that.
      --redis.addr string                                Address of your Redis server(ip:port). (default "127.0.0.1:6379")
      --redis.database int                               Database to be selected after connecting to the server.
      --redis.dial-timeout duration                      Dial timeout for establishing new connections. (default 5s)
      --redis.enable-trace                               Redis hook tracing (using open telemetry).
      --redis.max-retries int                            Maximum number of retries before giving up. (default 3)
      --redis.min-idle-conn int                          Minimum number of idle connections which is useful when establishing new connection is slow.
      --redis.password string                            Optional auth password for redis db.
      --redis.pool-size int                              Maximum number of socket connections. (default 10)
      --redis.pool-timeout duration                      Amount of time client waits for connection if all connections are busy before returning an error.
      --redis.read-timeout duration                      Timeout for socket reads. (default 3s)
      --redis.username string                            Username for access to redis service.
      --redis.write-timeout duration                     Timeout for socket writes.
      --tls.ca-cert string                               Path to ca cert for connecting to the server.
      --tls.cert string                                  Path to cert file for connecting to the server.
      --tls.insecure-skip-verify                         Controls whether a client verifies the server's certificate chain and host name.
      --tls.key string                                   Path to key file for connecting to the server.
      --tls.use-tls                                      Use tls transport to connect the server.
      --version version[=true]                           Print version information and quit
```

###### Auto generated by spf13/cobra on 21-Jul-2023
