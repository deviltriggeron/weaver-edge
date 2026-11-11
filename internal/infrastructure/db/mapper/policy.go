package mapper

import (
	"weaver-edge/internal/domain"
	"weaver-edge/internal/infrastructure/db/model"
)

func PolicyModelToDomain(policy model.PolicyModel) domain.Policy {
	return domain.Policy{
		ID:        policy.ID,
		Resource:  resourceModelToDomain(policy),
		CreatedAt: policy.CreatedAt,
	}
}

func resourceModelToDomain(policy model.PolicyModel) domain.Resource {
	return domain.Resource{
		ChannelID: policy.ChannelID,
		ThreadID:  policy.ThreadID,
	}
}
