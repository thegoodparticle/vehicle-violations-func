package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/thegoodparticle/vehicle-data-layer/internal/config"
	"github.com/thegoodparticle/vehicle-data-layer/internal/migration"
	"github.com/thegoodparticle/vehicle-data-layer/internal/routes"
	"github.com/thegoodparticle/vehicle-data-layer/internal/store"
)

func main() {
	configs := config.GetConfig()

	connection := GetConnection()
	repository := store.NewAdapter(connection)

	log.Print("Waiting service starting.... ", nil)

	errors := Migrate(connection)
	if len(errors) > 0 {
		for _, err := range errors {
			if err != nil {
				log.Panic("Error on migrate: ", err)
			}
		}
	}

	if err := checkTables(connection); err != nil {
		log.Panic("", err)
	}

	port := fmt.Sprintf(":%v", configs.Port)
	router := routes.NewRouter().SetRouters(repository)
	log.Print("Service running on port ", port)

	server := http.ListenAndServe(port, router)
	log.Fatal(server)
}

func Migrate(connection *dynamodb.DynamoDB) []error {
	var errors []error

	callMigrateAndAppendError(&errors, connection, migration.NewMigration())

	return errors
}

func callMigrateAndAppendError(errors *[]error, connection *dynamodb.DynamoDB, migration *migration.Migration) {
	err := migration.Migrate(connection)
	if err != nil {
		*errors = append(*errors, err)
	}
}

func checkTables(connection *dynamodb.DynamoDB) error {
	response, err := connection.ListTables(&dynamodb.ListTablesInput{})
	if response != nil {
		if len(response.TableNames) == 0 {
			log.Print("Tables not found: ", nil)
		}
		for _, tableName := range response.TableNames {
			log.Print("Table found: ", *tableName)
		}
	}
	return err
}

func GetConnection() *dynamodb.DynamoDB {
	// Initialize a session that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials
	// and region from the shared configuration file ~/.aws/config.
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	return dynamodb.New(sess)
}
