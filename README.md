
# Aoi (葵)

[中文说明](/README-cn.md)

Ghost in the Shell powered by ChatGPT

**Aoi** is an AI-based conversational agent powered by ChatGPT. With Aoi, you can have natural language conversations with an AI in the terminal that can understand your queries and execute appropriate commands.

## Features
You can use Aoi as a terminal version of ChatGPT, Besides, Aoi comes with several built-in features that can help you be more productive:

- `/code` - Generate code snippets and copy to the clipboard , e.g. `/code go generate random numbers`
- `/db` - Generate SQL and execute it on the database, e.g. `/db {url} list tables`
- `/shell` - Generate shell command and execute it, e.g. `/shell view listening ports`
- `/ssh` - Generate shell command and execute it on the remote host, e.g. `/ssh {host} view listening tcp ports`
- `/tldr` - Get a tl;dr explanation of a command
- `/trans` - Translate text to a specified language


## Getting Started
You can download Aoi from the GitHub [release page](https://github.com/shellfly/aoi/releases). Alternatively, you can use Go to install Aoi on your system:

```bash
go install github.com/shellfly/aoi@latest
```

Set your OpenAI API key as an environment variable, and then run the `aoi` command.

```bash
export OPENAI_API_KEY=<your_api_key>
aoi
```

## Demos
[![asciicast](https://asciinema.org/a/XjCGaMNf8Qp2nQ1UDlehjm5AN.svg)](https://asciinema.org/a/XjCGaMNf8Qp2nQ1UDlehjm5AN)

## Contributing
If you find any issues with Aoi or have suggestions for new features, please feel free to create an issue or submit a pull request on the GitHub repository. Contributions from anyone and everyone are welcome!
