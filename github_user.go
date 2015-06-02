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

func (u *GithubUser) saveGithubUser(tx *bolt.Tx) error {

	if u.user.Anonymous == "true" {
		return errors.New("Anonymous User!")
	}

	var b *bolt.Bucket

	b = tx.Bucket(s2b(AccessTokenToGithubUserIDBucket))
	b.Put(s2b(u.AccessToken), s2b(u.GithubUserID))

	b = tx.Bucket(s2b(GithubUserIDToUserIDBucket))
	b.Put(s2b(u.GithubUserID), s2b(u.user.ID))

	b = tx.Bucket(s2b(UserIDToGithubUserIDBucket))
	b.Put(s2b(u.user.ID), s2b(u.GithubUserID))

	b = tx.Bucket(s2b(GithubUserNameBucket))
	b.Put(s2b(u.user.ID), s2b(u.GithubUserName))

	return nil
}

func (u *GithubUser) loadGithubUser(tx *bolt.Tx) error {

	if u.user.Anonymous == "true" {
		return errors.New("Anonymous User!")
	}

	var b *bolt.Bucket

	b = tx.Bucket(s2b(UserIDToGithubUserIDBucket))
	u.GithubUserID = b2s(b.Get(s2b(u.user.ID)))

	b = tx.Bucket(s2b(UserNameBucket))
	u.GithubUserName = b2s(b.Get(s2b(u.user.ID)))

	b = tx.Bucket(s2b(GithubUserNameBucket))
	b.Put(s2b(u.user.ID), s2b(u.GithubUserName))

	return nil
}

func (u *GithubUser) ProvideAuth(accessToken string) {

	u.user.Anonymous = "false"
	u.AccessToken = accessToken

}
