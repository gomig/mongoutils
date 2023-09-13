package mongoutils

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Array generate primitive.A
func Array(args ...any) primitive.A {
	return args
}

// Map generate primitive.M
//
// Args count must even
func Map(args ...any) primitive.M {
	if len(args)%2 == 0 {
		res := make(primitive.M, len(args)/2)
		for i := 0; i < len(args); i++ {
			if i%2 == 0 {
				if k, ok := args[i].(string); ok {
					res[k] = args[i+1]
				}
			}
		}
		return res
	}
	return primitive.M{}
}

// Maps generate []primitive.M
//
// Args count must even
func Maps(args ...any) []primitive.M {
	res := make([]primitive.M, 0)
	if len(args)%2 == 0 {
		for i := 0; i < len(args); i++ {
			if i%2 == 0 {
				if k, ok := args[i].(string); ok {
					res = append(res, primitive.M{k: args[i+1]})
				}
			}
		}
	}
	return res
}

// Doc generate primitive.D from args
//
// Args count must even
// Example: Doc("_id", 1, "name", "John")
func Doc(args ...any) primitive.D {
	res := make([]primitive.E, 0)
	if len(args)%2 == 0 {
		for i := 0; i < len(args); i++ {
			if i%2 == 0 {
				if k, ok := args[i].(string); ok {
					res = append(res, primitive.E{Key: k, Value: args[i+1]})
				}
			}
		}
	}
	return res
}

// Regex generate Regex
//
// { pattern: "John.*", options: "i" }
func Regex(pattern string, opt string) primitive.Regex {
	return primitive.Regex{Pattern: pattern, Options: opt}
}

// RegexFor generate map with regex parameter
//
// { "name": { pattern: "John.*", options: "i" } }
func RegexFor(k string, pattern string, opt string) primitive.M {
	return primitive.M{k: Regex(pattern, opt)}
}

// In generate $in map
//
// {k: {$in: v}}
func In(k string, v ...any) primitive.M {
	return primitive.M{k: primitive.M{"$in": v}}
}

// Set generate simple set map
//
// {$set: v}
func Set(v any) primitive.M {
	return primitive.M{"$set": v}
}

// SetNested generate nested set map
//
// {$set: {k: v}}
func SetNested(k string, v any) primitive.M {
	return primitive.M{"$set": primitive.M{k: v}}
}

// Match generate nested set map
//
// {$match: v}
func Match(v any) primitive.M {
	return primitive.M{"$match": v}
}
