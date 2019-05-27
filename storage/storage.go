package storage

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/news-maily/app/entities"
	"github.com/news-maily/app/utils/pagination"
)

const key = "storage"

// Storage is the central interface for accessing and
// writing data in the datastore.
type Storage interface {
	GetUser(int64) (*entities.User, error)
	GetUserByUUID(string) (*entities.User, error)
	GetUserByUsername(string) (*entities.User, error)
	GetActiveUserByUsername(string) (*entities.User, error)
	CreateUser(*entities.User) error
	UpdateUser(*entities.User) error

	GetCampaigns(int64, *pagination.Pagination)
	GetCampaign(int64, int64) (*entities.Campaign, error)
	GetCampaignByName(name string, userID int64) (*entities.Campaign, error)
	GetCampaignsByTemplateName(string, int64) ([]entities.Campaign, error)
	CreateCampaign(*entities.Campaign) error
	UpdateCampaign(*entities.Campaign) error
	DeleteCampaign(int64, int64) error

	GetLists(int64, *pagination.Pagination)
	GetListsByIDs(userID int64, ids []int64) ([]entities.List, error)
	GetList(int64, int64) (*entities.List, error)
	GetListByName(name string, userID int64) (*entities.List, error)
	CreateList(*entities.List) error
	UpdateList(*entities.List) error
	DeleteList(int64, int64) error
	AppendSubscribers(*entities.List) error
	DetachSubscribers(*entities.List) error

	GetSubscribers(int64, *pagination.Pagination)
	GetSubscribersByListID(int64, int64, *pagination.Pagination)
	GetSubscriber(int64, int64) (*entities.Subscriber, error)
	GetSubscribersByIDs([]int64, int64) ([]entities.Subscriber, error)
	GetSubscriberByEmail(string, int64) (*entities.Subscriber, error)
	GetAllSubscribersByListID(listID, userID int64) ([]entities.Subscriber, error)
	GetDistinctSubscribersByListIDs(listIDs []int64, userID int64, blacklisted, active bool, nextID, limit int64) ([]entities.Subscriber, error)
	CreateSubscriber(*entities.Subscriber) error
	UpdateSubscriber(*entities.Subscriber) error
	BlacklistSubscriber(userID int64, email string) error
	DeleteSubscriber(int64, int64) error

	GetSesKeys(userID int64) (*entities.SesKeys, error)
	CreateSesKeys(s *entities.SesKeys) error
	DeleteSesKeys(userID int64) error

	CreateSendBulkLog(l *entities.SendBulkLog) error
	CountLogsByUUID(uuid string) (int, error)

	CreateBounce(b *entities.Bounce) error
	CreateComplaint(c *entities.Complaint) error
	CreateClick(c *entities.Click) error
	CreateOpen(o *entities.Open) error
	CreateDelivery(d *entities.Delivery) error
}

// SetToContext sets the storage to the context
func SetToContext(c *gin.Context, storage Storage) {
	c.Set(key, storage)
}

// GetFromContext returns the Storage associated with the context
func GetFromContext(c context.Context) Storage {
	return c.Value(key).(Storage)
}

// GetUser returns a User entity from the specified id.
func GetUser(c context.Context, id int64) (*entities.User, error) {
	return GetFromContext(c).GetUser(id)
}

// GetUserByUUID returns a User entity from the specified uuid.
func GetUserByUUID(c context.Context, uuid string) (*entities.User, error) {
	return GetFromContext(c).GetUserByUUID(uuid)
}

// GetUserByUsername returns a User entity from the specified username.
func GetUserByUsername(c context.Context, username string) (*entities.User, error) {
	return GetFromContext(c).GetUserByUsername(username)
}

// GetActiveUserByUsername returns an active User entity from the specified username.
func GetActiveUserByUsername(c context.Context, username string) (*entities.User, error) {
	return GetFromContext(c).GetActiveUserByUsername(username)
}

