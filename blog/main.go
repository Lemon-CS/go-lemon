package main

import (
	"fmt"
	"github.com/Lemon-CS/go-lemon/lime"
	"github.com/Lemon-CS/go-lemon/lime/gopool"
	lLog "github.com/Lemon-CS/go-lemon/lime/log"
	"github.com/Lemon-CS/go-lemon/lime/token"
	"log"
	"net/http"
	"sync"
	"time"
)

type User struct {
	Name      string   `xml:"name" json:"name" lime:"required"`
	Age       int      `xml:"age" json:"age" validate:"required,max=50,min=18"`
	Addresses []string `json:"addresses"`
	Email     string   `json:"email" lime:"required"`
}

func Log(next lime.HandlerFunc) lime.HandlerFunc {
	return func(ctx *lime.Context) {
		fmt.Println("打印请求参数")
		next(ctx)
		fmt.Println("返回执行时间")
	}
}

func main() {
	/*http.HandleFunc("/hello", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, "%s 欢迎来到go-limi", "lemon.com")
	})
	err := http.ListenAndServe(":8111", nil)
	if err != nil {
		log.Fatal(err)
	}*/

	engine := lime.Default()
	g := engine.Group("user")
	g.Get("/hello", func(ctx *lime.Context) {
		fmt.Fprintf(ctx.W, "%s get 欢迎来到go-lime", "lemon.com")
	})

	g.Post("/hello", func(ctx *lime.Context) {
		fmt.Fprintf(ctx.W, "%s post 欢迎来到go-lime", "lemon.com")
	})

	g.Post("/info", func(ctx *lime.Context) {
		fmt.Fprintf(ctx.W, "%s info", "lemon.com")
	})
	g.Get("/get/:id", func(ctx *lime.Context) {
		fmt.Fprintf(ctx.W, "%s get user info path variable", "lemon.com")
	})

	g.Get("/html", func(ctx *lime.Context) {
		ctx.HTML(http.StatusOK, "<h1>码神之路</h1>")
	})
	g.Get("/htmlTemplate", func(ctx *lime.Context) {
		user := &User{
			Name: "码神之路",
		}
		err := ctx.HTMLTemplate("login.html", user, "tpl/login.html", "tpl/header.html")
		if err != nil {
			log.Println(err)
		}
	})
	g.Get("/htmlTemplateGlob", func(ctx *lime.Context) {
		user := &User{
			Name: "码神之路",
		}
		err := ctx.HTMLTemplateGlob("login.html", user, "tpl/*.html")
		if err != nil {
			log.Println(err)
		}
	})
	//engine.LoadTemplate("tpl/*.html")

	g.Get("/template", func(ctx *lime.Context) {
		user := &User{
			Name: "码神之路",
		}
		err := ctx.Template("login.html", user)
		if err != nil {
			log.Println(err)
		}
	})

	g.Get("/json", func(ctx *lime.Context) {
		user := &User{
			Name: "码神之路",
		}
		err := ctx.JSON(http.StatusOK, user)
		if err != nil {
			log.Println(err)
		}
	})

	g.Get("/xml", func(ctx *lime.Context) {
		user := &User{
			Name: "码神之路",
			Age:  10,
		}
		err := ctx.XML(http.StatusOK, user)
		if err != nil {
			log.Println(err)
		}
	})

	g.Get("/excel", func(ctx *lime.Context) {
		ctx.File("tpl/test.xlsx")
	})
	g.Get("/excelName", func(ctx *lime.Context) {
		ctx.FileAttachment("tpl/test.xlsx", "aaaa.xlsx")
	})
	g.Get("/fs", func(ctx *lime.Context) {
		ctx.FileFromFS("test.xlsx", http.Dir("tpl"))
	})
	g.Get("/redirect", func(ctx *lime.Context) {
		ctx.Redirect(http.StatusFound, "/user/template")
	})
	g.Get("/string", func(ctx *lime.Context) {
		ctx.String(http.StatusFound, "和 %s %s学习 goweb框架", "码神之路", "从零")
	})

	g.Get("/add", func(ctx *lime.Context) {
		name := ctx.GetDefaultQuery("name", "张三")
		fmt.Printf("name: %v , ok: %v \n", name, true)
	})
	g.Get("/queryMap", func(ctx *lime.Context) {
		m, _ := ctx.GetQueryMap("user")
		ctx.JSON(http.StatusOK, m)
	})
	g.Post("/formPost", func(ctx *lime.Context) {
		m, _ := ctx.GetPostFormMap("user")
		//file := ctx.FormFile("file")
		//err := ctx.SaveUploadedFile(file, "./upload/"+file.Filename)
		//if err != nil {
		//	logger.Println(err)
		//}
		files := ctx.FormFiles("file")
		for _, file := range files {
			ctx.SaveUploadedFile(file, "./upload/"+file.Filename)
		}
		ctx.JSON(http.StatusOK, m)
	})
	//g.Post("/file", func(ctx *lime.Context) {
	//	m, _ := ctx.GetPostFormMap("user")
	//
	//	ctx.JSON(http.StatusOK, m)
	//})

	g.Post("/jsonParam", func(ctx *lime.Context) {
		user := make([]User, 0)
		ctx.DisallowUnknownFields = true
		//ctx.IsValidate = true
		err := ctx.BindJson(&user)
		if err == nil {
			ctx.JSON(http.StatusOK, user)
		} else {
			log.Println(err)
		}
	})
	engine.Logger.Level = lLog.LevelDebug
	//logger.Outs = append(logger.Outs, msLog.FileWriter("./log/log.log"))
	engine.Logger.LogFileSize = 1 << 10
	g.Post("/xmlParam", func(ctx *lime.Context) {
		user := &User{}
		_ = ctx.BindXML(user)
		ctx.Logger.WithFields(lLog.Fields{
			"name": "码神之路",
			"id":   1000,
		}).Debug("我是debug日志")
		ctx.Logger.Info("我是info日志")
		ctx.Logger.Error("我是error日志")
		//err := mserror.Default()
		//err.Result(func(msError *mserror.MsError) {
		//	ctx.Logger.Info(msError.Error())
		//	ctx.JSON(http.StatusInternalServerError, user)
		//})
		//a(1, err)
		//b(1, err)
		//c(1, err)
		ctx.JSON(http.StatusOK, user)
		//err := login()
		//ctx.HandleWithError(http.StatusOK, user, err)
	})

	p, _ := gopool.NewPool(15)
	g.Post("/pool", func(ctx *lime.Context) {
		currentTime := time.Now().UnixMilli()
		var wg sync.WaitGroup
		wg.Add(5)
		p.Submit(func() {
			defer func() {
				wg.Done()
			}()
			fmt.Println("1111111")
			//panic("这是1111的panic")
			time.Sleep(3 * time.Second)

		})
		p.Submit(func() {
			fmt.Println("22222222")
			time.Sleep(3 * time.Second)
			wg.Done()
		})
		p.Submit(func() {
			fmt.Println("33333333")
			time.Sleep(3 * time.Second)
			wg.Done()
		})
		p.Submit(func() {
			fmt.Println("44444")
			time.Sleep(3 * time.Second)
			wg.Done()
		})
		p.Submit(func() {
			fmt.Println("55555555")
			time.Sleep(3 * time.Second)
			wg.Done()
		})
		wg.Wait()
		fmt.Printf("time: %v \n", time.Now().UnixMilli()-currentTime)
		ctx.JSON(http.StatusOK, "success")
	})
	g.Get("/login", func(ctx *lime.Context) {
		jwt := &token.JwtHandler{}
		jwt.Key = []byte("123456")
		jwt.SendCookie = true
		jwt.TimeOut = 10 * time.Minute
		jwt.RefreshTimeOut = 20 * time.Minute
		jwt.Authenticator = func(ctx *lime.Context) (map[string]any, error) {
			data := make(map[string]any)
			data["userId"] = 1
			return data, nil
		}
		token, err := jwt.LoginHandler(ctx)
		if err != nil {
			log.Println(err)
			ctx.JSON(http.StatusOK, err.Error())
			return
		}
		ctx.JSON(http.StatusOK, token)
	})

	g.Get("/refresh", func(ctx *lime.Context) {
		jwt := &token.JwtHandler{}
		jwt.Key = []byte("123456")
		jwt.SendCookie = true
		jwt.TimeOut = 10 * time.Minute
		jwt.RefreshTimeOut = 20 * time.Minute
		jwt.RefreshKey = "blog_refresh_token"
		ctx.Set(jwt.RefreshKey, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTY1OTU3NzUsImlhdCI6MTY1NjU5NDU3NSwidXNlcklkIjoxfQ.v5rMFD-3kScPrbv6YOPR0ec9mpp84cXA14ZShVCTwC0")
		token, err := jwt.RefreshHandler(ctx)
		if err != nil {
			log.Println(err)
			ctx.JSON(http.StatusOK, err.Error())
			return
		}
		ctx.JSON(http.StatusOK, token)
	})

	//engine.Run()
	engine.RunTLS(":8118", "key/server.pem", "key/server.key")
}
