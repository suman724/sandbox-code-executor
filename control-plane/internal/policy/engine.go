package policy

import (
	"context"
	"errors"
	"sync"

	"github.com/open-policy-agent/opa/rego"
)

type Decision struct {
	Allowed bool
	Reason  string
}

type Evaluator interface {
	Evaluate(ctx context.Context, input any) (Decision, error)
}

type RulesetResolver interface {
	Ruleset(ctx context.Context, input any) (string, error)
}

type StaticRulesetResolver struct {
	RulesetText string
}

func (r StaticRulesetResolver) Ruleset(ctx context.Context, input any) (string, error) {
	_ = ctx
	_ = input
	if r.RulesetText == "" {
		return "", errors.New("empty ruleset")
	}
	return r.RulesetText, nil
}

type OPAEvaluator struct {
	Resolver RulesetResolver
	Query    string
	mu       sync.RWMutex
	cache    map[string]rego.PreparedEvalQuery
}

func (e *OPAEvaluator) Evaluate(ctx context.Context, input any) (Decision, error) {
	if e.Resolver == nil {
		return Decision{}, errors.New("missing ruleset resolver")
	}
	query := e.Query
	if query == "" {
		query = "data.policy"
	}
	ruleset, err := e.Resolver.Ruleset(ctx, input)
	if err != nil {
		return Decision{}, err
	}
	prepared, err := e.prepare(query, ruleset)
	if err != nil {
		return Decision{}, err
	}
	results, err := prepared.Eval(ctx, rego.EvalInput(input))
	if err != nil {
		return Decision{}, err
	}
	if len(results) == 0 || len(results[0].Expressions) == 0 {
		return Decision{Allowed: false, Reason: "no decision"}, nil
	}
	obj, ok := results[0].Expressions[0].Value.(map[string]any)
	if !ok {
		return Decision{}, errors.New("unexpected policy result")
	}
	allowed, _ := obj["allow"].(bool)
	reason, _ := obj["reason"].(string)
	if reason == "" && !allowed {
		reason = "denied"
	}
	return Decision{Allowed: allowed, Reason: reason}, nil
}

func (e *OPAEvaluator) prepare(query string, ruleset string) (rego.PreparedEvalQuery, error) {
	e.mu.RLock()
	if e.cache != nil {
		if prepared, ok := e.cache[ruleset]; ok {
			e.mu.RUnlock()
			return prepared, nil
		}
	}
	e.mu.RUnlock()

	compiler, err := rego.New(
		rego.Query(query),
		rego.Module("policy.rego", ruleset),
	).PrepareForEval(context.Background())
	if err != nil {
		return rego.PreparedEvalQuery{}, err
	}

	e.mu.Lock()
	if e.cache == nil {
		e.cache = make(map[string]rego.PreparedEvalQuery)
	}
	e.cache[ruleset] = compiler
	e.mu.Unlock()

	return compiler, nil
}
