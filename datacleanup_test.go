// +build clean

// ***************
// CAUTION: Running this will delete all keys in your Redis database
// ***************

package main

import (
	"fmt"
	"testing"
)

func TestCleanTestData(t *testing.T) {

	iter := redisClient.Scan(0, "", 0).Iterator()
	for iter.Next() {
		fmt.Println("Deleting:", iter.Val())
		redisClient.Del(iter.Val())
	}
	if err := iter.Err(); err != nil {
		panic(err)
	}

}
