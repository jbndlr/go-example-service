FROM jbndlr/dev-go:0.0.4

RUN go get -u -v \
    gopkg.in/yaml.v2 \
    github.com/gin-gonic/gin \
    github.com/dgrijalva/jwt-go \
    github.com/kelseyhightower/envconfig \
    golang.org/x/sync/errgroup

EXPOSE 7000
EXPOSE 8000
EXPOSE 9000
