package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	bolt "go.etcd.io/bbolt"
	//"encoding/binary"
	//"bytes"
)

var (
	listCmd = &cobra.Command{
		Use:   "list",
		Short: "List all tasks on the task manager",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			list(cmd, args)

		},
	}
)

func init() {
	rootCmd.AddCommand(listCmd)
}

func list(cmd *cobra.Command, args []string) {
	//fmt.Println("List command called")
	//fmt.Printf("Args[] = %+v\n",args)

	dbTasks.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(dbBucketName)
		counter := 0
		b.ForEach(func(k, v []byte) error {
			t := &task{}
			err := json.Unmarshal(v, t)
			counter++

			//var id uint64
			//binary.Read(bytes.NewBuffer(k), binary.BigEndian, &id)
			//fmt.Printf("key=%d, value=%s\n", id, v)
			fmt.Printf("%d - %s\n", t.Id, t.Text)
			return err
		})
		return nil
	})
}
