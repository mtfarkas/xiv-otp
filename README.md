# xivotp

A tiny CLI companion for Final Fantasy XIV players who use [XIVLauncher](https://github.com/goatcorp/FFXIVQuickLauncher). It automates the tedious copy-paste dance with your one-time password during login.

## What problem does this solve?

If you have two-factor authentication enabled on your Square Enix account, XIVLauncher needs a 6-digit TOTP code every time you log in. XIVLauncher has a built-in "OTP macro" feature — it listens on `http://localhost:4646` for a code sent via HTTP, and automatically fills it into the login form. But you still need *something* to generate and send that code.

That's where `xivotp` comes in. It:

1. Stores your OTP Auth URL securely in your operating system's keyring (Windows Credential Manager, macOS Keychain, or Linux Secret Service).
2. Generates a fresh TOTP code from that secret.
3. Fires it over to XIVLauncher's local HTTP endpoint so you don't have to type or paste anything.

One command and you're logged in.

## Installation

Pre-built binaries for Windows, macOS, and Linux are published on the [GitHub Releases](https://github.com/mfarkas/xiv-otp/releases) page — built and uploaded automatically by CI. Pick the archive for your platform, extract it, and put `xivotp` (or `xivotp.exe` on Windows) somewhere on your PATH.

If you have Go installed and prefer to build from source:

```sh
go install github.com/mfarkas/xiv-otp@latest
```

Or clone and build manually:

```sh
git clone https://github.com/mfarkas/xiv-otp.git
cd xiv-otp
go build -o xivotp .
```

## How to get your OTP Auth URL

<TODO>

## Usage

### Storing a secret

Before you can log in, you need to save your OTP Auth URL to the keyring. You only need to do this once per account.

```sh
xivotp secret add "otpauth://totp/Square%20Enix%20ID:..." --account mymain
```

The `--account` flag lets you name the account so you can manage multiple characters or service accounts. It defaults to `default` if you leave it out, which is perfect if you only have one.

### Logging in

Make sure XIVLauncher is running and sitting at the login screen waiting for a one-time password. Then:

```sh
xivotp auth login
```

If you stored the secret under a different account name:

```sh
xivotp auth login --account mymain
```

The tool generates your code, sends it to XIVLauncher, and you're in. It retries up to 3 times with a 2-second gap between attempts — useful if XIVLauncher is still waking up.

If it can't reach XIVLauncher after all retries, it prints a fallback code you can paste manually.

### One-off login without storing anything

If you'd rather not save your secret at all, you can pass the OTP Auth URL directly:

```sh
xivotp auth login-url "otpauth://totp/Square%20Enix%20ID:..."
```

This generates a code and sends it to XIVLauncher immediately, without touching the keyring.

### Changing the port

If you've configured XIVLauncher's OTP macro to listen on a different port, use the `--port` flag on any `auth` subcommand:

```sh
xivotp auth login --port 12345
xivotp auth login-url "..." --port 12345
```

The default is `4646`, which matches XIVLauncher's default.

### Clearing secrets

Remove a single account's secret:

```sh
xivotp secret clear --account mymain
```

Or wipe everything:

```sh
xivotp secret clear all
```

## How it works under the hood

```
  ┌──────────┐     ┌─────────────┐     ┌──────────────────┐
  │  xivotp  │────▶│ OS Keyring  │     │  XIVLauncher     │
  │  (CLI)   │     │ (Win/Mac/   │     │  (OTP macro      │
  │          │     │  Linux)     │     │   HTTP server)   │
  └────┬─────┘     └─────────────┘     └────────┬─────────┘
       │                                         │
       │  1. Read OTP Auth URL from keyring      │
       │  2. Generate TOTP code from secret      │
       │  3. GET http://localhost:4646/           │
       │     ffxivlauncher/<code> ──────────────▶│
       │                                         │
       │                        XIVLauncher       │
       │                        fills the code    │
       │                        into the login    │
       │                        form              │
```

- **Secret storage** uses [`go-keyring`](https://github.com/zalando/go-keyring), which wraps the native credential store on each platform — Credential Manager on Windows, Keychain on macOS, and the Secret Service / D-Bus API on Linux.
- **TOTP generation** uses [`pquerna/otp`](https://github.com/pquerna/otp), a well-tested Go implementation of the TOTP standard (RFC 6238).
- The CLI framework is [`urfave/cli`](https://github.com/urfave/cli), which gives us nice subcommand routing, flag parsing, and help text.

## Security notes

- Your OTP secret is stored in your OS keyring, not in a plaintext config file. This is about as good as it gets for a local CLI tool — the keyring is encrypted at rest on most platforms and unlocked when you're logged in.
- The secret is only ever sent to your own machine (`localhost`), never over the network.
- That said, this tool is a convenience, not a fortress. If someone has access to your user session, they can potentially dump your keyring. Treat your machine accordingly.

## License

MIT
