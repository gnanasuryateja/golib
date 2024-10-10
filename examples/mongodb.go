package main

// import (
// 	"context"
// 	"fmt"

// 	"github.com/gnanasuryateja/golib/datastore/database/mongodb"
// 	"go.mongodb.org/mongo-driver/bson"
// )

// type d struct {
// 	Id   string `json:"_id" bson:"_id"`
// 	Name string `json:"given_name" bson:"given_name"`
// }

// func main() {
// 	ctx := context.Background()
// 	mongoStoreConfig := mongodb.MongoStoreConfig{
// 		Uri:      "",
// 		DbName:   "",
// 		Username: "",
// 		Password: "",
// 	}
// 	db, cancel, err := mongodb.NewMongoStoreClient(ctx, mongoStoreConfig)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	defer db.CloseDB(ctx, cancel)
// 	// testing db.HealthCheck
// 	err = db.HealthCheck(ctx)
// 	if err != nil {
// 		fmt.Println(err)
// 	} else {
// 		fmt.Println("HealthCheck successful...:)")
// 	}
// 	data1 := d{
// 		Id:   "123",
// 		Name: "surya",
// 	}
// 	// testing db.AddData
// 	resAddData, err := db.AddData(ctx, "sample", data1)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	fmt.Println("res:", resAddData)
// 	data2 := d{
// 		Id:   "abc",
// 		Name: "susee",
// 	}
// 	data3 := d{
// 		Id:   "abcd",
// 		Name: "surya",
// 	}
// 	var ds []d
// 	ds = append(ds, data1, data2, data3)
// 	// testing db.AddMultipleData
// 	resAddMultipleData, err := db.AddMultipleData(ctx, "sample", ds)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	fmt.Println("resM:", resAddMultipleData)
// 	filter := make(map[string]any)
// 	filter["_id"] = "123"
// 	filter["given_name"] = "surya"
// 	// testing db.GetData
// 	resGd, err := db.GetData(ctx, "sample", filter)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	fmt.Printf("resGd:%v\n", resGd)
// 	// testing db.GetMultipleData
// 	resGmd, err := db.GetMultipleData(ctx, "sample", filter)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	fmt.Printf("resGmd:%v\n", resGmd)
// 	// testing db.UpdateData
// 	// building the updateData
// 	key := "new_field"
// 	value := "new field added"
// 	updateData := bson.M{
// 		"$set": bson.M{
// 			key: value,
// 		},
// 	}
// 	resUd, err := db.UpdateData(ctx, "sample", filter, updateData)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	fmt.Printf("resUd:%v\n", resUd)
// 	// testing db.UpdateMultipleData
// 	resUmd, err := db.UpdateMultipleData(ctx, "sample", filter, updateData)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	fmt.Printf("resUmd:%v\n", resUmd)
// 	// testing db.DeleteData
// 	resDd, err := db.DeleteData(ctx, "sample", filter)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	fmt.Printf("resDd:%v\n", resDd)
// 	// testing db.DeleteMultipleData
// 	resDmd, err := db.DeleteMultipleData(ctx, "sample", filter)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	fmt.Printf("resDmd:%v\n", resDmd)
// }
