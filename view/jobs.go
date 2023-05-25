package view

type JobCounts struct {
	Active    int `json:"active"`
	Completed int `json:"completed"`
	Failed    int `json:"failed"`
	Delayed   int `json:"delayed"`
	Waiting   int `json:"waiting"`
	Parsed    int `json:"parsed"`
}

type QueueStatus struct {
	IsActive bool `json:"isActive"`
	IsPaused bool `json:"isPaused"`
}

type Job struct {
	JobCounts   JobCounts   `json:"jobCounts"`
	QueueStatus QueueStatus `json:"queueStatus"`
}

type AllJobs struct {
	MetadataExtractionQueue       Job `json:"metadata-extraction-queue"`
	StorageTemplateMigrationQueue Job `json:"storage-template-migration-queue"`
	ThumbnailGenerationQueue      Job `json:"thumbnail-generation-queue"`
	VideoConversionQueue          Job `json:"video-conversion-queue"`
	ObjectTaggingQueue            Job `json:"object-tagging-queue"`
	ClipEncodingQueue             Job `json:"clip-encoding-queue"`
	BackgroundTaskQueue           Job `json:"background-task-queue"`
	SearchQueue                   Job `json:"search-queue"`
	RecognizeFacesQueue           Job `json:"recognize-faces-queue"`
	SidecarQueue                  Job `json:"sidecar-queue"`
}
