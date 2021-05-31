# SQS Messages replayer 

## Overview

This is a simple tool that allows replaying messages from the AWS SQS DLQ. While the primary function of this tool is to replay messages present in the DLQ, it can also be used to push a message from one given queue to another regardless if the source queue is a DLQ or not.

## Usage

Below is the print of the following command: `./sqs-pub -help`, which describe the usage of this tool


```
USAGE
  sqs-pub [-from queue1 - to queue2 -filter text1,text2,...]

FLAGS
  -delete true                                          delete messages from source after successfuly pushed to destination queue
  -filters 10104211111292                               comma separted text that can be used a message body filter
  -from vf-cm-dev-marketplace-emea-bi-deadletter-queue  sqs queue from where messages will be sourced from
  -to vf-cm-dev-marketplace-emea-bi-orders              sqs queue from where messages will be sourced from

```

## Flags
1. `delete`: by default this is set to true and the messages will be deleted from the soruce queue once successfully published to the destination queue.
2. `filters`: a comma separated list of keys that can be used to identify a given message is meant to be processed or not. This is a very simple filtering system, in subsequent releases timestamp or messageId can be used to filter out messages.
3. `from`: The source queue name (exp, DLQ)
4. `to`: The destionation queue name.
5. `dryrun`: By the default this flag is to false. If it is set to true then it will run the replay process but withouth sending the actual messages to the destination queue and reports will be generated as it has been a real dry. This flag is useful to get an picture of what is this tool attempt to do when run it for real.

## How to run

1. Set the AWS Credentials in the terminal 

```
export AWS_ACCESS_KEY_ID="xxxxxxxV2M3L"
export AWS_SECRET_ACCESS_KEY="xxxxxxxxDK"
export AWS_SESSION_TOKEN="IQoJxxxxxxxxx"
```

2. Run the binary as follow
```
./sqs-pub -from=queue-name-source -to=queue-name-destination -filters=text1,text2

```

## Binary releases
[Releases](https://github.com/akhettar/sqs-pub/releases)

## License
[MIT](LICENSE)
