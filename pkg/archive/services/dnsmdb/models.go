package dnsmdb

import (
	"time"

	"golang.org/x/net/publicsuffix"

	"github.com/luids-io/api/dnsutil"
)

type mdbResolvData struct {
	ID        string        `bson:"_id"`
	Timestamp time.Time     `bson:"timestamp"`
	Duration  time.Duration `bson:"duration"`
	ServerIP  string        `bson:"serverIP"`
	ClientIP  string        `bson:"clientIP"`
	//query info
	QID        uint16              `bson:"qid"`
	Name       string              `bson:"name"`
	IsIPv6     bool                `bson:"isIPv6"`
	QueryFlags mdbResolvQueryFlags `bson:"queryFlags"`
	//response info
	ReturnCode     int                    `bson:"returnCode"`
	ResolvedIPs    []string               `bson:"resolvedIPs,omitempty"`
	ResolvedCNAMEs []string               `bson:"resolvedCNAMEs,omitempty"`
	ResponseFlags  mdbResolvResponseFlags `bson:"responseFlags"`
	//calculated info
	TLD        string `bson:"tld"`
	TLDPlusOne string `bson:"tldPlusOne"`
}

type mdbResolvQueryFlags struct {
	Do                bool `bson:"do"`
	AuthenticatedData bool `bson:"authenticatedData"`
	CheckingDisabled  bool `bson:"checkingDisabled"`
}

type mdbResolvResponseFlags struct {
	AuthenticatedData bool `bson:"authenticatedData"`
}

func toMongoData(data *dnsutil.ResolvData) mdbResolvData {
	mdbData := mdbResolvData{
		ID:        data.ID,
		Timestamp: data.Timestamp,
		Duration:  data.Duration,
		ServerIP:  data.Server.String(),
		ClientIP:  data.Client.String(),
		//query data
		QID:    data.QID,
		Name:   data.Name,
		IsIPv6: data.IsIPv6,
		QueryFlags: mdbResolvQueryFlags{
			Do:                data.QueryFlags.Do,
			AuthenticatedData: data.QueryFlags.AuthenticatedData,
			CheckingDisabled:  data.QueryFlags.CheckingDisabled,
		},
		//response data
		ReturnCode: data.ReturnCode,
		ResponseFlags: mdbResolvResponseFlags{
			AuthenticatedData: data.ResponseFlags.AuthenticatedData,
		},
	}
	if len(data.ResolvedIPs) > 0 {
		mdbData.ResolvedIPs = make([]string, 0, len(data.ResolvedIPs))
		for _, r := range data.ResolvedIPs {
			mdbData.ResolvedIPs = append(mdbData.ResolvedIPs, r.String())
		}
	}
	if len(data.ResolvedCNAMEs) > 0 {
		mdbData.ResolvedCNAMEs = make([]string, 0, len(data.ResolvedCNAMEs))
		for _, r := range data.ResolvedCNAMEs {
			mdbData.ResolvedCNAMEs = append(mdbData.ResolvedCNAMEs, r)
		}
	}
	mdbData.TLD, _ = publicsuffix.PublicSuffix(data.Name)
	mdbData.TLDPlusOne, _ = publicsuffix.EffectiveTLDPlusOne(data.Name)
	return mdbData
}
