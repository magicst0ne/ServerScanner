package main

import (
	"fmt"
	"github.com/magicst0ne/ServerScanner/hwinfo"
	"github.com/magicst0ne/golib/cidr"
	"github.com/magicst0ne/golib/dispatcher"
	"github.com/magicst0ne/golib/net"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
	"os"
)

var (
	Version       string
	BuildRevision string
	BuildBranch   string
	BuildTime     string
	BuildHost     string

	verbose = kingpin.Flag("verbose", "Verbose mode.").Short('v').Bool()

	cliOutput    = kingpin.Flag("output", "output file.").Short('o').String()
	cliNet       = kingpin.Arg("net", "net to scan 1.1.1.1/24").Required().String()
	cliCommunity = kingpin.Arg("community", "snmp community").Required().String()
	cliUser      = kingpin.Arg("user", "user to request snmp").Required().String()
	cliPassword1 = kingpin.Arg("password1", "password to request snmp").Required().String()
	cliPassword2 = kingpin.Arg("password2", "retry password to request snmp").String()
)

func init() {
}

func main() {
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	hosts, err := cidr.AddressRange(*cliNet)
	if err != nil {
		fmt.Printf("[error]%s parser error\n", *cliNet)
		os.Exit(1)
	}
	doScanHwInfo(hosts)

}

func doScanHwInfo(hosts []string) {

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
				aliveIcmp = true
			}
		}

		// check ports active
		ports := []string{
			"22",
			"23",
			"80",
			"161",
			"443",
			"445",
			"2198",
			"5900",
			"9666",
			"9999",
			"17988",
		}

		portsStatus, _ := net.TcpGather(host, ports)
		activePortsString := ""

		for k, v := range portsStatus {
			if v {
				alivePort = true
				activePortsString = fmt.Sprintf("%s %s", activePortsString, k)
			}
		}

		if aliveIcmp || alivePort {

			if alivePort {
				hwItem := hwinfo.GetHwInfo(host, *cliCommunity, *cliUser, *cliPassword1, *cliPassword2)
				if hwItem != nil {
					fmt.Printf("[record], %v, alive, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v\n", host, hwItem.Mfr,
						hwItem.Model, hwItem.AssetTag, hwItem.SerialNumber, hwItem.ExpressServiceCode,
						hwItem.MacAddress, hwItem.BiosVerName, hwItem.Community, hwItem.User, hwItem.Password, activePortsString)
					return
				}
			}

			mfr := "Unknown"
			p17988, ok := portsStatus["17988"]
			if ok && p17988 {
				mfr = "HPE"
			}

			p9100, ok := portsStatus["9100"]
			if ok && p9100 {
				mfr = "Huawei"
			}

			p9999, ok := portsStatus["9999"]
			if ok && p9999 {
				mfr = "Inspire"
			}

			fmt.Printf("[record], %v, alive, %v, Unknown, -, -, -, -, -, -, -, -, %s\n", host, mfr, activePortsString)

		} else {
			fmt.Printf("[record], %v, noresp, Unknown, Unknown, -, -, -, -, -, -, -, -, -\n", host)
		}

	})
	d.Start()

	for _, host := range hosts {
		d.Add(host)
	}
	d.Wait()

}
