package db

import (
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
)

var readConn, writeConn redis.Conn

func Init() {
	readConn = getConnect()
	writeConn = getConnect()
}

func getConnect() redis.Conn {
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		fmt.Printf("redis dial failed:%s\n", err.Error())
		return nil
	}
	return conn
}

func SetKeyInTable(key string, userName string, score float64) error {
	if writeConn == nil {
		return errors.New("redis connect failed")
	}
	if _, err := writeConn.Do("zadd", key, score, userName); err != nil {
		return err
	}
	return nil
}

func GetUserRank(key string, userName string) int {
	if readConn == nil {
		fmt.Printf("redis connect failed\n")
		return -1
	}
	rank, err := redis.Int(readConn.Do("zrevrank", key, userName))
	if err != nil {
		fmt.Printf("zrevrank failed:%s\n", err.Error())
		return -1
	}
	return rank
}

func GetUserScore(key string, userName string) int {
	if readConn == nil {
		fmt.Printf("redis connect failed\n")
		return -1
	}
	score, err := redis.Float64(readConn.Do("zscore", key, userName))
	if err != nil {
		fmt.Printf("zscore failed:%s\n", err.Error())
		return -1
	}
	return int(score)
}

func Close() {
	readConn.Close()
	writeConn.Close()
}
