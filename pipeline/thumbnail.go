package pipeline

import (
	"fmt"
	"time"
)

func (p *asset) generateThumbnail() {
	fmt.Printf("gen thumbnail asset %d\n", p.ID)
	time.Sleep(time.Second)
	fmt.Printf("done %d\n", p.ID)
}
