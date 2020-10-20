package main

import (
	"context"
	"errors"
	"log"

	p "github.com/Amirkin/product-service/proto"
)

type Handlers struct {
	Store     *Store
	Requester *Requester
}

func (h *Handlers) Fetch(_ context.Context, r *p.FetchRequest) (*p.Empty, error) {
	body, err := h.Requester.GetCSV(r.Url)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	products, err := ParseCSV(body)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	for _, product := range products {
		err = h.Store.Save(product.Name, product.Price)
		if err != nil {
			log.Println(err)
			continue
		}
	}

	return nil, nil
}

func (h *Handlers) List(_ context.Context, r *p.ListParams) (*p.ListResponse, error) {
	paginParams := PagingParams{
		Page:   int64(r.PagingParams.Page),
		Offset: int64(r.PagingParams.Offset),
		Limit:  int64(r.PagingParams.Limit),
	}
	sortParams := SortingParams{
		Name:       sorting(r.SortParams.Name),
		Price:      sorting(r.SortParams.Price),
		LastUpdate: sorting(r.SortParams.LastUpdate),
	}
	products, err := h.Store.List(paginParams, sortParams)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &p.ListResponse{
		ListProduct: productDtoToProductClient(products),
	}, errors.New("list")
}

func productDtoToProductClient(dto []ProductDto) []*p.Product {
	client := make([]*p.Product, len(dto))

	for i, model := range dto {
		client[i] = &p.Product{
			Name:    model.Name,
			Counter: int32(model.Counter),
			Price:   model.Price.Float64(),
		}
	}

	return client
}
