package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/nsqio/go-nsq"
	"github.com/segmentio/ksuid"
	"github.com/sirupsen/logrus"

	"github.com/mailbadger/app/consumers"
	"github.com/mailbadger/app/entities"
	"github.com/mailbadger/app/queue"
	"github.com/mailbadger/app/s3"
	"github.com/mailbadger/app/services/campaigns"
	"github.com/mailbadger/app/services/templates"
	"github.com/mailbadger/app/storage"
	"github.com/mailbadger/app/utils"
)

// MessageHandler implements the nsq handler interface.
type MessageHandler struct {
	s           storage.Storage
	templatesvc templates.Service
	p           queue.Producer
}

// HandleMessage is the only requirement needed to fulfill the
// nsq.Handler interface.
func (h *MessageHandler) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		logrus.Error("Empty message, unable to proceed.")
		return nil
	}

	msg := new(entities.CampaignerTopicParams)
	svc := campaigns.New(h.s, h.p)

	err := json.Unmarshal(m.Body, msg)
	if err != nil {
		logrus.WithField("body", string(m.Body)).WithError(err).Error("Malformed JSON message.")
		return nil
	}

	campaign, err := h.s.GetCampaign(msg.CampaignID, msg.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logrus.WithFields(logrus.Fields{
				"campaign_id": msg.CampaignID,
				"user_id":     msg.UserID,
			}).WithError(err).Warn("unable to find campaign")
			return nil
		}
		logrus.WithFields(logrus.Fields{
			"campaign_id": msg.CampaignID,
			"user_id":     msg.UserID,
		}).WithError(err).Error("unable to find campaign")
		return err
	}
	if campaign.Status != entities.StatusDraft {
		logrus.WithError(err).
			WithFields(logrus.Fields{
				"user_id":     msg.UserID,
				"campaign_id": msg.CampaignID,
			}).Errorf("potentially duplicate message: campaign status is '%s', it should be 'draft'", campaign.Status)
		return nil
	}

	campaign.StartedAt.SetValid(time.Now().UTC())
	campaign.Status = entities.StatusSending
	err = h.s.UpdateCampaign(campaign)
	if err != nil {
		logrus.WithError(err).
			WithFields(logrus.Fields{
				"user_id":     campaign.UserID,
				"campaign_id": campaign.ID,
				"status":      campaign.Status,
			}).Error("unable to update campaign")
		return nil
	}

	// fetching subs that are active and that have not been blacklisted
	var (
		timestamp time.Time
		nextID    int64
		limit     int64 = 1000
	)

	parsedTemplate, err := h.templatesvc.ParseTemplate(context.Background(), campaign.TemplateID, msg.UserID)
	if err != nil {
		logrus.WithError(err).
			WithFields(logrus.Fields{
				"user_id":     campaign.UserID,
				"template_id": campaign.TemplateID,
			}).Error("unable to prepare campaign template data")

		campaign.Status = entities.StatusFailed
		err = h.s.UpdateCampaign(campaign)
		if err != nil {
			logrus.WithError(err).
				WithFields(logrus.Fields{
					"user_id":     campaign.UserID,
					"campaign_id": campaign.ID,
					"status":      campaign.Status,
				}).Error("unable to update campaign")
			return nil
		}
		return nil
	}

	id := ksuid.New()

	for {
		subs, err := h.s.GetDistinctSubscribersBySegmentIDs(
			msg.SegmentIDs,
			msg.UserID,
			false,
			true,
			timestamp,
			nextID,
			limit,
		)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				logrus.WithFields(logrus.Fields{
					"user_id":     msg.UserID,
					"segment_ids": msg.SegmentIDs,
					"campaign_id": msg.CampaignID,
				}).WithError(err).Warn("unable to fetch subscribers.")
				return nil
			}
			logrus.WithFields(logrus.Fields{
				"user_id":     msg.UserID,
				"segment_ids": msg.SegmentIDs,
				"campaign_id": msg.CampaignID,
			}).WithError(err).Error("unable to fetch subscribers.")
			return nil
		}

		for _, s := range subs {
			params, err := svc.PrepareSubscriberEmailData(s, id, *msg, campaign.ID, parsedTemplate.HTMLPart, parsedTemplate.SubjectPart, parsedTemplate.TextPart)
			if err != nil {
				sendLog := &entities.SendLog{
					ID:           id,
					UserID:       msg.UserID,
					SubscriberID: s.ID,
					CampaignID:   msg.CampaignID,
					Status:       entities.SendLogStatusFailed,
					Description:  fmt.Sprintf("Failed to prepare subscriber email data error: %s", err),
				}
				err = h.s.CreateSendLog(sendLog)
				if err != nil {
					logrus.WithFields(logrus.Fields{
						"id":            id.String(),
						"user_id":       msg.UserID,
						"segment_ids":   msg.SegmentIDs,
						"campaign_id":   msg.CampaignID,
						"subscriber_id": s.ID,
					}).WithError(err).Error("unable to insert send logs for subscriber.")
				}
				return nil
			}

			err = svc.PublishSubscriberEmailParams(params)
			if err != nil {
				sendLog := &entities.SendLog{
					ID:           id,
					UserID:       msg.UserID,
					SubscriberID: s.ID,
					CampaignID:   msg.CampaignID,
					Status:       entities.SendLogStatusFailed,
					Description:  fmt.Sprintf("Failed to publish subscriber email data error: %s", err),
				}

				err = h.s.CreateSendLog(sendLog)
				if err != nil {
					logrus.WithFields(logrus.Fields{
						"id":            id.String(),
						"user_id":       msg.UserID,
						"segment_ids":   msg.SegmentIDs,
						"campaign_id":   msg.CampaignID,
						"subscriber_id": s.ID,
					}).WithError(err).Error("unable to insert send logs for subscriber.")
				}

				return nil
			}

			id = id.Next()
		}

		if len(subs) == 0 {
			break
		}

		// set  vars for next batches
		lastSub := subs[len(subs)-1]
		nextID = lastSub.ID
		timestamp = lastSub.CreatedAt

		// exit the loop if next batch is smaller then 1k -> it's the last batch of 1k.
		if len(subs) < 1000 {
			break
		}
	}

	campaign.Status = entities.StatusSent
	campaign.CompletedAt.SetValid(time.Now().UTC())

	err = h.s.UpdateCampaign(campaign)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"user_id":     msg.UserID,
			"segment_ids": msg.SegmentIDs,
			"campaign_id": msg.CampaignID,
		}).WithError(err).Error("unable to update campaign")
		return nil
	}

	return nil
}

