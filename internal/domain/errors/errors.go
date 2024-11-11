package errors

type ErrOpenFile struct{}

func (e ErrOpenFile) Error() string {
	return "open file error"
}

type ErrCloseFile struct{}

func (e ErrCloseFile) Error() string {
	return "close file error"
}

type ErrOpenURL struct{}

func (e ErrOpenURL) Error() string {
	return "open URL error"
}

type ErrCloseURL struct{}

func (e ErrCloseURL) Error() string {
	return "close URL error"
}

type ErrNoSource struct{}

func (e ErrNoSource) Error() string {
	return "no source error"
}

type ErrZeroTime struct{}

func (e ErrZeroTime) Error() string {
	return "zero time error"
}

type ErrTimeParsing struct{}

func (e ErrTimeParsing) Error() string {
	return "err parsing time"
}

type ErrWrongTimeBoundaries struct{}

func (e ErrWrongTimeBoundaries) Error() string { return "to before from" }

type ErrFileCreation struct{}

func (e ErrFileCreation) Error() string { return "file creation error" }

type ErrFileWrite struct{}

func (e ErrFileWrite) Error() string { return "file write error" }

type ErrInvalidURL struct{}

func (e ErrInvalidURL) Error() string { return "invalid URL" }

type ErrGetContentFromURL struct{}

func (e ErrGetContentFromURL) Error() string { return "Get content from URL error" }

type ErrNotOkHTTPAnswer struct{}

func (e ErrNotOkHTTPAnswer) Error() string { return "Not Ok HTTP Answer" }

type ErrNoDataWereProcessed struct{}

func (e ErrNoDataWereProcessed) Error() string {
	return "no data were processed"
}

type ErrSourceClosure struct{}

func (e ErrSourceClosure) Error() string {
	return "source closure error"
}

type ErrOutPut struct{}

func (e ErrOutPut) Error() string {
	return "output error"
}
