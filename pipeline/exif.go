package pipeline

import (
	"fmt"
	"time"
)

func (p *asset) extractExif() {
	fmt.Printf("get exif for %d\n", p.ID)
	time.Sleep(time.Second)
	fmt.Printf("done %d\n", p.ID)
}
