# README
## What's this?
AWS Athena command line tool helping you to casually execute your queries.

* Execute query given by arg or from stdin.
  - Templating.
  - Synchronously waiting for its completion.
  - If interrupted, stop its execution.
* List query-executions, not query-excution IDs.
* Get s3 object with execution-id.
* Show query result with exeuction-id.

## Getting Started
You have these.
* A s3 bucket to write outputs.
* Your AWS profile with proper permissions.
* ELB logs

Rename `.athena-cli.yaml.templ` to `.athena-cli.yaml` and fill placeholders with your values.
```
# verbose: true
# work-group: <your-workgroup>
output-location: s3://<your-bucket>/outputs
database-name: <database-name>
timeout: 30s
```
Let's say the both bucket and database is `sample_elb` in the getting started.

Then run `list`, you'll see the below one and ready to execute.
```
$ export AWS_PROFILE=<your-profile>
$ ./athena-cli list
[]
```

### Query execution
Let's create your database and check it's created.
```
$ athena-cli query "create database sample_elb"

$ athena-cli query --query.result-csv "create database sample_elb"
sample_elb
```

Let's see the current execution list.
```
$ athena-cli list --list.limit 1 | jq .
[
  {
    "Query": "create database sample_elb",
    "QueryExecutionContext": {
      "Database": "sample_elb"
    },
    "QueryExecutionId": "6cf8459e-c72a-4375-ab01-34f0244868b8",
    "ResultConfiguration": {
      "EncryptionConfiguration": null,
      "OutputLocation": "s3://sample_elb/outputs/6cf8459e-c72a-4375-ab01-34f0244868b8.txt"
    },
    ...
  }
]
```

That's a sample for a simple query.


## Querying ELB logs
Let's try to query ELB logs in a S3 bucket. Prepare these.
- S3 bucket name (w/ prefix)
- AWS account ID

First of all, create a table.
```
$ athena-cli query --dry-run \
  --query.values '{
    "s3buketname": "a_elb/foo",
    "accountId": "123456789123"
  }' < query/mk-table-elb.q
...
LOCATION 's3://a_elb/foo/elb/AWSLogs/123456789123/elasticloadbalancing/';
```

You can check executed query with `--dry-run`.
Then run it without `--dry-run`.
```
$ athena-cli query --query.values '{"s3buketname":"a_elb/foo","accountId":"123456789123"}' < query/mk-table-elb.q

$
```

You're ready to select. Just in case, put `limit 1` at the tail.
```
$ athena-cli query "select * from elb_logs limit 1" | jq .
{
  "execution_id": "b5022e5e-0fec-4782-9649-c35633f0328b",
  "rows": [
    {
      "backend_ip": "10.1.145.234",
      "backend_port": "8080",
...
```

# Examples
## CloudTrail Logs
Steps of [Querying AWS CloudTrail Logs](https://docs.aws.amazon.com/athena/latest/ug/cloudtrail-logs.html).

Creating the table.
```
athena-cli query --query.values '{
    "s3bucketname": "your-cloudtrail-logs-bucket",
    "accountId": "123456789123"
  }' < query/mk-table-ct.q
```

Load the partitions.
```
athena-cli query --query.values '{
    "s3bucketname": "your-cloudtrail-logs-bucket",
    "accountId": "123456789123"
    "region": "ap-northeast-1",
    "year": "2020",
    "month": "03",
    "day": "29"
  }' < query/mk-part.q
```

```
athena-cli query "select * from cloudtrail_logs limit 1"
```

```
athena-cli query --query.values '{"year":"2019","month":"08","day":"19"}' < query/drop-part.q
```
