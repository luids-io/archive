package dnsmdb

import (
	"net"
	"time"

	"github.com/globalsign/mgo/bson"

	"github.com/luids-io/api/dnsutil"
)

type mdbResolvData struct {
	ID        bson.ObjectId `bson:"_id"`
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

func toMData(src *dnsutil.ResolvData, dst *mdbResolvData) {
	if src.ID != "" {
		dst.ID = bson.ObjectIdHex(src.ID)
	}
	dst.Timestamp = src.Timestamp
	dst.Duration = src.Duration
	dst.ServerIP = src.Server.String()
	dst.ClientIP = src.Client.String()
	//query data
	dst.QID = src.QID
	dst.Name = src.Name
	dst.IsIPv6 = src.IsIPv6
	dst.QueryFlags.Do = src.QueryFlags.Do
	dst.QueryFlags.AuthenticatedData = src.QueryFlags.AuthenticatedData
	dst.QueryFlags.CheckingDisabled = src.QueryFlags.CheckingDisabled
	//response data
	dst.ReturnCode = src.ReturnCode
	dst.ResponseFlags.AuthenticatedData = src.ResponseFlags.AuthenticatedData
	if len(src.ResolvedIPs) > 0 {
		dst.ResolvedIPs = make([]string, 0, len(src.ResolvedIPs))
		for _, ip := range src.ResolvedIPs {
			dst.ResolvedIPs = append(dst.ResolvedIPs, ip.String())
		}
	}
	if len(src.ResolvedCNAMEs) > 0 {
		dst.ResolvedCNAMEs = make([]string, 0, len(src.ResolvedCNAMEs))
		for _, cname := range src.ResolvedCNAMEs {
			dst.ResolvedCNAMEs = append(dst.ResolvedCNAMEs, cname)
		}
	}
	dst.TLD = src.TLD
	dst.TLDPlusOne = src.TLDPlusOne
}

func fromMData(src *mdbResolvData, dst *dnsutil.ResolvData) {
	dst.ID = src.ID.Hex()
	dst.Timestamp = src.Timestamp
	dst.Duration = src.Duration
	dst.Server = net.ParseIP(src.ServerIP)
	dst.Client = net.ParseIP(src.ClientIP)
	//query data
	dst.QID = src.QID
	dst.Name = src.Name
	dst.IsIPv6 = src.IsIPv6
	dst.QueryFlags.Do = src.QueryFlags.Do
	dst.QueryFlags.AuthenticatedData = src.QueryFlags.AuthenticatedData
	dst.QueryFlags.CheckingDisabled = src.QueryFlags.CheckingDisabled
	//response data
	dst.ReturnCode = src.ReturnCode
	dst.ResponseFlags.AuthenticatedData = src.ResponseFlags.AuthenticatedData
	if len(src.ResolvedIPs) > 0 {
		dst.ResolvedIPs = make([]net.IP, 0, len(src.ResolvedIPs))
		for _, ip := range src.ResolvedIPs {
			dst.ResolvedIPs = append(dst.ResolvedIPs, net.ParseIP(ip))
		}
	}
	if len(src.ResolvedCNAMEs) > 0 {
		dst.ResolvedCNAMEs = make([]string, 0, len(src.ResolvedCNAMEs))
		for _, cname := range src.ResolvedCNAMEs {
			dst.ResolvedCNAMEs = append(dst.ResolvedCNAMEs, cname)
		}
	}
	//calculated data
	dst.TLD = src.TLD
	dst.TLDPlusOne = src.TLDPlusOne
}
