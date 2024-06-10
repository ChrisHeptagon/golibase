package server

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func StartServer(db *sql.DB) {
	r := gin.New()
	r.SetTrustedProxies(nil)
	var handler gin.HandlerFunc
	if os.Getenv("MODE") == "DEV" {
		devServer, err := url.Parse(fmt.Sprintf("http://localhost:%s", os.Getenv("DEV_PORT")))
		if err != nil {
			fmt.Println("Error parsing dev server URL: ", err)
		}
		handler = func(c *gin.Context) {
			(*c).Request.Host = devServer.Host
			(*c).Request.URL.Host = devServer.Host
			(*c).Request.URL.Scheme = devServer.Scheme
			(*c).Request.RequestURI = ""

			if (*c).Request.URL.Path == "/" && (*c).Request.URL.RawQuery == "" {
				(*c).Writer.WriteHeader(http.StatusSwitchingProtocols)
				var ws websocket.Upgrader = websocket.Upgrader{
					HandshakeTimeout: 10 * time.Second,
					CheckOrigin: func(r *http.Request) bool {
						return true
					},
				}
				conn, err := ws.Upgrade((*c).Writer, (*c).Request, nil)
				if err != nil {
					fmt.Println("Error upgrading websocket: ", err)
					return
				}
				defer conn.Close()
				for {
					msgT, msgB, err := conn.ReadMessage()
					if err != nil {
						fmt.Println("Error reading message: ", err)
					}
					fmt.Printf("Message Type: %d\n", msgT)
					fmt.Printf("Message: %s\n", msgB)
					err = conn.WriteMessage(websocket.TextMessage, []byte("Hello from server"))
					if err != nil {
						fmt.Println("Error writing message: ", err)
						return
					}
				}
			}

			devServerResponse, err := http.DefaultClient.Do((*c).Request)
			if err != nil {
				fmt.Println("Error sending request to dev server: ", err)
				(*c).Writer.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf((*c).Writer, "Error sending request to dev server: %v", err)
				return
			}
			(*c).Writer.WriteHeader(devServerResponse.StatusCode)
			(*c).Writer.Header().Set("Content-Type", devServerResponse.Header.Get("Content-Type"))
			io.Copy((*c).Writer, devServerResponse.Body)
		}
		r.GET("/src/*wildcard",
			handler,
		)
		r.GET("/@vite/client",
			handler,
		)
		r.GET("/@fs/*wildcard",
			handler,
		)
		r.GET("/node_modules/*wildcard",
			handler,
		)
		r.GET("/.svelte-kit/*wildcard",
			handler,
		)
		r.GET("/@id/*wildcard",
			handler,
		)

	} else if os.Getenv("MODE") == "PROD" {
		r.Use(gzip.Gzip(gzip.DefaultCompression))
		corsConfig := cors.DefaultConfig()
		corsConfig.AllowAllOrigins = true
		r.Use(cors.New(corsConfig))
		nodeServer, err := url.Parse(fmt.Sprintf("http://localhost:%s", os.Getenv("PORT")))
		if err != nil {
			fmt.Println("Error parsing node server URL: ", err)
		}
		handler = gin.HandlerFunc(func(c *gin.Context) {
			c.Request.Host = nodeServer.Host
			c.Request.URL.Host = nodeServer.Host
			c.Request.URL.Scheme = nodeServer.Scheme
			c.Request.RequestURI = ""

			nodeServerResponse, err := http.DefaultClient.Do(c.Request)
			if err != nil {
				fmt.Println("Error sending request to node server: ", err)
				c.Writer.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(c.Writer, "Error sending request to node server: %v", err)
				return
			}
			c.Writer.WriteHeader(nodeServerResponse.StatusCode)
			c.Writer.Header().Set("Content-Type", nodeServerResponse.Header.Get("Content-Type"))
			io.Copy(c.Writer, nodeServerResponse.Body)

		})
	}

	r.GET("/ui/*wildcard",
		handler,
	)

	r.GET("/api/v1/user_schema",
		func(c *gin.Context) {
			rows, err := db.Query(`SELECT schema FROM schemas WHERE schema_name = 'user_schema' LIMIT 1;`)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			defer rows.Close()
			if !rows.Next() {
				c.JSON(http.StatusNotFound, gin.H{"error": "Schema not found"})
				return
			}
			var schema string
			err = rows.Scan(&schema)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, schema)
		},
	)

	fmt.Println("Server running at http://localhost:6701")

	r.Run(":6701")

}
