# Date Processor Plugin

Use the `date` processor to either add a metric timestamp as a human readable
tag or field or to translate an existing timestamp from one format to another.

A common use is to add a tag that can be used to group by month or year.

A few example use cases include:

1) consumption data for utilities on per month basis
2) bandwidth capacity per month
3) compare energy production or sales on a yearly or monthly basis
4) translate an existing timestamp into a unix timestamp

## Global configuration options <!-- @/docs/includes/plugin_config.md -->

In addition to the plugin-specific configuration settings, plugins support
additional global and plugin configuration settings. These settings are used to
modify metrics, tags, and field or create aliases and configure ordering, etc.
See the [CONFIGURATION.md][CONFIGURATION.md] for more details.

[CONFIGURATION.md]: ../../../docs/CONFIGURATION.md#plugins

## Configuration

```toml @sample.conf
# Dates measurements, tags, and fields that pass through this filter.
[[processors.date]]
  ## Field or Tag Key
  ## Specify the new tag key or field key to create. If the tag or field already
  ## exists, then the tag or key value is assumed to be a timestamp to
  ## translate and the `existing_date_format` option is also required.
  # tag_key = "month"
  # field_key = "month"

  ## Existing timestamp format
  ## If translating an existing timestamp specify it's format in Go reference
  ## time below.
  # existing_date_format = "Mon Jan 2 15:04:05 -0700 MST 2006"

  ## New timestamp format
  ## Date format string, must be a representation of the Go "reference time"
  ## which is "Mon Jan 2 15:04:05 -0700 MST 2006". If using a field_key,
  ## date format can also be one of "unix", "unix_ms", "unix_us", or "unix_ns",
  ## which will insert an integer field.
  date_format = "Jan"

  ## Offset duration added to the date string when writing the new tag.
  # date_offset = "0s"

  ## Timezone to use when creating the tag or field using a reference time
  ## string.  This can be set to one of "UTC", "Local", or to a location name
  ## in the IANA Time Zone database.
  ##   example: timezone = "America/Los_Angeles"
  # timezone = "UTC"
```

## Example

Example of adding a timestamp:

```diff
- throughput lower=10i,upper=1000i,mean=500i 1560540094000000000
+ throughput,month=Jun lower=10i,upper=1000i,mean=500i 1560540094000000000
```

Example of translating a timestamp:

```diff
- metric value=42,ts="" 1560540094000000000
+ metric value=42,ts= 1560540094000000000
```
