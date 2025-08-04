package main

import (
    "GolandPro/push"
    "encoding/json"
    "fmt"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
    "net/http"
    "strings"
)

const prefix = "/kpl"

func main1() {


    push.Start()

    // 连接etcd
    connectEtcd()


    r := gin.Default()
    logger, err := zap.NewProduction()
    if err != nil {
        println("\n logger err : ", err)
    }

    r.Use(func(c *gin.Context) {
        c.Next()
        logger.Info("request info",
            zap.String("path", c.Request.URL.Path),
            zap.Int("status", c.Writer.Status()))
    })

    g1 := r.Group(prefix)


    g1.GET("/list", func(c *gin.Context) {
        token := c.DefaultQuery("Token", "noToken")
        fmt.Println("\n token : ", token)
        d := GetData(c.Request)
        var model interface{}
        err := json.Unmarshal(d, &model)
        if err != nil {
          fmt.Println("\n err : ", err)
        }
        c.JSON(200, model)
    })

    g1.GET("/setting", func(c *gin.Context) {

        c.Redirect(http.StatusMovedPermanently, "https://applhb.longhuvip.com/w1/api/index.php?Token=7d5a372bc287f93b35716bc7652a1e1d&PhoneOSNew=2&VerSion=5.4.0.0&a=GetArrangeIndex&apiv=w29&UserID=799015&c=SysAppVersion")

        //var parame SettingModel
        //if err := c.ShouldBindJSON(&parame); err != nil {
        //    c.XML(http.StatusBadRequest, gin.H{"error": err.Error()})
        //    //c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        //}
        //fmt.Println("\n parame : ", parame)
        //c.String(200, parame.Name)
        //c.XML(http.StatusBadRequest, gin.H{"error": "34"})
    })

    r.POST("/list", func(c *gin.Context) {
        token := c.DefaultPostForm("Token", "noToken")
        fmt.Println("\n token : ", token)
        c.String(200, "do")
    })

    r.Run(":8888")


    //http.HandleFunc("/", ErrWrok(handleWeb))
    //err := http.ListenAndServe(":8888", nil)
    //if err != nil {
    //    panic(err)
    //}
}

func handleWeb(writer http.ResponseWriter, request *http.Request) error {
    path := request.URL.Path

    if strings.Index(path, prefix) == -1 {
        fmt.Printf("\n path err : ", path)
        return UserError("无效路径")
    }
    fmt.Println("\n reqPath ==", path)

    writer.Write(GetData(request))
    return nil
}
