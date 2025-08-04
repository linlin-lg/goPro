package main

import (
	"GolandPro/Connect"
	"GolandPro/proto"
	"GolandPro/server"
	"GolandPro/storage"
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"
	"github.com/gorilla/websocket"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"net"
	"net/http"
)

var (
	rootCrt   = "./config/ca.crt"
	serverCrt = "./config/server.crt"
	serverKey = "./config/server.pem"
)

type weatherModel struct {
	Day   string `json:"day"`
	Weath string `json:"weath"`
}

func collyRequestHanld(r *colly.Request) {
	r.Headers.Set("Host", "https://tianqi.2345.com")
	fmt.Println("visiting", r.URL)
}

func internetWorm() {
	c := colly.NewCollector()
	c.OnRequest(collyRequestHanld)

	c.OnResponse(func(r *colly.Response) {
		htmlDoc, err := goquery.NewDocumentFromReader(bytes.NewReader(r.Body))
		if err != nil {
			panic(err)
		}
		htmlDoc.Find(".weaul li a").Each(func(i int, s *goquery.Selection) {
			day := s.Find(".weaul_q .fl").Text()
			w := s.Find(".weaul_z").Text()
			querySql := "select weather from weatherList where date=?"
			rows, _ := storage.Query(querySql, day)
			defer rows.Close()
			//for rows.Next() {
			//	rows.Scan()
			//}
			if rows.Next() {
				updateSql := "update weatherList set `weather`=?,`w_image`=? where date=?"
				storage.Exec(updateSql, w, "", day)
			} else {
				insertSql := "insert into weatherList(date, weather, w_image) values(?,?,?)"
				storage.Exec(insertSql, day, w, "")
			}
		})
	})
	rule := colly.LimitRule{
		Parallelism: 8,
	}
	c.Limit(&rule)
	c.Visit("https://www.tianqi.com/guangzhou/15/")
}

func main() {

	//storage.ConnectSql()
	//go internetWorm()
	//startGinServer()
	//startTlsServer()
	//startGprcServer()
	//str := AuthSign("kpl-2024:2:5.15.0.3:98:14c15b0eb35bce4c97876e44247351a6d42bc807:0")
	//fmt.Println(str)
}

func luint32(b []byte) uint32 {
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
}

func rotateLeft32(x uint32, k int) uint32 {
	const n = 32
	s := uint(k) & (n - 1)
	//fmt.Println("s=", s)
	ss := x<<s | x>>(n-s)
	//fmt.Println("ss=", ss)
	return ss
}

func putUint32(b []byte, v uint32) {
	//fmt.Println("v=", v)
	b[0] = byte(v)
	//fmt.Println("b0=", byte(v))
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 24)
	//fmt.Println("blist=", b)
}

func AuthSign(str string) string {
	bs := []byte(str)
	a, b, c, d := uint32(8), uint32(16), uint32(32), uint32(64)
	nn := len(bs)
	//fmt.Println("nn=", nn)
	for i := 0; i < nn; i += 64 {
		aa, bb, cc, dd := a, b, c, d
		p := bs[i:]
		if len(p) < 64 {
			tmp := make([]byte, 64)
			copy(tmp, p)
			tmp[63] = byte(nn)
			p = tmp
		}
		p = p[:64:64]
		x0 := luint32(p[:])
		x1 := luint32(p[4:])
		x2 := luint32(p[8:])
		x3 := luint32(p[12:])
		x4 := luint32(p[16:])
		x5 := luint32(p[20:])
		x6 := luint32(p[24:])
		x7 := luint32(p[28:])
		x8 := luint32(p[32:])
		x9 := luint32(p[36:])
		xa := luint32(p[40:])
		xb := luint32(p[44:])
		xc := luint32(p[48:])
		xd := luint32(p[52:])
		xe := luint32(p[56:])
		xf := luint32(p[60:])

		a = b + rotateLeft32((((c^d)&b)^d)+a+x0+0xd76aa478, 7)
		fmt.Println("a=", a)
		d = a + rotateLeft32((((b^c)&a)^c)+d+x1+0xe8c7b756, 12)
		fmt.Println("d=", d)
		c = d + rotateLeft32((((a^b)&d)^b)+c+x2+0x242070db, 17)
		fmt.Println("c=", c)
		b = c + rotateLeft32((((d^a)&c)^a)+b+x3+0xc1bdceee, 22)
		fmt.Println("c=", b)
		a = b + rotateLeft32((((c^d)&b)^d)+a+x4+0xf57c0faf, 7)
		fmt.Println("c=", a)
		d = a + rotateLeft32((((b^c)&a)^c)+d+x5+0x4787c62a, 12)
		fmt.Println("c=", d)
		c = d + rotateLeft32((((a^b)&d)^b)+c+x6+0xa8304613, 17)
		fmt.Println("c=", c)
		b = c + rotateLeft32((((d^a)&c)^a)+b+x7+0xfd469501, 22)
		fmt.Println("c=", b)
		a = b + rotateLeft32((((c^d)&b)^d)+a+x8+0x698098d8, 7)
		fmt.Println("c=", a)
		d = a + rotateLeft32((((b^c)&a)^c)+d+x9+0x8b44f7af, 12)
		fmt.Println("c=", d)
		c = d + rotateLeft32((((a^b)&d)^b)+c+xa+0xffff5bb1, 17)
		fmt.Println("c=", c)
		b = c + rotateLeft32((((d^a)&c)^a)+b+xb+0x895cd7be, 22)
		fmt.Println("c=", b)
		a = b + rotateLeft32((((c^d)&b)^d)+a+xc+0x6b901122, 7)
		fmt.Println("c=", a)
		d = a + rotateLeft32((((b^c)&a)^c)+d+xd+0xfd987193, 12)
		fmt.Println("c=", d)
		c = d + rotateLeft32((((a^b)&d)^b)+c+xe+0xa679438e, 17)
		fmt.Println("c=", c)
		b = c + rotateLeft32((((d^a)&c)^a)+b+xf+0x49b40821, 22)
		fmt.Println("c=", b)
		a += aa
		b += bb
		c += cc
		d += dd
	}
	fmt.Println("tmpa=", a)
	fmt.Println("tmpb=", b)
	fmt.Println("tmpc=", c)
	fmt.Println("tmpd=", d)
	var digest [16]byte
	putUint32(digest[0:], a)
	putUint32(digest[4:], b)
	putUint32(digest[8:], c)
	putUint32(digest[12:], d)
	fmt.Println("digest=", digest)
	return fmt.Sprintf("%x", digest)
}

