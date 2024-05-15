package mongoutils

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Find find records
//
// @param ctx operation context
// @param filter (ignored on nil)
// @param sorts (ignored on nil)
// @param skip (ignored on 0)
// @param limit (ignored on 0)
// @opts operation option
func FindCtx[T any](
	ctx context.Context,
	filter any,
	sorts any,
	skip int64,
	limit int64,
	opts ...MongoOption,
) ([]T, error) {
	res := make([]T, 0)
	model := modelSafe(new(T))
	var pipeline MongoPipeline
	opt := optionOf(opts...)
	if v, err := callMethod(model, opt.Pipeline, opt.Params...); err != nil {
		return res, err
	} else {
		pipeline = parsePipeline(v)
	}
	if pipeline == nil {
		return res, errors.New(opt.Pipeline + " method should return MongoPipeline!")
	}
	pipe := pipeline.
		Match(filter).
		Sort(sorts).
		Skip(skip).
		Limit(limit).
		Build()
	if opt.DebugPipe {
		fmt.Println("=============== FIND PIPE ===============")
		prettyLog(pipe)
		fmt.Println("=========================================")
	}

	if opt.DebugResult {
		fmt.Println("============== FIND DECODE ==============")
		if cur, err := model.Collection(opt.Database).Aggregate(ctx, pipe, AggregateOption()); err != nil {
			fmt.Println("ERROR: " + err.Error())
		} else {
			defer cur.Close(ctx)
			var _res []map[string]any
			cur.All(ctx, &_res)
			prettyLog(_res)
		}
		fmt.Println("=========================================")
	}

	if cur, err := model.Collection(opt.Database).Aggregate(ctx, pipe, AggregateOption()); err != nil {
		return res, err
	} else {
		defer cur.Close(ctx)
		if err := cur.All(ctx, &res); err != nil {
			return res, err
		}
	}
	return res, nil
}
func Find[T any](filter any, sorts any, skip int64, limit int64, opts ...MongoOption) ([]T, error) {
	ctx, cancel := MongoOperationCtx()
	defer cancel()
	return FindCtx[T](ctx, filter, sorts, skip, limit, opts...)
}

// FindRaw find records from pipeline
// option pipeline not effected
//
// @param ctx operation context
// @param pipeline aggregation pipeline
// @opts operation option
func FindRawCtx[T any](
	ctx context.Context,
	pipeline MongoPipeline,
	opts ...MongoOption,
) ([]T, error) {
	res := make([]T, 0)
	model := modelSafe(new(T))
	option := optionOf(opts...)
	if option.DebugPipe {
		fmt.Println("============= FIND RAW PIPE =============")
		prettyLog(pipeline.Build())
		fmt.Println("=========================================")
	}

	if option.DebugResult {
		fmt.Println("============ FIND RAW DECODE ============")
		if cur, err := model.Collection(option.Database).Aggregate(ctx, pipeline.Build(), AggregateOption()); err != nil {
			fmt.Println("ERROR: " + err.Error())
		} else {
			defer cur.Close(ctx)
			var _res []map[string]any
			cur.All(ctx, &_res)
			prettyLog(_res)
		}
		fmt.Println("=========================================")
	}

	if cur, err := model.Collection(option.Database).Aggregate(ctx, pipeline.Build(), AggregateOption()); err != nil {
		return res, err
	} else {
		defer cur.Close(ctx)
		if err := cur.All(ctx, &res); err != nil {
			return res, err
		}
	}
	return res, nil
}
func FindRaw[T any](pipeline MongoPipeline, opts ...MongoOption) ([]T, error) {
	ctx, cancel := MongoOperationCtx()
	defer cancel()
	return FindRawCtx[T](ctx, pipeline, opts...)
}

