package httransform

import (
	"bytes"
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

	connectRequestBufferPool = sync.Pool{
		New: func() interface{} {
			return &bytes.Buffer{}
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
