# CDC Status

Combines data from Octopus and Metricly to print a summary of CDC statuses.

## Prerequisites

This command read a config file from **~/.monitoring/cdc_status.hcl**.

Example:

```hcl
Metricly {
    Username = "<your metricly username>"
    Password = "<your metricly password>"
}

Octopus {
    instanceURL = "https://<your organization>.octopus.app"
    apiKey      = "<your API Key>"
    space       = "<the Octopus Space to query>"
}
```

## Invocation

```shell
$ cdc_status
Current CDC Install/Replication status:
...
```
