---
id: codefresh
title: Codefresh Context
sidebar_label: Codefresh
---

## What

Display the active Codefresh context from the `current-context` field in `$HOME/.cfconfig`.
This segment reads local YAML only; it does not require the Codefresh CLI or make cloud requests.
It is disabled when the file is absent, empty, malformed, or has no non-whitespace string context.

## Sample Configuration

```json
{
  "type": "codefresh",
  "style": "powerline",
  "powerline_symbol": "\uE0B0",
  "foreground": "#ffffff",
  "background": "#3e9022",
  "template": " \uf5f7 {{ .Context }} "
}
```

## Options

None.

## Template ([info][templates])

:::note default template

```template
 {{ .Context }}
```

:::

### Properties

| Name       | Type     | Description                       |
| ---------- | -------- | --------------------------------- |
| `.Context` | `string` | The active Codefresh context name |

:::caution
Debug logging redacts the contents of `.cfconfig`. Decoder errors that could contain credential values are not logged.
:::

[templates]: /docs/configuration/templates
