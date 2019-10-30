package mongodb

import (
	"time"

	"github.com/luids-io/core/dnsutil"
)

type mdbResolvData struct {
	ID          string    `json:"id" bson:"_id"`
	Timestamp   time.Time `json:"timestamp"`
	ServerIP    string    `json:"server_ip"`
	ClientIP    string    `json:"client_ip"`
	ResolvedIPs []string  `json:"resolved_ips"`
	Name        string    `json:"name"`
}

func toMongoData(data dnsutil.ResolvData) mdbResolvData {
	resolved := make([]string, 0, len(data.Resolved))
	for _, r := range data.Resolved {
		resolved = append(resolved, r.String())
	}
	return mdbResolvData{
		ID:          data.ID,
		Timestamp:   data.Timestamp,
		ServerIP:    data.Server.String(),
		ClientIP:    data.Client.String(),
		ResolvedIPs: resolved,
		Name:        data.Name,
	}
}
