package main

import (
	"log"
	pb "github.com/raissov/shippy/shippy-service-user/proto/user"
	"github.com/micro/go-micro/v2"
)

const schema = `
	create table if not exists users (
		id varchar(36) not null,
		name varchar(125) not null,
		email varchar(225) not null unique,
		password varchar(225) not null,
		company varchar(125),
		primary key (id)
	);
`

func main() {


	db, err := NewConnection()
	if err != nil {
		log.Panic(err)
	}

	defer db.Close()

	if err != nil {
		log.Fatalf("Could not connect to DB: %v", err)
	}


	db.MustExec(schema)

	repo := NewPostgresRepository(db)

	tokenService := &TokenService{repo}

	// Create a new service. Optionally include some options here.
	service := micro.NewService(
		micro.Name("shippy.service.user"),
		micro.Version("latest"),
	)


	service.Init()


	if err := pb.RegisterUserServiceHandler(service.Server(), &handler{repo, tokenService}); err != nil {
		log.Panic(err)
	}


	if err := service.Run(); err != nil {
		log.Panic(err)
	}
}
