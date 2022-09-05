package views

import (
	"io"
	"os"
)

//counterfeiter:generate -o fakes/ . ViewEngineInterface
type ViewEngineInterface interface {
	Draw(view ViewInterface) error
}

//counterfeiter:generate -o fakes/ . ViewInterface
type ViewInterface interface {
	Draw(out io.Writer) error
	SetData(data interface{})
	Data() interface{}
}

type StdOutViewEngine struct {
}

func (ve *StdOutViewEngine) Draw(view ViewInterface) error {
	return view.Draw(os.Stdout)
}
