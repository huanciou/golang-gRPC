package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"

	pb "go-app/server/proto"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	r := gin.Default()

	conn, err := grpc.Dial("localhost:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Failed to Connection: %v", err)
	}
	defer conn.Close()

	client := pb.NewServiceServerClient(conn)

	r.GET("/api", func(c *gin.Context) {
		var wg sync.WaitGroup

		responseCh := make(chan *pb.SerResponse, 10)

		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				request := &pb.CliRequest{
					RequestName: "Jimmy",
				}
				response, err := client.SayHi(context.Background(), request)
				if err != nil {
					log.Fatalf("Req Failed: %v", err)
					return
				}

				responseCh <- response
			}()
		}

		wg.Wait()

		close(responseCh)

		var responses []*pb.SerResponse
		for response := range responseCh {
			responses = append(responses, response)
		}

		c.JSON(200, gin.H{
			"data": responses,
		})
	})

	r.GET("/api2", func(c *gin.Context) {
		var wg sync.WaitGroup
		var mu sync.Mutex

		// 创建一个通道来收集 JSON 响应
		jsonResponseCh := make(chan map[string]interface{}, 10)

		for i := 0; i < 10; i++ {
			wg.Add(1)

			go func() {
				defer wg.Done()

				res, err := http.Get("http://localhost:8082/api2")
				if err != nil {
					c.JSON(500, gin.H{
						"error": err.Error(),
					})
					return
				}

				defer res.Body.Close()

				var jsonResponse map[string]interface{}
				decoder := json.NewDecoder(res.Body)

				if err := decoder.Decode(&jsonResponse); err != nil {
					c.JSON(500, gin.H{
						"error": err.Error(),
					})
					return
				}

				// 使用互斥锁以确保安全地向通道发送响应
				mu.Lock()
				jsonResponseCh <- jsonResponse

				// 如果已经收集了足够的响应，关闭通道
				if len(jsonResponseCh) == 10 {
					close(jsonResponseCh)
				}
				mu.Unlock()
			}()
		}

		// 等待所有请求完成
		wg.Wait()

		// 从通道中获取所有 JSON 响应
		var allJSONResponses []map[string]interface{}
		for jsonResponse := range jsonResponseCh {
			allJSONResponses = append(allJSONResponses, jsonResponse)
		}

		// 返回所有 JSON 响应
		c.JSON(200, allJSONResponses)
	})

	r.Run(":8080")
}
