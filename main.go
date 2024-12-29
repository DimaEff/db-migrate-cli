package main

import (
 "bufio"
 "context"
 "database/sql"
 "fmt"
 "os"
 "strings"

 "github.com/db-migrate-cli/pkg/console"
 "github.com/db-migrate-cli/pkg/presets"

 _ "github.com/lib/pq"
 "go.mongodb.org/mongo-driver/mongo"
 "go.mongodb.org/mongo-driver/mongo/options"
)

const newConnectionOptionValue = console.OptionValue(-1)

func main() {
 fmt.Println("Подключение к вспомогательной базе данных...")
 cliDb, err := presets.ConnectToCliDb("./main.db")
 if err != nil {
  fmt.Println("Не удалось подключиться к вспомогательной базе данных")
  fmt.Println(err)
  os.Exit(1)
 }
 fmt.Println("Успешное подключился к вспомогательной базе данных")

 connections, err := presets.GetConnections(cliDb)
 if err != nil {
  fmt.Println("Не удалось получить сохраненные подключения")
  fmt.Println(err)
  os.Exit(1)
 }

 reader := bufio.NewReader(os.Stdin)
 var selectedConnection *presets.Connection
 if len(connections) > 0 {
  fmt.Println("Выберите сохраненные подключения")
  options := []console.SelectListOption{
   {Value: newConnectionOptionValue, Text: "<Создать новое подключение>"},
  }
  for _, con := range connections {
   options = append(options, console.SelectListOption{Value: console.OptionValue(con.ID), Text: con.Name})
  }
  selector, err := console.NewSelectList(options)
  selectedOption, err := selector.Run()
  fmt.Printf("Выбрано %s\n", selectedOption.Text)
  if err != nil {
   fmt.Println(err)
   os.Exit(1)
  }
  if selectedOption.Value == newConnectionOptionValue {
   selectedConnection = handleCreateNewConnection(reader, cliDb)
  } else {
   for _, con := range connections {
    if con.ID == int(selectedOption.Value) {
     selectedConnection = &con
     break
    }
   }
  }
 } else {
  selectedConnection = handleCreateNewConnection(reader, cliDb)
 }

 sourceClient := connectToMongoDb(selectedConnection.SourceURL)
 targetClient := connectToPgDb(selectedConnection.TargetURL)

 defer cliDb.Close()
 defer sourceClient.Disconnect(context.TODO())
 defer targetClient.Close()

 isEnd := requestCLIParam(reader, "Закончить выполнение?(y/n)")
 if isEnd == "y" {
  fmt.Println("Выполнение завершено")
  return
 }

 fmt.Println(selectedConnection.SourceURL, selectedConnection.TargetURL)
}

func handleCreateNewConnection(reader *bufio.Reader, cliDb *sql.DB) *presets.Connection {
 sourcDbUrl := requestCLIParam(reader, "Введите url исходной базы данных")
 targetDbUrl := requestCLIParam(reader, "Введите url целевой базы данных")
 connectionName := requestCLIParam(reader, "Введите название подключения")

 fmt.Println("Сохранение настроек подключения...")
 con, err := presets.SaveConnectionsURLs(cliDb, connectionName, sourcDbUrl, targetDbUrl)
 if err != nil {
  fmt.Println("Ошибка сохранения настроек подключения")
  fmt.Println(err)
  os.Exit(1)
 }
 fmt.Println("Настройки подключения успешно сохранены")
 return con
}

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
 fmt.Println(fmt.Sprintf("Подключение к базе данных mongodb...\n%s", uri))
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
 fmt.Println(fmt.Sprintf("Подключение к базе данных postgres...\n%s", uri))
 db, err := sql.Open("postgres", uri)
 if err != nil {
  fmt.Println("Ошибка подключения к базе данных")
  fmt.Println(err)
  os.Exit(1)
 }
 fmt.Println("Подключение к базе данных postgresql успешно")
 return db
}
