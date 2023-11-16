# Hcfy-Deepl Translation Adapter

This project serves as an adapter that translates text from the Hcfy format to the DeepL format and vice versa. It's designed to work with a specific translation tool that takes highlighted text and sends it to the adapter for translation.

这个项目的目的是为划词翻译适配 DeepLX，从而可以在自定义翻译源中使用。[自定义翻译源](https://hcfy.app/docs/services/custom-api)

## Features

- HTTP server that accepts POST requests for translation.
- Maps language codes from Hcfy to DeepL standards.
- Error handling for various scenarios including network errors and JSON parsing.
- Connection reuse for frequent external API calls.
- Configurable via environment variables.

## Getting Started

### Prerequisites

- Go 1.x
- Access to DeepL API

### Installation

Clone the repository to your local machine and Build it

```bash
git clone https://github.com/yourusername/hcfy-deepl-adapter.git
cd hcfy-deepl-adapter
CGO_ENABLED=0 GOOS=linux go build -o api_transformer main.go
```

```bash
# set the endpoint
export DEEPLX_ENDPOINT=your_deeplx_api_url/translate
# set the endpoint name
export DEEPLX_NAME=your_deeplx_name
./main
```

Another choice is to use the Docker: `nerdneils/deeplx_adapter_for_hcfy:latest`

```
docker run -it -p 9911:9911 -e "DEEPLX_ENDPOINT=you_deeplx.com/translate" -e "DEEPLX_NAME=your_deeplx_name" nerdneils/deeplx_adapter_for_hcfy:latest
```

### Usage

Send a POST request to `/` with the Hcfy formatted JSON. The service will process the request, call the DeepL API, and return the translation in the Hcfy format.

## Contributing
Contributions are welcome. Please open an issue first to discuss what you would like to change or add.

## License
This project is licensed under the BSD 2-Clause License - see the LICENSE.md file for details.

## Acknowledgments
Thanks to the [DeepLX](https://github.com/OwO-Network/DeepLX) API for providing the translation services.
Thanks to the contributors who maintain the [Hcfy(划词翻译)](https://hcfy.app) translation tool.