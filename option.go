package genius

type Opt func(*Option)

type Option struct {
	delimiter string
}

func WithDelimiter(delimiter string) Opt {
	return func(opt *Option) {
		opt.delimiter = delimiter
	}
}

func defaultOption(option *Option) {
	if option.delimiter == "" {
		option.delimiter = "."
	}
}
