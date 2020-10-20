package main

import (
	"log"
	"testing"
)

func TestMongoSave(t *testing.T) {
	log.SetFlags(log.Flags() | log.Llongfile)
	store, err := NewStore()
	if err != nil {
		t.Fatal(err)
		return
	}
	//err = store.Save("kek", decimal.NewFromFloat32(11.22))
	//if err != nil {
	//	t.Fatal(err)
	//	return
	//}
	log.Println(store.List(
		PagingParams{
			Page:   0,
			Offset: 0,
			Limit:  0,
		},
		SortingParams{
			Name:       0,
			Price:      1,
			LastUpdate: 0,
		},
	))
}
