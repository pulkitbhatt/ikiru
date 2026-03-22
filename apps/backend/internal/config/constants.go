package config

const (
	AppName = "IKIRU"
	EnvDev  = "dev"
	EnvProd = "prod"
)

var RegionsToMonitor = []string{"us", "eu", "apac"}

const (
	WorkerMaxConcurrency = 50
)

const (
	EventIncidentCreated  = "incident.created"
	EventIncidentResolved = "incident.resolved"
)
