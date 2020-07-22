package servlist

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
)

var (
	ProjectName = "project_name_undefined"
)

func init() {
	go regisiter()
}

func cc() (redis.Conn, error) {
	// connection redis
	c, err := redis.Dial("tcp", ":6379")
	if err != nil {
		return nil, err
	}
	if _, err := c.Do("SELECT", 7); err != nil {
		return nil, err
	}
	return c, nil
}

func regisiter() {
	register_ip := ""
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip := ipnet.IP.String()
				fmt.Println(">> ", ip)
				if strings.HasPrefix(ip, "172") {
					register_ip = ip
					fmt.Println("register ip: ", ip)
				}
			}
		}
	}

	// connection redis
	c, err := cc()
	if err != nil {
		panic(err)
	}
	defer c.Close()

	// get project name
	// get self ip
	tk := time.NewTicker(13 * time.Second)

	for {
		// ttl 15 second
		if _, err := c.Do("SET", ProjectName, register_ip, "EX", 15); err != nil {
			log.Println(err, " [set/ttl failed] ", ProjectName)
		}
		if ProjectName == "project_name_undefined" {
			log.Println("project name: ", ProjectName)
		}
		<-tk.C
		// set ttl
	}
}

func Get(pn string) (string, error) {
	// connection redis
	c, err := cc()
	if err != nil {
		return "", err
	}
	defer c.Close()
	s, err := redis.String(c.Do("GET", pn))
	if err != nil && err != redis.ErrNil {
		return "", err
	}
	return s, nil
}
