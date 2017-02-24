## Reproduction of unexpected behavior for close, streams, and graceful stop

There are two test cases in `main_test.go`, one which closes a streaming client connection and then gracefully shuts down the server, and second which triggers graceful shutdown first before stopping the client connection.

The second case prints some errors to stderr before hanging forever.

It's unclear if this is specific to testing locally, or if this bug might be triggered in production settings as well.

The lockup seems to be alleviated by injecting a cancelable context into the client stream, which is canceled in `client.Close()`, but it's unclear from documentation that this is required for correct usage of streams.
