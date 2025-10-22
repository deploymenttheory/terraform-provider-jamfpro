package framework

import (
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// CreateResponseContainer wraps resource.CreateResponse to implement StateContainer
type CreateResponseContainer struct {
	*resource.CreateResponse
}

func (c *CreateResponseContainer) GetState() tfsdk.State {
	return c.State
}

func (c *CreateResponseContainer) SetState(state tfsdk.State) {
	c.State = state
}

// UpdateResponseContainer wraps resource.UpdateResponse to implement StateContainer
type UpdateResponseContainer struct {
	*resource.UpdateResponse
}

func (c *UpdateResponseContainer) GetState() tfsdk.State {
	return c.State
}

func (c *UpdateResponseContainer) SetState(state tfsdk.State) {
	c.State = state
}

// StateContainer interface for anything that has a State field
type StateContainer interface {
	GetState() tfsdk.State
	SetState(tfsdk.State)
}
