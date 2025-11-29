package requestdata

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc/metadata"
	"net/http"
	"strconv"
	"strings"
)

// HTTP
const (
	AccountHeaderKey        = "X-Account"
	AccountSchemaHeaderKey  = "X-Account-Schema"
	AccountIdHeaderKey      = "X-Account-Id"
	FileNameHeaderKey       = "X-File-Name"
	FileSizeHeaderKey       = "X-File-Size"
	FileTypeHeaderKey       = "X-File-Type"
	SessionIDHeaderKey      = "X-Session-Id"
	ForwardedForHeaderKey   = "X-Forwarded-For"
	HealthCheckerHeaderKey  = "X-Health-Checker"
	RobotsTagHeaderKey      = "X-Robots-Tag"
	RequestIDHeaderKey      = "X-Request-Id"
	ClientNameHeaderKey     = "X-Client-Name"
	ClientSchemaHeaderKey   = "X-Client-Schema"
	ClientIdHeaderKey       = "X-Client-Id"
	ClientTimeZoneHeaderKey = "X-Client-Time-Zone"
	ClientStatusHeaderKey   = "X-Client-Status"
	UserIdHeaderKey         = "X-User-Id"
	UserLoginHeaderKey      = "X-User-Login"
	UserStatusHeaderKey     = "X-User-Status"
	UserRolesHeaderKey      = "X-User-Roles"
	UserLanguageHeaderKey   = "X-User-Language"

	ContentTypeHeaderKey                = "Content-Type"
	AccessControlAllowOriginHeaderKey   = "Access-Control-Allow-Origin"
	AccessControlMaxAgeHeaderKey        = "Access-Control-Max-Age"
	AccessControlAllowMethodsHeaderKey  = "Access-Control-Allow-Methods"
	AccessControlAllowHeadersHeaderKey  = "Access-Control-Allow-Headers"
	AccessControlExposeHeadersHeaderKey = "Access-Control-Expose-Headers"

	ConnectionHeaderKey    = "Connection"
	AuthorizationHeaderKey = "Authorization"
	UserAgentHeaderKey     = "User-Agent"
	FromServiceHeaderKey   = "From-Service"
)

const (
	intBase     = 10
	intBitSize  = 64
	RDHeaderKey = "request_data"
)

type RequestData struct {
	SessionID    string `json:"sessionID,omitempty"`
	UserAgent    string `json:"userAgent,omitempty"`
	IPAddress    string `json:"ipAddr,omitempty"`
	RequestID    string `json:"requestID,omitempty"`
	UserID       int64  `json:"clientID,omitempty"`
	ClientName   string `json:"clientName,omitempty"`
	ClientSchema string `json:"clientSchema,omitempty"`
	UserId       int64  `json:"userId,omitempty"`
	UserEmail    string `json:"userLogin,omitempty"`
	UserRole     string `json:"userRoles,omitempty"`
	UserStatus   string `json:"userStatus,omitempty"`
	UserLanguage string `json:"userLanguage,omitempty"`
	ClientStatus string `json:"clientStatus,omitempty"`
	FromService  string `json:"fromService,omitempty"`
}

func (rd *RequestData) ToMD() metadata.MD {
	return metadata.New(map[string]string{
		"client_name":   rd.ClientName,
		"client_schema": rd.ClientSchema,
		"client_id":     strconv.FormatInt(rd.UserID, intBase),
		"user_login":    rd.UserEmail,
		"user_id":       strconv.FormatInt(rd.UserId, intBase),
		"user_roles":    rd.UserRole,
		"user_status":   rd.UserStatus,
		"user_agent":    rd.UserAgent,
		"ip_address":    rd.IPAddress,
		"request_id":    rd.RequestID,
		"session_id":    rd.SessionID,
		"client_status": rd.ClientStatus,
		"from_service":  rd.FromService,
	})
}

func (rd *RequestData) ToHeader(header http.Header) http.Header {
	header.Set(ClientNameHeaderKey, rd.ClientName)
	header.Set(ClientSchemaHeaderKey, rd.ClientSchema)
	header.Set(ClientIdHeaderKey, strconv.FormatInt(rd.UserID, intBase))
	header.Set(ClientStatusHeaderKey, rd.ClientStatus)

	header.Set(UserIdHeaderKey, strconv.FormatInt(rd.UserId, intBase))
	header.Set(UserLoginHeaderKey, rd.UserEmail)
	header.Set(UserStatusHeaderKey, rd.UserStatus)
	header.Set(UserRolesHeaderKey, rd.UserRole)
	header.Set(UserLanguageHeaderKey, rd.UserLanguage)

	header.Set(RequestIDHeaderKey, rd.RequestID)
	header.Set(UserAgentHeaderKey, rd.UserAgent)
	header.Set(SessionIDHeaderKey, rd.SessionID)
	header.Set(ForwardedForHeaderKey, rd.IPAddress)

	header.Set(FromServiceHeaderKey, rd.FromService)

	return header
}

