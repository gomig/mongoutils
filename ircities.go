package mongoutils

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// IrCity only implement model methods
type IrCity struct {
	EmptyModel `bson:"inline"`
	Code       uint                `bson:"code" json:"code"`
	Name       string              `bson:"name" json:"name"`
	CountyId   *primitive.ObjectID `bson:"county_id" json:"county_id"`
	ProvinceId *primitive.ObjectID `bson:"province_id" json:"province_id"`
	County     *IrCity             `bson:"county,omitempty" json:"county"`
	Province   *IrCity             `bson:"province,omitempty" json:"province"`
}

func (*IrCity) TypeName() string {
	return "IrCity"
}

func (*IrCity) Collection(db *mongo.Database) *mongo.Collection {
	return db.Collection("ir_cities")
}

func (ic *IrCity) NewId() {
	ic.ID = primitive.NewObjectID()
}

func (ic *IrCity) SetID(id primitive.ObjectID) {
	ic.ID = id
}

func (ic *IrCity) GetID() primitive.ObjectID {
	return ic.ID
}

func (*IrCity) IsEditable() bool {
	return false
}

func (ic *IrCity) Cleanup() {
	ic.County = nil
	ic.Province = nil
}

func (*IrCity) SinglePipeline() MongoPipeline {
	return NewPipe().
		LoadRelation("ir_cities", "county_id", "_id", "county").
		LoadRelation("ir_cities", "province_id", "_id", "province")
}

func (ic *IrCity) Index(db *mongo.Database) error {
	_, err := ic.Collection(db).Indexes().CreateMany(context.TODO(), []mongo.IndexModel{
		{Keys: primitive.M{"name": 1}},
		{Keys: primitive.M{"county_id": 1}},
		{Keys: primitive.M{"province_id": 1}},
	})
	return err
}

func (i *IrCity) Seed(db *mongo.Database) error {
	type city struct {
		ID       int    `json:"id"`
		Code     int    `json:"code"`
		Name     string `json:"name"`
		Type     int    `json:"type"`
		County   int    `json:"county"`
		Province int    `json:"province"`
	}

	var err error
	var provinceData []byte
	var countyData []byte
	var cityData []byte

	// fetch data
	if provinceData, err = fetchJson("https://raw.githubusercontent.com/mekramy/data/master/province.json"); err != nil {
		return err
	}
	if countyData, err = fetchJson("https://raw.githubusercontent.com/mekramy/data/master/county.json"); err != nil {
		return err
	}
	if cityData, err = fetchJson("https://raw.githubusercontent.com/mekramy/data/master/city.json"); err != nil {
		return err
	}

	// unmarshal data
	provinces := make([]city, 0)
	counties := make([]city, 0)
	cities := make([]city, 0)

	// unmarshal data
	if err := json.Unmarshal(provinceData, &provinces); err != nil {
		return err
	}
	if err := json.Unmarshal(countyData, &counties); err != nil {
		return err
	}
	if err := json.Unmarshal(cityData, &cities); err != nil {
		return err
	}

	// parse data
	records := make([]IrCity, 0)
	for _, _province := range provinces {
		province := IrCity{}
		province.ID = *uniqueID(3, _province.Code)
		province.Code = uint(_province.Code)
		province.Name = _province.Name
		records = append(records, province)
		for _, _county := range counties {
			if _county.Province != _province.ID {
				continue
			}
			county := IrCity{}
			county.ID = *uniqueID(7, _county.Code)
			county.Code = uint(_county.Code)
			county.Name = _county.Name
			county.ProvinceId = &province.ID
			records = append(records, county)
			for _, _city := range cities {
				if _city.Type != 0 || _city.County != _county.ID {
					continue
				}
				city := IrCity{}
				city.ID = *uniqueID(11, _city.Code)
				city.Code = uint(_city.Code)
				city.Name = _city.Name
				city.CountyId = &county.ID
				city.ProvinceId = &province.ID
				records = append(records, city)
			}
		}
	}

	for _, rec := range records {
		if count, err := i.Collection(db).CountDocuments(context.TODO(), primitive.M{"code": rec.Code}); err != nil {
			return err
		} else if count == 0 {
			if _, err = i.Collection(db).InsertOne(context.TODO(), rec); err != nil {
				return err
			}
		}
	}
	return nil
}

// utils
func uniqueID(n, code int) *primitive.ObjectID {
	pattern := fmt.Sprintf("%%0%dd", n)
	suffix := fmt.Sprintf(pattern, code)
	id := "60486e0000040cc14406ca54"
	id = id[:len(id)-n]
	id = id + suffix
	return ParseObjectID(id)
}

func fetchJson(url string) ([]byte, error) {
	if resp, err := http.Get(url); err != nil {
		return nil, err
	} else {
		defer resp.Body.Close()
		return io.ReadAll(resp.Body)
	}
}
