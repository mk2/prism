package prism

import "github.com/boltdb/bolt"

const (
	UserAuthProviderGithub = "UserAuthProviderGithub"
)

const (
	UserIDBucket        = "UserIDBucket"
	UserAnonymousBucket = "UserAnonymousBucket"
	UserNameBucket      = "UserNameBucket"
	UserCreatedBucket   = "UserCreatedBucket"
	UserUpdatedBucket   = "UserUpdatedBucket"
)

/*
AutoProvider type alias for AuthProvider
*/
type AuthProvider string

type UserInterface interface {
	EntityIface
}

type User struct {
	Entity
	GithubUser

	name             string
	userAuthProvider string
	anonymous        string
}

func CreateUserBuckets(db *bolt.DB) error {

	tx, err := db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	requiredBuckets := []string{
		UserIDBucket,
		UserAnonymousBucket,
		UserNameBucket,
		UserCreatedBucket,
		UserUpdatedBucket,
		AccessTokenToGithubUserIDBucket,
		GithubUserIDToUserIDBucket,
		UserIDToGithubUserIDBucket,
		GithubUserNameBucket,
	}

	for _, requiredBucket := range requiredBuckets {
		_, err := tx.CreateBucketIfNotExists(s2b(requiredBucket))

		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func DeleteUserBuckets(db *bolt.DB) error {

	tx, err := db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	requiredBuckets := []string{
		UserIDBucket,
		UserAnonymousBucket,
		UserNameBucket,
		UserCreatedBucket,
		UserUpdatedBucket,
	}

	for _, requiredBucket := range requiredBuckets {
		err := tx.DeleteBucket(s2b(requiredBucket))

		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func NewUser(db *bolt.DB) *User {

	var u User

	db.Update(func(tx *bolt.Tx) error {

		u.newUserID(tx)

		created := u.Created(tx, UserCreatedBucket)
		updated := u.Updated(tx, UserUpdatedBucket)

		dbg.Printf("Created: %v", created)
		dbg.Printf("Updated: %v", updated)

		return nil

	})

	return &u
}

func LoadUser(db *bolt.DB, id string) *User {

	var u User
	u.id = id

	db.View(func(tx *bolt.Tx) error {

		switch u.userAuthProvider {

		case UserAuthProviderGithub:
			u.loadGithubUser(tx)

		}

		return nil
	})

	return &u
}

func (u *User) SaveUser(db *bolt.DB) error {

	return db.Batch(func(tx *bolt.Tx) error {

		switch u.userAuthProvider {

		case UserAuthProviderGithub:
			u.saveGithubUser(tx)

		}

		return nil
	})
}

func (u *User) newUserID(tx *bolt.Tx) error {

	return u.newID(tx, UserIDBucket, "userID")
}
