package push

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"golang.org/x/crypto/pkcs12"
	"golang.org/x/net/http2"
	"io/ioutil"
	"net/http"
)

const userToken = "7b332c3697a83159bd67afff1f65130120712ce420899488811b0f13f8cb9ee4"
const userToken1 = "3b59ccad062fd39b029ade2219c6a195dbafde6e5c77a28f69cef8ddf9615d50"

const url = "https://api.development.push.apple.com"
const url1 = "https://api.push.apple.com"


var pool *redis.Pool


func Start()  {

	//pool = storage.RedisInitConn()
	//
	//for {
	//	conn := pool.Get()
	//	defer conn.Close()
	//
	//	tokenStr, err := conn.Do("BLPOP", "list1", 10)
	//	if err != nil {
	//		panic(fmt.Sprintf("redis BLPOP error: %v\n", err))
	//	}
	//
	//	//tokenStr, err := redis.String(re, err)
	//	//if err != nil {
	//	//	panic(fmt.Sprintf("redisToString error: %v\n", err))
	//	//}
	//	println("token == ", tokenStr)
	//
	//}


	//cert, err := FromP12File("push/certificate/iosPush_development.p12", "000000")
	cert, err := FromP12File("push/certificate/push_pro_12.p12", "000000")
	if err != nil {
		panic(err)
	}
	host := fmt.Sprintf("%v/3/device/%v", url1, userToken1)


	payload := Payload{
		Aps: APS{Alert: Alert{
							Title: "this is title",
							Body: "body call ",
							ActionLocKey: "PLAY",
			},
			Badge: 8,
			Sound: "bingbong.aiff",

		},
	}
	bodyData, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", bodyData)

	httpReq, err := http.NewRequest(
		"POST",
		host,
		bytes.NewBuffer(bodyData))
	if err != nil {
		panic(fmt.Sprintf("newRequest error: %v\n", err))
	}

	//httpReq.Header.Set(":apns-id", userToken)
	//httpReq.Header.Set(":path", "/3/device/fd9880b6c0b9651c64917c7481d7cac198e2e39ef663fe594b1adbb1c698ed99")
	httpReq.Header.Set("Content-Type", "application/json; charset=utf-8")
	httpReq.Header.Set("apns-topic", "com.kaipanla.www")
	//httpReq.Header.Set("apns-collapse-id", "1")
	httpReq.Header.Set("apns-priority", "10")
	//httpReq.Header.Set("apns-expiration","0"))

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	transport := &http2.Transport{
		TLSClientConfig:tlsConfig,
		//DialTLS: defu
	}

	client := http.Client{
		Transport: transport,

		CheckRedirect: func(
			req *http.Request,
			via []*http.Request) error {
			fmt.Println("Redirect:", req)
			return nil
		},
	}

	response, err := client.Do(httpReq)
	if err != nil {
		panic(fmt.Sprintf("http error: %v\n", err))
	}
	defer  httpReq.Body.Close()
	respData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(fmt.Sprintf("read error: %v\n", err))
	}
	fmt.Println("response :",string(respData))

}


// cert
func FromP12File(filename string, password string) (tls.Certificate, error) {
	p12bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return tls.Certificate{}, err
	}
	return FromP12Bytes(p12bytes, password)
}

// FromP12File loads a PKCS#12 certificate from an in memory byte array
// and returns a tls.Certificate.
func FromP12Bytes(bytes []byte, password string) (tls.Certificate, error) {
	key, cert, err := pkcs12.Decode(bytes, password)
	if err != nil {
		return tls.Certificate{}, err
	}
	return tls.Certificate{
		Certificate: [][]byte{cert.Raw},
		PrivateKey:  key,
		Leaf:        cert,
	}, nil
}


type Alert struct {
	// Title is a short string describing the purpose of the notification.
	// This key was added in iOS 8.2 or newer.
	Title        string   `json:"title,omitempty"`
	TitleLocKey  string   `json:"title-loc-key,omitempty"`
	TitleLocArgs []string `json:"title-loc-args,omitempty"`

	// Subtitle added in iOS 10
	Subtitle string `json:"subtitle,omitempty"`

	// Body is the text of the alert message.
	Body    string   `json:"body,omitempty"`
	LocKey  string   `json:"loc-key,omitempty"`
	LocArgs []string `json:"loc-args,omitempty"`

	// If a string is specified, the system displays an alert that includes the Close and View buttons.
	// The string is used as a key to get a localized string in the current localization
	// to use for the right button’s title instead of “View”.
	ActionLocKey string `json:"action-loc-key,omitempty"`

	// LaunchImage is the filename of an image file in the app bundle,
	// with or without the filename extension.
	// The image is used as the launch image when users tap the action button or move the action slider.
	LaunchImage string `json:"launch-image,omitempty"`

	Badge int `json:"badge,omitempty"`
	Sound string `json:"sound,omitempty"`
}

type APS struct {
	// Alert is a string or alert dictionary.
	Alert Alert `json:"alert,omitempty"`

	// Badge used for define the app icon display rules.
	Badge interface{} `json:"badge,omitempty"`

	// Include this key when you want the system to play a sound.
	// The value of this key is the name of a sound file in your app’s main bundle
	// or in the Library/Sounds folder of your app’s data container.
	// If the sound file cannot be found, or if you specify default for the value,
	// the system plays the default alert sound.
	Sound string `json:"sound,omitempty"`

	// Include this key with a value of 1 to configure a background update notification.
	// When this key is present, the system wakes up your app in the background
	// and delivers the notification to its app delegate.
	ContentAvailable int `json:"content-available,omitempty"`

	// Category represents the notification’s type.
	// This value corresponds to the value in the
	// identifier property of one of your app’s registered categories.
	Category string `json:"category,omitempty"`

	// ThreadID identifier for grouping notifications.
	// If you provide a Notification Content app extension,
	// you can use this value to group your notifications together.
	ThreadID string `json:"thread-id,omitempty"`

	// Mutable is used for Service Extensions introduced in iOS 10.
	// It's value must be one if APNs wants to achieve service extensions.
	MutableContent int `json:"mutable-content,omitempty"`
}

type Payload struct {
	Aps APS `json:"aps,omitempty"`
	Acme1 string `json:"acme1,omitempty"`
}

