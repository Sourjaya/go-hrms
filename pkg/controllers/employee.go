package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Sourjaya/go-hrms/pkg/models"
	"github.com/gofiber/fiber"
	"github.com/gorilla/mux"
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
	if err := godotenv.Load("../.env"); err != nil {
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

func insertOneUser(user models.Employee) interface{} {
	inserted, err := collection.InsertOne(context.Background(), user)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted 1 user in db with id: ", inserted.InsertedID)
	return inserted.InsertedID
}

func deleteOneUser(Id string, w http.ResponseWriter) {
	id, err := primitive.ObjectIDFromHex(Id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	filter := bson.M{"_id": id}
	deleteCount, err := collection.DeleteOne(context.Background(), filter)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("User got deleted with delete count: ", deleteCount)
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
	return
}

func NewEmployee(ctx *fiber.Ctx) {

	//w.Header().Set("Content-Type", "application/json")
	//w.Header().Set("Allow-Control-Allow-Methods", "POST")

	ctx.Append("Content-Type", "application/json")
	ctx.Append("Allow-Control-Allow-Methods", "POST")
	var employee models.Employee

	if err := ctx.BodyParser(employee); err != nil {
		ctx.Status(400).SendString(err.Error())
		return
	}
	ob := insertOneUser(employee)
	employee.ID = ob.(primitive.ObjectID)
	ctx.JSON(employee)
	return
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "DELETE")

	params := mux.Vars(r)
	deleteOneUser(params["id"], w)
	json.NewEncoder(w).Encode(params["id"])
}
