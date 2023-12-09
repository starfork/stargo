package server

/*
/可以通过grpc直接转换的方法路由
var Routers = map[string][][]any{
	"serverA": {
		{"GET", "/test1", "Test", &ps.TestRequest{}, &ps.TestResponse{}, "contx/test1"},
		{"GET", "/test2", "Test", &ps.TestRequest{}, &ps.TestResponse{}, "contx/test2"},
		{"GET", "/test3", "Test", &ps.TestRequest{}, &ps.TestResponse{}, "contx/index"},
	},
	"serverB": {
		{"GET", "/test3", "Test", &ps.TestRequest{}, &ps.TestResponse{}, "article2/read"},
	},
}

*/
