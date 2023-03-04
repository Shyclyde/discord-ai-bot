# Discord Admin and AI bot

This Discord bot written in Go is currently meant for:

- Assisting in managing another bot (such as restarting when it goes wrong).
- OpenAI chat and OpenAI image generation, if enabled.
- Game server information and management.
- Other general Discord administration functions.

## Setup

### Environment

Copy both the `.env_sample` and the `config_sample.json` to non `_sample` files.

Tweak the values in `config.json` as you see fit. If you need assistance with the OpenAI values, see the OpenAI docs for more information:

[OpenAI Text Completion API Docs](https://platform.openai.com/docs/api-reference/completions)
[OpenAI Image Creation API Docs](https://platform.openai.com/docs/api-reference/images)

The `.env` values require you to have setup a Discord App and obtained the token for your Discord App, as well as have an OpenAI account and have your API key ready.
