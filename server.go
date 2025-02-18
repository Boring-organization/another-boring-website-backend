package main

import (
	"TestGoLandProject/cli"
	"TestGoLandProject/core/database"
	"TestGoLandProject/core/router"
	"TestGoLandProject/core/utils/common"
	"TestGoLandProject/core/validation"
	"TestGoLandProject/global_const"
	graph "TestGoLandProject/graph/generated"
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
	sq "github.com/Masterminds/squirrel"
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

func graphqlHandler(resolver graph.ResolverRoot, databaseInstance sq.StatementBuilderType) gin.HandlerFunc {
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
	playgroundHandler := playground.Handler("GraphQL", globalConst.ApiUrl)

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

	initializedDatabase, err := database.InitDb()
	if err != nil {
		panic(err)
	}

	defer initializedDatabase.Close()
	sqDb := sq.StatementBuilder.RunWith(initializedDatabase)

	resolverInstance := &resolver.Resolver{sqDb}

	ginRouter := gin.Default()
	ginRouter.Use(GinContextToContextMiddleware())

	ginRouter.POST(globalConst.QueryUrl, graphqlHandler(resolverInstance, sqDb))
	ginRouter.GET("/", playgroundHandler())
	router.InitHttpRoutes(ginRouter)

	ctx, cancel := context.WithCancel(context.Background())
	go commonUtils.PeriodicImageCleaner(ctx, sqDb)

	defer cancel()

	ginRouter.Run()
}
