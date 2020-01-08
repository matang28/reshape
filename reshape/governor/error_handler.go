package governor

// Error Handler will be used to handle errors thrown while governor runs
// you should return false only if you want the governor to ignore that error.
// For example, the governor has limits on errors count, you can override this logic by returning false
// on your handler and handle the error the way you want.
type ErrorHandler = func(err error) bool

// The default handler returns true for any error reported
var defaultErrorHandler = func(err error) bool {
	if err == nil {
		return false
	}
	return true
}
