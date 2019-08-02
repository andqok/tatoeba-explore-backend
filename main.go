package main

import (
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gomodule/redigo/redis"
	_ "github.com/lib/pq"
)

var (
	Pool *redis.Pool
)

func init() {
	redisHost := ":6379"
	Pool = newPool(redisHost)
	cleanupHook()
}

func newPool(server string) *redis.Pool {

	return &redis.Pool{

		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,

		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			return c, err
		},

		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func cleanupHook() {

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	signal.Notify(c, syscall.SIGKILL)
	go func() {
		<-c
		Pool.Close()
		os.Exit(0)
	}()
}

func RedisSet(key string, value []byte) error {

	conn := Pool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", key, value)
	if err != nil {
		v := string(value)
		if len(v) > 15 {
			v = v[0:12] + "..."
		}
		return fmt.Errorf("error setting key %s to %s: %v", key, v, err)
	}
	return err
}

const connStr = "user=postgres password=ajoutee dbname=tatoeba_explore sslmode=disable"

// Sentence my-f-ckomment
type Sentence struct {
	ID     int    `json:"id"`
	Number int    `json:"number"`
	Text   string `json:"text"`
	Lang   string `json:"lang"`
}

func setToRedis() {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	rows, err1 := db.Query("SELECT * FROM sentences")
	if err1 != nil {
		fmt.Println("Error 1 happened", err1)
		return
	}
	defer rows.Close()

	// insertStatement := `INSERT INTO words (word, sentence_number, lang) VALUES ($1, $2, $3)`
	for rows.Next() {
		snt := Sentence{}
		fmt.Println(rows)
		err2 := rows.Scan(
			&snt.ID,
			&snt.Number,
			&snt.Text,
			&snt.Lang)
		if err2 != nil {
			fmt.Println("Error 2 happened", err2)
			return
		}
		err3 := RedisSet(strconv.Itoa(snt.Number), []byte(snt.Text))
		fmt.Println(snt.Number)
		if err3 != nil {
			fmt.Println("Error 3 happened", err3)
			return
		}
		// for _, word := range strings.Fields(snt.Text) {
		// 	lower := strings.ToLower(word)
		// 	removed := removePunctuation(lower)
		// 	fmt.Println(lower, removed)
		// 	_, err = db.Exec(insertStatement, removed, snt.Number, currLang)
		// 	if err != nil {
		// 		panic(err)
		// 	}
		// }
	}
}

func main() {
	setToRedis()
	// getSplittedWords("eng")
	// runServer()
}
