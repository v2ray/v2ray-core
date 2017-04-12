package mux

type statusHandler func(meta *FrameMetadata) error
