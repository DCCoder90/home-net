# OpenClaw Voice Call (Twilio)

:::info
Last Updated 3/25/2026
:::

:::warning
Manual NPM and Cloudflare configuration required before voice calls will work. See Section 4.
:::

## **1. Service Overview**

- **Service Name:** OpenClaw Voice Call Plugin
- **Purpose:** Enables OpenClaw to accept and conduct incoming voice calls via Twilio. Callers speak to the AI agent in real-time using Twilio Programmable Voice + Media Streams for audio and ElevenLabs for text-to-speech.
- **Status:** Pending Setup
- **Parent Service:** [OpenClaw](../../../stacks/ai.yaml) (`claw.dcapi.app` / `192.168.5.50`)

---

## **2. Architecture**

The voice-call plugin runs a **separate HTTP/WebSocket listener on port 3334** inside the OpenClaw container, independent of the main Gateway on port 18789. This allows exposing only the voice endpoints to the internet while keeping the Gateway internal-only.

```
Caller (PSTN)
  -> Twilio
    -> Cloudflare (voice.dcapi.app)
      -> EdgeRouter (port 443)
        -> Nginx Proxy Manager (192.168.4.2)
          -> OpenClaw voice listener (192.168.5.50:3334)
```

**Endpoints served on port 3334:**

| Path | Type | Purpose |
| --- | --- | --- |
| `/voice/webhook` | HTTP POST | Twilio call event webhooks |
| `/voice/stream` | WebSocket | Twilio media stream (real-time audio) |

---

## **3. Secrets & Configuration**

### Infisical Secrets

| Infisical Key | Environment Variable | Purpose |
| --- | --- | --- |
| `twilio_account_sid` | `TWILIO_ACCOUNT_SID` | Twilio Account SID |
| `twilio_auth_token` | `TWILIO_AUTH_TOKEN` | Twilio Auth Token |
| `elevenlabs_api_key` | `ELEVENLABS_API_KEY` | ElevenLabs TTS API key |
| `openclaw_openai_api_key` | `OPENAI_API_KEY` | Required for Twilio Media Streams speech-to-text |

These are defined in `config/stacks/ai.yaml` under the openclaw service's `secrets` block.

### OpenClaw Plugin Config

The voice-call plugin is configured in `/mnt/cache/appdata/openclaw/openclaw.json` (JSON5 format). The Gateway auto-reloads on file changes.

```json5
{
  plugins: {
    entries: {
      "voice-call": {
        enabled: true,
        config: {
          provider: "twilio",
          fromNumber: "+1YOUR_TWILIO_NUMBER",

          twilio: {
            accountSid: "${TWILIO_ACCOUNT_SID}",
            authToken: "${TWILIO_AUTH_TOKEN}",
          },

          serve: {
            port: 3334,
            path: "/voice/webhook",
          },

          publicUrl: "https://voice.dcapi.app/voice/webhook",

          webhookSecurity: {
            allowedHosts: ["voice.dcapi.app"],
          },

          // Inbound calls
          inboundPolicy: "allowlist",
          allowFrom: ["+1ALLOWED_CALLER_NUMBER"],
          inboundGreeting: "Hello! How can I help you today?",
          responseTimeoutMs: 10000,

          // Media streaming (real-time audio over WebSocket)
          streaming: {
            enabled: true,
            streamPath: "/voice/stream",
            preStartTimeoutMs: 5000,
          },

          // Text-to-speech (Microsoft TTS does NOT work for voice calls)
          tts: {
            provider: "elevenlabs",
            elevenlabs: {
              apiKey: "${ELEVENLABS_API_KEY}",
              voiceId: "pMsXgVXv3BLzUgSXRplE",
              modelId: "eleven_multilingual_v2",
            },
          },

          maxDurationSeconds: 300,
          staleCallReaperSeconds: 360,
        },
      },
    },
  },
}
```

### Plugin Installation

```bash
docker exec -it openclaw openclaw plugins install @openclaw/voice-call
```

Restart the container after installation.

---

## **4. Manual Networking Setup (NOT managed by Terraform)**

The voice webhook must be reachable from the internet for Twilio to deliver call events. To avoid exposing the main OpenClaw Gateway (port 18789), a separate domain (`voice.dcapi.app`) routes only to the voice listener (port 3334).

### Step 1: Cloudflare DNS

Add a DNS record for `voice.dcapi.app`:

- **Type:** CNAME (to existing dcapi.app record) or A record (to WAN IP)
- **Proxied:** Yes (orange cloud)
- **Name:** `voice`

### Step 2: Nginx Proxy Manager

Create a new proxy host in [Nginx Proxy Manager](../Networking/Nginx-Proxy.md) (`http://192.168.4.2:81`):

- **Domain Names:** `voice.dcapi.app`
- **Scheme:** `http`
- **Forward Hostname / IP:** `192.168.5.50`
- **Forward Port:** `3334`
- **Websockets Support:** Enabled (required for `/voice/stream`)
- **Access List:** Cloudflare Only
- **SSL:**
  - Use the `*.dcapi.app` wildcard certificate, OR
  - Request a new Let's Encrypt certificate

### Step 3: Twilio Console

Configure the Twilio phone number:

1. Go to **Phone Numbers** -> select your number
2. Under **Voice Configuration**:
   - **"A Call Comes In"** -> Webhook
   - **URL:** `https://voice.dcapi.app/voice/webhook`
   - **HTTP Method:** POST

---

## **5. Twilio-Specific Behavior**

- **Signature verification** is always enforced in production. The plugin reconstructs the public URL from forwarded headers using `webhookSecurity.allowedHosts`.
- **Replay protection** is built in. Each conversation turn gets a per-turn token in Gather callbacks.
- **Stream disconnect grace period:** 2000ms before auto-ending a call when the media stream disconnects. Reconnection during this window cancels auto-end.
- **TTS:** When streaming is enabled, core TTS (ElevenLabs) handles speech, not Twilio's native voices. If the stream is active and TTS is unavailable, playback fails rather than falling back.

---

## **6. Verification & Debugging**

```bash
# Check for config issues
docker exec -it openclaw openclaw doctor

# Confirm plugin is loaded
docker exec -it openclaw openclaw plugins list

# Watch call events in real time
docker exec -it openclaw openclaw voicecall tail

# Check latency stats
docker exec -it openclaw openclaw voicecall latency --last 20

# Get status of a specific call
docker exec -it openclaw openclaw voicecall status --call-id <id>
```

---

## **7. Dependencies**

- **[OpenClaw](../../../stacks/ai.yaml):** Parent service (`192.168.5.50:18789`)
- **[Nginx Proxy Manager](../Networking/Nginx-Proxy.md):** Reverse proxy for `voice.dcapi.app` -> `192.168.5.50:3334`
- **[EdgeRouter Pro](../../Hardware/Routers/EdgeRouter%20Pro.md):** Port forwarding (80/443) to NPM
- **[Cloudflare](../External%20Services/Cloudflare.md):** DNS for `voice.dcapi.app`
- **Twilio:** Telephony provider (external)
- **ElevenLabs:** Text-to-speech provider (external)

---

## **8. Useful Links**

- **OpenClaw Voice Call Docs:** https://docs.openclaw.ai/plugins/voice-call
- **Twilio Console:** https://console.twilio.com
- **ElevenLabs:** https://elevenlabs.io
- **OpenClaw Plugin Config Docs:** https://docs.openclaw.ai/configuration
