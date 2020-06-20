package contract

// EntityTrait contains basic functions shared by multiple different app components,
// such as logging.
type EntityTrait struct {
	Log Logger
}

// Entity creates a new EntityTrait instance.
func Entity(log Logger) EntityTrait {
	return EntityTrait{
		Log: log,
	}
}

// Error is an error handler.
func (e EntityTrait) Error(err error) {
	e.Log.Error(err)
	panic("AAAAAAAAAAAA")
}
