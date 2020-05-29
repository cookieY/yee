package yee

//type group struct {
//	prefix     string
//	middleware []HandlerFunc
//	core       *Core
//	root       bool
//}






//func (g *group) createDistHandler(relativePath string, fs http.FileSystem) HandlerFunc {
//	ab := path.Join(g.prefix, relativePath)
//	fs_1 := http.StripPrefix(ab, http.FileServer(fs))
//	return func(c Context) (err error) {
//		file := c.Params("filepath")
//		if _, err := fs.Open(file); err != nil {
//			return c.String(http.StatusNotFound, fmt.Sprintf("404 NOT FOUND!"))
//		}
//		fs_1.ServeHTTP(c.Response(), c.Request())
//		return
//	}
//
//}

//func (g *group) Static(relativePath, dist string) {
//	handler := g.createDistHandler(relativePath, http.Dir(dist))
//	url := path.Join(relativePath, "/*filepath")
//	g.GET(url, handler)
//}
