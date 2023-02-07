package repo

type OptionFunc func(*MetricsRepo)

func StoreFilePath(path string) OptionFunc {
	return func(repo *MetricsRepo) {
		repo.StoreFilePath = path
	}
}

func Restore() OptionFunc {
	return func(repo *MetricsRepo) {
		repo.Restore = true
	}
}