// FindOne find one record
//
// @param ctx operation context
// @param filter (ignored on nil)
// @param sorts  (ignored on nil)
// @opts operation option
func FindOneCtx[T any](
	ctx context.Context,
	filter any,
	sorts any,
	opts ...MongoOption,
) (*T, error) {
	res := new(T)
	model := modelSafe(new(T))
	var pipeline MongoPipeline
	opt := optionOf(opts...)
	if v, err := callMethod(model, opt.Pipeline, opt.Params...); err != nil {
		return res, err
	} else {
		pipeline = parsePipeline(v)
	}
	if pipeline == nil {
		return res, errors.New(opt.Pipeline + " method should return MongoPipeline!")
	}
	pipe := pipeline.
		Match(filter).
		Sort(sorts).
		Limit(1).
		Build()
	if opt.DebugPipe {
		fmt.Println("============= FIND ONE PIPE =============")
		prettyLog(pipe)
		fmt.Println("=========================================")
	}

	if cur, err := model.Collection(opt.Database).Aggregate(ctx, pipe, AggregateOption()); err != nil {
		return res, err
	} else {
		defer cur.Close(ctx)
		for cur.Next(ctx) {
			if opt.DebugResult {
				var _res map[string]any
				cur.Decode(&_res)
				fmt.Println("============ FIND ONE DECODE ============")
				prettyLog(_res)
				fmt.Println("=========================================")
			}
			if err := cur.Decode(res); err != nil {
				return res, err
			} else {
				return res, nil
			}
		}
	}
	return nil, nil
}
func FindOne[T any](filter any, sorts any, opts ...MongoOption) (*T, error) {
	ctx, cancel := MongoOperationCtx()
	defer cancel()
	return FindOneCtx[T](ctx, filter, sorts, opts...)
}

// Insert insert new record
// this function use FindOne to find old record
//
// @param ctx operation context
// @param v model
// @opts operation option
func InsertCtx[T any](
	ctx context.Context,
	v *T,
	opts ...MongoOption,
) (*mongo.InsertOneResult, error) {
	model := modelSafe(v)
	opt := optionOf(opts...)
	model.Cleanup()
	model.FillCreatedAt()
	FillBackupFields(v)
	model.OnInsert(ctx, opts...)
	if res, err := model.Collection(opt.Database).InsertOne(ctx, model); err != nil {
		return res, err
	} else {
		if opt.DebugResult {
			fmt.Println("============= INSERT RESULT =============")
			prettyLog(res)
			fmt.Println("=========================================")
		}
		if id, ok := res.InsertedID.(primitive.ObjectID); !ok {
			return res, errors.New("no ObjectId returned")
		} else {
			model.SetID(id)
			model.OnInserted(ctx, opts...)
			return res, nil
		}
	}
}
func Insert[T any](v *T, opts ...MongoOption) (*mongo.InsertOneResult, error) {
	ctx, cancel := MongoOperationCtx()
	defer cancel()
	return InsertCtx(ctx, v, opts...)
}

// Update update one record
//
// @param ctx operation context
// @param v model
// @param isSilent disable update meta (updated_at)
// @opts operation option
func UpdateCtx[T any](
	ctx context.Context,
	v *T,
	isSilent bool,
	opts ...MongoOption,
) (*mongo.UpdateResult, error) {
	model := modelSafe(v)
	opt := optionOf(opts...)
	old, err := FindOneCtx[T](ctx, primitive.M{"_id": model.GetID()}, nil, opts...)
	if err != nil {
		return nil, err
	}
	// Handle model changes
	model.Cleanup()
	isChanged := true
	oldCS, _ := modelChecksum(old)
	if cs, backup := modelChecksum(v); cs != "" {
		if cs != oldCS {
			backup.SetChecksum(cs)
			backup.UnMarkBackup()
		}
		isChanged = cs != oldCS
	}
	if !isSilent && isChanged {
		model.FillUpdatedAt()
	}
	model.OnUpdate(ctx, opts...)
	if res, err := model.Collection(opt.Database).UpdateByID(ctx, model.GetID(), Set(model)); err != nil {
		return nil, err
	} else {
		if opt.DebugResult {
			fmt.Println("============= UPDATE RESULT =============")
			prettyLog(res)
			fmt.Println("=========================================")
		}
		if res.ModifiedCount+res.UpsertedCount > 0 {
			model.OnUpdated(old, ctx, opts...)
		}
		return res, nil
	}
}
func Update[T any](v *T, silent bool, opts ...MongoOption) (*mongo.UpdateResult, error) {
	ctx, cancel := MongoOperationCtx()
	defer cancel()
	return UpdateCtx(ctx, v, silent, opts...)
}

