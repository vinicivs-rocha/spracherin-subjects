# gRPC Server

**How To Run**
1. Build a server binary that calls `grpc.ListenAndServe` with a `SubjectRepository` and `Messager`.
2. Provide the DB environment variables described in `docs/db.md`.

**Notes**
1. The handlers use the existing `DescribeSubject` application flow and call the repository directly for list/detail/remove.
