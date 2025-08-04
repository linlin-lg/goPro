package server

import (
	"GolandPro/Client"
	"GolandPro/proto"
	"golang.org/x/net/context"
	"io"
	"log"
	"strconv"
	"time"
)

var MyServer *Server

type Server struct {
	AllConnect []Client.Clinet
}

func InitServer() *Server {
	MyServer = &Server{}
	MyServer.AllConnect = make([]Client.Clinet, 1000)
	MyServer.sendGuangbo()
	return  MyServer
}

func (g *Server) SayHello(ctx context.Context, in *proto.HelloRequest) (*proto.HelloReply, error) {
	return &proto.HelloReply{Message: "hi," + in.Name}, nil
}

func (g *Server) SayGb(ctx context.Context, in *proto.HelloReply) (*proto.HelloRequest, error) {
	return &proto.HelloRequest{Name: "bye" + in.Message}, nil
}

// 双向流式 ：集合请求，集合响应
func (s *Server) MultiReqMultiResp(reqStream proto.IGrpcStremService_MultiReqMultiRespServer) error {
	// 简单处理，对于收到的每一条记录都返回一个响应
	for {
		singleRequest, err := reqStream.Recv()

		// 不等于io.EOS表示这是条有效记录
		if err == io.EOF {
			log.Println("4. 接收完毕")
			return nil
		} else if err != nil {
			log.Fatalln("4. 接收时发生异常", err)
			return err
		} else {
			log.Println("4. 接收到数据", singleRequest.GetId())

			id := singleRequest.GetId()

			if sendErr := reqStream.Send(&proto.SingleResponse{Id: id, Name: "4. name-" + strconv.Itoa(int(id))}); sendErr != nil {
				log.Println("4. 返回数据异常数据", sendErr)
				return sendErr
			}
		}
	}
}

func (g *Server)sendGuangbo() {
	time := time.Tick(5 * time.Second)
	for {
		select  {
		case <- time:
			if len(MyServer.AllConnect) == 0 {
				continue
			}
			for _, cli := range MyServer.AllConnect {
				cli.Send()
		  	}
		}

	}
}