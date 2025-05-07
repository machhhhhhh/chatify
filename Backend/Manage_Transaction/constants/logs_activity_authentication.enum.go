package constants

import "slices"

// enum type
type LogsActivityAuthentication string

// enum values
const (
	LogsActivityAuthenticationLogIn  LogsActivityAuthentication = "LOGIN"
	LogsActivityAuthenticationLogOut LogsActivityAuthentication = "LOGOUT"
)

var ValidatorLogsActivityAuthentication []LogsActivityAuthentication = []LogsActivityAuthentication{
	LogsActivityAuthenticationLogIn,
	LogsActivityAuthenticationLogOut,
}

func IsLogsActivityAuthenticationExist(role LogsActivityAuthentication) bool {
	return slices.Contains(ValidatorLogsActivityAuthentication, role) == true
}

func GetAllLogsActivityAuthentication() []string {
	var logs_activity_authentication []string
	for i := range ValidatorLogsActivityAuthentication {
		logs_activity_authentication = append(logs_activity_authentication, string(ValidatorLogsActivityAuthentication[i]))
	}
	return logs_activity_authentication
}
