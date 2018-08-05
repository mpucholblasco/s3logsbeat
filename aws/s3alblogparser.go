package aws

// S3ObjectSQSMessage contains information of S3 objects (sqsMessage not
// null implies that this object is extracted from an SQS message)
type S3ALBLogParser struct{}

func NewS3ALBLogParser() *S3ALBLogParser {
	return &S3ALBLogParser{}
}

// ExtractNewS3Objects extracts those new S3 objects present on an SQS message
func (l *S3ALBLogParser) parse(logs *string, mh s3LogParserMessageHandler) error {
	return nil
}
