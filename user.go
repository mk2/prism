package prism

type User struct {
	EntityInterface
	ID        int
	Name      string
	Anonymous bool
	Entity
	GithubUser
}
