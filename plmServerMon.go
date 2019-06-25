package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

//定义server结构体
type server struct {
	cas int
	ip  string
}

//定义result结构体
type result struct {
	cas    int
	secs   float64
	status string
}

//定义结构体切片
type servers []server
type results []result

//makeRequest 检查http请求状态并记录响应时间
func makeRequest(s server, ch chan<- result) {
	var status string
	start := time.Now()
	url := "http://" + s.ip + "/3dspace/services/NioPortalWebService"
	resp, err := http.Get(url)
	if err != nil {
		status = "Unavailable"
	} else if resp.StatusCode == 200 {
		status = "OK"
	} else {
		status = "Down"
	}
	secs := time.Since(start).Seconds()
	ch <- result{s.cas, secs, status} //发送结果至ch通道
}

func main() {
	start := time.Now()
	ch := make(chan result)
	//list := readLines("server.config")
	list := readLines(filepath.Dir(os.Args[0]) + "/server.config")
	for _, s := range list {
		go makeRequest(s, ch) //并发执行makeRequest
	}
	//rs := results{}
	for range list {
		fmt.Println(<-ch) //显示ch通道接收到的结果
		//rs = append(rs, <-ch)
	}

	/*sort.Sort(rs)
	for _, r := range rs {
		fmt.Println(r)
	}*/
	fmt.Printf("%.2fms elapsed total\n", time.Since(start).Seconds()*1000)
}

func readLines(file string) servers {
	fi, err := os.Open(file)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
	defer fi.Close()
	br := bufio.NewReader(fi)
	count := 1
	list := servers{}
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		list = append(list, server{count, string(a)})
		count++
	}
	return list
}

/*func (r results) Len() int {
	return len(r)
}

func (r results) Less(i, j int) bool {
	return r[i].cas < r[j].cas
}

func (r results) Swap(i, j int) {
	r[i].cas, r[j].cas = r[j].cas, r[i].cas
	r[i].secs, r[j].secs = r[j].secs, r[i].secs
	r[i].status, r[j].status = r[j].status, r[i].status
}*/

func (r result) String() string {
	return fmt.Sprintf("%.2fms elapsed with response : CAS%02d %s", r.secs*1000, r.cas, r.status)
}
