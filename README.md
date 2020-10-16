# Google Compute Engine Graceful Shutdown 

This tool monitors an Instance Group for changes. The moment the given Instance Template
is no longer part of the monitored Instance Group, `on-shutdown` is run. Afterwards
the Instance is deleted via an API call.

It does not currently trigger if:

* Instance is deleted
* Instance is restarted


## Motivation

Google Shutdown Scripts are unreliable and only run for 60-90s max.


## Usage

```
./gce-graceful-shutdown -on-shutdown 'docker stop --time 30 my-container' -timeout 30s 
```

| Variable                         | Description                                                                                                                                                                                                                                          |
|----------------------------------|---------------------------------------------------------------------------------------|
| `-creds`                         | Path to Google Service Account Credentials. Defaults to Application Credentials.      |
| `-on-shutdown`                   | ***Required*** Command to run when instance is about to be deleted.                   |
| `-timeout`                       | Time to wait for -`on-shutdown` to finish. Default `10m`.                             |
| `-project`                       | Manually set Google Project.                                                          |
| `-region`                        | Manually set Region.                                                                  |
| `-zone`                          | Manually set Zone.                                                                    |
| `-instance-group`                | Manually set Instance Group.                                                          |
| `-instance-template`             | Manually set Instance Template.                                                       |
| `-instance`                      | Manually set Instance.                                                                |


