package main

/**
 * write 100w record in zset
 * 2 goroutine, one update random the other read
 * check the elasp time
 */

import (
	"explore/redis/db"
	"fmt"
	"git.ixiaochuan.cn/xclib/common/elapsed"
	_ "github.com/gomodule/redigo"
	"math/rand"
	"time"
)
import "explore/redis/common"

func writeRedis(finish chan int) {
	et := elapsed.NewElapsedTime()
	for i := 0; i < common.RECORD_NUM; i++ {
		userName := fmt.Sprintf("%s_%d", common.BASE_NAME, i)
		if err := db.SetKeyInTable(common.TABLE_NAME, userName, rand.Intn(100)); err != nil {
			fmt.Printf("SetKeyInTable failed:%s, param userName:%s, %d) \n", err.Error(), userName, rand.Intn(100))
			continue
		}
	}
	fmt.Printf("elaptime:%d \n", et.Elapsed())
	finish <- 1
}

func updateUserData() {
	userName := fmt.Sprintf("%s_%d", common.BASE_NAME, rand.Intn(common.RECORD_NUM))
	et := elapsed.NewElapsedTime()
	db.SetKeyInTable(common.TABLE_NAME, userName, rand.Intn(100))
	t := et.Elapsed()
	fmt.Printf("user:%s elasp_time:%d\n", userName, t)
}

func readUserData() {
	userName := fmt.Sprintf("%s_%d", common.BASE_NAME, rand.Intn(common.RECORD_NUM))
	et := elapsed.NewElapsedTime()
	rank := db.GetUserRank(common.TABLE_NAME, userName)
	score := db.GetUserRank(common.TABLE_NAME, userName)
	t := et.Elapsed()
	fmt.Printf("user:%s rank:%d score:%d  elasp_time:%d\n", userName, rank, score, t)
}

func done(finish chan int) {
	d := time.Duration(time.Second * 2)
	t := time.NewTicker(d)
	defer t.Stop()
	for {
		<-t.C
		finish <- 1
	}
}

func main() {
	finish := make(chan int)
	go writeRedis(finish)
	<-finish
	go func() {
		updateUserData()
		time.Sleep(1 * time.Second)
	}()
	go func() {
		readUserData()
		time.Sleep(2 * time.Second)
	}()
	go done(finish)
	<-finish
	close(finish)
}
