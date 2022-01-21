package breaker

type Config struct {
	// 熔断器 名称唯一表示
	Name string `json:"name"`
	//  执行command的超时时间 ms
	Timeout int `json:"timeout"`
	//  command的最大并发量
	MaxConcurrentRequests int `json:"max_concurrent_requests"`
	// 统计窗口10s内的请求数量，达到这个请求数量后才去判断是否要开启熔断
	RequestVolumeThreshold int `json:"request_volume_threshold"`
	// 当熔断器被打开后，SleepWindow的时间就是控制过多久后去尝试服务是否可用了
	SleepWindow int `json:"sleep_window"`
	// 错误百分比，请求数量大于等于RequestVolumeThreshold并且错误率到达这个百分比后就会启动熔断
	ErrorPercentThreshold int `json:"error_percent_threshold"`
}
