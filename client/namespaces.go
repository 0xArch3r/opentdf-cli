package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/opentdf/platform/protocol/go/common"
	"github.com/opentdf/platform/protocol/go/policy"
	"github.com/opentdf/platform/protocol/go/policy/namespaces"
	"github.com/opentdf/platform/protocol/go/policy/subjectmapping"
)

// NamespaceState represents the desired state of a namespace.
type NamespaceState struct {
	Name       string           `json:"name"`
	Attributes []AttributeState `json:"attributes"`
}

// AttributeState represents the desired state of an attribute.
type AttributeState struct {
	Name     string                `json:"name"`
	Rule     string                `json:"rule"`
	Values   []string              `json:"values"`
	Mappings []SubjectMappingState `json:"mappings"`
}

// SubjectMappingState represents the desired state of a subject mapping.
type SubjectMappingState struct {
	AttributeValue string                `json:"attributeValue"`
	Actions        []string              `json:"actions"`
	Subjects       []SubjectConditionSet `json:"subjects"`
}

// SubjectConditionSet represents the desired state of a subject condition set.
type SubjectConditionSet struct {
	ConditionGroups []ConditionGroup `json:"conditionGroups"`
}

// ConditionGroup represents the desired state of a condition group.
type ConditionGroup struct {
	BooleanOperator string      `json:"booleanOperator"`
	Conditions      []Condition `json:"conditions"`
}

// Condition represents the desired state of a condition.
type Condition struct {
	SubjectExternalSelectorValue string   `json:"subjectExternalSelectorValue"`
	Operator                     string   `json:"operator"`
	SubjectExternalValues        []string `json:"subjectExternalValues"`
}

func (s *Sdk) CreateNamespace(name string) (*namespaces.CreateNamespaceResponse, error) {
	request := &namespaces.CreateNamespaceRequest{
		Name: name,
	}
	return s.sdk.Namespaces.CreateNamespace(context.Background(), request)
}

func (s *Sdk) ListNamespaces() (*namespaces.ListNamespacesResponse, error) {
	request := &namespaces.ListNamespacesRequest{
		State: common.ActiveStateEnum_ACTIVE_STATE_ENUM_ANY,
	}
	return s.sdk.Namespaces.ListNamespaces(context.Background(), request)
}

