package agent

type WebAPIAgent interface {
	SendMetrics(string, string, interface{}) error
}
