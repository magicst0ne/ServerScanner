package hwinfo

import (
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"github.com/gosnmp/gosnmp"
	"github.com/magicst0ne/HttpRequest"
	"strings"
)

type HwInfo struct {
	Host               string
	User               string
	Password           string
	Community          string
	Mfr                string
	Model              string
	AssetTag           string
	SerialNumber       string
	ExpressServiceCode string
	MacAddress         string
	BiosVerName        string
}

func GetHwInfo(host string, community string, user string, password1 string, password2 string) *HwInfo {

	req := HttpRequest.NewRequest()
	req.SetTLSClient(&tls.Config{InsecureSkipVerify: true})

	resp, err := req.Get(fmt.Sprintf("https://%s", host), nil)
	if err != nil {
		return nil
	}

	//unknow 0
	//dell   1
	//hpe    2
	//inspur 3
	//lenovo 4
	//huawei 5
	//Inpire 7
	mfr := 0

	if resp.StatusCode() == 200 {
		body, _ := resp.Body()

		if strings.Contains(string(body), "iDRAC.Embedded.1") {
			mfr = 1
		} else if strings.Contains(string(body), "isSCenabled") {
			mfr = 1
		} else if strings.Contains(string(body), "idrac") {
			mfr = 1
		} else if strings.Contains(string(body), "Hewlett") {
			mfr = 2
		} else if strings.Contains(string(body), "login.asp?lang=") {
			mfr = 5
		} else if strings.Contains(string(body), "ibmc") {
			mfr = 5
		} else if strings.Contains(string(body), "ATEN International") {
			mfr = 6
			return &HwInfo{
				User:               "-",
				Password:           "-",
				Host:               host,
				Community:          "-",
				Mfr:                "citrix",
				Model:              "Unknown",
				AssetTag:           "-",
				SerialNumber:       "-",
				ExpressServiceCode: "-",
				BiosVerName:        "-",
				MacAddress:         "00",
			}

		} else if strings.Contains(string(body), "COPYRIGHT.manufacturer") && strings.Contains(string(body), "Huawei") {

			mfr = 11
			return &HwInfo{
				User:               "-",
				Password:           "-",
				Host:               host,
				Community:          "-",
				Mfr:                "LanSwitch",
				Model:              "Unknown",
				AssetTag:           "-",
				SerialNumber:       "-",
				ExpressServiceCode: "-",
				BiosVerName:        "-",
				MacAddress:         "00",
			}

		} else if strings.Contains(string(body), "signString.copyrightStringSuffix") {
			mfr = 7
			return getInpireHwInfo(host, community)

		} else if strings.Contains(string(body), "Lenovo") {
			mfr = 8
			return &HwInfo{
				User:               "-",
				Password:           "-",
				Host:               host,
				Community:          "-",
				Mfr:                "Lenovo",
				Model:              "Unknown",
				AssetTag:           "-",
				SerialNumber:       "-",
				ExpressServiceCode: "-",
				BiosVerName:        "-",
				MacAddress:         "00",
			}

		} else {

			fmt.Println(host, string(body))
			return nil

		}

		if mfr == 1 {
			return getDellHwInfo(host, community)
		}
		if mfr == 2 {
			return getHpeHwInfo(host, community)
		}
		if mfr == 5 {
			return getHuaweiHwInfo(host, user, password1, password2)
		}
		if mfr == 3 {
			return getInpireHwInfo(host, community)
		}
		fmt.Println(host, "mfr=", mfr)

	} else {
		body, _ := resp.Body()
		fmt.Println(host, resp.StatusCode(), string(body))
	}

	return nil

}

