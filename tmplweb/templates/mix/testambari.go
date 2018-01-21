package main

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/rakanalh/scheduler"
	"github.com/rakanalh/scheduler/storage"

	tk "github.com/eaciit/toolkit"
	"gopkg.in/resty.v1"
)

func main() {
	s := scheduler.New(storage.NewMemoryStorage())
	s.RunEvery(15*time.Second, GetResult01)
	s.Start()
	s.Wait()
	log.Println("Tasks ended...")

}

func GetResult01() {
	tstamp := time.Now()
	t1 := tstamp.Unix()
	t0 := tstamp.Add(-15 * time.Second).Unix()

	resp, _ := GetAmbariResponse("http://10.20.172.190:8081", "STANCDEV1TDH", "a1577318", "a1577318", "", "YARN/components/NODEMANAGER", "", "", []string{"metrics/disk/read_bps._sum", "metrics/disk/write_bps._sum"}, "", t0, t1)
	data := make(map[string]interface{})
	json.Unmarshal(resp, &data)
	tk.Printfn("RESP %v", data)
}

func buildClusterQuery(fieldQuery []string, t0, t1 int64) string {
	// fields := []string{"metrics/cpu/System._avg", "metrics/cpu/User._avg", "metrics/memory/Total._avg", "metrics/memory/Use._avg"}
	flist := make([]string, 0)
	for _, fl := range fieldQuery {
		fq := tk.Sprintf("%s[%d,%d,%d]", fl, t0, t1, 15)
		flist = append(flist, fq)
	}
	fqs := strings.Join(flist, ",")
	return "fields=" + fqs
}

func GetAmbariResponse(apiURL, cluster, username, passwd, extraPath, service, hostname, hostcomp string, fieldQuery []string, queryString string, t0, t1 int64) ([]byte, error) {
	qryString := ""
	if len(fieldQuery) > 0 {
		qryString = buildClusterQuery(fieldQuery, t0, t1)
	} else {
		qryString = queryString
	}
	url := apiURL + "/clusters/" + cluster
	if strings.TrimSpace(extraPath) != "" {
		url += "/" + extraPath
	}
	if strings.TrimSpace(hostname) != "" {
		url += "/hosts/" + hostname
	}
	if strings.TrimSpace(service) != "" {
		url += "/services/" + service
	}
	if strings.TrimSpace(hostcomp) != "" {
		url += "/host_components/" + hostcomp
	}
	tk.Println("URL: ", url+"?"+qryString)
	resp, err := resty.R().
		SetBasicAuth(username, passwd).
		SetQueryString(qryString).
		Get(url)
	tk.Printfn("Error: %v", err)
	tk.Printfn("Response Status Code: %v", resp.StatusCode())
	tk.Printfn("Response Status: %v", resp.Status())
	tk.Printfn("Response Time: %v", resp.Time())
	tk.Printfn("Response Received At: %v", resp.ReceivedAt())
	if err != nil {
		return nil, err
	}
	return resp.Body(), nil
}
