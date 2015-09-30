package pq

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"github.com/oursky/ourd/oddb"
)

var subscribeListenOnce sync.Once
var appEventChannelsMap map[string][]chan oddb.RecordEvent

// Assume all app resist on one Database
func (c *conn) Subscribe(recordEventChan chan oddb.RecordEvent) error {
	channels := appEventChannelsMap[c.appName]
	appEventChannelsMap[c.appName] = append(channels, recordEventChan)

	// TODO(limouren): Seems a start-up time config would be better?
	subscribeListenOnce.Do(func() {
		go newRecordListener(c.option).Listen()
	})

	return nil
}

func emit(n *notification) {
	channels := appEventChannelsMap[n.AppName]
	for _, channel := range channels {
		go func(ch chan oddb.RecordEvent) {
			ch <- oddb.RecordEvent{
				Record: &n.Record,
				Event:  n.ChangeEvent,
			}
		}(channel)
	}
}

// the channel to listen for record changes
const recordChangeChannel = "record_change"

type notification struct {
	AppName     string
	ChangeEvent oddb.RecordHookEvent
	Record      oddb.Record
}

type rawNotification struct {
	AppName    string
	Op         string
	RecordType string
	Record     []byte
}

type recordListener struct {
	option string
	db     *sqlx.DB
}

func newRecordListener(option string) *recordListener {
	return &recordListener{
		option: option,
		db:     sqlx.MustOpen("postgres", option),
	}
}

func (l *recordListener) Listen() {
	eventCallback := func(event pq.ListenerEventType, err error) {
		if err != nil {
			log.WithField("err", err).Errorf("pq/listener: Received an error")
		} else {
			log.WithField("event", event).Infof("pq/listener: Received an event")
		}
	}

	listener := pq.NewListener(
		l.option,
		10*time.Second,
		time.Minute,
		eventCallback)

	if err := listener.Listen(recordChangeChannel); err != nil {
		log.WithFields(log.Fields{
			"channel": recordChangeChannel,
			"err":     err,
		}).Errorln("pq/listener: got an err while trying to listen")
		return
	}

	log.Infof("pq/listener: Listening to %s...", recordChangeChannel)

	for {
		select {
		case pqNotification := <-listener.Notify:
			log.WithField("pqNotification", pqNotification).Infoln("Received a notify")

			n := notification{}
			if err := l.fetchNotification(pqNotification.Extra, &n); err != nil {
				log.WithFields(log.Fields{
					"pqNotification": pqNotification,
					"err":            err,
				}).Errorln("pq/listener: failed to fetch notification")

				continue
			}

			emit(&n)

			l.deleteNotification(pqNotification.Extra)
		case <-time.After(60 * time.Second):
			go func() {
				if err := listener.Ping(); err != nil {
					log.WithField("err", err).Errorln("pq/listener: got an err while pinging connection")
				}
			}()
		}
	}
}

// NOTE(limouren): pending_notification.id is integer in database.
func (l *recordListener) fetchNotification(notificationID string, n *notification) error {
	var rawNoti rawNotification
	err := l.db.QueryRowx("SELECT op, appname, recordtype, record FROM pending_notification WHERE id = $1", notificationID).
		StructScan(&rawNoti)
	if err != nil {
		log.WithFields(log.Fields{
			"notificationID": notificationID,
			"err":            err,
		}).Errorln("Failed to fetch pending notification")
		return err
	}

	if err := parseNotification(&rawNoti, n); err != nil {
		return err
	}

	return nil
}

func (l *recordListener) deleteNotification(notificationID string) {
	result, err := l.db.Exec("DELETE FROM pending_notification WHERE id = $1", notificationID)
	if err != nil {
		log.WithFields(log.Fields{
			"notificationID": notificationID,
			"err":            err,
		}).Errorln("Failed to delete notification")

		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.WithFields(log.Fields{
			"notificationID": notificationID,
			"err":            err,
			"rowsAffected":   rowsAffected,
		}).Errorln("More than one notification deleted")

		return
	}

	if rowsAffected != 1 {
		log.WithFields(log.Fields{
			"notificationID": notificationID,
			"rowsAffected":   rowsAffected,
		}).Errorln("Zero or more than one notification deleted")
	}
}

func parseNotification(raw *rawNotification, n *notification) error {
	if err := parseAppName(raw.AppName, &n.AppName); err != nil {
		return err
	}

	if err := parseChangeEvent(raw.Op, &n.ChangeEvent); err != nil {
		return err
	}

	if err := parseRecordData(raw.Record, &n.Record); err != nil {
		return err
	}
	n.Record.ID.Type = raw.RecordType

	return nil
}

func parseAppName(rawAppName string, appName *string) error {
	if !strings.HasPrefix(rawAppName, "app_") {
		return fmt.Errorf("Invalid AppName = %v", rawAppName)
	}
	*appName = rawAppName[4:]
	return nil
}

func parseChangeEvent(rawOp string, changeEvent *oddb.RecordHookEvent) error {
	switch rawOp {
	case "INSERT":
		*changeEvent = oddb.RecordCreated
	case "UPDATE":
		*changeEvent = oddb.RecordUpdated
	case "DELETE":
		*changeEvent = oddb.RecordDeleted
	default:
		return fmt.Errorf("Unrecongized Op = %v", rawOp)
	}

	return nil
}

func parseRecordData(data []byte, record *oddb.Record) error {
	recordData := map[string]interface{}{}
	if err := json.Unmarshal(data, &recordData); err != nil {
		return fmt.Errorf("invalid json: %v", err)
	}

	recordID, _ := recordData["_id"].(string)
	rawDatabaseID, _ := recordData["_database_id"].(string)
	rawOwnerID, _ := recordData["_owner_id"].(string)

	if recordID == "" || rawDatabaseID == "" || rawOwnerID == "" {
		return errors.New(`missing key "_id", "_database_id" or "_owner_id"`)
	}

	for key := range recordData {
		if key[0] == '_' {
			delete(recordData, key)
		}
	}

	record.ID.Key = recordID
	record.Data = recordData
	record.DatabaseID = rawDatabaseID
	record.OwnerID = rawOwnerID

	return nil
}

func init() {
	appEventChannelsMap = map[string][]chan oddb.RecordEvent{}
}