func getDellHwInfo(host string, community string) *HwInfo {

	hwItem := HwInfo{
		User:               "-",
		Password:           "-",
		Host:               host,
		Community:          community,
		Mfr:                "Dell",
		Model:              "Unknown",
		AssetTag:           "-",
		SerialNumber:       "-",
		ExpressServiceCode: "-",
		BiosVerName:        "-",
		MacAddress:         "00",
	}

	//Dell 0  systemServiceTag                 .1.3.6.1.4.1.674.10892.5.1.3.2.0
	//Dell 1  systemExpressServiceCode         .1.3.6.1.4.1.674.10892.5.1.3.3.0
	//Dell 2  systemAssetTag                   .1.3.6.1.4.1.674.10892.5.1.3.4.0
	//Dell 3  systemDataCenterName             .1.3.6.1.4.1.674.10892.5.1.3.8.0
	//Dell 4  systemAisleName                  .1.3.6.1.4.1.674.10892.5.1.3.9.0
	//Dell 5  systemRackName                   .1.3.6.1.4.1.674.10892.5.1.3.10.0
	//Dell 6  systemRackSlot                   .1.3.6.1.4.1.674.10892.5.1.3.11.0
	//Dell 7  sysModelName                     .1.3.6.1.4.1.674.10892.5.1.3.12.0
	//Dell 8  systemSystemID                   .1.3.6.1.4.1.674.10892.5.1.3.13.0
	//Dell 9  systemBIOSVersionName            .1.3.6.1.4.1.674.10892.5.4.300.50.1.8.1.1
	//Dell 10 systemBIOSManufacturerName       .1.3.6.1.4.1.674.10892.5.4.300.50.1.11.1.1
	//Dell 11 PhysAddress                      .1.3.6.1.2.1.2.2.1.6.6

	oids := []string{
		".1.3.6.1.4.1.674.10892.5.1.3.2.0",
		".1.3.6.1.4.1.674.10892.5.1.3.3.0",
		".1.3.6.1.4.1.674.10892.5.1.3.4.0",
		".1.3.6.1.4.1.674.10892.5.1.3.8.0",
		".1.3.6.1.4.1.674.10892.5.1.3.9.0",
		".1.3.6.1.4.1.674.10892.5.1.3.10.0",
		".1.3.6.1.4.1.674.10892.5.1.3.11.0",
		".1.3.6.1.4.1.674.10892.5.1.3.12.0",
		".1.3.6.1.4.1.674.10892.5.1.3.13.0",
		".1.3.6.1.4.1.674.10892.5.4.300.50.1.8.1.1",
		".1.3.6.1.4.1.674.10892.5.4.300.50.1.11.1.1",
		".1.3.6.1.2.1.2.2.1.6.6",
		".1.3.6.1.2.1.2.2.1.6.7",
	}

	snmpData, err := getSnmpV2(host, community, oids)
	if err != nil && strings.Contains(err.Error(), "request timeout") {
		snmpData, err = getSnmpV2(host, "public", oids)
		if err != nil {
			hwItem.Community = "-"
			return &hwItem
		}
		hwItem.Community = "public"
	}

	if len(snmpData) == 13 {
		if snmpData[0].Type == gosnmp.OctetString {
			hwItem.SerialNumber = string(snmpData[0].Value.([]byte))
		}
		if snmpData[1].Type == gosnmp.OctetString {
			hwItem.ExpressServiceCode = string(snmpData[1].Value.([]byte))
		}
		if snmpData[2].Type == gosnmp.OctetString {
			hwItem.AssetTag = string(snmpData[2].Value.([]byte))
			if hwItem.AssetTag == "" {
				hwItem.AssetTag = "-"
			}
		}

		if snmpData[7].Type == gosnmp.OctetString {
			hwItem.Model = string(snmpData[7].Value.([]byte))
		}

		if snmpData[9].Type == gosnmp.OctetString {
			hwItem.BiosVerName = string(snmpData[9].Value.([]byte))
		}

		if snmpData[11].Type == gosnmp.OctetString {
			macAddress := snmpData[11].Value.([]byte)
			if len(macAddress) == 6 {
				hwItem.MacAddress = fmt.Sprintf("%v:%v:%v:%v:%v:%v",
					hex.EncodeToString(([]byte{macAddress[0]})),
					hex.EncodeToString(([]byte{macAddress[1]})),
					hex.EncodeToString(([]byte{macAddress[2]})),
					hex.EncodeToString(([]byte{macAddress[3]})),
					hex.EncodeToString(([]byte{macAddress[4]})),
					hex.EncodeToString(([]byte{macAddress[5]})),
				)
			}
		}
		if snmpData[12].Type == gosnmp.OctetString {
			macAddress := snmpData[12].Value.([]byte)
			if len(macAddress) == 6 {
				hwItem.MacAddress = fmt.Sprintf("%v:%v:%v:%v:%v:%v",
					hex.EncodeToString(([]byte{macAddress[0]})),
					hex.EncodeToString(([]byte{macAddress[1]})),
					hex.EncodeToString(([]byte{macAddress[2]})),
					hex.EncodeToString(([]byte{macAddress[3]})),
					hex.EncodeToString(([]byte{macAddress[4]})),
					hex.EncodeToString(([]byte{macAddress[5]})),
				)
			}
		}
	}

	return &hwItem
}

