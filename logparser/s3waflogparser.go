package logparser

var (
	// S3WAFLogParser S3 WAF logs parser
	S3WAFLogParser = NewJSONLogParser("timestamp", mustKindFromString("timeUnixMilliseconds"))
)
