package logging

func BrandscreenLogger(name string, filename string) (*Logger, error) {
	return FileLogger(
		name,
		DEBUG,
		"%s [%s] %s[%d] %s:%d:%s: %s\ntime,levelname,name,process,filename,lineno,funcname,message",
		"2006-01-02 15:04:05.000000",
		filename,
		false)
}
