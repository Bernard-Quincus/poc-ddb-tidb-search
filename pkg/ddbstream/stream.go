package ddbstream

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// ConvertStreamRecord converts a ddb stream record to a type
func ConvertStreamRecord[T any](record *events.DynamoDBStreamRecord, avType string, out *T) error {
	sMap, err := streamAttributesToMap(record.NewImage)
	if err != nil {
		return err
	}

	err = attributevalue.UnmarshalMapWithOptions(sMap, out, func(o *attributevalue.DecoderOptions) {
		o.TagKey = avType
	})
	if err != nil {
		return err
	}
	return nil
}

func streamAttributesToMap(attributes map[string]events.DynamoDBAttributeValue) (map[string]types.AttributeValue, error) {
	dbAttrMap := make(map[string]types.AttributeValue)

	for k, v := range attributes {
		nv, err := convertAttribute(v)
		if err != nil {
			return nil, err
		}
		dbAttrMap[k] = nv
	}
	return dbAttrMap, nil
}

func convertAttribute(av events.DynamoDBAttributeValue) (types.AttributeValue, error) {
	switch av.DataType() {
	default:
		return nil, fmt.Errorf("unhandled attribute type %d", av.DataType())
	case events.DataTypeBinary:
		return &types.AttributeValueMemberB{Value: av.Binary()}, nil
	case events.DataTypeString:
		return &types.AttributeValueMemberS{Value: av.String()}, nil
	case events.DataTypeNumber:
		return &types.AttributeValueMemberN{Value: av.Number()}, nil
	case events.DataTypeBoolean:
		return &types.AttributeValueMemberBOOL{Value: av.Boolean()}, nil
	case events.DataTypeNull:
		return &types.AttributeValueMemberNULL{Value: av.IsNull()}, nil
	case events.DataTypeList:
		list := make([]types.AttributeValue, 0)
		for _, lv := range av.List() {
			o, err := convertAttribute(lv)
			if err != nil {
				return nil, err
			}

			list = append(list, o)
		}
		return &types.AttributeValueMemberL{Value: list}, nil
	case events.DataTypeMap:
		m := make(map[string]types.AttributeValue)
		for k, v := range av.Map() {
			mv, err := convertAttribute(v)
			if err != nil {
				return nil, err
			}
			m[k] = mv
		}
		return &types.AttributeValueMemberM{Value: m}, nil
	}
}