// Delete delete record
//
// @param ctx operation context
// @param v model
// @opts operation option
func DeleteCtx[T any](
	ctx context.Context,
	v *T,
	opts ...MongoOption,
) (*mongo.DeleteResult, error) {
	model := modelSafe(v)
	opt := optionOf(opts...)
	model.OnDelete(ctx, opts...)
	if res, err := model.Collection(opt.Database).DeleteOne(ctx, primitive.M{"_id": model.GetID()}); err != nil {
		return nil, err
	} else {
		if opt.DebugResult {
			fmt.Println("============= DELETE RESULT =============")
			prettyLog(res)
			fmt.Println("=========================================")
		}
		model.OnDeleted(ctx, opts...)
		return res, nil
	}
}
func Delete[T any](v *T, opts ...MongoOption) (*mongo.DeleteResult, error) {
	ctx, cancel := MongoOperationCtx()
	defer cancel()
	return DeleteCtx(ctx, v, opts...)
}

// Count get records count
//
// @param ctx operation context
// @param filter (ignored on nil)
// @opts operation option
func CountCtx[T any](
	ctx context.Context,
	filter any,
	opts ...MongoOption,
) (int64, error) {
	model := typeModelSafe[T]()
	var pipeline MongoPipeline
	opt := optionOf(opts...)
	if v, err := callMethod(model, opt.Pipeline, opt.Params...); err != nil {
		return 0, err
	} else {
		pipeline = parsePipeline(v)
	}
	if pipeline == nil {
		return 0, errors.New(opt.Pipeline + " method should return MongoPipeline!")
	}
	pipe := pipeline.
		Match(filter).
		Add(func(d MongoDoc) MongoDoc {
			return d.Add("$count", "count")
		}).
		Build()
	if opt.DebugPipe {
		fmt.Println("=============== COUNT PIPE ===============")
		prettyLog(pipe)
		fmt.Println("==========================================")
	}
	if cur, err := model.Collection(opt.Database).Aggregate(ctx, pipe, AggregateOption()); err != nil {
		return 0, err
	} else {
		defer cur.Close(ctx)
		for cur.Next(ctx) {
			if opt.DebugResult {
				var _res map[string]any
				cur.Decode(&_res)
				fmt.Println("=============== COUNT DECODE ===============")
				prettyLog(_res)
				fmt.Println("============================================")
			}
			rec := new(countResult)
			if err := cur.Decode(rec); err != nil {
				return 0, err
			} else {
				return rec.Count, nil
			}
		}
	}
	return 0, nil
}
func Count[T any](filter any, opts ...MongoOption) (int64, error) {
	ctx, cancel := MongoOperationCtx()
	defer cancel()
	return CountCtx[T](ctx, filter, opts...)
}

