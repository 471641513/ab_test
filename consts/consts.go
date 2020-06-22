package consts

type ExperimentGroup string

func (group ExperimentGroup) IsControlGroup() bool {
	return group == CONTROL_GROUP
}

func (group ExperimentGroup) IsTreatmentGroup() bool {
	return group == TREATMENT_GROUP
}

const (
	CONTROL_GROUP   = "control_group"
	TREATMENT_GROUP = "treatment_group"
)
