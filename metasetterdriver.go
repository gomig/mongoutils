package mongoutils

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type metaV struct {
	ID    primitive.ObjectID
	Meta  string
	Value any
}

type metaSetter struct {
	Data map[string][]metaV
}

func (ms *metaSetter) addCol(col string) {
	if _, ok := ms.Data[col]; !ok {
		ms.Data[col] = make([]metaV, 0)
	}
}

func (ms *metaSetter) Add(_col, _meta string, id *primitive.ObjectID, value any) MetaSetter {
	if id != nil {
		ms.addCol(_col)
		for i, mt := range ms.Data[_col] {
			if mt.ID == *id && mt.Meta == _meta {
				ms.Data[_col][i].Value = value
				return ms
			}
		}
		ms.Data[_col] = append(ms.Data[_col], metaV{Meta: _meta, ID: *id, Value: value})
	}
	return ms
}

func (ms *metaSetter) Build() []MetaSetterResult {
	result := make([]MetaSetterResult, 0)
	ignores := make(map[string]map[string]any)
	addIgnore := func(_col, _meta string, value any) {
		if _, ok := ignores[_col]; !ok {
			ignores[_col] = make(map[string]any)
		}
		ignores[_col][_meta] = value
	}
	isAdded := func(_col, _meta string, value any) bool {
		for k, i := range ignores {
			if k == _col {
				for _k, v := range i {
					if _k == _meta && v == value {
						return true
					}
				}
			}
		}
		return false
	}
	foundIds := func(_meta string, value any, data []metaV) []primitive.ObjectID {
		ids := make([]primitive.ObjectID, 0)
		for _, m := range data {
			if m.Meta == _meta && value == m.Value {
				ids = append(ids, m.ID)
			}
		}
		return ids
	}
	for _col, _meta := range ms.Data {
		for _, m := range _meta {
			if !isAdded(_col, m.Meta, m.Value) {
				result = append(result, MetaSetterResult{
					Col:    _col,
					Ids:    foundIds(m.Meta, m.Value, _meta),
					Values: map[string]any{m.Meta: m.Value},
				})
				addIgnore(_col, m.Meta, m.Value)
			}
		}
	}
	return result
}
