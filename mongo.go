package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/mgocompat"

	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Store struct {
	client *mongo.Client
}

func NewStore(url string) (store *Store, err error) {
	cl, err := mongo.NewClient(options.Client().ApplyURI(url).SetRegistry(mgocompat.NewRegistryBuilder().Build()))
	if err != nil {
		return
	}
	err = cl.Connect(context.Background())
	if err != nil {
		return
	}

	return &Store{
		client: cl,
	}, nil
}

func (s *Store) Save(name string, price decimal.Decimal) (err error) {
	priceFloat, _ := price.BigFloat().Float64()
	_, err = s.client.Database("Product").Collection("Products").UpdateOne(context.Background(), bson.M{
		"Name": name,
	}, bson.M{
		"$set": bson.M{
			"Price":      priceFloat,
			"LastUpdate": time.Now(),
		},
		"$inc": bson.M{
			"Counter": 1,
		},
	}, options.Update().SetUpsert(true))

	if err != nil {
		return
	}
	return
}

type SuperFloat decimal.Decimal

func (sf SuperFloat) String() string {
	return decimal.Decimal(sf).String()
}

func (sf SuperFloat) GetBSON() (interface{}, error) {
	d := decimal.Decimal(sf)
	f, ok := d.Float64()
	if !ok {
		return nil, errors.New("error parse decimal to float")
	}

	return f, nil
}

func (sf *SuperFloat) SetBSON(raw bson.RawValue) error {
	if d, ok := raw.DoubleOK(); !ok {
		return fmt.Errorf("cant unmarshalling value %v", raw)
	} else {
		*sf = SuperFloat(decimal.NewFromFloat(d))
		return nil
	}
}

func (sf *SuperFloat) Float64() float64 {
	f, _ := decimal.Decimal(*sf).Float64()
	return f
}

type ProductDto struct {
	Name    string     `bson:"Name"`
	Counter int        `bson:"Counter"`
	Price   SuperFloat `bson:"Price"`
}

type sorting int

var (
	nothing sorting = 0
	asc     sorting = -1
	desc    sorting = 1
)

type SortingParams struct {
	Name       sorting
	Price      sorting
	LastUpdate sorting
}

func (s *SortingParams) ToBson() bson.M {
	bsonModel := make(map[string]interface{}, 0)

	if s.Name != nothing {
		bsonModel["Name"] = s.Name
	}

	if s.Price != nothing {
		bsonModel["Price"] = s.Price
	}

	if s.LastUpdate != nothing {
		bsonModel["LastUpdate"] = s.LastUpdate
	}
	return bsonModel
}

type PagingParams struct {
	Page   int64
	Offset int64
	Limit  int64
}

func (p *PagingParams) GetPage() int64 {
	if p.Page <= 0 {
		return 0 // by default
	}
	return p.Page
}

func (p *PagingParams) GetOffset() int64 {
	if p.Offset <= 0 {
		return 0 // by default
	}
	return p.Offset
}

func (p *PagingParams) GetLimit() int64 {
	if p.Limit <= 0 {
		return 10 // by default
	}
	return p.Limit
}

func (s *Store) List(pagingParams PagingParams, sortingParams SortingParams) (products []ProductDto, err error) {
	offset := pagingParams.GetOffset()
	limit := pagingParams.GetLimit()
	opts := options.FindOptions{
		Skip:  &offset,
		Limit: &limit,
	}

	opts.SetSort(sortingParams.ToBson())

	cursor, err := s.client.Database("Product").
		Collection("Products").
		Find(context.Background(),
			bson.M{}, &opts)
	if err != nil {
		return
	}
	err = cursor.All(context.Background(), &products)
	if err != nil {
		return
	}
	return
}
