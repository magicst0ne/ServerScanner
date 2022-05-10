package main

import (
	"fmt"
	"github.com/magicst0ne/ServerScanner/dispatcher"
	"github.com/magicst0ne/gofish"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
	"os"
	"strings"
)

var (
	Version       string
	BuildRevision string
	BuildBranch   string
	BuildTime     string
	BuildHost     string

	verbose = kingpin.Flag("verbose", "Verbose mode.").Short('v').Bool()

	redfishTimeout   = kingpin.Flag("timeout", "Timeout waiting for redfish api.").Default("60s").Short('t').Duration()
	redfishOutput    = kingpin.Flag("output", "output file.").Short('o').String()
	redfishNet       = kingpin.Arg("net", "net to scan 1.1.1.1/24").Required().String()
	redfishUser      = kingpin.Arg("user", "user to request redfish api").Required().String()
	redfishPassword1 = kingpin.Arg("password1", "password to request redfish api").Required().String()
	redfishPassword2 = kingpin.Arg("password2", "retry password to request redfish api").String()
	redfishPassword3 = kingpin.Arg("password3", "retry password to request redfish api").String()
)

func init() {
}

func main() {
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	hosts, err := cidrHosts(*redfishNet)
	if err != nil {
		fmt.Printf("[error]%s parser error\n", *redfishNet)
		os.Exit(1)
	}
	doScanByRedfish(hosts)
}

type ServerInfo struct {
	User               string
	Password           string
	Host               string
	AssetTag           string
	ServiceTag         string
	SerialNumber       string
	ManagerMACAddress  string
	SystemManufacturer string
	SystemModel        string
}

func getInfoByRedFish(host string, user string, password string) (int, ServerInfo) {

	record := ServerInfo{
		User:               user,
		Password:           password,
		Host:               host,
		AssetTag:           "null",
		SerialNumber:       "null",
		ServiceTag:         "null",
		ManagerMACAddress:  "00",
		SystemManufacturer: "Unknown",
		SystemModel:        "Unknown",
	}

	url := fmt.Sprintf("https://%s", host)
	config := gofish.ClientConfig{
		Endpoint:  url,
		Username:  user,
		Password:  password,
		Insecure:  true,
		BasicAuth: false,
	}
	c, err := gofish.Connect(config)
	if err != nil {
		//fmt.Printf("[error]connect %s faile  %v\n", host, err.Error())
		return -1, record
	}

	// Retrieve the service root
	service := c.Service

	// Query the chassis data using the session token
	systems, err := service.Systems()
	if err != nil {
		//fmt.Printf("[error]auth failed %s %v\n", host, err.Error())
		return -2, record
	}

	systemModel := "Unknown"
	systemManufacturer := "Unknown"
	ManagerMACAddress := "00"
	AssetTag := "null"
	ServiceTag := "null"
	SerialNumber := "null"

	for _, system := range systems {
		// server info
		SerialNumber = system.SerialNumber
		systemModel = strings.Replace(system.Model, " ", "", -1)

		if system.AssetTag != "" {
			AssetTag = system.AssetTag
		}
		if system.Manufacturer != "" {
			tmpStr := strings.Split(system.Manufacturer, " ")
			systemManufacturer = tmpStr[0]
		}
		if systemManufacturer == "Dell" {
			SerialNumber = system.SKU
		}
	}

	if ServiceTag == "null" {
		//ServiceTag = SerialNumber
	}

	record.AssetTag = AssetTag
	record.ServiceTag = ServiceTag
	record.SystemManufacturer = systemManufacturer
	record.SystemModel = systemModel
	record.SerialNumber = SerialNumber

	// Query the chassis data using the session token
	managers, err := service.Managers()
	if err != nil {
		//fmt.Printf("[error]auth failed %s", err.Error())
		return -3, record
	}

	for _, item := range managers {
		eths, _ := item.EthernetInterfaces()
		for _, eth := range eths {
			ManagerMACAddress = eth.MACAddress
			if eth.MACAddress == "" {
				ManagerMACAddress = eth.PermanentMACAddress
			}
		}
	}
	if ManagerMACAddress == "" {
		ManagerMACAddress = "00"
	}
	record.ManagerMACAddress = ManagerMACAddress
	return 0, record
}

func doScanByRedfish(hosts []string) {

	d := dispatcher.NewDispatcher(100, 10240, func(v interface{}) {
		host := v.(string)

		aliveIcmp := false
		alivePort := false

		rtt, err := doPing(host)
		if err != nil {
			//scanLoggerCtx.Info("ping failed")
		} else {
			//scanLoggerCtx.Info(fmt.Sprintf("ping alive rtt=%v", rtt))
			if rtt > 0 {
				alivePort = true
			}
		}

		ports := []string{"17988", "9666", "9999", "443", "22", "161", "2198", "5900"}

		portsStatus := tcpGather(host, ports)

		for _, item := range portsStatus {
			if item == true {
				alivePort = true
				break
			}
		}

		if aliveIcmp || alivePort {

			ret, sInfo := getInfoByRedFish(host, *redfishUser, *redfishPassword1)
			if *redfishPassword2 != "" && ret < 0 {
				ret, sInfo = getInfoByRedFish(host, *redfishUser, *redfishPassword2)
			} else if *redfishPassword3 != "" && ret < 0 {
				ret, sInfo = getInfoByRedFish(host, *redfishUser, *redfishPassword2)
			}
			_ = ret

			if sInfo.SystemManufacturer == "Unknown" {
				if portsStatus["443"] && portsStatus["9666"] && portsStatus["9999"] {
					sInfo.SystemManufacturer = "Inspire"
				}

				if portsStatus["443"] && portsStatus["17988"] {
					sInfo.SystemManufacturer = "HPE"
				}
				if portsStatus["443"] && portsStatus["9100"] {
					sInfo.SystemManufacturer = "Huawei"
				}
				sInfo.User = "-"
				sInfo.Password = "-"
			}

			fmt.Printf("[record] %v alive %v %v %v %v %v %v %v %v\n", host, sInfo.SystemManufacturer,
				sInfo.SystemModel, sInfo.AssetTag, sInfo.ServiceTag, sInfo.SerialNumber,
				sInfo.ManagerMACAddress, sInfo.User, sInfo.Password)

		} else {
			fmt.Printf("[record] %v noresp Unknown Unknown null null null 00 - - \n", host)
		}

	})
	d.Start()

	for _, host := range hosts {
		//rootLoggerCtx.Info(fmt.Sprintf("add host %v to scan list", host))
		d.Add(host)
	}
	d.Wait()

}
