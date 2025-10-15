package main

import (
    "github.com/gin-gonic/gin"
    "github.com/99designs/gqlgen/graphql/handler"
    "github.com/99designs/gqlgen/graphql/playground"
    "chatear-backend/graph" // ajuste o import conforme o módulo do seu go.mod
)

func main() {
    r := gin.Default()

    srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

    r.POST("/query", gin.WrapH(srv))
    r.GET("/", gin.WrapH(playground.Handler("GraphQL playground", "/query")))

    r.Run(":8080")
}

