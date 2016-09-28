# BigQuery Table Autocreator
This tool can be used to create BigQuery tables in several projects with names
in the following format: `<prefix><date_tomorrow>`.

`prefix` can be changed by editing the YAML configuration file in `./configs`.
`date_tomorrow` is tomorrow's date which looks like this: `20060102`.

This table name format allows you to use BigQuery's `TABLE_DATE_RANGE()`
function to query data partitioned by day.

## Usage
### Development
Install [Go](https://golang.org/dl/) and [gb](https://getgb.io/) first, then edit the configuration file
and schema in `./configs` and run this:
```sh
make run
```

### Production
This utility was created to be executed every day using cron. Configure it to
run every day at a time when you are comfortable receiving Slack mentions in
case of errors.
