# DLQ Message Replayer

This is a utility to replay a message from DLQ to a target SQS queue. If there are any failures to push to the target queue the body of the message will be saved into `failed-messages.json` file.

## How to run

1. Set the AWS Credentials in the terminal 

```
export AWS_ACCESS_KEY_ID="ASIAQIYCJYNNGFOV2M3L"
export AWS_SECRET_ACCESS_KEY="yqezFhH0Q+a3yYi6uSQSxGinMyThaAecxZdh8BDK"
export AWS_SESSION_TOKEN="IQoJb3JpZ2luX2VjEKn//////////wEaCXVzLWVhc3QtMSJIMEYCIQDawR3t48PP73Mm9/BMksRgZlbADNKrqIXsO3celPRAFAIhALBjYL6W7mSQcc/IeTNrTeVFTFK9Hqz/M+tgduA0z00kKpUDCEIQAhoMMDE4Nzk1MzE2MDU4Igz1JgvYH7SrIA2FQ0cq8gIgFU/Ww8ZwiCumJaewzrQtcvS1vd44Mtu6nSdw4WM8lGzeeQU9EIZCD0O5xCGqLtMSTyjQX8m5xJEpQGgZ/r/z3u/nnipnLr0z1pthR0I/UwQsCK49mMIWYGG8BBt/TuMEuXeCPgt0+oY5DqPlQAR+uaVyp34vxfzYUckaePabj2VGSQfU1hJZloBn446z7Stkz4t8qn8TB5iPEOXmzOc0lbVZnUIUp909SYYWk8U2RDT9+xHsmKZlOIofwwxPwF8StNdIE4Cmc76Z8rEcLbM2FnxeZIGrQSpK5z39OHA3ue9BZfnisUpBpdgO9KvmCEosahirk7OHpFh5auhve1mGyszfToq65HxDkqzAyILLoDwfkrM3Nkk6dkvvYg7plq/lCMH3tmoUsePu4+xeTLhaNDgvvgll9q8SVOZJc/fN9HwVlDG7Keq5CzcXVdHtkDEhoBVi1xQ8QlOb5+7SeW3W8NNrJt+dmaVJUkFGck/HFOngMNHiiIUGOqUBQLlI+ZWLs2jDfpu1sqE+CJND9dUXcauFnQLZ26I2Ib2bewx8TQVVT4WNA95C1+RSllUCy4RySXQSXec8bebGKbvIh+GvW5x6f9i8z7h6OEpr7AjeOImWPideDVDOhr9j4ftOJgXSOUG0WayjJXuaFiLeoI6FvLH+L0ro3tntZmjOI7UN0p6ZV6HuEtAn0A2sPxl28FIhi8JZ2CM0lHo1rWxSvery"
```

2. Run the binary as follow
```
./sqs-dlq-replayer -from=https://sqs.eu-west-1.amazonaws.com/453318629493/vf-cm-prod-marketplace-emea-bi-deadletter-queue -to=https://sqs.eu-west-1.amazonaws.com/453318629493/vf-cm-prod-marketplace-emea-bi-orders

```

## What happens when all messages have been pushed?

Once all the messages have been replayed successfully, check the DB to make sure the recrods have been successfully persisted. Then purge the DLQ and that will delete all the messages