func startGprcServer() {

	listen, err := net.Listen("tcp", "127.0.0.1:8888")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	g := server.InitServer()
	// 实例化grpc服务端
	s := grpc.NewServer()
	proto.RegisterGreetServer(s, g)
	proto.RegisterIGrpcStremServiceServer(s, g)
	if err := s.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	for {
		//conn, err := listen.Accept()
		//if err != nil {
		//	fmt.Println("accept err =", err)
		//	continue
		//}
		//Client.NewClinet(conn)
	}
}

func startTlsServer() {

	pem, err := ioutil.ReadFile(rootCrt)
	if err != nil {
		panic(fmt.Sprintf("rootCrt err = %v", err))
	}
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pem) {
		panic(fmt.Sprintf("AppendCertsToPoolErr"))
	}
	crt, err := tls.LoadX509KeyPair(serverCrt, serverKey)
	if err != nil {
		panic(fmt.Sprintf("tlsCer err = %v", err))
	}
	tlsConfig := tls.Config{
		Certificates: []tls.Certificate{crt},
		RootCAs:      certPool,
		ClientCAs:    certPool,
		//ClientAuth:   tls.RequireAndVerifyClientCert,
		//MinVersion:   tls.VersionTLS12,
		//InsecureSkipVerify: true,

		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		},
	}
	listen, err := tls.Listen("tcp", "127.0.0.1:8888", &tlsConfig)
	if err != nil {
		panic(fmt.Sprintf("tlslisten err = %v", err))
	}

	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("accept err =", err)
			continue
		}
		go connetHanld(conn)
	}

}

func connetHanld(c net.Conn) {
	for {
		readBuf := make([]byte, 1024)
		bsCount, err := c.Read(readBuf)
		if err != nil {
			fmt.Println("conn read err =", err)
			return
		}
		fmt.Println("Read suc =", string(readBuf[:bsCount]))
		code, err := c.Write(readBuf[:bsCount])
		if err != nil {
			fmt.Println("conn write err =", err)
			return
		}
		fmt.Println("write suc = ", code)
	}
}

func startGinServer() {
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Next()
	})
	//g1 := r.Group("/kpl")

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		// Allow connections from any Origin
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	r.GET("/ws", func(c *gin.Context) {
		if conn, err := upgrader.Upgrade(c.Writer, c.Request, c.Request.Header); err != nil {
			fmt.Println("Upgrade err = ", err)
		} else {
			achan := make(chan interface{}, 1000)
			c := Connect.InitConnect(conn, achan)
			c.SendMessage(200, []byte("发送成功"))
		}
	})

	//r.GET("/socket/ws", func(c *gin.Context) {
	//	if conn, err := upgrader.Upgrade(c.Writer, c.Request, c.Request.Header); err != nil {
	//		fmt.Println("Upgrade err = ", err)
	//	} else {
	//		returnData(conn, "socket data")
	//	}
	//})

	//r.GET("/list", func(c *gin.Context) {
	//	userName, _ := c.Params.Get("userName")
	//	userId, _ := c.Params.Get("userId")
	//	fmt.Println("userName = %d, userId = %d", userName, userId)
	//
	//	querySql := "select date,weather from weatherList"
	//	rows, _ := storage.Query(querySql)
	//	defer rows.Close()
	//	var date, weather string
	//	var group []weatherModel
	//	for rows.Next() {
	//		err := rows.Scan(&date, &weather)
	//		if err != nil {
	//			fmt.Println("rows scan err =", err)
	//			return
	//		}
	//		group = append(group, weatherModel{date, weather})
	//	}
	//	returnSuccess(c, 200, group)
	//})
	//r.POST("/update", func(c *gin.Context) {
	//	userName, _ := c.Params.Get("userName")
	//	userId, _ := c.Params.Get("userId")
	//	fmt.Println("userName = %s, userId = %s", userName, userId)
	//	returnSuccess(c, 200, nil)
	//})
	r.Run(":8888")
}

func returnSuccess(g *gin.Context, code int, items interface{}) {
	fmt.Println(items)
	result := &JsonStruct{Code: code, Items: items, Msg: ""}
	//变成json输出
	if dataByte, err := json.Marshal(&result); err != nil {
		g.Writer.WriteString(err.Error())
	} else {
		g.Writer.WriteString(string(dataByte))
	}
}

func returnData(conn *websocket.Conn, items interface{}) {
	fmt.Println(items)
	result := &JsonStruct{Code: 200, Items: items, Msg: ""}
	//变成json输出
	if dataByte, err := json.Marshal(&result); err != nil {
		conn.WriteMessage(0, nil)
	} else {
		conn.WriteMessage(200, dataByte)
	}

}

type JsonStruct struct {
	Code  int         `json:"code"`
	Msg   string      `json:"msg"`
	Items interface{} `json:"items"`
}
