package models

type ChangeSet struct {
	Changes map[string]interface{}
}

func (c *ChangeSet) IsEmpty() bool {
	return len(c.Changes) == 0
}

func (c *ChangeSet) RegisterChange(k string, val interface{}) {
	c.Changes[k] = val
}

func NewChangeSet() *ChangeSet {
	return &ChangeSet{
		Changes: map[string]interface{}{},
	}
}
