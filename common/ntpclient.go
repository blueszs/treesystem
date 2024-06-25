package common

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"syscall"
	"time"
)

type NTPServerConf struct {
	ntpHost        string  `default:"time.windows.com"` // NTP服务器地址ntp.tencent.com
	ntpPort        int     `default:"123"`              // NTP服务端口123
	ntpEpochOffset float64 `default:"2208988800"`       //NTP 秒数是从 1900 年开始计算的，因此必须使用一个纪元偏移量来进行校正，才能将 NTP 秒数转换为 Unix 时间，具体做法是减去 70 年的秒数（1970-1900），即 2208988800 秒。
}

// /获取默认的NTPServerConf实体
func DefaultNTPServerConf() NTPServerConf {
	return NTPServerConf{
		ntpHost:        "time.windows.com",
		ntpPort:        123,
		ntpEpochOffset: 2208988800,
	}
}

// NTP 数据包格式（v3，移除了可选的 v4 字段）
//
// 0                   1                   2                   3
// 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |LI | VN  |Mode |    Stratum     |     Poll      |  Precision   |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                         Root Delay                            |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                         Root Dispersion                       |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                          Reference ID                         |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                                                               |
// +                     Reference Timestamp (64)                  +
// |                                                               |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                                                               |
// +                      Origin Timestamp (64)                    +
// |                                                               |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                                                               |
// +                      Receive Timestamp (64)                   +
// |                                                               |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                                                               |
// +                      Transmit Timestamp (64)                  +
// |                                                               |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//LI（LeapIndicator）：闰秒标识器，占用2个 bit。
//VN（VersionNumber）：版本号，占用3个 bits，表示 NTP 的版本号，现在为3。
//Mode：模式，占用3个 bits，表示模式，包括预留、对称行为、客户端、服务器、广播、NTP 控制信息等。
//Stratum（层）：占用8个 bits。
//Poll：测试间隔，占用8个 bits，表示连续信息之间的最大间隔。
//Precision：精度，占用8个 bits，表示本地时钟精度。
//Root Delay：根时延，占用8个 bits，表示在主参考源之间往返的总共时延。
//Root Dispersion：根离散，占用8个 bits，表示在主参考源有关的名义错误。
//Reference Identifier：参考时钟标识符，占用8个 bits，用来标识特殊的参考源。
//Reference Timestamp：参考时间戳，64bits，时间戳，表示本地时钟被修改的最新时间。
//Originate Timestamp：原始时间戳，客户端发送的时间，64bits。
//Receive Timestamp：接受时间戳，服务端接受到的时间，64bits。
//Transmit Timestamp：传送时间戳，服务端送出应答的时间，64bits。

type packet struct {
	Settings       uint8  // 闰秒标识器，占用2个 bit。版本号，占用3个 bits，表示 NTP 的版本号。模式，占用3个 bits，表示模式，包括预留、对称行为、客户端、服务器、广播、NTP 控制信息等。
	Stratum        uint8  // （层）：占用8个 bits。
	Poll           int8   // 测试间隔，占用8个 bits，表示连续信息之间的最大间隔。
	Precision      int8   // 精度，占用8个 bits，表示本地时钟精度。
	RootDelay      uint32 // 根时延，占用8个 bits，表示在主参考源之间往返的总共时延。
	RootDispersion uint32 // 根离散，占用8个 bits，表示在主参考源有关的名义错误。
	ReferenceID    uint32 // 参考时钟标识符，占用8个 bits，用来标识特殊的参考源。
	RefTimeSec     uint32 // 参考时间戳秒
	RefTimeFrac    uint32 // 参考时间戳小数部分
	OrigTimeSec    uint32 // 原始时间戳，客户端发送的时间秒。
	OrigTimeFrac   uint32 // 原始时间戳，客户端发送的时间小数部分
	RxTimeSec      uint32 // 接受时间戳，服务端接受到的时间秒。
	RxTimeFrac     uint32 // 接受时间戳，客户端发送的时间小数部分
	TxTimeSec      uint32 // 传送时间戳，服务端送出应答的时间秒。
	TxTimeFrac     uint32 // 传送时间戳，客户端发送的时间小数部分
}

func Npt_Time(cf NTPServerConf) (time.Time, error) {
	// 解析NTP服务器地址
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", cf.ntpHost, cf.ntpPort))
	if err != nil {
		fmt.Println(err)
		return time.Time{}, err
	}
	// 创建UDP连接并发送请求数据包
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Println(err)
		return time.Time{}, err
	}
	defer conn.Close()
	// 发送请求数据包
	req := &packet{Settings: 0x1B}
	// send time request
	if err := binary.Write(conn, binary.BigEndian, req); err != nil {
		fmt.Printf("failed to send request: %v\n", err)
		return time.Time{}, err
	}
	// 接收响应数据包并解析时间戳字段
	rsp := &packet{}
	if err := binary.Read(conn, binary.BigEndian, rsp); err != nil {
		fmt.Printf("failed to read server response: %v\n", err)
		return time.Time{}, err
	}
	//在遵循 POSIX 标准的操作系统中，时间是用 Unix 时间纪元（自 1970 年以来的秒数）来表示的。
	//NTP 秒数是从 1900 年开始计算的，因此必须使用一个纪元偏移量来进行校正，才能将 NTP 秒数转换为 Unix 时间，具体做法是减去 70 年的秒数（1970-1900），即 2208988800 秒。
	secs := float64(rsp.TxTimeSec) - cf.ntpEpochOffset
	nanos := (int64(rsp.TxTimeFrac) * 1e9) >> 32 // 将小数部分转换为纳秒
	t := time.Unix(int64(secs), nanos)
	fmt.Printf("%v\n", t)
	return t, nil
}

func SetSystemDate(newTime time.Time) error {
	systemName := runtime.GOOS
	switch systemName {
	case "windows":
		{
			// 创建cmd命令
			cmd := exec.Command("cmd", "/C", "date", newTime.Format("2006-01-02"))
			err := cmd.Run()
			if err != nil {
				return err
			} else {
				return nil
			}
		}
	case "linux":
		{
			binary, lookErr := exec.LookPath("date")
			if lookErr != nil {
				fmt.Printf("Date binary not found, cannot set system date: %s", lookErr.Error())
				return lookErr
			} else {
				//dateString := newTime.Format("2006-01-2 15:4:5")
				dateString := newTime.Format("2 Jan 2006 15:04:05")
				fmt.Printf("Setting system date to: %s", dateString)
				args := []string{"--set", dateString}
				env := os.Environ()
				return syscall.Exec(binary, args, env)
			}
		}
	default:
		return nil
	}
}

func SetLinuxDate(newTime time.Time) error {
	_, lookErr := exec.LookPath("date")
	if lookErr != nil {
		fmt.Printf("Date binary not found, cannot set system date: %s\n", lookErr.Error())
		return lookErr
	} else {
		//dateString := newTime.Format("2006-01-2 15:4:5")
		dateString := newTime.Format("2 Jan 2006 15:04:05")
		fmt.Printf("Setting system date to: %s\n", dateString)
		args := []string{"--set", dateString}
		return exec.Command("date", args...).Run()
	}
}
