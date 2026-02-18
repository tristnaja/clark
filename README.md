# Clark

A personal AI butler for your WhatsApp.

Clark is a command-line application that runs a sophisticated AI-powered butler for your WhatsApp account. It uses OpenRouter to generate intelligent, context-aware responses, acting as a gatekeeper for your messages while you're away. Clark only interacts with a pre-approved list of "VIP" contacts, ensuring your privacy and focus.

![License](https://img.shields.io/badge/License-MIT-blue.svg)
![Go Version](https://img.shields.io/badge/Go-1.21+-brightgreen.svg)

## How It Works

Clark connects to your WhatsApp account as a client and listens for incoming messages. When a message is received from a recognized VIP, it forwards the conversation to a powerful AI model via OpenRouter. The AI, acting as a professional butler, formulates a response based on your current status and a predefined persona.

**Flow:**
`WhatsApp Message (from VIP) -> Clark (CLI) -> OpenRouter (AI) -> Clark (CLI) -> WhatsApp Reply`

## Features

- **AI-Powered Responses:** Leverages large language models via OpenRouter for natural and intelligent conversations.
- **WhatsApp Integration:** Seamlessly connects to your WhatsApp account using the `whatsmeow` library.
- **Configurable Persona:** The AI operates based on a "Butler Protocol," ensuring all responses are professional and in character.
- **VIP Management:** You control which contacts the bot interacts with through a simple command.
- **Context-Aware:** Set a "master context" (e.g., "In a meeting," "On vacation") to inform the AI of your status.
- **Persistent History:** Stores conversation history in a local SQLite database for context continuity.
- **Easy-to-Use CLI:** Manage the assistant through a straightforward command-line interface.

## Getting Started

Follow these steps to get your personal butler up and running.

### Prerequisites

- **Go (Version 1.21+):** [Installation Guide](https://go.dev/doc/install)
- **OpenRouter API Key:** Get one from the [OpenRouter website](https://openrouter.ai/).
- **WhatsApp Account:** The account you wish to run the butler on.

### Installation & Configuration

1.  **Clone the Repository:**
    ```sh
    git clone https://github.com/tristnaja/clark.git
    cd clark
    ```

2.  **Install Dependencies:**
    ```sh
    go mod tidy
    ```

3.  **Build the Binary:**
    ```sh
    go build .
    ```

4.  **Set Up Environment:**
    Create a `.env` file in the project root and add your OpenRouter API key:
    ```
    OPENROUTER_API="your_openrouter_api_key"
    ```

5.  **Initialize the Assistant:**
    This creates the necessary database and default settings.
    ```sh
    ./clark init
    ```

## Usage

Clark is managed via a set of simple commands.

-   **`run`**: Starts the assistant. On the first run, a QR code will be displayed in your terminal. Scan it with your WhatsApp mobile app (in `Settings > Linked Devices`) to connect your account.
    ```sh
    ./clark run
    ```

-   **`add`**: Adds a contact to the VIP list. The bot will only respond to contacts on this list.
    -   **Format:** `"[number],[name],[relation]"`
    -   `number`: The contact's phone number with country code (e.g., `11234567890`).
    -   `name`: The contact's name.
    -   `relation`: Your relationship to them (e.g., "colleague," "family").
    ```sh
    ./clark add -v "11234567890,John Doe,Colleague"
    ```

-   **`ctx`**: Sets the master context for the AI. This tells the butler your current status.
    ```sh
    ./clark ctx -c "Currently in a board meeting until 5 PM."
    ```

-   **`toggle`**: Toggles the assistant's active status (on/off). The bot will not respond if toggled off.
    ```sh
    ./clark toggle
    ```

-   **`view`**: Displays the current assistant settings, including name, model, status, context, and the VIP list.
    ```sh
    ./clark view
    ```

## License

This project is licensed under the MIT License. See the [LICENSE](./LICENSE) file for details.
