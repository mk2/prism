package prism

import (
	"errors"

	"github.com/boltdb/bolt"
)

const (
	AccessTokenToGithubUserIDBucket = "AccessTokenToGithubUserIDBucket"
	GithubUserIDToUserIDBucket      = "GithubUserIDToUserIDBucket"
	UserIDToGithubUserIDBucket      = "GithubUserIDToUserIDBucket"
	GithubUserNameBucket            = "GithubUserNameBucket"
)

type GithubUser struct {
	user           *User
	GithubUserID   string
	GithubUserName string
	AccessToken    string
}

func LoadGithubUser(db *bolt.DB, accessToken string) *User {

	var userID, githubUserID string

	db.View(func(tx *bolt.Tx) error {

		var b *bolt.Bucket

		b = tx.Bucket(s2b(AccessTokenToGithubUserIDBucket))
		githubUserID = b2s(b.Get(s2b(accessToken)))

		b = tx.Bucket(s2b(GithubUserIDToUserIDBucket))
		userID = b2s(b.Get(s2b(githubUserID)))

		return nil
	})

	return LoadUser(db, userID)
}

func (u *GithubUser) saveGithubUser(tx *bolt.Tx) error {

	if u.user.anonymous == "true" {
		return errors.New("Anonymous User!")
	}

	var b *bolt.Bucket

	b = tx.Bucket(s2b(AccessTokenToGithubUserIDBucket))
	b.Put(s2b(u.AccessToken), s2b(u.GithubUserID))

	b = tx.Bucket(s2b(GithubUserIDToUserIDBucket))
	b.Put(s2b(u.GithubUserID), s2b(u.user.id))

	b = tx.Bucket(s2b(UserIDToGithubUserIDBucket))
	b.Put(s2b(u.user.id), s2b(u.GithubUserID))

	b = tx.Bucket(s2b(GithubUserNameBucket))
	b.Put(s2b(u.user.id), s2b(u.GithubUserName))

	return nil
}

func (u *GithubUser) loadGithubUser(tx *bolt.Tx) error {

	if u.user.anonymous == "true" {
		return errors.New("Anonymous User!")
	}

	var b *bolt.Bucket

	b = tx.Bucket(s2b(UserIDToGithubUserIDBucket))
	u.GithubUserID = b2s(b.Get(s2b(u.user.id)))

	b = tx.Bucket(s2b(GithubUserNameBucket))
	u.GithubUserName = b2s(b.Get(s2b(u.user.id)))

	return nil
}

func (u *GithubUser) provideGithubAuth(accessToken string) {

	u.user.anonymous = "false"
	u.AccessToken = accessToken

}
