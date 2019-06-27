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
	"math/rand"
	"time"
)
import "explore/redis/common"

func writeRedis(finish chan int) {
	et := elapsed.NewElapsedTime()
	for i := 0; i < common.RECORD_NUM; i++ {
		userName := fmt.Sprintf("%s_%d", common.BASE_NAME, i)
		if err := db.SetKeyInTable(common.TABLE_NAME, userName, float64(rand.Intn(100))); err != nil {
			fmt.Printf("SetKeyInTable failed:%s, param userName:%s, %d) \n", err.Error(), userName, rand.Intn(100))
			continue
		}
	}
	fmt.Printf("write data finished:  elaptime:%d \n", et.Elapsed())
	finish <- 1
}

func updateUserData() {
	userName := fmt.Sprintf("%s_%d", common.BASE_NAME, rand.Intn(common.RECORD_NUM))
	et := elapsed.NewElapsedTime()
	db.SetKeyInTable(common.TABLE_NAME, userName, float64(rand.Intn(100)))
	t := et.Elapsed()
	fmt.Printf("[write] user:%s elasp_time:%d\n", userName, t)
}

func readUserData() {
	userName := fmt.Sprintf("%s_%d", common.BASE_NAME, rand.Intn(common.RECORD_NUM))
	et := elapsed.NewElapsedTime()
	rank := db.GetUserRank(common.TABLE_NAME, userName)
	score := db.GetUserScore(common.TABLE_NAME, userName)
	t := et.Elapsed()
	fmt.Printf("[read] user:%s rank:%d score:%d  elasp_time:%d\n", userName, rank, score, t)
}

func done(finish chan int) {
	d := time.Duration(time.Minute * 2)
	t := time.NewTicker(d)
	defer t.Stop()
	for {
		<-t.C
		finish <- 1
	}
}

func main() {
	db.Init()
	finish := make(chan int)
	go writeRedis(finish)
	<-finish
	go func() {
		for {
			updateUserData()
			time.Sleep(1 * time.Second)
		}
	}()
	go func() {
		for {
			readUserData()
			time.Sleep(2 * time.Second)
		}
	}()
	go done(finish)
	<-finish
	close(finish)
	db.Close()
}
