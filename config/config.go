package config

import (
    "context"
    "fmt"
    "log"
    "time"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Client

func ConnectDB() *mongo.Client {
    uri := "mongodb+srv://skillbarter_rupa:rupaskillbarter@barter.gdyiob8.mongodb.net/?retryWrites=true&w=majority&appName=Barter"

    client, err := mongo.NewClient(options.Client().ApplyURI(uri))
    if err != nil {
        fmt.Println("MongoDB connection error1:", err)

        log.Fatal(err)
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    err = client.Connect(ctx)
    if err != nil {
        fmt.Println("MongoDB connection error2:", err)

        log.Fatal(err)
    }

    // Optional: Ping the database to test connection
    err = client.Ping(ctx, nil)
    if err != nil {
        fmt.Println("MongoDB connection error3:", err)

        log.Fatal(err)
    }

    fmt.Println("MongoDB connected ")
    DB = client
    return client
}
