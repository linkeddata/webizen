package webizen

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDBSync(t *testing.T) {
	err := db.Sync(
		new(User),
		new(UserName),
		new(UserImage),
		new(UserMbox),
	)
	assert.NoError(t, err)

	results, err := db.Query("SHOW TABLES")
	assert.NoError(t, err)
	assert.Equal(t, len(results), 4)
}

var (
	testUser = &User{
		Uri: "https://webid.mit.edu/presbrey#",
	}
	testName = &UserName{
		Name: "Test User",
	}
	testImage = &UserImage{
		Image: "https://0.gravatar.com/avatar/39e047043fbfdf600dfe0230d92c32e5",
	}
	testMbox = &UserMbox{
		Local:  "test",
		Domain: "test.com",
	}
)

func TestDBSave(t *testing.T) {
	n, err := db.InsertOne(testUser)
	assert.Equal(t, n, 1)
	assert.NoError(t, err)

	testName.User,
		testImage.User,
		testMbox.User =
		testUser.Id,
		testUser.Id,
		testUser.Id

	n, err = db.Insert(testName, testMbox, testImage)
	assert.Equal(t, n, 3)
	assert.NoError(t, err)
}

func TestDBSearch(t *testing.T) {
	res1 := make([]User, 0)
	err := db.Cols("id").Where("uri LIKE ?", `%`+testUser.Uri+`%`).Find(&res1)
	assert.NoError(t, err)
	assert.Equal(t, res1[0].Id, testUser.Id)

	res2 := make([]UserName, 0)
	err = db.Cols("user").Where("name LIKE ?", `%test%`).Find(&res2)
	assert.NoError(t, err)
	assert.Equal(t, res2[0].User, testUser.Id)

	res3 := make([]UserMbox, 0)
	err = db.Cols("user").Where("local LIKE ?", `%test%`).Find(&res3)
	assert.NoError(t, err)
	assert.Equal(t, res3[0].User, testUser.Id)
}

func TestDBDelete(t *testing.T) {
	var (
		n   int64
		err error
	)

	n, err = db.Delete(testUser)
	assert.Equal(t, n, 1)
	assert.NoError(t, err)

	n, err = db.Delete(testName)
	assert.Equal(t, n, 1)
	assert.NoError(t, err)

	n, err = db.Delete(testImage)
	assert.Equal(t, n, 1)
	assert.NoError(t, err)

	n, err = db.Delete(testMbox)
	assert.Equal(t, n, 1)
	assert.NoError(t, err)
}
