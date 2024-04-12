# Agent

For a complete list of configuration options and details about the agent, please
see the [configuration][] document's agent section.

[configuration]: ../docs/CONFIGURATION.md#agent

## Configuration for telegraf agent

```toml
[agent]
  interval = "10s"
  round_interval = true

  # metric_batch_size = 1000
  # metric_buffer_limit = 10000

  # collection_jitter = "0s"
  # collection_offset = "0s"

  # flush_interval = "10s"
  # flush_jitter = "0s"

  # precision = "0s"

  # debug = false
  # quiet = false
  # logtarget = "file"
  # logfile = ""
  # logfile_rotation_interval = "0h"
  # logfile_rotation_max_size = "0MB"
  # logfile_rotation_max_archives = 5
  # log_with_timezone = ""

  # hostname = ""
  # omit_hostname = false

  # snmp_translator = "netsnmp"

  # statefile = ""

  # skip_processors_after_aggregators = false
```