func (rd *RequestData) ToHeaderMap() map[string]interface{} {
	return map[string]interface{}{
		"client_name":   rd.ClientName,
		"client_schema": rd.ClientSchema,
		"client_id":     rd.UserID,
		"user_login":    rd.UserEmail,
		"user_id":       rd.UserId,
		"user_roles":    rd.UserRole,
		"user_status":   rd.UserStatus,
		"user_agent":    rd.UserAgent,
		"ip_address":    rd.IPAddress,
		"request_id":    rd.RequestID,
		"session_id":    rd.SessionID,
		"client_status": rd.ClientStatus,
		"from_service":  rd.FromService,
	}
}

func (rd *RequestData) LogrusFields() map[string]interface{} {
	return map[string]interface{}{
		"user":         rd.UserEmail,
		"client_name":  rd.ClientName,
		"request_id":   rd.RequestID,
		"session_id":   rd.SessionID,
		"user_agent":   rd.UserAgent,
		"ip_address":   rd.IPAddress,
		"from_service": rd.FromService,
	}
}

func FromHeader(header http.Header) (rd RequestData) {
	rd.ClientName = header.Get(ClientNameHeaderKey)
	rd.ClientSchema = header.Get(ClientSchemaHeaderKey)
	rd.UserID, _ = strconv.ParseInt(header.Get(ClientIdHeaderKey), intBase, intBitSize)
	rd.ClientStatus = header.Get(ClientStatusHeaderKey)

	rd.UserId, _ = strconv.ParseInt(header.Get(UserIdHeaderKey), intBase, intBitSize)
	rd.UserEmail = header.Get(UserLoginHeaderKey)
	rd.UserStatus = header.Get(UserStatusHeaderKey)
	rd.UserRole = header.Get(UserRolesHeaderKey)
	rd.UserLanguage = header.Get(UserLanguageHeaderKey)

	rd.RequestID = header.Get(RequestIDHeaderKey)
	rd.UserAgent = header.Get(UserAgentHeaderKey)
	rd.SessionID = header.Get(SessionIDHeaderKey)
	rd.IPAddress = header.Get(ForwardedForHeaderKey)

	return rd
}

func (rd *RequestData) RD() RequestData {
	if rd == nil {
		return RequestData{RequestID: uuid.Must(uuid.NewV4(), nil).String()}
	}
	return *rd
}

func (rd *RequestData) String() string {
	var parts []string

	// Форматирование строковых полей с кавычками
	appendPart := func(fieldName, format string, value interface{}) {
		parts = append(parts, fmt.Sprintf("%s: %"+format, fieldName, value))
	}

	appendPart("SessionID", "q", rd.SessionID)
	appendPart("UserAgent", "q", rd.UserAgent)
	appendPart("IPAddress", "q", rd.IPAddress)
	appendPart("RequestID", "q", rd.RequestID)
	appendPart("ClientName", "q", rd.ClientName)
	appendPart("ClientSchema", "q", rd.ClientSchema)
	appendPart("UserEmail", "q", rd.UserEmail)
	appendPart("UserStatus", "q", rd.UserStatus)
	appendPart("UserLanguage", "q", rd.UserLanguage)
	appendPart("ClientStatus", "q", rd.ClientStatus)
	appendPart("FromService", "q", rd.FromService)
	appendPart("Route", "q", rd.FromService)
	appendPart("UserRole", "q", rd.UserRole)
	// Числовые поля
	appendPart("UserID", "d", rd.UserID)
	appendPart("UserId", "d", rd.UserId)

	// Объединение всех частей
	return fmt.Sprintf("{ %s }", strings.Join(parts, ", "))
}

func (rd *RequestData) Copy() RequestData {
	return RequestData{
		SessionID:    rd.SessionID,
		UserAgent:    rd.UserAgent,
		IPAddress:    rd.IPAddress,
		RequestID:    rd.RequestID,
		UserID:       rd.UserID,
		ClientName:   rd.ClientName,
		ClientSchema: rd.ClientSchema,
		UserId:       rd.UserId,
		UserEmail:    rd.UserEmail,
		UserRole:     rd.UserRole,
		UserStatus:   rd.UserStatus,
		UserLanguage: rd.UserLanguage,
		ClientStatus: rd.ClientStatus,
		FromService:  rd.FromService,
	}
}