func getHpeHwInfo(host string, community string) *HwInfo {

	hwItem := HwInfo{
		User:               "-",
		Password:           "-",
		Host:               host,
		Community:          community,
		Mfr:                "HPE",
		Model:              "Unknown",
		AssetTag:           "-",
		SerialNumber:       "-",
		ExpressServiceCode: "-",
		BiosVerName:        "-",
		MacAddress:         "00",
	}

	//HPE 0 cpqSiSysSerialNum                        .1.3.6.1.4.1.232.2.2.2.1.0
	//HPE 1 cpqSiAssetTag                            .1.3.6.1.4.1.232.2.2.2.3.0
	//HPE 2 cpqSiProductName                         .1.3.6.1.4.1.232.2.2.4.2.0
	//HPE 3 BiosVerName                              .1.3.6.1.4.1.232.11.2.14.1.1.5.1
	//HPE 4 MacAddress                               .1.3.6.1.4.1.232.9.2.5.1.1.4.2

	oids := []string{
		".1.3.6.1.4.1.232.2.2.2.1.0",
		".1.3.6.1.4.1.232.2.2.2.3.0",
		".1.3.6.1.4.1.232.2.2.4.2.0",
		".1.3.6.1.4.1.232.11.2.14.1.1.5.1",
		".1.3.6.1.4.1.232.9.2.5.1.1.4.2",
	}

	snmpData, err := getSnmpV2(host, community, oids)
	if err != nil && strings.Contains(err.Error(), "request timeout") {
		snmpData, err = getSnmpV2(host, "public", oids)
		if err != nil {
			hwItem.Community = "-"
			return &hwItem
		}
		hwItem.Community = "public"
	}

	if len(snmpData) == 5 {
		if snmpData[0].Type == gosnmp.OctetString {
			hwItem.SerialNumber = string(snmpData[0].Value.([]byte))
		}
		if snmpData[1].Type == gosnmp.OctetString {
			hwItem.AssetTag = string(snmpData[1].Value.([]byte))
		}

		if snmpData[2].Type == gosnmp.OctetString {
			hwItem.Model = string(snmpData[2].Value.([]byte))
		}

		if snmpData[3].Type == gosnmp.OctetString {
			hwItem.BiosVerName = string(snmpData[3].Value.([]byte))
		}

		if snmpData[4].Type == gosnmp.OctetString {
			macAddress := snmpData[4].Value.([]byte)
			if len(macAddress) == 6 {
				hwItem.MacAddress = fmt.Sprintf("%v:%v:%v:%v:%v:%v",
					hex.EncodeToString(([]byte{macAddress[0]})),
					hex.EncodeToString(([]byte{macAddress[1]})),
					hex.EncodeToString(([]byte{macAddress[2]})),
					hex.EncodeToString(([]byte{macAddress[3]})),
					hex.EncodeToString(([]byte{macAddress[4]})),
					hex.EncodeToString(([]byte{macAddress[5]})),
				)
			}
		}
	}

	return &hwItem
}

func getHuaweiHwInfo(host string, user string, password1 string, password2 string) *HwInfo {

	hwItem := HwInfo{
		User:               "-",
		Password:           "-",
		Host:               host,
		Community:          "-",
		Mfr:                "Huawei",
		Model:              "Unknown",
		AssetTag:           "-",
		SerialNumber:       "-",
		ExpressServiceCode: "-",
		BiosVerName:        "-",
		MacAddress:         "00",
	}

	//Huawei 0 sysModelName                   .1.3.6.1.4.1.2011.2.235.1.1.1.6.0
	//Huawei 1 sysSerialNo                    .1.3.6.1.4.1.2011.2.235.1.1.1.7.0
	//Huawei 2 Guid                           .1.3.6.1.4.1.2011.2.235.1.1.1.10.0
	//Huawei 3 BiosName                       .1.3.6.1.4.1.2011.2.235.1.1.11.50.1.5.4.66.73.79.83
	//Huawei 4 MacAddress                     .1.3.6.1.2.1.2.2.1.6.2
	//Huawei 5 MacAddress                     .1.3.6.1.2.1.2.2.1.6.4

	oids := []string{
		".1.3.6.1.4.1.2011.2.235.1.1.1.6.0",
		".1.3.6.1.4.1.2011.2.235.1.1.1.7.0",
		".1.3.6.1.4.1.2011.2.235.1.1.1.10.0",
		".1.3.6.1.4.1.2011.2.235.1.1.11.50.1.5.4.66.73.79.83",
		".1.3.6.1.2.1.2.2.1.6.2",
		".1.3.6.1.2.1.2.2.1.6.4",
	}

	snmpData, err := getSnmpV3(host, oids, user, password1)
	if err != nil {
		snmpData, err = getSnmpV3(host, oids, user, password2)
		if err != nil {
			return &hwItem
		}
	}

	if len(snmpData) == 6 {
		if snmpData[1].Type == gosnmp.OctetString {
			hwItem.SerialNumber = string(snmpData[1].Value.([]byte))
		}
		if snmpData[0].Type == gosnmp.OctetString {
			hwItem.Model = string(snmpData[0].Value.([]byte))
		}

		if snmpData[3].Type == gosnmp.OctetString {
			hwItem.BiosVerName = string(snmpData[3].Value.([]byte))
		}

		if snmpData[4].Type == gosnmp.OctetString {
			macAddress := snmpData[4].Value.([]byte)
			if len(macAddress) == 6 {
				hwItem.MacAddress = fmt.Sprintf("%v:%v:%v:%v:%v:%v",
					hex.EncodeToString(([]byte{macAddress[0]})),
					hex.EncodeToString(([]byte{macAddress[1]})),
					hex.EncodeToString(([]byte{macAddress[2]})),
					hex.EncodeToString(([]byte{macAddress[3]})),
					hex.EncodeToString(([]byte{macAddress[4]})),
					hex.EncodeToString(([]byte{macAddress[5]})),
				)
			}
		}

		if snmpData[5].Type == gosnmp.OctetString {
			macAddress := snmpData[5].Value.([]byte)
			if len(macAddress) == 6 {
				hwItem.MacAddress = fmt.Sprintf("%v:%v:%v:%v:%v:%v",
					hex.EncodeToString(([]byte{macAddress[0]})),
					hex.EncodeToString(([]byte{macAddress[1]})),
					hex.EncodeToString(([]byte{macAddress[2]})),
					hex.EncodeToString(([]byte{macAddress[3]})),
					hex.EncodeToString(([]byte{macAddress[4]})),
					hex.EncodeToString(([]byte{macAddress[5]})),
				)
			}
		}
	}

	return &hwItem
}

