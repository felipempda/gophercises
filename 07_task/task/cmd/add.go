package cmd

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	bolt "go.etcd.io/bbolt"
	"strings"
)

var (
	addCmd = &cobra.Command{
		Use:   "add [TASK DESCRIPTION]",
		Short: "Add a task to the task manager",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			add(cmd, args)
		},
	}
)

func init() {
	rootCmd.AddCommand(addCmd)
}

func add(cmd *cobra.Command, args []string) {
	//fmt.Println("Command add called")
	//fmt.Printf("Args[] = %+v\n", args)

	err := dbTasks.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket(dbBucketName)
		id, _ := b.NextSequence()

		task := task{
			Text: strings.Join(args, " "),
			Id:   int(id),
		}
		buf, err := json.Marshal(task)
		if err != nil {
			return err
		}
		//fmt.Println("Id: ", idBinary)
		return b.Put(intToBinary(task.Id), buf)
	})
	if err != nil {
		fmt.Println("Error adding task")
		panic(err)
	}
	fmt.Println("task added successfully")
}

// convert int into []byte using bigEndian format
func intToBinary(i int) []byte {
	idBinary := make([]byte, 8)
	binary.BigEndian.PutUint64(idBinary, uint64(i))
	return idBinary
}
