package web

type node struct {
	pattern  string  // router path to be matched
	part     string  // matched part in router path
	children []*node // sub path of mathch router
	isWild   bool    // match wildcard */: true
}
