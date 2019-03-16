package beater

import "flag"

var (
	once            = flag.Bool("once", false, "Run s3logsbeat only once until all inputs will be read")
	keepSQSMessages = flag.Bool("keepsqsmessages", false, "Do not delete SQS messages when processed (set for testing)")
)
