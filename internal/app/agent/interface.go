package agent

type WebAgentAPI interface {
	SendMetrics(string, string, interface{}) error
}
