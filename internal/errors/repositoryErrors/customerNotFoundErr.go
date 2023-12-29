package repositoryErrors

type CustomerNotFoundErr struct{}

const (
	msg = "There's no customer!"
)

func (e CustomerNotFoundErr) Error() string {
	return msg
}
