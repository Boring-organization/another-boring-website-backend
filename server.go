package main

import (
	"TestGoLandProject/cli"
	"TestGoLandProject/core/database"
	"TestGoLandProject/core/router"
	"TestGoLandProject/core/validation"
	"TestGoLandProject/global_consts"
	"TestGoLandProject/graph"
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
	"github.com/gin-gonic/gin"
	"github.com/vektah/gqlparser/v2/ast"
)

func generateSchema() {
	cfg, err := config.LoadConfigFromDefaultLocations()

	if err != nil {
		panic(err)
	}

	err = api.Generate(cfg)

	if err != nil {
		panic(err)
	}
}

func graphqlHandler(resolver graph.ResolverRoot, databaseInstance database.Database) gin.HandlerFunc {
	graphQlConfig := graph.Config{Resolvers: resolver}

	err := validation.ImplementDirectives(&graphQlConfig.Directives, databaseInstance)

	if err != nil {
		panic(err)
	}

	serverHandler := handler.New(graph.NewExecutableSchema(graphQlConfig))

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

	databaseInstance := database.Database{database.InitDb()}

	resolverInstance := &resolver.Resolver{&databaseInstance}
	defer resolverInstance.CloseDb()

	ginRouter := gin.Default()
	ginRouter.Use(GinContextToContextMiddleware())

	ginRouter.POST(global_consts.QueryUrl, graphqlHandler(resolverInstance, databaseInstance))
	ginRouter.GET("/", playgroundHandler())
	router.InitHttpRoutes(ginRouter)

	ginRouter.Run()
}
