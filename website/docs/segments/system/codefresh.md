---
id: codefresh
title: Codefresh Context
sidebar_label: Codefresh
---

## What

Display the currently active Codefresh context name.

## Sample Configuration

```json
{
  "type": "codefresh",
  "style": "powerline",
  "powerline_symbol": "\uE0B0",
  "foreground": "#ffffff",
  "background": "#3e9022",
  "template": " \uf5f7 {{.Context}} "
}
```

## Properties

- None

## Template ([info][templates])

:::note default template

``` template
{{ .Context }}
```

:::

### Properties

- `.Context`: `string` - the current codefresh context

## Tips

This segment assumes that the .cfconfig which stores the context for the Codefresh CLI is stores in the HOME directory.

[templates]: /docs/configuration/templates