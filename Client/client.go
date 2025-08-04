package Client

import (
	"GolandPro/proto"
	"context"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"time"
)



type Clinet struct {
	conn net.Conn
}

func (c *Clinet) Send() {
	c.conn.Write([]byte("headeData"))
}

func NewClinet(conn net.Conn) {
	//c := Clinet{conn: conn}
	//server.MyServer.AllConnect = append(server.MyServer.AllConnect, c)
}

func main() {

	//conn, err := net.Dial("tcp","127.0.0.1:8888")
	//if err != nil {
	//	log.Fatalf("did not connect: %v", err)
	//}
	//defer conn.Close()
	//for {
	//	readBuf := make([]byte, 1024)
	//	bsCount, err := conn.Read(readBuf)Dial
	//	if err != nil {
	//		fmt.Println("conn read err =", err)
	//		return
	//	}
	//	fmt.Println("Read suc =", string(readBuf[:bsCount]))
	//}


	//conn, err := grpc.Dial("127.0.0.1:8888", grpc.WithInsecure(), grpc.WithBlock())
	//if err != nil {
	//	log.Fatalf("did not connect: %v", err)
	//}
	//defer conn.Close()
	//Greet(conn)
	//MultiReqMultiResp(conn)
}

func MultiReqMultiResp(conn *grpc.ClientConn) {
	client := proto.NewIGrpcStremServiceClient(conn)
	// 超时设置
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// 远程调用
	intOutStream, err := client.MultiReqMultiResp(ctx)
	if err != nil {
		log.Fatalf(" 远程调用异常 : %v", err)
		return
	}

	// 发送多条记录到服务端
	for i := 0; i < 2; i++ {
		if err = intOutStream.Send(&proto.SingleRequest{Id: int32(400 + i)}); err != nil {
			log.Fatalf("通过流发送数据异常 : %v", err)
			return
		}
	}

	// 服务端一直在接收，直到收到io.EOF为止
	// 因此，这里必须发送io.EOF到服务端，让服务端知道发送已经结束(很重要)
	intOutStream.CloseSend()

	// 接收服务端发来的数据
	for {
		singleResponse, err := intOutStream.Recv()
		if err == io.EOF {
			log.Printf("4. 获取数据完毕")
			break
		} else if err != nil {
			log.Fatalf("4. 接收服务端数据异常 : %v", err)
			break
		}

		log.Printf("4. 收到服务端响应, id : %d, name : %s", singleResponse.GetId(), singleResponse.GetName())
	}

}

func multiReqMultiResp(ctx context.Context, client proto.IGrpcStremServiceClient) error {

	return nil
}

func Greet(conn *grpc.ClientConn) {
	client := proto.NewGreeterClient(conn)
	// 超时设置
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := client.SayHello(ctx, &proto.HelloRequest{Name: "jerry"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.Message)

	h, err := client.SayGb(ctx, &proto.HelloReply{Message: "结束对话"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", h.Name)
}
