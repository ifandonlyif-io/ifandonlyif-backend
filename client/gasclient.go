package client

import (
	"fmt"

	"github.com/robfig/cron"
)

func RunCronFetchGas() {
	fmt.Println("AAAAAAAAAAAAAAAAA")
	i := 0
	c := cron.New()
	spec := "*/1 * * * * ?"
	err := c.AddFunc(spec, func() {
		i++
		fmt.Println("cron times : ", i)
	})
	if err != nil {
		fmt.Errorf("AddFunc error : %v", err)
		return
	}
	c.Start()

	defer c.Stop()
	select {}
}
