# micro-peers

Micro peers is a side project where multiple services are connected to network through consul.
These services communicate with each other and also discover other services. In future we will add more services that will discover other services and share data arcoss data centers

# Consul Node :

    consul agent -dev -node=machine 
    Consul agent initate on your system.

    Consul run through configuration files
    consul agent -dev -config-dir=./web.json -node=machine

    Checkout your consul dns information
    dig @127.0.0.1 -p 8600 service_name.service.consul

# How to run app

   - 👀 forking the project
   - 👽 Another terminal Ctrl+T and type cd client/ && go run main.go 
   - 💻 Open terminal Ctrl+T and run go run main.go


# RUN CONSUL :

   - 👀 fork the project and type ctrl+ t ; then cd consul
   - 😘 consul agent -dev -config-dir=./web.json -node=machine

    check on 127.0.0.1:8500
