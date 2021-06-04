package router

import (
	"context"
	"github.com/v2fly/v2ray-core/v4/features/routing"

	"github.com/v2fly/v2ray-core/v4/features/extension"
	"github.com/v2fly/v2ray-core/v4/features/outbound"
)

type BalancingStrategy interface {
	PickOutbound([]string) string
}

type Balancer struct {
	selectors   []string
	strategy    routing.BalancingStrategy
	ohm         outbound.Manager
	fallbackTag string

	override override
}

// PickOutbound picks the tag of a outbound
func (b *Balancer) PickOutbound() (string, error) {
	candidates, err := b.SelectOutbounds()
	if err != nil {
		if b.fallbackTag != "" {
			newError("fallback to [", b.fallbackTag, "], due to error: ", err).AtInfo().WriteToLog()
			return b.fallbackTag, nil
		}
		return "", err
	}
	var tag string
	if o := b.override.Get(); o != nil {
		tag = b.strategy.Pick(o.selects)
	} else {
		tag = b.strategy.SelectAndPick(candidates)
	}
	if tag == "" {
		if b.fallbackTag != "" {
			newError("fallback to [", b.fallbackTag, "], due to empty tag returned").AtInfo().WriteToLog()
			return b.fallbackTag, nil
		}
		// will use default handler
		return "", newError("balancing strategy returns empty tag")
	}
	return tag, nil
}

func (b *Balancer) InjectContext(ctx context.Context) {
	if contextReceiver, ok := b.strategy.(extension.ContextReceiver); ok {
		contextReceiver.InjectContext(ctx)
	}
}

// SelectOutbounds select outbounds with selectors of the Balancer
func (b *Balancer) SelectOutbounds() ([]string, error) {
	hs, ok := b.ohm.(outbound.HandlerSelector)
	if !ok {
		return nil, newError("outbound.Manager is not a HandlerSelector")
	}
	tags := hs.Select(b.selectors)
	return tags, nil
}
