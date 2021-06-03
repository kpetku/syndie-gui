package database

import (
	"bytes"
	"image"
	"image/jpeg"
	"io/ioutil"
	"sync"

	"github.com/kpetku/libsyndie/syndieutil"
	"github.com/kpetku/syndie-core/data"
	bolt "go.etcd.io/bbolt"
)

type Database struct {
	ChanList map[string][]data.Message
	Avatars  map[string]*image.Image

	Channels []data.Channel
	Messages []data.Message
	sync.Mutex
}

func New() *Database {
	db := new(Database)
	db.Lock()
	db.ChanList = make(map[string][]data.Message)
	db.Avatars = make(map[string]*image.Image)
	db.Unlock()
	return db
}

func (db *Database) Open(path string) error {
	err := data.OpenDB(path)
	if err != nil {
		return err
	}
	return data.InitDB()
}

func (db *Database) GetAvatar(identhash string) image.Image {
	if _, ok := db.Avatars[identhash]; ok {
		return *db.Avatars[identhash]
	}
	return *db.Avatars["empty"]
}

func (db *Database) NameFromChanIdentHash(s string) string {
	for _, channel := range db.Channels {
		if channel.IdentHash == s {
			if channel.Name == "" {
				return "Anonymous"
			}
			return channel.Name
		}
	}
	return "Unknown"
}

func (db *Database) Reload() {
	db.ChanList = make(map[string][]data.Message)
	db.loadChannels()
	db.loadMessages()
	db.loadAvatars()
}

func (db *Database) loadChannels() {
	db.Lock()
	db.Channels = []data.Channel{}
	data.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("channels"))
		b.ForEach(func(k, v []byte) error {
			c := data.Channel{}
			c.Decode(v)
			cHash, _ := syndieutil.ChanHash(c.Identity)
			c.IdentHash = cHash
			db.Channels = append(db.Channels, c)
			return nil
		})
		return nil
	})
	db.Unlock()
}

func (db *Database) loadMessages() {
	db.Lock()
	db.Messages = []data.Message{}
	data.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("messages"))
		b.ForEach(func(k, v []byte) error {
			m := data.Message{}
			m.Decode(v)
			db.Messages = append(db.Messages, m)
			db.ChanList[m.TargetChannel] = append(db.ChanList[m.TargetChannel], m)
			return nil
		})
		return nil
	})
	db.Unlock()
}

func (db *Database) loadAvatars() {
	db.Lock()
	empty, _ := ioutil.ReadFile("resources/jpeg.jpg")
	emptyAvatar, _ := jpeg.Decode(bytes.NewReader(empty))
	db.Avatars = make(map[string]*image.Image)
	db.Avatars["empty"] = &emptyAvatar

	for _, channel := range db.Channels {
		if len(channel.Avatar) > 0 {
			avatar, _ := jpeg.Decode(bytes.NewReader(channel.Avatar))
			db.Avatars[channel.IdentHash] = &avatar
		} else {
			db.Avatars[channel.IdentHash] = &emptyAvatar
		}
	}
	db.Unlock()
}
