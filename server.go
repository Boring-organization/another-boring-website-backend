package main

import (
	"TestGoLandProject/cli"
	"TestGoLandProject/core/database"
	"TestGoLandProject/core/router"
	"TestGoLandProject/global_consts"
	"TestGoLandProject/graph"
	directiveValidationHooks "TestGoLandProject/graph/directive_hooks/validation"
	"TestGoLandProject/graph/resolvers"
	"context"
	"fmt"
	"github.com/99designs/gqlgen/api"
	"github.com/99designs/gqlgen/codegen/config"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/99designs/gqlgen/plugin/modelgen"
	"github.com/gin-gonic/gin"
	"github.com/vektah/gqlparser/v2/ast"
)

func generateSchema() {
	cfg, err := config.LoadConfigFromDefaultLocations()

	if err != nil {
		panic(err)
	}

	p := modelgen.Plugin{
		FieldHook: directiveValidationHooks.CustomFieldHook,
	}

	err = api.Generate(cfg, api.ReplacePlugin(&p))

	if err != nil {
		panic(err)
	}

}

func graphqlHandler(resolver graph.ResolverRoot) gin.HandlerFunc {

	serverHandler := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	serverHandler.AddTransport(transport.Options{})
	serverHandler.AddTransport(transport.GET{})
	serverHandler.AddTransport(transport.POST{})

	serverHandler.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	serverHandler.Use(extension.Introspection{})
	serverHandler.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	return func(c *gin.Context) {
		serverHandler.ServeHTTP(c.Writer, c.Request)
	}
}

func playgroundHandler() gin.HandlerFunc {
	playgroundHandler := playground.Handler("GraphQL", global_consts.ApiUrl)

	return func(c *gin.Context) {
		playgroundHandler.ServeHTTP(c.Writer, c.Request)
	}
}
func GinContextToContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), "GinContextKey", c)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func main() {
	if cli.GetArgs(cli.GenerateSchema) == true {
		fmt.Println("[GraphQL] Start schema generation...")
		generateSchema()
		fmt.Println("[GraphQL] Successfully generated schema")
		fmt.Println()
	} else {
		fmt.Println("[GraphQL] Skip schema generation\n")
	}

	resolverInstance := &resolver.Resolver{&database.Database{database.InitDb()}}
	defer resolverInstance.CloseDb()

	ginRouter := gin.Default()
	ginRouter.Use(GinContextToContextMiddleware())

	ginRouter.POST(global_consts.QueryUrl, graphqlHandler(resolverInstance))
	ginRouter.GET("/", playgroundHandler())
	router.InitHttpRoutes(ginRouter)

	ginRouter.Run()
}
