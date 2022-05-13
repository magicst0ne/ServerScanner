package hwinfo

import (
	"github.com/gosnmp/gosnmp"
	"time"
)

func getSnmpV2(host string, community string, oids []string) ([]gosnmp.SnmpPDU, error) {

	// build our own GoSNMP struct, rather than using g.Default
	g := &gosnmp.GoSNMP{
		Target:    host,
		Port:      161,
		Version:   gosnmp.Version2c,
		Community: community,
		Timeout:   time.Duration(10) * time.Second,
	}

	err := g.Connect()
	if err != nil {
		return nil, err
	}
	defer g.Conn.Close()

	result, err2 := g.Get(oids) // Get() accepts up to g.MAX_OIDS
	if err2 != nil {
		return nil, err2
	}

	return result.Variables, nil
}

func getSnmpV3(host string, oids []string, user string, password string) ([]gosnmp.SnmpPDU, error) {

	// build our own GoSNMP struct, rather than using g.Default
	g := &gosnmp.GoSNMP{
		Target:        host,
		Port:          161,
		Version:       gosnmp.Version3,
		SecurityModel: gosnmp.UserSecurityModel,
		MsgFlags:      gosnmp.AuthPriv,
		Timeout:       time.Duration(10) * time.Second,
		SecurityParameters: &gosnmp.UsmSecurityParameters{
			UserName:                 user,
			AuthenticationProtocol:   gosnmp.SHA,
			AuthenticationPassphrase: password,
			PrivacyProtocol:          gosnmp.AES,
			PrivacyPassphrase:        password,
		},
	}

	err := g.Connect()
	if err != nil {
		return nil, err
	}
	defer g.Conn.Close()

	result, err2 := g.Get(oids) // Get() accepts up to g.MAX_OIDS
	if err2 != nil {
		return nil, err2
	}

	return result.Variables, nil
}
