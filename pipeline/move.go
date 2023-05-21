package pipeline

import (
	"fmt"
	"time"
)

func (p *asset) moveToLibrary() {
	fmt.Printf("move asset %d\n", p.ID)
	time.Sleep(time.Second)
	fmt.Printf("done %d\n", p.ID)
}
