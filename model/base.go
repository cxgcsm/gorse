package model

import "github.com/zhenghaoz/gorse/core"
import "github.com/zhenghaoz/gorse/base"

/* Base */

// Base model must be included by every recommendation model. Hyper-parameters,
// ID sets, random generator and fitting options are managed the Base model.
type Base struct {
	Params      base.Params          // Hyper-parameters
	UserIndexer *base.Indexer        // Users' ID set
	ItemIndexer *base.Indexer        // Items' ID set
	rng         base.RandomGenerator // Random generator
	randState   int64                // Random seed
	// Tracker
	isSetParamsCalled bool // Check whether SetParams called
}

// SetParams sets hyper-parameters for the Base model.
func (model *Base) SetParams(params base.Params) {
	model.isSetParamsCalled = true
	model.Params = params
	model.randState = model.Params.GetInt64(base.RandomState, 0)
}

// GetParams returns all hyper-parameters.
func (model *Base) GetParams() base.Params {
	return model.Params
}

// Predict has not been implemented.
func (model *Base) Predict(userId, itemId int) float64 {
	panic("Predict() not implemented")
}

// Fit has not been implemented,
func (model *Base) Fit(trainSet core.DataSet, options *base.RuntimeOptions) {
	panic("Fit() not implemented")
}

// Init the Base model. The method must be called at the beginning of Fit.
func (model *Base) Init(trainSet core.DataSetInterface) {
	// Check Base.SetParams() called
	if model.isSetParamsCalled == false {
		panic("Base.SetParams() not called")
	}
	// Setup ID set
	model.UserIndexer = trainSet.UserIndexer()
	model.ItemIndexer = trainSet.ItemIndexer()
	// Setup random state
	model.rng = base.NewRandomGenerator(model.randState)
}

// ItemPop recommends items by their popularity. The popularity of a item is
// defined as the occurrence frequency of the item in the training data set.
type ItemPop struct {
	Base
	Pop []float64
}

// NewItemPop creates an ItemPop model.
func NewItemPop(params base.Params) *ItemPop {
	pop := new(ItemPop)
	pop.SetParams(params)
	return pop
}

// Fit the ItemPop model.
func (pop *ItemPop) Fit(set core.DataSetInterface, options *base.RuntimeOptions) {
	pop.Init(set)
	// Get items' popularity
	pop.Pop = make([]float64, set.ItemCount())
	for i := 0; i < set.ItemCount(); i++ {
		pop.Pop[i] = float64(set.ItemByIndex(i).Len())
	}
}

// Predict by the ItemPop model.
func (pop *ItemPop) Predict(userId, itemId string) float64 {
	// Return items' popularity
	denseItemId := pop.ItemIndexer.ToIndex(itemId)
	if denseItemId == base.NotId {
		return 0
	}
	return pop.Pop[denseItemId]
}
