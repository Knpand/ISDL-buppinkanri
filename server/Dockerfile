# 2020/10/14最新versionを取得
FROM golang:1.17.9-alpine
# アップデートとgitのインストール！！
RUN apk update && apk add git
# ワーキングディレクトリの設定
WORKDIR /app
# CMD ["/bin/ash", "../init/init.sh"]
RUN go install github.com/cosmtrek/air@latest

# RUN mod tidy -go=1.17
# RUN go mod download
# RUN go build -o /go/bin/air github.com/cosmtrek/air

#air導入によりgoのホットリロードが可能
CMD ["air"]

# ADD . /app  
#  <= ymlでバインドしてßるから必要ない

# /bin/ash に入る