func getInpireHwInfo(host string, community string) *HwInfo {

	hwItem := HwInfo{
		User:               "-",
		Password:           "-",
		Host:               host,
		Community:          community,
		Mfr:                "Inspur",
		Model:              "Unknown",
		AssetTag:           "-",
		SerialNumber:       "-",
		ExpressServiceCode: "-",
		BiosVerName:        "-",
		MacAddress:         "00",
	}

	//Inpire 0 mfr                            .1.3.6.1.4.1.37945.2.1.5.1.1.3.14.17.77.97.110.117.102.97.99.116.117.114.101.114.32.78.97.109.101
	//Inpire 1 Serial Number                  .1.3.6.1.4.1.37945.2.1.5.1.1.3.19.9.65.115.115.101.116.32.84.97.103
	//Inpire 2 model                          .1.3.6.1.4.1.37945.2.1.5.1.1.3.12.17.66.111.97.114.100.32.80.97.114.116.32.78.117.109.98.101.114
	//Inpire 3 PhysAddress                    .1.3.6.1.2.1.2.2.1.6.4

	oids := []string{
		".1.3.6.1.4.1.37945.2.1.5.1.1.3.14.17.77.97.110.117.102.97.99.116.117.114.101.114.32.78.97.109.101",
		".1.3.6.1.4.1.37945.2.1.5.1.1.3.19.9.65.115.115.101.116.32.84.97.103",
		".1.3.6.1.4.1.37945.2.1.5.1.1.3.12.17.66.111.97.114.100.32.80.97.114.116.32.78.117.109.98.101.114",
		".1.3.6.1.4.1.2011.2.235.1.1.1.10",
	}

	snmpData, err := getSnmpV2(host, community, oids)
	if err != nil && strings.Contains(err.Error(), "request timeout") {
		snmpData, err = getSnmpV2(host, "public", oids)
		if err != nil {
			return &hwItem
		}
		hwItem.Community = "public"
	}

	if len(snmpData) == 4 {
		if snmpData[1].Type == gosnmp.OctetString {
			hwItem.SerialNumber = string(snmpData[1].Value.([]byte))
		}
		if snmpData[2].Type == gosnmp.OctetString {
			hwItem.Model = string(snmpData[2].Value.([]byte))
		}

		if snmpData[3].Type == gosnmp.OctetString {
			macAddress := snmpData[3].Value.([]byte)
			if len(macAddress) == 6 {
				hwItem.MacAddress = fmt.Sprintf("%v:%v:%v:%v:%v:%v",
					hex.EncodeToString(([]byte{macAddress[0]})),
					hex.EncodeToString(([]byte{macAddress[1]})),
					hex.EncodeToString(([]byte{macAddress[2]})),
					hex.EncodeToString(([]byte{macAddress[3]})),
					hex.EncodeToString(([]byte{macAddress[4]})),
					hex.EncodeToString(([]byte{macAddress[5]})),
				)
			}
		}
	}

	/*
		for i, item := range snmpData {
			fmt.Printf("%d: oid: %s ", i, item.Name)

			// the Value of each variable returned by Get() implements
			// interface{}. You could do a type switch...
			switch item.Type {
			case gosnmp.OctetString:
				fmt.Printf("string: %s\n", string(item.Value.([]byte)))
			default:
				// ... or often you're just interested in numeric values.
				// ToBigInt() will return the Value as a BigInt, for plugging
				// into your calculations.
				fmt.Printf("number: %d\n", gosnmp.ToBigInt(item.Value))
			}
		}
	*/
	return &hwItem
}
