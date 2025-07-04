package metrics

const (
	MetricPostShortUrlMissingParam  = "PostShortUrl_MissingParam"
	MetricPostShortUrlError         = "PostShortUrl_Error"
	MetricPostShortUrlSuccess       = "PostShortUrl_Success"
	MetricPostShortUrlCreatedNew    = "PostShortUrl_CreatedNew"
	MetricPostShortUrlFoundExisting = "PostShortUrl_FoundExisting"

	MetricDeleteShortUrlMissingParam = "DeleteShortUrl_MissingParam"
	MetricDeleteShortUrlError        = "DeleteShortUrl_Error"
	MetricDeleteShortUrlSuccess      = "DeleteShortUrl_Success"

	MetricResolveShortUrlError    = "ResolveShortUrl_Error"
	MetricResolveShortUrlNotFound = "ResolveShortUrl_NotFound"
	MetricResolveShortUrlSuccess  = "ResolveShortUrl_Success"

	MetricPostShortUrlDuration    = "PostShortUrl_Duration"
	MetricDeleteShortUrlDuration  = "DeleteShortUrl_Duration"
	MetricResolveShortUrlDuration = "ResolveShortUrl_Duration"
)
