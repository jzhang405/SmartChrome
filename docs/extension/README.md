# ChromLLM Chrome Extension Documentation

This directory contains documentation for the ChromLLM Chrome extension.

## Installation

1. Clone the repository
2. Navigate to the `chrome-extension` directory
3. Install dependencies: `npm install`
4. Build the extension: `npm run build`
5. Load the extension in Chrome:
   - Open `chrome://extensions/`
   - Enable "Developer mode"
   - Click "Load unpacked"
   - Select the `dist` directory

## Usage

1. Navigate to any webpage with text content
2. Click the ChromLLM extension icon in the toolbar
3. Click "Extract Page Content" to analyze the page
4. Ask questions about the content in the chat interface
5. Receive real-time responses from the LLM

## Configuration

Configure backend URL and settings in the extension options page.

## Development

- `npm run dev` - Development build with watch
- `npm run test` - Run tests
- `npm run lint` - Run linting
- `npm run build` - Production build