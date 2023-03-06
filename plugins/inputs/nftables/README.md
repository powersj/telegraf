# nftables Input Plugin

TODO

## Global configuration options <!-- @/docs/includes/plugin_config.md -->

In addition to the plugin-specific configuration settings, plugins support
additional global and plugin configuration settings. These settings are used to
modify metrics, tags, and field or create aliases and configure ordering, etc.
See the [CONFIGURATION.md][CONFIGURATION.md] for more details.

[CONFIGURATION.md]: ../../../docs/CONFIGURATION.md#plugins

## Configuration

```toml @sample.conf
# Gather packets and bytes throughput from nftables
# This plugin ONLY supports Linux
[[inputs.nftables]]
  ## Use sudo to run nft commands
  # use_sudo = false
```

### Using sudo

If your account does not already have the ability to run commands
with passwordless sudo then updates to the sudoers file are required. Below
is an example to allow the requires LVM commands:

First, use the `visudo` command to start editing the sudoers file. Then add
the following content, where `<username>` is the username of the user that
needs this access:

```text
Cmnd_Alias NFT = /usr/bin/nft *
<username>  ALL=(root) NOPASSWD: NFT
Defaults!NFT !logfile, !syslog, !pam_session
```

## Metrics

### Measurements & Fields

* nftables
  * packets (integer)
  * bytes (integer)

### Tags

* All measurements have the following tags:
  * family
  * table
  * chain

## Example Output


```text
nftables,table=filter,chain=INPUT,ruleid=ssh pkts=100i,bytes=1024i 1453831884664956455
nftables,table=filter,chain=INPUT,ruleid=httpd pkts=42i,bytes=2048i 1453831884664956455
```
