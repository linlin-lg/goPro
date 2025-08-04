package test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

const mUrl = "https://www.zhaohaowang.com"

func Pachong() {

	request, err := http.NewRequest(
		http.MethodGet,
		mUrl,
		nil)
	if err != nil {
		panic(err)
	}
	//cookie := "ci_session=u6qi9l7urhof3kv4plftp5t5pqb5tebd; Hm_lpvt_13995e0e78c2dc80389e4d4526856615=1647244479; Hm_lvt_13995e0e78c2dc80389e4d4526856615=1646379112,1646902495,1646961449,1647239756; usto_scrollTop=0"
	//request.Header.Add("Cookie", cookie)
	//request.Header.Add("Referer", "https://www.kaipanla.com/")
	//request.Header.Add("Cache-Control", "max-age=0")
	//request.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.3 Safari/605.1.15")

	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return nil
		},
	}
	resp, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	dataStr := string(data)
	//fmt.Println(dataStr)

	rePhone := `1[3456789]\d\s?\d{4}\s?\d{4}`
	regexp, _ := regexp.Compile(rePhone)
	res := regexp.FindAllStringSubmatch(dataStr, -1)
	for _, result := range res {
		fmt.Println(result)
	}
}


func SqlSet()  {


}