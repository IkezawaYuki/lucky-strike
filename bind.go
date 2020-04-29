package echo

type (
	Binder interface {
		Bind(i interface{}, c Context) error
	}
	DefaultBinder   struct{}
	BindUnmarshaler interface {
		UnmarshalParam(param string) error
	}
)

func (b *DefaultBinder) Bind(i interface{}, c Context) (err error) {

}
