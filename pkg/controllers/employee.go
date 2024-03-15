package controllers

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Sourjaya/go-hrms/pkg/models"
	"github.com/gofiber/fiber"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const dbName = "hrms"
const colName = "watchlist"

// Creating a instance
var collection *mongo.Collection

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading env variables")
	}
	connectionString := os.Getenv("MONGODB_URI")
	//client option
	clientOption := options.Client().ApplyURI(connectionString)

	//connect to mongodb
	client, err := mongo.Connect(context.TODO(), clientOption)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("MongoDB connection success")

	collection = client.Database(dbName).Collection(colName)
	fmt.Println("Collection created")
}

func insertOneEmployee(employee *models.Employee) interface{} {
	inserted, err := collection.InsertOne(context.Background(), employee)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted 1 user in db with id: ", inserted.InsertedID)
	return inserted.InsertedID
}

func deleteOneEmployee(id primitive.ObjectID, ctx *fiber.Ctx) {
	filter := bson.M{"_id": id}
	deleteCount, err := collection.DeleteOne(ctx.Context(), filter)

	if err != nil {
		ctx.SendStatus(500)
		return
	}
	if deleteCount.DeletedCount < 1 {
		ctx.SendStatus(404)
		return
	}
	log.Println("User got deleted with delete count: ", deleteCount.DeletedCount)
	ctx.Status(200).JSON("record deleted")
}

func GetAllEmployees(ctx *fiber.Ctx) {
	var employees []models.Employee
	query := bson.D{{}}
	cursor, err := collection.Find(ctx.Context(), query)
	if err != nil {
		ctx.Status(500).SendString(err.Error())
	}
	if err := cursor.All(ctx.Context(), &employees); err != nil {
		ctx.Status(500).SendString(err.Error())
	}
	ctx.JSON(employees)
}

func GetEmployeeByID(ctx *fiber.Ctx) {
	id := ctx.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.Status(500).SendString(err.Error())
		return
	}
	employee := models.Employee{}
	err = collection.FindOne(context.Background(), bson.M{"_id": oid}).Decode(&employee)
	if err != nil {
		ctx.Status(500).SendString(err.Error())
		return
	}
	ctx.JSON(employee)
}

func NewEmployee(ctx *fiber.Ctx) {

	//w.Header().Set("Content-Type", "application/json")
	//w.Header().Set("Allow-Control-Allow-Methods", "POST")

	ctx.Append("Content-Type", "application/json")
	ctx.Append("Allow-Control-Allow-Methods", "POST")
	var employee *models.Employee

	if err := ctx.BodyParser(&employee); err != nil {
		ctx.Status(400).SendString(err.Error())
		return
	}
	ob := insertOneEmployee(employee)
	employee.ID = ob.(primitive.ObjectID)
	ctx.JSON(employee)
}

func DeleteEmployee(ctx *fiber.Ctx) {
	ctx.Append("Content-Type", "application/json")
	ctx.Append("Allow-Control-Allow-Methods", "DELETE")
	id := ctx.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.SendStatus(400)
		return
	}
	deleteOneEmployee(oid, ctx)
}