// CreateUser persists a new User entity in the datastore.
func CreateUser(c context.Context, user *entities.User) error {
	return GetFromContext(c).CreateUser(user)
}

// UpdateUser updates the User entity.
func UpdateUser(c context.Context, user *entities.User) error {
	return GetFromContext(c).UpdateUser(user)
}

// GetCampaigns populates a pagination object with a collection of
// campaigns by the specified user id.
func GetCampaigns(c context.Context, userID int64, p *pagination.Pagination) {
	GetFromContext(c).GetCampaigns(userID, p)
}

// GetCampaign returns a Campaign entity by the given id and user id.
func GetCampaign(c context.Context, id, userID int64) (*entities.Campaign, error) {
	return GetFromContext(c).GetCampaign(id, userID)
}

// GetCampaignByName returns a Campaign entity by the given name and user id.
func GetCampaignByName(c context.Context, name string, userID int64) (*entities.Campaign, error) {
	return GetFromContext(c).GetCampaignByName(name, userID)
}

// GetCampaignsByTemplateName returns a collection of campaigns by the given template name and user id.
func GetCampaignsByTemplateName(c context.Context, templateName string, userID int64) ([]entities.Campaign, error) {
	return GetFromContext(c).GetCampaignsByTemplateName(templateName, userID)
}

// CreateCampaign persists a new Campaign entity in the datastore.
func CreateCampaign(c context.Context, campaign *entities.Campaign) error {
	return GetFromContext(c).CreateCampaign(campaign)
}

// UpdateCampaign updates a Campaign entity.
func UpdateCampaign(c context.Context, campaign *entities.Campaign) error {
	return GetFromContext(c).UpdateCampaign(campaign)
}

// DeleteCampaign deletes a Campaign entity by the given id.
func DeleteCampaign(c context.Context, id, userID int64) error {
	return GetFromContext(c).DeleteCampaign(id, userID)
}

// GetLists populates a pagination object with a collection of
// lists by the specified user id.
func GetLists(c context.Context, userID int64, p *pagination.Pagination) {
	GetFromContext(c).GetLists(userID, p)
}

// GetListsByIDs fetches lists by user id and the given ids
func GetListsByIDs(c context.Context, userID int64, ids []int64) ([]entities.List, error) {
	return GetFromContext(c).GetListsByIDs(userID, ids)
}

// GetList returns a List entity by the given id and user id.
func GetList(c context.Context, id, userID int64) (*entities.List, error) {
	return GetFromContext(c).GetList(id, userID)
}

// GetListByName returns a Campaign entity by the given name and user id.
func GetListByName(c context.Context, name string, userID int64) (*entities.List, error) {
	return GetFromContext(c).GetListByName(name, userID)
}

// CreateList persists a new List entity in the datastore.
func CreateList(c context.Context, l *entities.List) error {
	return GetFromContext(c).CreateList(l)
}

// UpdateList updates a List entity.
func UpdateList(c context.Context, l *entities.List) error {
	return GetFromContext(c).UpdateList(l)
}

// DeleteList deletes a List entity by the given id.
func DeleteList(c context.Context, id, userID int64) error {
	return GetFromContext(c).DeleteList(id, userID)
}

// AppendSubscribers appends subscribers to the existing association.
func AppendSubscribers(c context.Context, l *entities.List) error {
	return GetFromContext(c).AppendSubscribers(l)
}

// DetachSubscribers deletes subscribers from the list.
func DetachSubscribers(c context.Context, l *entities.List) error {
	return GetFromContext(c).DetachSubscribers(l)
}

// GetSubscribers populates a pagination object with a collection of
// subscribers by the specified user id.
func GetSubscribers(c context.Context, userID int64, p *pagination.Pagination) {
	GetFromContext(c).GetSubscribers(userID, p)
}

