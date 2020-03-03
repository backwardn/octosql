package logical

import (
	"context"
	"github.com/cube2222/octosql"
	"github.com/cube2222/octosql/physical"
	"github.com/pkg/errors"
)

type Offset struct {
	data       Node
	offsetExpr Expression
}

func NewOffset(data Node, expr Expression) Node {
	return &Offset{data: data, offsetExpr: expr}
}

func (node *Offset) Physical(ctx context.Context, physicalCreator *PhysicalPlanCreator) ([]physical.Node, octosql.Variables, error) {
	sourceNodes, variables, err := node.data.Physical(ctx, physicalCreator)
	if err != nil {
		return nil, nil, errors.Wrap(err, "couldn't get physical plan for source nodes")
	}

	offsetExpr, offsetVariables, err := node.offsetExpr.Physical(ctx, physicalCreator)
	if err != nil {
		return nil, nil, errors.Wrap(err, "couldn't get physical plan for offset expression")
	}
	variables, err = variables.MergeWith(offsetVariables)
	if err != nil {
		return nil, nil, errors.Wrap(err, "couldn't get offset node variables")
	}

	return []physical.Node{physical.NewOffset(physical.NewUnionAll(sourceNodes...), offsetExpr)}, variables, nil
}