func (s *Sdk) ApplyNamespace(f io.Reader) error {

	data, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	var desiredState NamespaceState
	err = json.Unmarshal(data, &desiredState)
	if err != nil {
		return err
	}

	existingNamespaces, err := s.ListNamespaces()
	if err != nil {
		return err
	}

	var ns *policy.Namespace
	for _, n := range existingNamespaces.Namespaces {
		if n.Name == desiredState.Name {
			ns = n
			break
		}
	}

	if ns == nil {
		x, err := s.CreateNamespace(desiredState.Name)
		if err != nil {
			return err
		}
		ns = x.GetNamespace()
	}

	existingAttributes, err := s.ListAttributes(ns.Id)
	if err != nil {
		return err
	}

	for _, da := range desiredState.Attributes {
		var ex *policy.Attribute
		for _, a := range existingAttributes.GetAttributes() {
			if a.Name == da.Name {
				ex = a
				break
			}
		}
		if ex == nil {
			attr, err := s.CreateAttribute(ns.Id, da.Name, da.Rule)
			if err != nil {
				return err
			}
			ex = attr.Attribute
		}
		attr_values, err := s.ListAttributeValues(ex.Id)
		if err != nil {
			return err
		}

		valueMap := make(map[string]*policy.Value)
		for _, v := range attr_values.GetValues() {
			valueMap[v.Value] = v
		}

		for _, dv := range da.Values {
			if _, ok := valueMap[dv]; !ok {
				_, err := s.CreateAttributeValue(ex.Id, dv)
				if err != nil {
					return err
				}
			}
		}

		for _, ev := range attr_values.GetValues() {
			var found bool
			for _, dv := range da.Values {
				if ev.Value == dv {
					found = true
					break
				}
			}
			if !found {
				_, err := s.DeactivateAttributeValue(ev.Id)
				if err != nil {
					return err
				}
			}
		}

		for _, desired_mapping := range da.Mappings {
			attr_values, err = s.ListAttributeValues(ex.Id)
			if err != nil {
				return fmt.Errorf("unable to list attribute values: %w", err)
			}
			var pol_val *policy.Value
			for _, val := range attr_values.GetValues() {
				if val.Value == desired_mapping.AttributeValue {
					pol_val = val
				}
			}
			if pol_val == nil {
				return errors.New("could not get attribute value tied to mapping")
			}

			actions, err := buildActionsFromStrings(desired_mapping.Actions)
			if err != nil {
				return err
			}

			var subjectConditionSets []*policy.SubjectConditionSet
			for _, subject := range desired_mapping.Subjects {
				var ConditionSet *subjectmapping.SubjectConditionSetCreate
				if subject.ConditionGroups == nil {
					ConditionSet = &subjectmapping.SubjectConditionSetCreate{
						SubjectSets: []*policy.SubjectSet{
							{
								ConditionGroups: []*policy.ConditionGroup{
									{
										BooleanOperator: policy.ConditionBooleanTypeEnum_CONDITION_BOOLEAN_TYPE_ENUM_AND,
										Conditions:      []*policy.Condition{}, // Allow empty conditions for ALL_USERS
									},
								},
							},
						},
					}
				} else {
					ConditionSet = &subjectmapping.SubjectConditionSetCreate{
						SubjectSets: []*policy.SubjectSet{
							{
								ConditionGroups: buildConditionGroups(subject.ConditionGroups),
							},
						},
					}
				}
				resp, err := s.CreateSubjectConditionSet(ConditionSet)
				if err != nil {
					return err
				}
				subjectConditionSets = append(subjectConditionSets, resp.SubjectConditionSet)
			}

			for _, subjectConditionSet := range subjectConditionSets {
				_, err := s.CreateSubjectMapping(pol_val.Id, subjectConditionSet, actions)
				if err != nil {
					return err
				}
			}
		}

		// Delete Subject Mappings
		s_mappings, err := s.ListSubjectMappings()
		if err != nil {
			return err
		}

		for _, s_mapping := range s_mappings.GetSubjectMappings() {
			if s_mapping.AttributeValue.Attribute == nil {
				continue
			}
			if s_mapping.AttributeValue.Id != ex.Id {
				continue
			}
			var found bool
			for _, mapping := range da.Mappings {
				if s_mapping.AttributeValue.Value == mapping.AttributeValue {
					found = true
					break
				}
			}

			if !found {
				_, err := s.DeleteSubjectMapping(s_mapping.Id)
				if err != nil {
					return err
				}
			}
		}
	}

	for _, existing_attr := range existingAttributes.GetAttributes() {
		var found bool
		for _, desired_attr := range desiredState.Attributes {
			if existing_attr.Name == desired_attr.Name {
				found = true
				break
			}
		}
		if !found {
			_, err := s.DeactivateAttribute(existing_attr.Id)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *Sdk) ExportNamespace(name string) (*NamespaceState, error) {

	out := &NamespaceState{}

	namespaces, err := s.ListNamespaces()
	if err != nil {
		return nil, err
	}
	var ns *policy.Namespace
	for _, n := range namespaces.GetNamespaces() {
		if n.Name == name {
			ns = n
		}
	}
	if ns == nil {
		return nil, fmt.Errorf("unable to find namespace %q", name)
	}
	out.Name = ns.Name

	attributes, err := s.ListAttributes(ns.Id)
	if err != nil {
		return nil, err
	}

	subject_mappings_list, err := s.ListSubjectMappings()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve subject mappings: %w", err)
	}

	attributeStates := make([]AttributeState, 0)
	for _, attr := range attributes.GetAttributes() {
		attribute_values, err := s.ListAttributeValues(attr.Id)
		if err != nil {
			return nil, err
		}
		valuesList := make([]string, len(attribute_values.Values))
		mappings := make([]SubjectMappingState, 0)
		for i, val := range attribute_values.GetValues() {
			valuesList[i] = val.Value
			subjects := make([]SubjectConditionSet, 0)
			actions := make(map[string]bool, 0)
			for _, sm := range subject_mappings_list.GetSubjectMappings() {
				if sm.AttributeValue.Value == val.Value {
					for _, action := range sm.Actions {
						switch a := action.GetValue().(type) {
						case *policy.Action_Standard:
							actions[a.Standard.String()] = true
						case *policy.Action_Custom:
							actions[a.Custom] = true
						}
					}
					subjects = append(subjects, getSubjectsFromMapping(sm)...)
				}
			}
			action_slice := make([]string, 0)
			for action, _ := range actions {
				action_slice = append(action_slice, action)
			}
			mappings = append(mappings, SubjectMappingState{
				AttributeValue: val.Value,
				Actions:        action_slice,
				Subjects:       subjects,
			})

		}
		attributeStates = append(attributeStates, AttributeState{
			Name:     attr.Name,
			Rule:     attr.Rule.String(),
			Values:   valuesList, // This is unsorted
			Mappings: mappings,
		})
	}

	out.Attributes = attributeStates

	return out, nil
}
