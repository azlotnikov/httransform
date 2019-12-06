package httransform

import (
	"sync"
)

var (
	layerStatePool = sync.Pool{
		New: func() interface{} {
			return &LayerState{
				ctx: make(map[string]interface{}),
			}
		},
	}
)

func getLayerState() *LayerState {
	return layerStatePool.Get().(*LayerState)
}

func releaseLayerState(state *LayerState) {
	for k := range state.ctx {
		delete(state.ctx, k)
	}

	layerStatePool.Put(state)
}
