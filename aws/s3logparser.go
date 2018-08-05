package aws

type s3LogParserMessageHandler func(*struct{})

// S3LogParser interface to inherit on each type of S3 log parsers
type S3LogParser interface {
	parse(*string, s3LogParserMessageHandler)
}

// TODO: decide if it should be a string or can be s streaming
