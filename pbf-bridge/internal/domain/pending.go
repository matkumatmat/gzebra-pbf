package domain

type PendingJobRepository interface {
	SaveFailedJob(jobType string, payload interface{}) error
	ListPendingFiles() ([]string, error)
	DeletePendingFile(filename string) error
	ReadPendingFile(filename string) (string, []byte, error)
}

type PendingUseCase interface {
	HandleFailedPrint(jobType string, payload interface{}, originalErr error) error
	RetryAllPending() error
	SetUsecases(shipping ShippingUseCase, identity IdentityUseCase)
}
