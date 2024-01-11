package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	bolt "go.etcd.io/bbolt"
)

var (
	rootCmd = &cobra.Command{
		Use:   "task",
		Short: "task is a program for managing your TODO tasks.",
	}
	dbTasks      *bolt.DB
	dbBucketName = []byte("TASKS_BUCKET")
)

type task struct {
	Id   int
	Text string
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	openDB()
	createBucket()
}

func openDB() {
	var err error
	dbTasks, err = bolt.Open("tasks.db", 0600, nil)
	if err != nil {
		panic(err)
	}
}

func createBucket() {
	err := dbTasks.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(dbBucketName)
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}