// CountRaw get records count
// option Pipeline not effected
//
// @param ctx operation context
// @param filter (ignored on nil)
// @opts operation option
func CountRawCtx[T any](
	ctx context.Context,
	pipeline MongoPipeline,
	opts ...MongoOption,
) (int64, error) {
	model := typeModelSafe[T]()
	option := optionOf(opts...)
	pipeline.Add(func(d MongoDoc) MongoDoc { return d.Add("$count", "count") })
	if option.DebugPipe {
		fmt.Println("=============== COUNT PIPE ===============")
		prettyLog(pipeline.Build())
		fmt.Println("==========================================")
	}
	if cur, err := model.Collection(option.Database).Aggregate(ctx, pipeline.Build(), AggregateOption()); err != nil {
		return 0, err
	} else {
		defer cur.Close(ctx)
		for cur.Next(ctx) {
			if option.DebugResult {
				var _res map[string]any
				cur.Decode(&_res)
				fmt.Println("=============== COUNT DECODE ===============")
				prettyLog(_res)
				fmt.Println("============================================")
			}
			rec := new(countResult)
			if err := cur.Decode(rec); err != nil {
				return 0, err
			} else {
				return rec.Count, nil
			}
		}
	}
	return 0, nil
}
func CountRaw[T any](filter any, opts ...MongoOption) (int64, error) {
	ctx, cancel := MongoOperationCtx()
	defer cancel()
	return CountCtx[T](ctx, filter, opts...)
}

// BatchUpdate update multiple records
//
// @param ctx operation context
// @param condition update condition
// @param updates update value
// @opts operation option
func BatchUpdateCtx[T any](
	ctx context.Context,
	condition any,
	updates any,
	opts ...MongoOption,
) (*mongo.UpdateResult, error) {
	model := typeModelSafe[T]()
	opt := optionOf(opts...)
	if res, err := model.Collection(opt.Database).UpdateMany(ctx, condition, updates); err != nil {
		return nil, err
	} else {
		if opt.DebugResult {
			fmt.Println("========== BATCH UPDATE RESULT ==========")
			prettyLog(res)
			fmt.Println("=========================================")
		}
		return res, nil
	}
}
func BatchUpdate[T any](condition any, updates any, opts ...MongoOption) (*mongo.UpdateResult, error) {
	ctx, cancel := MongoOperationCtx()
	defer cancel()
	return BatchUpdateCtx[T](ctx, condition, updates, opts...)
}

// Patch partial update multiple records using $set
//
// @param ctx operation context
// @param condition update condition
// @param data update value
// @param silent disable update meta (updated_at)
// @opts operation option
func PatchCtx[T any](
	ctx context.Context,
	condition any,
	data primitive.M,
	silent bool,
	opts ...MongoOption,
) (*mongo.UpdateResult, error) {
	model := typeModelSafe[T]()
	opt := optionOf(opts...)
	if !silent {
		data["updated_at"] = time.Now().UTC()
	}
	if res, err := model.Collection(opt.Database).UpdateMany(ctx, condition, Set(data)); err != nil {
		return nil, err
	} else {
		if opt.DebugResult {
			fmt.Println("============== PATCH RESULT ==============")
			prettyLog(res)
			fmt.Println("==========================================")
		}
		return res, nil
	}
}
func Patch[T any](condition any, data primitive.M, silent bool, opts ...MongoOption) (*mongo.UpdateResult, error) {
	ctx, cancel := MongoOperationCtx()
	defer cancel()
	return PatchCtx[T](ctx, condition, data, silent, opts...)
}

// Increment increment numeric data
// pass negative value for decrement
// increment run on silent mode
//
// @param ctx operation context
// @param condition update condition
// @param data update value
// @opts operation option
func IncrementCtx[T any](
	ctx context.Context,
	condition any,
	data any,
	opts ...MongoOption,
) (*mongo.UpdateResult, error) {
	model := typeModelSafe[T]()
	opt := optionOf(opts...)
	if res, err := model.Collection(opt.Database).UpdateMany(ctx, condition, primitive.M{"$inc": data}); err != nil {
		return nil, err
	} else {
		if opt.DebugResult {
			fmt.Println("============ INCREMENT RESULT ============")
			prettyLog(res)
			fmt.Println("==========================================")
		}
		return res, nil
	}
}
func Increment[T any](condition any, data any, opts ...MongoOption) (*mongo.UpdateResult, error) {
	ctx, cancel := MongoOperationCtx()
	defer cancel()
	return IncrementCtx[T](ctx, condition, data, opts...)
}