// GetSubscribersByListID populates a pagination object with a collection of
// subscribers by the specified user id and list id.
func GetSubscribersByListID(c context.Context, listID, userID int64, p *pagination.Pagination) {
	GetFromContext(c).GetSubscribersByListID(listID, userID, p)
}

// GetSubscriber returns a Subscriber entity by the given id and user id.
func GetSubscriber(c context.Context, id, userID int64) (*entities.Subscriber, error) {
	return GetFromContext(c).GetSubscriber(id, userID)
}

// GetSubscribersByIDs returns a Subscriber entity by the given ids and user id.
func GetSubscribersByIDs(c context.Context, ids []int64, userID int64) ([]entities.Subscriber, error) {
	return GetFromContext(c).GetSubscribersByIDs(ids, userID)
}

// GetSubscriberByEmail returns a Subscriber entity by the given email and user id.
func GetSubscriberByEmail(c context.Context, email string, userID int64) (*entities.Subscriber, error) {
	return GetFromContext(c).GetSubscriberByEmail(email, userID)
}

// GetAllSubscribersByListID fetches all subscribers by user id and list id
func GetAllSubscribersByListID(c context.Context, listID, userID int64) ([]entities.Subscriber, error) {
	return GetFromContext(c).GetAllSubscribersByListID(listID, userID)
}

// GetDistinctSubscribersByListIDs fetches all distinct subscribers by user id and list ids
func GetDistinctSubscribersByListIDs(
	c context.Context,
	listIDs []int64,
	userID int64,
	blacklisted, active bool,
	nextID, limit int64,
) ([]entities.Subscriber, error) {
	return GetFromContext(c).GetDistinctSubscribersByListIDs(listIDs, userID, blacklisted, active, nextID, limit)
}

// CreateSubscriber persists a new Subscriber entity in the datastore.
func CreateSubscriber(c context.Context, s *entities.Subscriber) error {
	return GetFromContext(c).CreateSubscriber(s)
}

// UpdateSubscriber updates a Subscriber entity.
func UpdateSubscriber(c context.Context, s *entities.Subscriber) error {
	return GetFromContext(c).UpdateSubscriber(s)
}

// BlacklistSubscriber blacklists a Subscriber entity by the given email.
func BlacklistSubscriber(c context.Context, userID int64, email string) error {
	return GetFromContext(c).BlacklistSubscriber(userID, email)
}

// DeleteSubscriber deletes a Subscriber entity by the given id.
func DeleteSubscriber(c context.Context, id, userID int64) error {
	return GetFromContext(c).DeleteSubscriber(id, userID)
}

// GetSesKeys returns the SES keys by the given user id
func GetSesKeys(c context.Context, userID int64) (*entities.SesKeys, error) {
	return GetFromContext(c).GetSesKeys(userID)
}

// CreateSesKeys adds new SES keys in the database.
func CreateSesKeys(c context.Context, s *entities.SesKeys) error {
	return GetFromContext(c).CreateSesKeys(s)
}

// DeleteSesKeys deletes SES keys configuration by the given user ID.
func DeleteSesKeys(c context.Context, userID int64) error {
	return GetFromContext(c).DeleteSesKeys(userID)
}

// CreateBounce adds new bounce in the database.
func CreateBounce(c context.Context, b *entities.Bounce) error {
	return GetFromContext(c).CreateBounce(b)
}

// CreateComplaint adds new complaint in the database.
func CreateComplaint(c context.Context, compl *entities.Complaint) error {
	return GetFromContext(c).CreateComplaint(compl)
}

// CreateClick adds new click in the database.
func CreateClick(c context.Context, click *entities.Click) error {
	return GetFromContext(c).CreateClick(click)
}

// CreateOpen adds new open in the database.
func CreateOpen(c context.Context, open *entities.Open) error {
	return GetFromContext(c).CreateOpen(open)
}

// CreateDelivery adds new delivery in the database.
func CreateDelivery(c context.Context, d *entities.Delivery) error {
	return GetFromContext(c).CreateDelivery(d)
}
