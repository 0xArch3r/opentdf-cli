package client

import (
	"context"

	"github.com/opentdf/platform/protocol/go/policy"
	"github.com/opentdf/platform/protocol/go/policy/subjectmapping"
)

func (s *Sdk) ListSubjectMappings() (*subjectmapping.ListSubjectMappingsResponse, error) {
	request := &subjectmapping.ListSubjectMappingsRequest{}
	return s.sdk.SubjectMapping.ListSubjectMappings(context.Background(), request)
}

func (s *Sdk) DeleteSubjectMapping(mapping_id string) (*subjectmapping.DeleteSubjectMappingResponse, error) {
	request := &subjectmapping.DeleteSubjectMappingRequest{
		Id: mapping_id,
	}

	return s.sdk.SubjectMapping.DeleteSubjectMapping(context.Background(), request)
}

func (s *Sdk) CreateSubjectConditionSet(input *subjectmapping.SubjectConditionSetCreate) (*subjectmapping.CreateSubjectConditionSetResponse, error) {
	request := &subjectmapping.CreateSubjectConditionSetRequest{
		SubjectConditionSet: input,
	}
	return s.sdk.SubjectMapping.CreateSubjectConditionSet(context.Background(), request)
}

func (s *Sdk) CreateSubjectMapping(attribute_value_id string, condition_set *policy.SubjectConditionSet, actions []*policy.Action) (*subjectmapping.CreateSubjectMappingResponse, error) {
	request := &subjectmapping.CreateSubjectMappingRequest{
		AttributeValueId:              attribute_value_id,
		ExistingSubjectConditionSetId: condition_set.Id,
		Actions:                       actions,
	}
	return s.sdk.SubjectMapping.CreateSubjectMapping(context.Background(), request)
}

func buildActionsFromStrings(actionStrings []string) ([]*policy.Action, error) {
	actions := make([]*policy.Action, len(actionStrings))
	for i, actionStr := range actionStrings {
		switch actionStr {
		case "DECRYPT":
			actions[i] = &policy.Action{Value: &policy.Action_Standard{Standard: policy.Action_STANDARD_ACTION_DECRYPT}}
		case "TRANSMIT":
			actions[i] = &policy.Action{Value: &policy.Action_Standard{Standard: policy.Action_STANDARD_ACTION_TRANSMIT}}
		default:
			actions[i] = &policy.Action{Value: &policy.Action_Custom{Custom: actionStr}}
		}
	}
	return actions, nil
}

func buildConditionGroups(conditionGroups []ConditionGroup) []*policy.ConditionGroup {
	groups := make([]*policy.ConditionGroup, len(conditionGroups))
	for i, cg := range conditionGroups {
		groups[i] = &policy.ConditionGroup{
			BooleanOperator: getBooleanOperatorEnum(cg.BooleanOperator),
			Conditions:      buildConditions(cg.Conditions),
		}
	}
	return groups
}

func buildConditions(conditions []Condition) []*policy.Condition {
	builtConditions := make([]*policy.Condition, len(conditions))
	for i, cond := range conditions {
		builtConditions[i] = &policy.Condition{
			SubjectExternalSelectorValue: cond.SubjectExternalSelectorValue,
			Operator:                     getOperatorEnum(cond.Operator),
			SubjectExternalValues:        cond.SubjectExternalValues,
		}
	}
	return builtConditions
}

func getBooleanOperatorEnum(operator string) policy.ConditionBooleanTypeEnum {
	switch operator {
	case "AND":
		return policy.ConditionBooleanTypeEnum_CONDITION_BOOLEAN_TYPE_ENUM_AND
	case "OR":
		return policy.ConditionBooleanTypeEnum_CONDITION_BOOLEAN_TYPE_ENUM_OR
	default:
		return policy.ConditionBooleanTypeEnum_CONDITION_BOOLEAN_TYPE_ENUM_UNSPECIFIED
	}
}

func getOperatorEnum(operator string) policy.SubjectMappingOperatorEnum {
	switch operator {
	case "IN":
		return policy.SubjectMappingOperatorEnum_SUBJECT_MAPPING_OPERATOR_ENUM_IN
	case "NOT_IN":
		return policy.SubjectMappingOperatorEnum_SUBJECT_MAPPING_OPERATOR_ENUM_NOT_IN
	default:
		return policy.SubjectMappingOperatorEnum_SUBJECT_MAPPING_OPERATOR_ENUM_UNSPECIFIED
	}
}

// Helper function to get subjects from a subject mapping
func getSubjectsFromMapping(sm *policy.SubjectMapping) []SubjectConditionSet {
	subjectSets := []SubjectConditionSet{}
	if len(sm.SubjectConditionSet.SubjectSets) == 0 {
		subjectSets = append(subjectSets, SubjectConditionSet{})
	} else {
		for _, subjectSet := range sm.SubjectConditionSet.SubjectSets {
			subjectSets = append(subjectSets, SubjectConditionSet{
				ConditionGroups: getConditionGroups(subjectSet),
			})
		}
	}
	return subjectSets
}

func getConditionGroups(subjectSet *policy.SubjectSet) []ConditionGroup {
	groups := []ConditionGroup{}
	for _, conditionGroup := range subjectSet.ConditionGroups {
		conditions := []Condition{}
		for _, condition := range conditionGroup.Conditions {
			conditions = append(conditions, Condition{
				SubjectExternalSelectorValue: condition.SubjectExternalSelectorValue,
				Operator:                     condition.Operator.String(),
				SubjectExternalValues:        condition.SubjectExternalValues,
			})
		}

		groups = append(groups, ConditionGroup{
			BooleanOperator: conditionGroup.BooleanOperator.String(),
			Conditions:      conditions,
		})
	}
	return groups
}
