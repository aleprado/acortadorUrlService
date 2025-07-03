package metrics

const (
	MetricGetShortUrlMissingParam  = "GetShortUrl_MissingParam"
	MetricGetShortUrlError         = "GetShortUrl_Error"
	MetricGetShortUrlSuccess       = "GetShortUrl_Success"
	MetricGetShortUrlCreatedNew    = "GetShortUrl_CreatedNew"
	MetricGetShortUrlFoundExisting = "GetShortUrl_FoundExisting"

	MetricDeleteShortUrlMissingParam = "DeleteShortUrl_MissingParam"
	MetricDeleteShortUrlError        = "DeleteShortUrl_Error"
	MetricDeleteShortUrlSuccess      = "DeleteShortUrl_Success"

	MetricResolveShortUrlError    = "ResolveShortUrl_Error"
	MetricResolveShortUrlNotFound = "ResolveShortUrl_NotFound"
	MetricResolveShortUrlSuccess  = "ResolveShortUrl_Success"

	MetricGetShortUrlDuration     = "GetShortUrl_Duration"
	MetricDeleteShortUrlDuration  = "DeleteShortUrl_Duration"
	MetricResolveShortUrlDuration = "ResolveShortUrl_Duration"
)
