package pipeline

import (
	"fmt"
	"time"
)

func (p *asset) encodeVideo() {
	fmt.Printf("encode video %d\n", p.ID)
	time.Sleep(time.Second)
	fmt.Printf("done %d\n", p.ID)
}
