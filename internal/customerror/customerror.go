package customerror

var (
	ErrIDNotFound = New("ID not found", false)
	ErrIDExists   = New("ID exists", false)
	ErrUnknown    = New("Unknown error", true)
	ErrDelete     = New("Error while deleting", true)
	ErrSet        = New("Error while setting", true)
	ErrUpdate     = New("Error while updating", true)
	ErrGetAll     = New("Error while getting all", true)
)

type CustomError interface {
	Wrap(err error) CustomError
	Unwrap() error
	AddData(any) CustomError
	DestroyData() CustomError
	Error() string
}

type Error struct {
	Err      error
	Message  string
	Data     any `json:"-"`
	Loggable bool
}

func (e *Error) Wrap(err error) CustomError {
	e.Err = err
	return e
}

func (e *Error) Unwrap() error {
	return e.Err
}

func (e *Error) AddData(d any) CustomError {
	e.Data = d
	return e
}

func (e *Error) DestroyData() CustomError {
	e.Data = nil
	return e
}

func (e *Error) Error() string {
	if e.Err != nil {
		return e.Err.Error() + ", " + e.Message
	}
	return e.Message
}

func New(message string, l bool) CustomError {
	return &Error{
		Message:  message,
		Loggable: l,
	}
}
