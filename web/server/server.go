package server

import (
	"apps/spiderpool/funcs"
	sf "apps/spiderpool/funcs/spiderpool"
	"context"
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/starfork/stargo/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Server struct {
	g    *gin.Engine
	conn map[string]grpc.ClientConnInterface
	ctx  context.Context
}

func NewServer(conf map[string]*config.Config) (*Server, error) {

	s := &Server{
		conn: make(map[string]grpc.ClientConnInterface),
	}

	for k, v := range conf {
		s.conn[k], _ = grpc.Dial(v.ServerPort,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)

	}
	return s, nil
}

func (s *Server) Start(g *gin.Engine, efs embed.FS, viewPath string) {

	s.ctx = context.Background()
	fcs := &funcs.Funcs{
		Sp: sf.NewFuncs(s.ctx, s.conn["spiderpool"]),
	}

	tpl := template.Must(template.New("").Funcs(funcs.Register(fcs)).
		ParseFS(efs, "views/**/**/**/*"))
	g.SetHTMLTemplate(tpl)

	s.g = g
	routers := make(map[string][][]any)
	s.Register(routers)

	// //handler 应该不会有很多才对
	// sh := &sh.Handler{
	// 	Tpl:  tpl,
	// 	Ns:   "/go.park.spider.SpiderpoolHandler/",
	// 	Conn: s.conn["spiderpool"],
	// 	Ctx:  s.ctx,
	// }

	// s.g.Handle("GET", "/article/:name", sh.ReadArticle)

}

func (s *Server) Register(routers map[string][][]any) {
	for k, rts := range routers {
		for _, r := range rts {
			tplName := strings.ToLower(r[2].(string))
			if len(r) == 6 {
				tplName = r[5].(string)
			}
			s.g.Handle(r[0].(string), r[1].(string), func(c *gin.Context) {
				out := r[4]
				if err := s.Invoke(k, r[2].(string), r[3], out); err != nil {
					fmt.Println(err) //throw http error
				}

				//添加其他变量

				c.HTML(http.StatusOK, tplName, out)
			})
		}
	}
}

func (s *Server) Invoke(serviceName, method string, args any, reply any, opts ...grpc.CallOption) error {
	ns := fmt.Sprintf("/go.park.%s.%sHandler/", strings.ToLower(serviceName), serviceName)
	return s.conn[serviceName].Invoke(s.ctx, ns+method, args, reply, opts...)
}
