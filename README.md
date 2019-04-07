# S3logsbeat

S3logsbeat is a [beat](https://www.elastic.co/products/beats) to read logs from AWS S3 and send them to
ElasticSearch. AWS uses S3 as destination for several internal logs: ALB, CloudFront, CloudTrail, etc.
This beat is based on [S3 event notifications](https://docs.aws.amazon.com/AmazonS3/latest/dev/NotificationHowTo.html)
to send notifications to an SQS queue when a new object is created on S3. Then, S3logsbeat polls
these SQS queues, reads new objects messages (ignoring others), downloads S3 objects, parses logs to convert
them into events, and finally publishes to ElasticSearch. If all events are published correctly, SQS message
is deleted from SQS queue.

## Whys

### Why should I use S3logsbeat if I can use Lambdas?
Lambdas are perfect when you have a few log entries. However, its price is based on how much time it takes to send
events to ElasticSearch, if you have many log entries, its price could be eleveted.

As S3logsbeat consumes few resources (~20MB RAM on my tests), it can be installed on your ELK or on a container.

For instance, assume we have 10 million of new S3 objects of ALBs logs per month. Each object takes 1 second
to be downloaded, parsed, and sent to our ElasticSearch (I'm assuming an average of 100 events per object). If
we use the minimum amount of memory allowed by Lamba (128MB), our cost is (according to [AWS Lambda pricing
calculator](https://s3.amazonaws.com/lambda-tools/pricing-calculator.html)): $2 for requests + $20.84 for execution = $22.84 / month.

As S3logsbeat can read until 10 SQS messages per request, we need to perform 1,000,000 request to obtain these 10 million of new S3 objects to SQS. As we read the message and delete it after processing, we need 2,000,000 request,
which implies an SQS cost of $0.80/month. S3 costs would be: 10,000,000 GETs = $0.4 (as the instance in which
S3logsbeat is running in the same region, I'm ignoring data transfer cost). The total cost with S3logs beat is: $1.20/month.

## Features
S3logsbeat has the following features:
* Limited workers to poll from SQS and download objects from S3 to avoid exceeding AWS request limits
* Usage of internal bounded queues to avoid overloading outputs
* If output is overloaded or inaccessible, no more messages are read from SQS
* High availability: you can have several S3logsbeat running in parallel
* Reliability: SQS messages are only deleted when output contains all events
* Avoid duplicates on supported outputs
* Supported S3 log parsers: ELB, ALB, CloudFront
* Extra fields based on S3 key
* Delayed shutdown based on timout and pending messages to be acked by outputs
* Limited amount of resources: ~20MB RAM in my tests

However, it has some use cases not supported yet:
* S3 objects already present on bucket when S3 event notifications is activated are not processed
* CloudTrail logs not supported yet
* Custom logs not supported yet

Unsupported features are on the roadmap, so just wait for them.

### Reduce overloading on outputs
We have to configure internal beat queue in order to reduce overloading on outputs. Edit your configuration
file and add the following configuration on root:
```yaml
queue.mem:
  events: 128 # default: 4096
  flush.min_events: 64 # default: 2048
```

In this way, if we try to send 128 messages to outputs and output is outline or saturated, publishing is blocked. As
publishing is blocked, all internal bounded queues are filled eventually and no more SQS messages are read (avoiding
extra requests to SQS or S3). Once output is available again, messages present on internal queues are sent to
output and new messages are then read from SQS.

### SQS messages deleted when events are acked
S3logsbeat is based on [SQS Visibility timeout feature](https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/sqs-visibility-timeout.html). When a message is read from SQS, it just "dissapears" from the queue. However, if message is not deleted, it
reappears on the queue after visitiblity timeout expires, making it available to be read again. This is perfect
for those cases in which a consumer deads unexpectedly.

S3logsbeat uses this principle. The workflow to process an SQS is the following:
1) SQS message is read from an SQS queue
2) New S3 objects appearing on SQS message are downloaded from S3 and parsed to generate events
3) Each event is sent to output
4) When output confirms the reception of all events related to an SQS message, SQS message is deleted from the queue

### Avoid duplicates
S3logsbeat can avoid duplicates on ElasticSearch output transparently by adding an event (document on ES) identifier
based on its content.

It can also be used on Logstash output because a `@metadata` field called `_id` is added with this information. You
can configure your [Logstash ES output](https://www.elastic.co/guide/en/logstash/current/plugins-outputs-elasticsearch.html)
in the following way to take this value and use it as document identifier:
```
output {
  elasticsearch {
    ···
    document_id => "%{[@metadata][_id]}"
    ···
  }
}
```

### Fields based on S3 key
If you are sending several origin logs to the same S3 bucket and you want to distinguish them on ElasticSearch,
you can set a regular expression on `key_regex_fields` in order to parse S3 keys and add extracted fields to
each event extracted from it.

For instance, imagine you are sending production and testing logs to the same S3 bucket from your ALBs. In order
to distinguish them, you added a prefix based on pattern `{environment}-{application}/` to your ALBs configuration.

By default S3logsbeat is aware of this and generates events based on content, not from where them have been
extracted. Due to that, you couldn't distinguish between production and testing logs. In order to make S3logsbeat
aware of it, you have to configure your input as follows:
```yaml
s3logsbeat:
  inputs:
    - type: sqs
      queues_url:
        - https://sqs.{aws-region}.amazonaws.com/{account ID}/{queue name}
      log_format: alb
      key_regex_fields: ^(?P<environment>[^\-]+)-(?P<application>[^/\-]+)
      poll_frequency: 1m
```

### Delayed shutdown
By default, when S3logsbeat is stopped, SQS messages being processed are cancelled. It is not problematic because
SQS message is not deleted until all events are present on output. Due to that, when you starts S3logsbeat again,
SQS message cancelled is processed again. That combined with no duplicate events produces the same events on
ElasticSearch as if initial message would be processed entirely.

However, we are doing the same process again on the same SQS message and the same S3 objects associated with it. In
order to avoid it, we can configure a shutdown timeout as follows:
```yaml
s3logsbeat:
  shutdown_timeout: 5s
```

Then, when S3logsbeat is stopped, it stops to process new SQS messages and waits until (what happens first):
1) All current SQS messages are processed, or
2) Timeout expires

First case is the typical one when everything is ok because an SQS message is processed in less than a second and we have not to repeat the same process.

Second case can happen when we have a lot of S3 objects on pending SQS messages. This case is similar to the one exposed initially when `shutdown_timeout` is not configured. It means, unprocessed S3 objects will be processed when S3logsbeat starts again.

### AWS IAM
S3logsbeat requires the following IAM permissions:
* SQS permissions: `sqs:ReceiveMessage` and `sqs:DeleteMessage`.
* S3 permissions: `s3:GetObject`.

IAM policy:
```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "sqs:ReceiveMessage",
        "sqs:DeleteMessage"
      ],
      "Resource": "arn:aws:sqs:*:123456789012:<QUEUE_NAME>"
    },
    {
      "Effect": "Allow",
      "Action": [
        "s3:GetObject"
      ],
      "Resource": "arn:aws:s3:::<BUCKET-NAME>/*"
    }
  ]
}
```

### Initial import
You may already have S3 log files when you configure an SQS queue to import new files via `s3logsbeat`. If this is the case,
you can import those files by using the command `s3imports` and a configuration file as this:
```yaml
s3logsbeat:
  inputs:
    # S3 inputs (only taken into account when command `s3import` is executed)
    -
      type: s3
      # S3
      buckets:
        - s3://mybucket/mypath
      log_format: alb
      # Optional fields extractor from key. E.g. key=staging-myapp/eu-west-1/2018/06/01/
      key_regex_fields: ^(?P<environment>[^\-]+)-(?P<application>[^/]+)/(?P<awsregion>[^/]+)
      since: 2018-10-15T01:00 # ISO8601 format - optional
      to: 2018-11-20T01:00 # ISO8601 format - optional
```

Using command `./s3logsbeat s3imports -c config.yml` you can import all those S3 files that you already have on S3. This command
is not executed as a daemon and exits when all S3 objects are processed. Those SQS inputs present on configuration will be
ignored when command `s3imports` is executed.

This command is useful on first import, however, you should take care because in combination with standard mode of `s3logsbeat`
can generate duplicates. In order to avoid this problem you can:
* Use `@metadata._id` in order to avoid duplicates on ElasticSearch (see section [Avoid duplicates](#avoid-duplicates)).
* Configure the SQS queue and S3 event notifications. Wait until the first element is present on the queue. Via console or
  cli, analyse the element present on the SQS queue without deleting it (it will reappear later). Then edit yaml configuration
  and set the `to` property to just one second before the one obtained and execute `s3imports` command.

### Supported log formats
`s3logsbeat` supports the following log formats:
* `elb`: parses Elastic Load Balancer (classic ELB) log.
* `alb`: parses Application Load Balancer (ALB) log.
* `cloudfront`: parses CloudFront logs.
* `json`: parses JSON logs. Requires the following options (set via parameter `log_format_options`):
    * `timestamp_field`: field that represents the timestamp of log event. Mandatory.
    * `timestamp_format`: format in which timestamp is represented and from which should be converted into Date/Time. See [Suported timestamp formats](#supported-timestamp-formats). Mandatory.

### Supported timestamp formats
The following timestamp formats are supported:
* `timeUnixMilliseconds`: long or string with epoc millis.
* `timeISO8601`: string with ISO8601 format.
* `time:layout`: string with layout format present after prefix `time:`. Valid layouts correspond to ones parsed by [time.Parse](https://golang.org/pkg/time/#Parse).

### Example of events

#### ALB
The following log event example is generated when `log_format: alb` is present:
```
{
  "_index" : "yourindex-2018.09.14",
  "_type" : "doc",
  "_id" : "678d4e087033da5bc429c1834941b71868b1cb6c",
  "_score" : 1.0,
  "_source" : {
    "@timestamp" : "2018-09-14T19:11:02.796Z",
    "trace_id" : "Root=1-5b9c07c6-686bc18f187f75f58d06cd63",
    "target_group_arn" : "arn:aws:elasticloadbalancing:region:123456789012:targetgroup/yourtg/f3b5cb16d14453e9",
    "application" : "yourapp",
    "beat" : {
      "hostname" : "macbook-pro-4.home",
      "version" : "6.0.2",
      "name" : "macbook-pro-4.home"
    },
    "environment" : "yourenvironment",
    "client_port" : 37054,
    "received_bytes" : 272,
    "target_port" : 80,
    "user_agent" : "cURL",
    "request_url" : "https://www.example.com/",
    "ssl_protocol" : "TLSv1.2",
    "elb" : "app/yourenvironment-yourapp/ad4ceee8a897566c",
    "request_proto" : "HTTP/1.0",
    "response_processing_time" : 0,
    "ssl_cipher" : "ECDHE-RSA-AES128-GCM-SHA256",
    "client_ip" : "1.2.3.4",
    "target_status_code" : 200,
    "sent_bytes" : 11964,
    "request_processing_time" : 0,
    "elb_status_code" : 200,
    "target_processing_time" : 0.049,
    "request_verb" : "GET",
    "type" : "https",
    "target_ip" : "2.3.4.5"
  }
}
```

These fields corresond to the ones established by AWS on [this page](https://docs.aws.amazon.com/es_es/elasticloadbalancing/latest/application/load-balancer-access-logs.html).

#### CloudFront
The following log event example is generated when `log_format: cloudfront` is present:
```
{
  "_index" : "yourindex-2018.09.02",
  "_type" : "doc",
  "_id" : "4d2f4dceb5ed5b4a4300eccc657608a7cc57923d",
  "_score" : 1.0,
  "_source" : {
    "@timestamp" : "2018-09-02T10:10:11.000Z",
    "ssl_protocol" : "TLSv1.2",
    "cs_protocol_version" : "HTTP/1.1",
    "x_edge_response_result_type" : "Hit",
    "cs_user_agent" : "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:61.0) Gecko/20100101 Firefox/61.0",
    "x_edge_request_id" : "9QuHsIZRJxPbF_F-3to6KberARZ7Ddd6wthSnivYJEWalJWqlYCe7A==",
    "ssl_cipher" : "ECDHE-RSA-AES128-GCM-SHA256",
    "x_edge_result_type" : "Hit",
    "x_edge_location" : "MAD50",
    "x_host_header" : "www.example.com",
    "beat" : {
      "name" : "macbook-pro-4.home",
      "hostname" : "macbook-pro-4.home",
      "version" : "6.0.2"
    },
    "cs_bytes" : 457,
    "sc_status" : 200,
    "cs_cookie" : "-",
    "cs_protocol" : "https",
    "cs_method" : "GET",
    "cs_uri_stem" : "/",
    "cs_host" : "d111111abcdef8.cloudfront.net",
    "c_ip" : "1.2.3.4",
    "time_taken" : 0.004,
    "sc_bytes" : 1166
  }
}
```

These fields corresond to the ones established by AWS on [this page](https://docs.aws.amazon.com/AmazonCloudFront/latest/DeveloperGuide/AccessLogs.html#LogFileFormat).

## Development

### Requirements

* [Golang](https://golang.org/dl/) 1.7

### Init Project
To get running with S3logsbeat and also install the
dependencies, run the following command:

```
make setup
```

It will create a clean git history for each major step. Note that you can always rewrite the history if you wish before pushing your changes.

To push S3logsbeat in the git repository, run the following commands:

```
git remote set-url origin https://github.com/sequra/s3logsbeat
git push origin master
```

For further development, check out the [beat developer guide](https://www.elastic.co/guide/en/beats/libbeat/current/new-beat.html).

### Build

To build the binary for S3logsbeat run the command below. This will generate a binary
in the same directory with the name s3logsbeat.

```
make
```


### Run

To run S3logsbeat with debugging output enabled, run:

```
./s3logsbeat -c s3logsbeat.yml -e -d "*"
```


### Test

To test S3logsbeat, run the following command:

```
make testsuite
```

alternatively:
```
make unit-tests
make system-tests
make integration-tests
make coverage-report
```

The test coverage is reported in the folder `./build/coverage/`

### Update

Each beat has a template for the mapping in elasticsearch and a documentation for the fields
which is automatically generated based on `fields.yml` by running the following command.

```
make update
```


### Cleanup

To clean  S3logsbeat source code, run the following commands:

```
make fmt
```

To clean up the build directory and generated artifacts, run:

```
make clean
```


### Clone

To clone S3logsbeat from the git repository, run the following commands:

```
mkdir -p ${GOPATH}/src/github.com/sequra/s3logsbeat
git clone https://github.com/sequra/s3logsbeat ${GOPATH}/src/github.com/sequra/s3logsbeat
```


For further development, check out the [beat developer guide](https://www.elastic.co/guide/en/beats/libbeat/current/new-beat.html).


## Packaging

The beat frameworks provides tools to crosscompile and package your beat for different platforms. This requires [docker](https://www.docker.com/) and vendoring as described above. To build packages of your beat, run the following command:

```
make package
```

This will fetch and create all images required for the build process. The hole process to finish can take several minutes.
