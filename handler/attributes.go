package handler

import (
	"context"
	"fmt"

	"github.com/opentdf/platform/protocol/go/policy"
	"github.com/opentdf/platform/protocol/go/policy/attributes"
	"github.com/opentdf/platform/sdk"
)

func (s *Sdk) ListAttributes(namespace string) (*attributes.ListAttributesResponse, error) {
	request := &attributes.ListAttributesRequest{
		Namespace: namespace,
	}
	return s.sdk.Attributes.ListAttributes(context.Background(), request)
}

func (s *Sdk) ListAttributeValues(attribute_id string) (*attributes.ListAttributeValuesResponse, error) {
	request := &attributes.ListAttributeValuesRequest{
		AttributeId: attribute_id,
	}

	return s.sdk.Attributes.ListAttributeValues(context.Background(), request)
}

func (s *Sdk) CreateAttribute(namespace, name, rule string) (*attributes.CreateAttributeResponse, error) {
	ruleEnum, err := getAttributeRuleEnum(rule)
	if err != nil {
		return nil, err
	}
	request := &attributes.CreateAttributeRequest{
		Name:        name,
		NamespaceId: namespace,
		Rule:        ruleEnum,
	}

	return s.sdk.Attributes.CreateAttribute(context.Background(), request)
}

func (s *Sdk) CreateAttributeValue(attribute_id, value string) (*attributes.CreateAttributeValueResponse, error) {
	request := &attributes.CreateAttributeValueRequest{
		AttributeId: attribute_id,
		Value:       value,
	}
	return s.sdk.Attributes.CreateAttributeValue(context.Background(), request)
}

func (s *Sdk) DeactivateAttribute(attribute_id string) (*attributes.DeactivateAttributeResponse, error) {
	request := &attributes.DeactivateAttributeRequest{
		Id: attribute_id,
	}

	return s.sdk.Attributes.DeactivateAttribute(context.Background(), request)
}

func (s *Sdk) DeactivateAttributeValue(attribute_id string) (*attributes.DeactivateAttributeValueResponse, error) {
	request := &attributes.DeactivateAttributeValueRequest{
		Id: attribute_id,
	}
	return s.sdk.Attributes.DeactivateAttributeValue(context.Background(), request)
}

func getAttributeRuleEnum(ruleStr string) (policy.AttributeRuleTypeEnum, error) {
	switch ruleStr {
	case "ANY_OF":
		return policy.AttributeRuleTypeEnum_ATTRIBUTE_RULE_TYPE_ENUM_ANY_OF, nil
	case "ALL_OF":
		return policy.AttributeRuleTypeEnum_ATTRIBUTE_RULE_TYPE_ENUM_ALL_OF, nil
	case "HIERARCHY":
		return policy.AttributeRuleTypeEnum_ATTRIBUTE_RULE_TYPE_ENUM_HIERARCHY, nil
	default:
		return policy.AttributeRuleTypeEnum_ATTRIBUTE_RULE_TYPE_ENUM_UNSPECIFIED, fmt.Errorf("invalid attribute rule: %s", ruleStr)
	}
}

func syncAttributeValues(client *sdk.SDK, attribute *policy.Attribute, existingValues []*policy.Value, desiredValues []string) error {
	ctx := context.Background()

	existingValueMap := make(map[string]*policy.Value)
	for _, v := range existingValues {
		existingValueMap[v.Value] = v
	}

	// Create missing values
	for _, desiredVal := range desiredValues {
		if _, ok := existingValueMap[desiredVal]; !ok {
			_, err := client.Attributes.CreateAttributeValue(ctx, &attributes.CreateAttributeValueRequest{
				AttributeId: attribute.Id,
				Value:       desiredVal,
			})
			if err != nil {
				return fmt.Errorf("could not create attribute value: %w", err)
			}
		}
	}

	// Delete extra values
	for _, existingVal := range existingValues {
		var found bool
		for _, desiredVal := range desiredValues {
			if existingVal.Value == desiredVal {
				found = true
				break
			}
		}
		if !found {
			_, err := client.Attributes.DeactivateAttributeValue(ctx, &attributes.DeactivateAttributeValueRequest{
				Id: existingVal.Id,
			})
			if err != nil {
				return fmt.Errorf("could not delete attribute value: %w", err)
			}
		}
	}

	return nil
}
