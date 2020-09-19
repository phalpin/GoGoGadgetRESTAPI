package pherr

type options struct {
	errorType ErrorType
}

type Option = func(opts *options)

func WithErrorType(errType ErrorType) Option {
	return func(opts *options) {
		opts.errorType = errType
	}
}
