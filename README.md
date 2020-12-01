## Example-Service

Start VSCode with the folder containing the ``.devcontainer.json`` as the root folder; VSCode will bring up a panel in bottom right corner asking whether to restart in container. Choose yes.

VSCode then brings up the desired development container and asks to add extensions where required. It's up to you to install them up front via ``.devcontainer.json``'s ``extensions`` list.

When developing and in container mode (check bottom left corner of info bar saying "Dev Container"), hit ^F5 (Ctrl + F5) to run without debugging and VSCode will start your application within the dev container. If the application already listens on a port (and that port is exposed in ``appPorts`` list), direct your local browser to ``localhost:<port>`` and see your application's output.

# Build Protocol Buffers

```bash
protoc --proto_path=example/api/grpc/proto --go_out=example/api/grpc/pb --go_opt=paths=source_relative --go-grpc_out=example/api/grpc/pb --go-grpc_opt=paths=source_relative example/api/grpc/proto/*.proto
```

# Access gRPC Methods Using ``grpcurl``

```bash
grpcurl -plaintext localhost:9000 grpc.Example/Info
```

# Authenticate with ``curl``

Authenticate using credentials and store the resulting cookie to ``/tmp/cookie``:

```bash
curl -iL -d '{"identity": "user1", "secret": "pw1"}' -c /tmp/cookie -X POST localhost:8000/auth
```

Access a restricted resource providing the previously saved cookie:

```bash
curl -iL -b /tmp/cookie localhost:8000/secret
```
