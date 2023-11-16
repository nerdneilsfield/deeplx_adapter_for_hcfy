# Hcfy-Deepl Translation Adapter

This project serves as an adapter that translates text from the Hcfy format to the DeepL format and vice versa. It's designed to work with a specific translation tool that takes highlighted text and sends it to the adapter for translation.

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

Clone the repository to your local machine:

```bash
git clone https://github.com/yourusername/hcfy-deepl-adapter.git
cd hcfy-deepl-adapter
```