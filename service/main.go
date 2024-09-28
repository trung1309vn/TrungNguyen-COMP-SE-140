package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	ps "github.com/mitchellh/go-ps"
	"golang.org/x/sys/unix"
)

type SystemInfo struct {
	IP           string
	Processes    [][2]string
	DiskUsage    uint64
	LastBootTime string
}

func getLastBootTime() string {
	cmd := exec.Command("uptime", "-s")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return ""
	}
	t := strings.TrimSpace(out.String())
	return t
}

func getDiskUsage() uint64 {
	var stat unix.Statfs_t
	wd, _ := os.Getwd()
	unix.Statfs(wd, &stat)
	// Available blocks * size per block = available space in bytes
	return stat.Bavail * uint64(stat.Bsize)
}

func getProccesses() [][2]string {
	processList, err := ps.Processes()
	var proccesses [][2]string
	if err != nil {
		log.Println("ps.Processes() Failed, are you using windows?")
		return proccesses
	}

	// map ages
	for x := range processList {
		var process ps.Process
		var proc [2]string
		process = processList[x]
		proc[0] = strconv.Itoa(process.Pid())
		proc[1] = process.Executable()
		proccesses = append(proccesses, proc)
		// do os.* stuff on the pid
	}

	return proccesses
}

func GetLocalIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddress := conn.LocalAddr().(*net.UDPAddr)

	return localAddress.IP
}

func helloWorld(w http.ResponseWriter, r *http.Request) {
	var passkey string = r.URL.Query().Get("passkey")
	if passkey != "" {
		if passkey == "test" {
			var IPAdress = GetLocalIP()
			var proccesses = getProccesses()
			// var procs = ""
			// for _, proc := range proccesses {
			// 	procs += fmt.Sprintf("\tid: %s, name: %s\n", proc[0], proc[1])
			// }
			t := getLastBootTime()
			var disk_usage = getDiskUsage()

			response := SystemInfo{}
			response.IP = IPAdress.String()
			response.Processes = proccesses
			response.DiskUsage = disk_usage
			response.LastBootTime = t

			responseJson, err := json.Marshal(response)
			if err != nil {
				panic(err)
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(responseJson)
		} else {
			fmt.Fprintf(w, "Wrong passkey!")
		}
	} else {
		fmt.Fprintf(w, "Access denied!")
	}
}
func main() {
	http.HandleFunc("/", helloWorld)
	http.ListenAndServe(":8080", nil)
}
