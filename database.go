package main

import (
	"bytes"
	"image"
	"image/jpeg"
	"io/ioutil"

	"github.com/boltdb/bolt"

	"github.com/kpetku/libsyndie/syndieutil"
	"github.com/kpetku/syndie-core/data"
)

type database struct {
	chanList map[string][]data.Message
	avatars  map[string]*image.Image

	Channels []data.Channel
	Messages []data.Message
}

func NewDatabase() *database {
	return new(database)
}

func (db *database) openDB(path string) error {
	db.chanList = make(map[string][]data.Message)
	err := data.OpenDB(path)
	if err != nil {
		return err
	}
	return data.InitDB()
}

func (db *database) loadChannels() {
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
}

func (db *database) loadMessages() {
	data.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("messages"))
		b.ForEach(func(k, v []byte) error {
			m := data.Message{}
			m.Decode(v)
			db.Messages = append(db.Messages, m)
			db.chanList[m.TargetChannel] = append(db.chanList[m.TargetChannel], m)
			return nil
		})
		return nil
	})
}

func (db *database) loadAvatars() {
	db.avatars = make(map[string]*image.Image)

	empty, _ := ioutil.ReadFile("resources/jpeg.jpg")
	emptyAvatar, _ := jpeg.Decode(bytes.NewReader(empty))
	db.avatars["empty"] = &emptyAvatar

	for _, channel := range db.Channels {
		if len(channel.Avatar) > 0 {
			avatar, _ := jpeg.Decode(bytes.NewReader(channel.Avatar))
			db.avatars[channel.IdentHash] = &avatar
		} else {
			db.avatars[channel.IdentHash] = &emptyAvatar
		}
	}
}

func (db *database) getAvatar(identhash string) image.Image {
	if _, ok := db.avatars[identhash]; ok {
		return *db.avatars[identhash]
	}
	return *db.avatars["empty"]
}