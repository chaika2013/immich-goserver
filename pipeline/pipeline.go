package pipeline

import (
	"github.com/chaika2013/immich-goserver/config"
	"golang.design/x/chann"
)

const (
	JobGetExif = 1 << iota
	JobMoveToLibrary
	JobGenerateThumbnail
	JobEncodeVideo
)

const AllJobs = JobGetExif | JobMoveToLibrary | JobGenerateThumbnail | JobEncodeVideo
const TerminalJobs = JobGenerateThumbnail | JobEncodeVideo

var inst *pipeline

func Setup() {
	inst = &pipeline{
		semaphore:     make(chan struct{}, *config.ConcurrentFiles),
		exifChan:      chann.New[*asset](),
		moveChan:      chann.New[*asset](),
		thumbnailChan: chann.New[*asset](),
		encodeChan:    chann.New[*asset](),
	}

	for i := 0; i < *config.ConcurrentFiles; i++ {
		inst.semaphore <- struct{}{}
	}

	go inst.processQueues()
}

type pipeline struct {
	semaphore     chan struct{}
	exifChan      *chann.Chann[*asset]
	moveChan      *chann.Chann[*asset]
	thumbnailChan *chann.Chann[*asset]
	encodeChan    *chann.Chann[*asset]
}

// we use only id not to have the memory occupied with data which can be taken from db
type asset struct {
	ID   uint
	jobs uint32
}

func Enqueue(assetID uint, jobs uint32) {
	asset := &asset{
		ID:   assetID,
		jobs: jobs,
	}
	inst.enqueue(asset, 0, nil)
}

func (p *pipeline) enqueue(asset *asset, job uint32, err error) {
	if job > 0 {
		p.semaphore <- struct{}{}
	}

	if err != nil {
		// TODO log job failed with error
		return
	}

	if job&TerminalJobs > 0 {
		return
	}

	asset.jobs &= ^job

	if asset.jobs&JobGetExif > 0 {
		p.exifChan.In() <- asset
		return
	}
	if asset.jobs&JobMoveToLibrary > 0 {
		p.moveChan.In() <- asset
		return
	}
	if asset.jobs&JobGenerateThumbnail > 0 {
		p.thumbnailChan.In() <- asset
	}
	if asset.jobs&JobEncodeVideo > 0 {
		p.encodeChan.In() <- asset
	}
}

// ugly solution for prioritizing channels
func (p *pipeline) processQueues() {
	for {
		<-p.semaphore
		select {
		case asset := <-p.exifChan.Out():
			go func() {
				err := asset.extractExif()
				p.enqueue(asset, JobGetExif, err)
			}()
			continue
		default:
		}
		select {
		case asset := <-p.exifChan.Out():
			go func() {
				err := asset.extractExif()
				p.enqueue(asset, JobGetExif, err)
			}()
			continue
		case asset := <-p.moveChan.Out():
			go func() {
				err := asset.moveToLibrary()
				p.enqueue(asset, JobMoveToLibrary, err)
			}()
			continue
		default:
		}
		select {
		case asset := <-p.exifChan.Out():
			go func() {
				err := asset.extractExif()
				p.enqueue(asset, JobGetExif, err)
			}()
			continue
		case asset := <-p.moveChan.Out():
			go func() {
				err := asset.moveToLibrary()
				p.enqueue(asset, JobMoveToLibrary, err)
			}()
			continue
		case asset := <-p.thumbnailChan.Out():
			go func() {
				err := asset.generateThumbnail()
				p.enqueue(asset, JobGenerateThumbnail, err)
			}()
			continue
		default:
		}
		select {
		case asset := <-p.exifChan.Out():
			go func() {
				err := asset.extractExif()
				p.enqueue(asset, JobGetExif, err)
			}()
			continue
		case asset := <-p.moveChan.Out():
			go func() {
				err := asset.moveToLibrary()
				p.enqueue(asset, JobMoveToLibrary, err)
			}()
			continue
		case asset := <-p.thumbnailChan.Out():
			go func() {
				err := asset.generateThumbnail()
				p.enqueue(asset, JobGenerateThumbnail, err)
			}()
			continue
		case asset := <-p.encodeChan.Out():
			go func() {
				err := asset.encodeVideo()
				p.enqueue(asset, JobEncodeVideo, err)
			}()
			continue
		}
	}
}
