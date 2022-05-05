package main

import (
	"bufio"
	"embed"
	"errors"
	"log"
	"math/rand"
	"net/http"
	"os"
	redisPool "pastetor/redis"
	"strings"
	"time"
)

//go:embed styles/*
var styles embed.FS

//go:embed templates/*
var index embed.FS

type User struct {
	Username string
	Password string
}

var users []User

func main() {
	rand.Seed(time.Now().UnixNano())
	usrs, err := parseUsers()
	if err != nil {
		log.Fatal(err)
		return
	}
	users = usrs
	redisPool.GetPool().InitPool()
	r := http.NewServeMux()
	routes(r)
	m := HeaderHandler(r)
	log.Println("Starting server on port 9000")
	err = http.ListenAndServe(":9000", m)
	if err != nil {
		panic(err)
	}
}

func parseUsers() ([]User, error) {
	usersfile := os.Getenv("USERS")
	if usersfile == "" {
		return nil, errors.New("no USERS env found")
	}
	fh, err := os.Open(usersfile)
	if err != nil {
		return nil, errors.New(usersfile + " file not found")
	}
	defer fh.Close()
	var userstmp []User
	scanner := bufio.NewScanner(fh)
	for scanner.Scan() {
		line := strings.SplitN(scanner.Text(), ":", 2)
		if len(line) == 2 {
			userstmp = append(userstmp, User{Username: line[0], Password: line[1]})
		}
	}
	if err = scanner.Err(); err != nil {
		return nil, err
	}
	return userstmp, nil
}
