package models

import wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"

func (gs *GameState) HasDiscovery(d ResearchDiscovery) bool {
	for _, x := range gs.Discoveries {
		if x == d {
			return true
		}
	}

	return false
}

func (gs *GameStatePatch) HasDiscovery(d ResearchDiscovery) bool {
	for _, x := range gs.Discoveries {
		if x == d {
			return true
		}
	}

	return false
}

func (m *Mouse) ToPatch(isNew bool) *GameStatePatch_MousePatch {
	return &GameStatePatch_MousePatch{
		IsNew:           isNew,
		Name:            wrapperspb.String(m.Name),
		Education:       &GameStatePatch_EducationPatch{Value: m.Education},
		IsEducated:      wrapperspb.Bool(m.IsEducated),
		IsBeingEducated: wrapperspb.Bool(m.IsBeingEducated),
	}
}

func (p *GameStatePatch_MousePatch) ToMouse() *Mouse {
	m := &Mouse{}

	if p.Name != nil {
		m.Name = p.Name.Value
	}
	if p.Education != nil {
		m.Education = p.Education.Value
	}
	if p.IsBeingEducated != nil {
		m.IsBeingEducated = p.IsBeingEducated.Value
	}
	if p.IsEducated != nil {
		m.IsEducated = p.IsEducated.Value
	}

	return m
}
