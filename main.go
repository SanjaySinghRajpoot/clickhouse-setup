package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"math/rand"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/google/uuid"
	"github.com/icrowley/fake"
)

func main() {
	conn, err := connect()
	if err != nil {
		panic((err))
	}

	ctx := context.Background()
	rows, err := conn.Query(ctx, "SELECT event,event_type FROM event_data LIMIT 5")
	if err != nil {
		log.Fatal(err)
	}

	// Insert 100 dummy rows
	for i := 0; i < 100; i++ {

		event := fake.Word()
		eventType := fake.Word()
		userID := rand.Intn(1000) + 1
		source := rand.Intn(3) + 1
		utmContent := fake.Word()
		sessionID := uuid.New().String()

		// Insert the dummy data into the events table
		_, err := conn.Query(ctx, `INSERT INTO event_data (event, event_type, user_id, source, utm_content, session_id) VALUES (?, ?, ?, ?, ?, ?)`, event, eventType, userID, source, utmContent, sessionID)
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	for rows.Next() {
		var (
			event, event_type string
		)
		if err := rows.Scan(
			&event,
			&event_type,
		); err != nil {
			log.Fatal(err)
		}
		log.Printf("name: %s, uuid: %s",
			event, event_type)
	}

}
func connect() (driver.Conn, error) {
	var (
		ctx       = context.Background()
		conn, err = clickhouse.Open(&clickhouse.Options{
			Addr: []string{"gecem3ya47.ap-south-1.aws.clickhouse.cloud:9440"},
			Auth: clickhouse.Auth{
				Database: "default",
				Username: "default",
				Password: "c~avAy0F5Tc8M",
			},
			ClientInfo: clickhouse.ClientInfo{
				Products: []struct {
					Name    string
					Version string
				}{
					{Name: "an-example-go-client", Version: "0.1"},
				},
			},

			Debugf: func(format string, v ...interface{}) {
				fmt.Printf(format, v)
			},
			TLS: &tls.Config{
				InsecureSkipVerify: true,
			},
		})
	)

	if err != nil {
		return nil, err
	}

	if err := conn.Ping(ctx); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("Exception [%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		}
		return nil, err
	}
	return conn, nil
}
