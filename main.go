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
	pool        *redis.Pool
	registerIp  = ""
)

func init() {
	pool = newPool()
	go regisiter()
}

func newPool() *redis.Pool {
	if registerIp == "" {
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
						registerIp = ip
						fmt.Println("register ip: ", ip)
					}
				}
			}
		}
	}

	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 120 * time.Second,
		// Dial or DialContext must be set. When both are set, DialContext takes precedence over Dial.
		Dial: func() (redis.Conn, error) {
			if strings.HasPrefix(registerIp, "172.17") {
				return redis.Dial("tcp", "172.17.0.1:6379")
			}
			return redis.Dial("tcp", ":6379")
		},
	}
}

func regisiter() {

	// get self ip
	tk := time.NewTicker(13 * time.Second)
	for {
		set(ProjectName, registerIp)
		<-tk.C
	}
}

func set(pn, ip string) {
	// connection redis
	c := pool.Get()
	defer c.Close()

	if _, err := c.Do("SELECT", 7); err != nil {
		log.Println(err, " [redis select failed] ", pn)
		return
	}
	// ttl 15 second
	if _, err := c.Do("SET", pn, ip, "EX", 15); err != nil {
		log.Println(err, " [set/ttl failed] ", pn)
		return
	}
	if pn == "project_name_undefined" {
		log.Println("project name: ", pn)
	}
}

func Get(pn string) (string, error) {
	// connection redis
	c := pool.Get()
	defer c.Close()
	s, err := redis.String(c.Do("GET", pn))
	if err != nil && err != redis.ErrNil {
		return "", err
	}
	return s, nil
}

func Start() {
	fmt.Println("Auto Register Sevelist >>> ", ProjectName, " <<<")
}
