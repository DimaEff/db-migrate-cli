package main

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func requestCLIParam(reader *bufio.Reader, message string) string {
	fmt.Println(message)
	input, err := reader.ReadString('\n')

	if err != nil {
		fmt.Println("Ошибка чтения параметра")
		os.Exit(1)
	}

	return strings.TrimSpace(input)
}

func connectToMongoDb(uri string) *mongo.Client {
	infoMessage := fmt.Sprintf("Подключение к базе данных mongodb...\n%s", uri)
	fmt.Println(infoMessage)
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		fmt.Println("Ошибка подключения к базе данных")
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Подключение к базе данных mongodb успешно")
	return client
}

func connectToPgDb(uri string) *sql.DB {
	infoMessage := fmt.Sprintf("Подключение к базе данных mongodb...\n%s", uri)
	fmt.Println(infoMessage)
	db, err := sql.Open("postgres", uri)
	if err != nil {
		fmt.Println("Ошибка подключения к базе данных")
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Подключение к базе данных postgresql успешно")
	return db
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	sourcDbUrl := requestCLIParam(reader, "Введите url исходной базы данных")
	targetDbUrl := requestCLIParam(reader, "Введите url целевой базы данных")

	sourceClient := connectToMongoDb(sourcDbUrl)
	targetClient := connectToPgDb(targetDbUrl)

	defer sourceClient.Disconnect(context.TODO())
	defer targetClient.Close()

	isEnd := requestCLIParam(reader, "Закончить выполнение?(y/n)")
	if isEnd == "y" {
		fmt.Println("Выполнение завершено")
		return
	}

	fmt.Println(sourcDbUrl, targetDbUrl)
}
