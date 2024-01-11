package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	bolt "go.etcd.io/bbolt"
	"strconv"
)

var (
	delCmd = &cobra.Command{
		Use:   "delete [TASK NUMBER]",
		Short: "Delete a task from the task manager",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			delete(cmd, args)
		},
	}
)

func init() {
	rootCmd.AddCommand(delCmd)
}

func delete(cmd *cobra.Command, args []string) {
	//fmt.Println("Command add called")
	//fmt.Printf("Args[] = %+v\n", args)

	taskId, _ := strconv.Atoi(args[0])

	err := dbTasks.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket(dbBucketName)
		err := b.Delete(intToBinary(taskId)) // doesn't return error if entry doesn't exist!!

		return err
	})
	if err != nil {
		fmt.Printf("Could not delete task %s\n", taskId)
		panic(err)
	}
	fmt.Println("task deleted successfully")
}
