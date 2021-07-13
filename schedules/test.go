package schedules

import (
	"fmt"
	"time"
)

func (c *Schedules) Test() {
	fmt.Println("current time:", time.Now().Format("2006-01-02 15:04:05"))
}
