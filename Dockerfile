FROM scratch
MAINTAINER Sope Shen "shenshouer51@gmail.com"

EXPOSE 8080

COPY k8s-ui /
COPY templates /templates
COPY static /static

ENV PORT=8080

CMD ["/k8s-ui", "-alsologtostderr=true -logtostderr=true -addr=:8080"]