func main() {
	lvl, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		lvl = logrus.InfoLevel
	}

	logrus.SetLevel(lvl)
	if utils.IsProductionMode() {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}
	logrus.SetOutput(os.Stdout)

	driver := os.Getenv("DATABASE_DRIVER")
	conf := storage.MakeConfigFromEnv(driver)
	s := storage.New(driver, conf)

	s3Client, err := s3.NewS3Client(
		os.Getenv("AWS_S3_ACCESS_KEY"),
		os.Getenv("AWS_S3_SECRET_KEY"),
		os.Getenv("AWS_S3_REGION"),
	)
	if err != nil {
		logrus.Fatal(err)
	}

	svc := templates.New(s, s3Client)

	p, err := queue.NewProducer(os.Getenv("NSQD_HOST"), os.Getenv("NSQD_PORT"))
	if err != nil {
		logrus.Fatal(err)
	}

	config := nsq.NewConfig()

	consumer, err := nsq.NewConsumer(entities.CampaignerTopic, entities.CampaignerTopic, config)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to create consumer")
	}

	consumer.ChangeMaxInFlight(200)

	consumer.SetLogger(
		&consumers.NoopLogger{},
		nsq.LogLevelError,
	)

	consumer.AddConcurrentHandlers(
		&MessageHandler{s, svc, p},
		20,
	)

	addr := fmt.Sprintf("%s:%s", os.Getenv("NSQLOOKUPD_HOST"), os.Getenv("NSQLOOKUPD_PORT"))
	nsqlds := []string{addr}

	logrus.Infoln("Connecting to NSQlookup...")
	if err := consumer.ConnectToNSQLookupds(nsqlds); err != nil {
		logrus.Fatal(err)
	}

	logrus.Infoln("Connected to NSQlookup")

	shutdown := make(chan os.Signal, 2)
	signal.Notify(shutdown, os.Interrupt)
	signal.Notify(shutdown, syscall.SIGINT)
	signal.Notify(shutdown, syscall.SIGTERM)

	for {
		select {
		case <-consumer.StopChan:
			return // consumer disconnected. Time to quit.
		case <-shutdown:
			// Synchronously drain the queue before falling out of main
			logrus.Infoln("Stopping consumer...")
			consumer.Stop()
		}
	}
}
