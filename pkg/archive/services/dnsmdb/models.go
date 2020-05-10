package dnsmdb

import (
	"time"

	"github.com/luids-io/api/dnsutil"
)

type mdbResolvData struct {
	ID        string        `json:"id" bson:"_id"`
	Timestamp time.Time     `json:"timestamp"`
	Duration  time.Duration `json:"duration"`
	ServerIP  string        `json:"server_ip"`
	ClientIP  string        `json:"client_ip"`
	//query info
	QID              uint16 `json:"qid"`
	Name             string `json:"name"`
	CheckingDisabled bool   `json:"checking_disabled"`
	//response info
	ReturnCode        int      `json:"return_code"`
	AuthenticatedData bool     `json:"authenticated_data"`
	ResolvedIPs       []string `json:"resolved_ips,omitempty" bson:",omitempty"`
}

func toMongoData(data dnsutil.ResolvData) mdbResolvData {
	mdbData := mdbResolvData{
		ID:        data.ID,
		Timestamp: data.Timestamp,
		Duration:  data.Duration,
		ServerIP:  data.Server.String(),
		ClientIP:  data.Client.String(),
		//query data
		QID:              data.QID,
		Name:             data.Name,
		CheckingDisabled: data.CheckingDisabled,
		//response data
		ReturnCode:        data.ReturnCode,
		AuthenticatedData: data.AuthenticatedData,
	}
	if len(data.Resolved) > 0 {
		mdbData.ResolvedIPs = make([]string, 0, len(data.Resolved))
		for _, r := range data.Resolved {
			mdbData.ResolvedIPs = append(mdbData.ResolvedIPs, r.String())
		}
	}
	return mdbData
}
