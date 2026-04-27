package bo

type OperationAuditLogCreateInput struct {
	UserId      int64
	Username    string
	UserType    int
	Module      string
	Action      string
	LogType     string
	Method      string
	Path        string
	RequestData string
	Ip          string
	UserAgent   string
	TraceId     string
	DeviceInfo  string
}
