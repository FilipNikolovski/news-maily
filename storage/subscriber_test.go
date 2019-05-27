package storage

import (
	"testing"

	"github.com/news-maily/app/entities"
	"github.com/news-maily/app/utils/pagination"
	"github.com/stretchr/testify/assert"
)

func TestSubscriber(t *testing.T) {
	db := openTestDb()
	defer db.Close()

	store := From(db)

	l := &entities.List{
		Name:   "foo",
		UserId: 1,
	}

	err := store.CreateList(l)
	assert.Nil(t, err)

	//Test create subscriber
	s := &entities.Subscriber{
		Name:   "foo",
		Email:  "john@example.com",
		UserId: 1,
		Metadata: []entities.SubscriberMetadata{
			{Key: "key", Value: "val"},
		},
		Blacklisted: false,
		Active:      true,
	}
	s.Lists = append(s.Lists, *l)

	err = store.CreateSubscriber(s)
	assert.Nil(t, err)

	//Test get subscriber
	s, err = store.GetSubscriber(s.Id, 1)
	assert.Nil(t, err)
	assert.Equal(t, s.Name, "foo")
	assert.NotEmpty(t, s.Metadata)
	assert.Equal(t, s.Metadata[0].Key, "key")
	assert.Equal(t, s.Metadata[0].Value, "val")

	//Test get subscriber by email
	s, err = store.GetSubscriberByEmail("john@example.com", 1)
	assert.Nil(t, err)
	assert.Equal(t, s.Name, "foo")

	//Test update subscriber
	s.Name = "bar"
	err = store.UpdateSubscriber(s)
	assert.Nil(t, err)
	assert.Equal(t, s.Name, "bar")

	//Test subscriber validation when name and email are invalid
	s.Name = ""
	s.Email = "foo bar"
	s.Validate()
	assert.Equal(t, s.Errors["name"], entities.ErrSubscriberNameEmpty.Error())
	assert.Equal(t, s.Errors["email"], entities.ErrEmailInvalid.Error())

	//Test get subs
	p := &pagination.Pagination{PerPage: 10}
	store.GetSubscribers(1, p)
	assert.NotEmpty(t, p.Collection)
	assert.Equal(t, len(p.Collection), int(p.Total))

	//Test get subs by ids
	subs, err := store.GetSubscribersByIDs([]int64{1}, 1)
	assert.Nil(t, err)
	assert.NotEmpty(t, subs)

	//Test get subs by list id
	p = &pagination.Pagination{PerPage: 10}
	store.GetSubscribersByListID(l.Id, 1, p)
	assert.NotEmpty(t, p.Collection)
	assert.Equal(t, len(p.Collection), int(p.Total))

	subs, err = store.GetAllSubscribersByListID(l.Id, 1)
	assert.Equal(t, 1, len(subs))

	subs, err = store.GetDistinctSubscribersByListIDs([]int64{l.Id}, 1, false, true, 0, 10)
	assert.Equal(t, 1, len(subs))

	err = store.DeleteSubscriber(1, 1)
	assert.Nil(t, err)
}
