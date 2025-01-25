package main

import (
	"TestGoLandProject/database"
	"TestGoLandProject/graph"
	"TestGoLandProject/graph/resolvers"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/vektah/gqlparser/v2/ast"
)

const apiUrl = "/api"

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
	playgroundHandler := playground.Handler("GraphQL", apiUrl)

	return func(c *gin.Context) {
		playgroundHandler.ServeHTTP(c.Writer, c.Request)
	}
}

func main() {
	resolverInstance := &resolver.Resolver{&database.Database{database.InitDb()}}
	defer resolverInstance.CloseDb()

	ginRouter := gin.Default()

	ginRouter.POST(apiUrl, graphqlHandler(resolverInstance))
	ginRouter.GET("/", playgroundHandler())

	ginRouter.Run()
}
