package prism

type UserInterface interface {
	EntityInterface
}

type User struct {
	UserInterface
	ID        int
	Name      string
	Anonymous bool
	Entity
	GithubUser
}
