package prism

type UserInterface interface {
	EntityInterface
}

type User struct {
	ID        int
	Name      string
	Anonymous bool
	Entity
	GithubUser
}
