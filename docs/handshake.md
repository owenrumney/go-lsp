# Handshake and Message Order

Every LSP session runs in the same order: negotiate, work, shut down. The order is not advisory. A server must reject work before `initialize`, and a client that skips `initialized` leaves the server dropping its notifications.

## The sequence

```mermaid
sequenceDiagram
    autonumber
    participant C as Client (package client)
    participant S as Server (package server)

    Note over C,S: Negotiate
    C->>S: spawn process, connect stdin/stdout
    C->>S: initialize (clientInfo, capabilities, rootUri + workspaceFolders, processId, initializationOptions)
    S->>S: handler.Initialize, then merge auto-detected capabilities
    S-->>C: InitializeResult (serverInfo, capabilities, positionEncoding)
    C->>S: initialized
    Note over C,S: Server now accepts requests, readiness grace starts

    Note over C,S: Work
    C->>S: textDocument/didOpen
    S--)C: textDocument/publishDiagnostics
    opt dynamic registration
        S->>C: client/registerCapability
        C-->>S: null
    end
    opt long startup work
        S->>C: window/workDoneProgress/create
        C-->>S: null
        S--)C: $/progress begin ... end
    end
    C->>S: textDocument/hover (and other requests)
    S-->>C: result
    opt server needs the client
        S->>C: workspace/configuration, workspace/applyEdit, window/showMessageRequest, window/showDocument
        C-->>S: result
    end

    Note over C,S: Shut down
    C->>S: shutdown
    S-->>C: null
    C->>S: exit
    S->>S: process exits
```

## What negotiation actually decides

| Decided at | What |
| --- | --- |
| Client capabilities in `initialize` | Which server→client requests the client will honour, and which dynamic registrations it will act on. `DefaultCapabilities` advertises `dynamicRegistration` only where `Client` honours it, because over-declaring makes a server withhold static capabilities. |
| Server capabilities in `InitializeResult` | Which requests the client may send. go-lsp builds these from the interfaces your handler implements and merges them into whatever your `Initialize` returned, so explicit values win. |
| `positionEncoding` | How offsets in every later position are read. The client defaults to UTF-16 and switches only if the server names another encoding. `DefaultCapabilities` offers UTF-16 alone. |
| `client/registerCapability`, after the fact | Capabilities the server adds later. `Client.Supports(method, uri, languageID)` answers over both the static set and the live registrations. |

## Rules the server enforces

- A request before `initialize` completes fails with `-32002` (`ServerNotInitialized`).
- A second `initialize` fails with `InvalidRequest`. A failed `initialize` resets the state, so the client may retry.
- Notifications outside the initialized window are dropped, except `exit`.
- After `shutdown`, every request fails with `InvalidRequest`. The diagnostics engine stops before `handler.Shutdown` runs.

## Readiness is not part of the handshake

A server answers `initialize` long before it has indexed anything. Ask for references in that window and the empty list looks like a real answer.

LSP has no readiness signal, so `Client.Ready` approximates one: it tracks work-done progress tokens from `window/workDoneProgress/create` and `$/progress`, and reports ready when none are outstanding. For a server that reports no progress at all, it waits out `Config.ReadyGrace` (2s by default), timed from the end of `initialize`. Requests are never blocked on it.

See [`client.md`](./client.md) for the client API and [`documents.md`](./documents.md) for position encoding in practice.
