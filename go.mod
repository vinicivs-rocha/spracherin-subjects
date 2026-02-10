module github.com/vinicivs-rocha/spracherin-subjects

go 1.25.6

require (
	github.com/go-sql-driver/mysql v1.9.0
	github.com/google/uuid v1.6.0
	github.com/joserocha/spracherin/spracherin-proto v0.0.0
	google.golang.org/grpc v1.70.0
)

replace github.com/joserocha/spracherin/spracherin-proto => ../spracherin-proto

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	golang.org/x/net v0.32.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241202173237-19429a94021a // indirect
	google.golang.org/protobuf v1.36.5 // indirect
)
