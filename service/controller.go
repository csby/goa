package main

type controllers struct {
	opt  controllerOpt
	app  controllerApp
	auth controllerAuth
}

func (s *Handler) initController() {
	s.ctrl.opt.initController(s)
	s.ctrl.app.initController(s)
	s.ctrl.auth.initController(s)
}
