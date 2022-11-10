package main

import "github.com/csby/gwsf/gtype"

type controllerOpt struct {
}

func (s *controllerOpt) initController(h *Handler) {
}

func (s *controllerOpt) initRouter(router gtype.Router, path *gtype.Path, preHandle gtype.HttpHandle) {

}
