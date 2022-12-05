package mongo_test

// import (
// 	"context"
// 	"fmt"
// 	"os"
// 	"reflect"
// 	"testing"
// 	"time"

// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/bson/bsoncodec"
// 	"go.mongodb.org/mongo-driver/bson/bsonrw"
// 	"go.mongodb.org/mongo-driver/bson/bsontype"
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// 	"go.mongodb.org/mongo-driver/mongo"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// )

// type FlatStruct struct {
// 	DT time.Time
// }

// type NestedMapStruct struct {
// 	Data map[string]interface{}
// }

// func getConn() *mongo.Database {
// 	conOpts := options.Client().ApplyURI("mongodb://artem:hometest@localhost:27017/")
// 	if err := conOpts.Validate(); err != nil {
// 		os.Exit(1)
// 	}
// 	cl, err := mongo.Connect(context.TODO(), conOpts)
// 	if err != nil {
// 		os.Exit(1)
// 	}
// 	db := cl.Database("admin")

// 	return db
// }

// func TestWriteFlat(t *testing.T) {
// 	db := getConn()
// 	c := db.Collection("testing")

// 	_, err := c.InsertOne(context.Background(), FlatStruct{DT: time.Now()}, nil)

// 	if err != nil {
// 		t.Fatal(err)
// 	}
// }

// func TestWriteNested(t *testing.T) {
// 	db := getConn()
// 	c := db.Collection("testing")

// 	_, err := c.InsertOne(context.Background(), NestedMapStruct{Data: map[string]interface{}{"dt": time.Now()}}, nil)

// 	if err != nil {
// 		t.Fatal(err)
// 	}
// }

// func TestReadFlat(t *testing.T) {
// 	db := getConn()
// 	c := db.Collection("testing")

// 	id, err := primitive.ObjectIDFromHex("6388ce9384462c91137b3a32")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	var res FlatStruct
// 	err = c.FindOne(context.Background(), primitive.M{"_id": id}, nil).Decode(&res)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	fmt.Printf(`
// 		Type: %T
// 		Value: %v
// 		Inner type: %T
// 		Inner Valie: %v
// 	`, res, res, res.DT, res.DT)

// }

// func TestReadNested(t *testing.T) {
// 	db := getConn()
// 	c := db.Collection("testing")

// 	id, err := primitive.ObjectIDFromHex("6388ce96f4cb4048bbfc297f")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	var res NestedMapStruct
// 	err = c.FindOne(context.Background(), primitive.M{"_id": id}, nil).Decode(&res)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	fmt.Printf(`
// 		Type: %T
// 		Value: %v
// 		Inner type: %T
// 		Inner Valie: %v
// 		Map value type": %T
// 		Map value value: %v
// 	`, res, res, res.Data, res.Data, res.Data["dt"], res.Data["dt"])

// }

// func TestReadNestedWithCustomRegistry(t *testing.T) {
// 	db := getConn()
// 	c := db.Collection("testing")

// 	cr := CustomRegistry()

// 	id, err := primitive.ObjectIDFromHex("6388ce96f4cb4048bbfc297f")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	var res NestedMapStruct
// 	sr := c.FindOne(context.Background(), primitive.M{"_id": id}, nil)
// 	if sr.Err() != nil {
// 		t.Fatal(err)
// 	}

// 	raw, err := sr.DecodeBytes()
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	bson.UnmarshalWithRegistry(cr, raw, &res)

// 	fmt.Printf(`
// 		Type: %T
// 		Value: %v
// 		Inner type: %T
// 		Inner Valie: %v
// 		Map value type": %T
// 		Map value value: %v
// 	`, res, res, res.Data, res.Data, res.Data["dt"], res.Data["dt"])

// }

// func CustomRegistry() *bsoncodec.Registry {
// 	rb := bsoncodec.NewRegistryBuilder()

// 	bsoncodec.DefaultValueEncoders{}.RegisterDefaultEncoders(rb)
// 	bsoncodec.DefaultValueDecoders{}.RegisterDefaultDecoders(rb)

// 	// rb.RegisterTypeDecoder(reflect.TypeOf(primitive.DateTime(0)), CustomMapDecoder{})
// 	rb.RegisterTypeMapEntry(bsontype.DateTime, reflect.TypeOf(time.Time{}))
// 	return rb.Build()
// }

// //DecodeValue(DecodeContext, bsonrw.ValueReader, reflect.Value) error

// type CustomMapDecoder struct{}

// func (cmd CustomMapDecoder) DecodeValue(dc bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
// 	if !val.CanSet() || val.Type() != reflect.TypeOf(primitive.DateTime(0)) {
// 		return bsoncodec.ValueDecoderError{Name: "AAAAA", Types: []reflect.Type{reflect.TypeOf(primitive.DateTime(0))}, Received: val}
// 	}

// 	elem, err := cmd.DecodeDT(dc, vr, reflect.TypeOf(time.Time{}))
// 	if err != nil {
// 		return err
// 	}

// 	val.Set(elem)
// 	return nil
// }

// func (CustomMapDecoder) DecodeDT(dc bsoncodec.DecodeContext, vr bsonrw.ValueReader, t reflect.Type) (reflect.Value, error) {
// 	if t != reflect.TypeOf(primitive.DateTime(0)) {
// 		return reflect.Value{}, bsoncodec.ValueDecoderError{
// 			Name:     "AAAA",
// 			Types:    []reflect.Type{reflect.TypeOf(primitive.DateTime(0))},
// 			Received: reflect.Zero(t),
// 		}
// 	}

// 	var dt int64
// 	var err error
// 	switch vrType := vr.Type(); vrType {
// 	case bsontype.DateTime:
// 		dt, err = vr.ReadDateTime()
// 	case bsontype.Null:
// 		err = vr.ReadNull()
// 	case bsontype.Undefined:
// 		err = vr.ReadUndefined()
// 	default:
// 		return reflect.Value{}, fmt.Errorf("cannot decode %v into a DateTime", vrType)
// 	}
// 	if err != nil {
// 		return reflect.Value{}, err
// 	}

// 	return reflect.ValueOf(time.UnixMilli(dt)), nil
// }
