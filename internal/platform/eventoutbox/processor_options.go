package eventoutbox

func optionUnix(options PublishOptions) int64 {
	if options.Now == nil {
		return 0
	}
	return options.Now().Unix()
}
