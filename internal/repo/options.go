package repo

type OptionFunc func(*MetricsRepo)

// StoreFilePath sets the store file path for the repository.
func StoreFilePath(path string) OptionFunc {
	return func(repo *MetricsRepo) {
		repo.StoreFilePath = path
	}
}

// Restore sets the restore flag for the repository.
func Restore() OptionFunc {
	return func(repo *MetricsRepo) {
		repo.Restore = true
	}
}
