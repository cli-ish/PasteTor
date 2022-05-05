package logic

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/go-redis/redis"
	"math/rand"
	redisPool "pastetor/redis"
	"regexp"
	"strings"
)

type Note struct {
	Data  string
	State int
	Id    string
}

const notePreField = "note:"
const noteDataField = ":data"
const noteReportedField = ":reported"
const noteAllowedField = ":allowed"

var idvalidregex = regexp.MustCompile("^[\\da-z]+$")

func generateId() string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	s := make([]rune, 64)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	h := sha256.New()
	h.Write([]byte(string(s)))
	return hex.EncodeToString(h.Sum(nil))
}

func AddNote(data string, retry int) (string, error) {
	if retry >= 3 {
		return "", errors.New("max retries reached")
	}
	if data == "" {
		return "", errors.New("empty text")
	}
	pool := redisPool.GetPool()
	rdb, i := pool.GetClient()
	id := ""
	for {
		id = generateId()
		_, err := GetNote(id, 0)
		if err == redis.Nil { // Exit when the id does not exist.
			break
		}
	}
	setm := rdb.Set(notePreField+id+noteDataField, data, 0)
	if setm.Err() != nil {
		pool.Reconnect(i)
		return AddNote(data, retry+1)
	}
	return id, nil
}

func GetNote(id string, retry int) (Note, error) {
	if retry >= 3 {
		return Note{}, errors.New("max retries reached")
	}
	if !IsValidId(id) {
		return Note{}, errors.New("invalid node id")
	}
	pool := redisPool.GetPool()
	rdb, i := pool.GetClient()
	key := notePreField + id
	data := rdb.MGet(key+noteDataField, key+noteReportedField, key+noteAllowedField)
	if data.Err() != nil {
		if data.Err() != redis.Nil {
			pool.Reconnect(i)
			return GetNote(id, retry+1)
		}
		return Note{}, data.Err() // Exit anyway.
	}
	values := data.Val()
	if values[0] == nil {
		return Note{}, redis.Nil
	}
	state := 0
	if values[1] != nil {
		state = 1
	}
	if values[2] != nil {
		state = 2
	}
	return Note{values[0].(string), state, id}, nil
}

func ReportNote(note Note) error {
	pool := redisPool.GetPool()
	rdb, _ := pool.GetClient()
	rdb.Set(notePreField+note.Id+noteReportedField, 1, 0)
	return nil
}

func AllowNote(note Note, retry int) error {
	if retry >= 3 {
		return errors.New("max retries reached")
	}
	pool := redisPool.GetPool()
	rdb, _ := pool.GetClient()
	err := rdb.Del(notePreField + note.Id + noteReportedField).Err()
	if err != nil {
		return AllowNote(note, retry+1)
	}
	err = rdb.Set(notePreField+note.Id+noteAllowedField, 1, 0).Err()
	if err != nil {
		return AllowNote(note, retry+1)
	}
	return nil
}

func UnallowNote(note Note, retry int) error {
	if retry >= 3 {
		return errors.New("max retries reached")
	}
	pool := redisPool.GetPool()
	rdb, _ := pool.GetClient()
	err := rdb.Del(notePreField + note.Id + noteAllowedField).Err()
	if err != nil {
		return AllowNote(note, retry+1)
	}
	err = rdb.Set(notePreField+note.Id+noteReportedField, 1, 0).Err()
	if err != nil {
		return AllowNote(note, retry+1)
	}
	return nil
}

func DeleteNote(note Note, retry int) error {
	if retry >= 3 {
		return errors.New("max retries reached")
	}
	pool := redisPool.GetPool()
	rdb, _ := pool.GetClient()
	key := notePreField + note.Id
	err := rdb.Del(key+noteDataField, key+noteReportedField, key+noteAllowedField).Err()
	if err != nil {
		return AllowNote(note, retry+1)
	}
	return nil
}

func GetReportedNotes(retry int) ([]string, error) {
	if retry >= 3 {
		return []string{}, errors.New("max retries reached")
	}
	var note []string
	pool := redisPool.GetPool()
	rdb, _ := pool.GetClient()
	keys := rdb.Keys(notePreField + "*" + noteReportedField)
	if keys.Err() != nil {
		return GetReportedNotes(retry + 1)
	}
	for _, entry := range keys.Val() {
		id := strings.TrimSuffix(strings.TrimPrefix(entry, notePreField), noteReportedField)
		note = append(note, id)
	}
	return note, nil
}

func GetAllowedNotes(retry int) ([]string, error) {
	if retry >= 3 {
		return []string{}, errors.New("max retries reached")
	}
	var note []string
	pool := redisPool.GetPool()
	rdb, _ := pool.GetClient()
	keys := rdb.Keys(notePreField + "*" + noteAllowedField)
	if keys.Err() != nil {
		return GetReportedNotes(retry + 1)
	}
	for _, entry := range keys.Val() {
		id := strings.TrimSuffix(strings.TrimPrefix(entry, notePreField), noteAllowedField)
		note = append(note, id)
	}
	return note, nil
}

func IsValidId(id string) bool {
	if len(id) != 64 {
		return false
	}
	return idvalidregex.Match([]byte(id))
}
