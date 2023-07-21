## zeroctl new

Generate demo command code

### Synopsis

Used to generate demo command source code file.

 Can use this command generate a command template file, and do some modify based on your needs. This can improve your R &D efficiency.

```
zeroctl new CMD_NAME | CMD_NAME CMD_DESCRIPTION
```

### Examples

```
  # Create a default 'test' command file without a description
  zeroctl new test
  
  # Create a default 'test' command file in /tmp/
  zeroctl new test -d /tmp/
  
  # Create a default 'test' command file with a description
  zeroctl new test "This is a test command"
  
  # Create command 'test' with two subcommands
  zeroctl new -g test "This is a test command with two subcommands"
```

### Options

```
  -g, --group           Generate two subcommands.
  -h, --help            help for new
  -d, --outdir string   Where to create demo command files. (default ".")
```

### Options inherited from parent commands

```
      --alsologtostderr                           log to standard error as well as files
  -c, --config FILE                               Read configuration from specified FILE, support JSON, TOML, YAML, HCL, or Java properties formats.
      --gateway.address string                    The address and port of the Zero API server
      --gateway.certificate-authority string      Path to a cert file for the certificate authority
      --gateway.insecure-skip-tls-verify          If true, the server's certificate will not be checked for validity. This will make your HTTPS connections insecure
      --gateway.max-retries int                   Maximum number of retries.
      --gateway.retry-interval duration           The interval time between each attempt.
      --gateway.timeout duration                  The length of time to wait before giving up on a single server request. Non-zero values should contain a corresponding time unit (e.g. 1s, 2m, 3h). A value of zero means don't timeout requests.
      --kubeconfig string                         Paths to a kubeconfig. Only required if out-of-cluster.
      --log-backtrace-at traceLocations           when logging hits line file:N, emit a stack trace
      --log-dir string                            If non-empty, write log files in this directory
      --log-link string                           If non-empty, add symbolic links in this directory to the log files
      --logbuflevel int                           Buffer log messages logged at this level or lower (-1 means don't buffer; 0 means buffer INFO only; ...). Has limited applicability on non-prod platforms.
      --logtostderr                               log to standard error instead of files
      --profile string                            Name of profile to capture. One of (none|cpu|heap|goroutine|threadcreate|block|mutex) (default "none")
      --profile-output string                     Name of the file to write the profile to (default "profile.pprof")
      --stderrthreshold severityFlag              logs at or above this threshold go to stderr (default 2)
      --user.client-certificate string            Path to a client certificate file for TLS
      --user.client-key string                    Path to a client key file for TLS
      --user.password string                      Password for basic authentication to the API server
      --user.secret-id string                     SecretID for JWT authentication to the API server
      --user.secret-key string                    SecretKey for jwt authentication to the API server
      --user.token string                         Bearer token for authentication to the API server
      --user.username string                      Username for basic authentication to the API server
      --usercenter.address string                 The address and port of the Zero API server
      --usercenter.certificate-authority string   Path to a cert file for the certificate authority
      --usercenter.insecure-skip-tls-verify       If true, the server's certificate will not be checked for validity. This will make your HTTPS connections insecure
      --usercenter.max-retries int                Maximum number of retries.
      --usercenter.retry-interval duration        The interval time between each attempt.
      --usercenter.timeout duration               The length of time to wait before giving up on a single server request. Non-zero values should contain a corresponding time unit (e.g. 1s, 2m, 3h). A value of zero means don't timeout requests.
  -v, --v Level                                   log level for V logs
      --version version[=true]                    Print version information and quit
      --vmodule vModuleFlag                       comma-separated list of pattern=N settings for file-filtered logging
      --warnings-as-errors                        Treat warnings received from the server as errors and exit with a non-zero exit code
      --zeroconfig string                         Path to the zeroconfig file to use for CLI.
```

### SEE ALSO

* [zeroctl](zeroctl.md)	 - zeroctl controls the zero cloud platform

###### Auto generated by spf13/cobra on 21-Jul-2023
