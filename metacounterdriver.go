package mongoutils

import "go.mongodb.org/mongo-driver/bson/primitive"

type meta struct {
	ID     primitive.ObjectID
	Meta   string
	Amount int64
}

type metaCounter struct {
	Data map[string][]meta
}

func (mc *metaCounter) addCol(col string) {
	if _, ok := mc.Data[col]; !ok {
		mc.Data[col] = make([]meta, 0)
	}
}

func (mc *metaCounter) Add(_col string, _meta string, id *primitive.ObjectID, amount int64) MetaCounter {
	if id != nil {
		mc.addCol(_col)
		for i, mt := range mc.Data[_col] {
			if mt.ID == *id && mt.Meta == _meta {
				mc.Data[_col][i].Amount += amount
				return mc
			}
		}
		mc.Data[_col] = append(mc.Data[_col], meta{Meta: _meta, ID: *id, Amount: amount})
	}
	return mc
}

func (mc *metaCounter) Build() []MetaCounterResult {
	result := make([]MetaCounterResult, 0)
	ignores := make(map[string]map[string]int64)
	addIgnore := func(_col, _meta string, amount int64) {
		if _, ok := ignores[_col]; !ok {
			ignores[_col] = make(map[string]int64)
		}
		ignores[_col][_meta] = amount
	}
	isAdded := func(_col, _meta string, amount int64) bool {
		for k, i := range ignores {
			if k == _col {
				for _k, _a := range i {
					if _k == _meta && _a == amount {
						return true
					}
				}
			}
		}
		return false
	}
	foundIds := func(_meta string, amount int64, data []meta) []primitive.ObjectID {
		ids := make([]primitive.ObjectID, 0)
		for _, m := range data {
			if m.Meta == _meta && amount == m.Amount {
				ids = append(ids, m.ID)
			}
		}
		return ids
	}
	for _col, _meta := range mc.Data {
		for _, m := range _meta {
			if m.Amount != 0 && !isAdded(_col, m.Meta, m.Amount) {
				result = append(result, MetaCounterResult{
					Col:    _col,
					Ids:    foundIds(m.Meta, m.Amount, _meta),
					Values: map[string]int64{m.Meta: m.Amount},
				})
				addIgnore(_col, m.Meta, m.Amount)
			}
		}
	}
	return result
}
