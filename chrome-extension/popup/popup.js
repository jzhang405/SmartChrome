class ChromLLMPopup {
  constructor() {
    this.backendURL = 'http://localhost:8080';
    this.sessionID = null;
    this.conversationID = null;
    this.extractedContent = null;
    this.websocket = null;
    
    this.initializeElements();
    this.setupEventListeners();
    this.checkBackendConnection();
    this.loadStoredSession();
  }

  initializeElements() {
    this.extractBtn = document.getElementById('extractBtn');
    this.sendBtn = document.getElementById('sendBtn');
    this.messageInput = document.getElementById('messageInput');
    this.chatMessages = document.getElementById('chatMessages');
    this.extractionStatus = document.getElementById('extractionStatus');
    this.statusIndicator = document.getElementById('statusIndicator');
    this.statusDot = this.statusIndicator.querySelector('.status-dot');
    this.statusText = this.statusIndicator.querySelector('.status-text');
  }

  setupEventListeners() {
    this.extractBtn.addEventListener('click', () => this.extractContent());
    this.sendBtn.addEventListener('click', () => this.sendMessage());
    this.messageInput.addEventListener('input', () => this.updateSendButton());
    this.messageInput.addEventListener('keydown', (e) => {
      if (e.key === 'Enter' && !e.shiftKey) {
        e.preventDefault();
        this.sendMessage();
      }
    });
  }

  async checkBackendConnection() {
    try {
      const response = await fetch(`${this.backendURL}/v1/health`);
      if (response.ok) {
        this.updateConnectionStatus(true);
      } else {
        this.updateConnectionStatus(false);
      }
    } catch (error) {
      this.updateConnectionStatus(false);
    }
  }

  updateConnectionStatus(connected) {
    if (connected) {
      this.statusDot.classList.add('connected');
      this.statusText.textContent = 'Connected';
    } else {
      this.statusDot.classList.remove('connected');
      this.statusText.textContent = 'Disconnected';
    }
  }

  async loadStoredSession() {
    try {
      const result = await chrome.storage.local.get(['sessionID', 'conversationID']);
      if (result.sessionID) {
        this.sessionID = result.sessionID;
      }
      if (result.conversationID) {
        this.conversationID = result.conversationID;
      }
    } catch (error) {
      console.error('Failed to load stored session:', error);
    }
  }

  async extractContent() {
    this.extractBtn.disabled = true;
    this.extractionStatus.innerHTML = '<div class="loading"></div> Extracting content...';

    try {
      const [tab] = await chrome.tabs.query({ active: true, currentWindow: true });
      
      const results = await chrome.scripting.executeScript({
        target: { tabId: tab.id },
        function: this.extractPageContent,
      });

      if (results && results[0] && results[0].result) {
        const { title, content } = results[0].result;
        this.extractedContent = { title, content };
        this.extractionStatus.textContent = `Extracted: ${title.substring(0, 50)}...`;
        this.updateSendButton();
        
        // Create new conversation
        await this.createConversation(tab.url, title);
      } else {
        this.extractionStatus.textContent = 'Failed to extract content';
      }
    } catch (error) {
      console.error('Content extraction failed:', error);
      this.extractionStatus.textContent = 'Extraction failed';
    } finally {
      this.extractBtn.disabled = false;
    }
  }

  extractPageContent() {
    const title = document.title;
    
    // Remove unwanted elements
    const unwantedSelectors = [
      'script', 'style', 'noscript', 'iframe', 'nav', 'header', 'footer',
      '.advertisement', '.ads', '.sidebar', '.menu', '.navigation'
    ];
    
    unwantedSelectors.forEach(selector => {
      document.querySelectorAll(selector).forEach(el => el.remove());
    });

    // Extract main content
    const content = document.body.innerText || document.body.textContent || '';
    
    return {
      title,
      content: content.trim().substring(0, 5000) // Limit content length
    };
  }

  async createConversation(url, title) {
    try {
      if (!this.sessionID) {
        await this.createSession();
      }

      const response = await fetch(`${this.backendURL}/v1/conversations`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${this.getAuthToken()}`
        },
        body: JSON.stringify({
          url,
          title,
          webpage_content: {
            extracted_text: this.extractedContent.content,
            content_hash: this.generateHash(this.extractedContent.content)
          }
        })
      });

      if (response.ok) {
        const conversation = await response.json();
        this.conversationID = conversation.id;
        await chrome.storage.local.set({ conversationID: conversation.id });
        
        this.addMessage('system', `New conversation started for: ${title}`);
      } else {
        console.error('Failed to create conversation');
      }
    } catch (error) {
      console.error('Conversation creation failed:', error);
    }
  }

  async createSession() {
    try {
      const response = await fetch(`${this.backendURL}/v1/sessions`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({})
      });

      if (response.ok) {
        const session = await response.json();
        this.sessionID = session.session.id;
        await chrome.storage.local.set({ 
          sessionID: session.session.id,
          authToken: session.token 
        });
      }
    } catch (error) {
      console.error('Session creation failed:', error);
    }
  }

  async sendMessage() {
    const message = this.messageInput.value.trim();
    if (!message || !this.conversationID) return;

    this.addMessage('user', message);
    this.messageInput.value = '';
    this.updateSendButton();
    this.sendBtn.disabled = true;

    try {
      const response = await fetch(
        `${this.backendURL}/v1/conversations/${this.conversationID}/messages`,
        {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${this.getAuthToken()}`
          },
          body: JSON.stringify({
            content: message,
            type: 'user_question'
          })
        }
      );

      if (response.ok) {
        const messageData = await response.json();
        this.connectWebSocket(messageData.id);
      } else {
        this.addMessage('assistant', 'Sorry, I failed to process your message.');
      }
    } catch (error) {
      console.error('Message sending failed:', error);
      this.addMessage('assistant', 'Sorry, an error occurred while sending your message.');
    } finally {
      this.sendBtn.disabled = false;
    }
  }

  connectWebSocket(messageID) {
    const wsURL = `ws://localhost:8080/v1/stream?conversationId=${this.conversationID}&messageId=${messageID}`;
    
    this.websocket = new WebSocket(wsURL);
    
    this.websocket.onopen = () => {
      console.log('WebSocket connected');
    };

    this.websocket.onmessage = (event) => {
      const data = JSON.parse(event.data);
      if (data.type === 'stream') {
        this.updateStreamingMessage(data.content, data.data.is_complete);
      } else if (data.type === 'error') {
        this.addMessage('assistant', `Error: ${data.content}`);
      }
    };

    this.websocket.onerror = (error) => {
      console.error('WebSocket error:', error);
      this.addMessage('assistant', 'Sorry, a streaming error occurred.');
    };

    this.websocket.onclose = () => {
      console.log('WebSocket closed');
    };
  }

  updateStreamingMessage(content, isComplete) {
    let lastMessage = this.chatMessages.lastElementChild;
    
    if (!lastMessage || !lastMessage.classList.contains('streaming')) {
      lastMessage = this.addMessage('assistant', '', true);
    }

    const contentElement = lastMessage.querySelector('.content');
    contentElement.textContent += content;

    if (isComplete) {
      lastMessage.classList.remove('streaming');
      const timestamp = new Date().toLocaleTimeString();
      lastMessage.querySelector('.timestamp').textContent = timestamp;
    }
  }

  addMessage(type, content, isStreaming = false) {
    const messageDiv = document.createElement('div');
    messageDiv.className = `message ${type === 'user' ? 'user' : 'assistant'}`;
    if (isStreaming) {
      messageDiv.classList.add('streaming');
    }

    const timestamp = new Date().toLocaleTimeString();
    
    messageDiv.innerHTML = `
      <div class="timestamp">${timestamp}</div>
      <div class="content">${content}</div>
    `;

    this.chatMessages.appendChild(messageDiv);
    this.chatMessages.scrollTop = this.chatMessages.scrollHeight;

    return messageDiv;
  }

  updateSendButton() {
    const hasContent = this.messageInput.value.trim().length > 0;
    const hasExtractedContent = this.extractedContent !== null;
    this.sendBtn.disabled = !(hasContent && hasExtractedContent);
  }

  getAuthToken() {
    // In a real implementation, this would retrieve the stored auth token
    return 'placeholder-token';
  }

  generateHash(content) {
    // Simple hash function for content deduplication
    let hash = 0;
    for (let i = 0; i < content.length; i++) {
      const char = content.charCodeAt(i);
      hash = ((hash << 5) - hash) + char;
      hash = hash & hash; // Convert to 32-bit integer
    }
    return hash.toString(16);
  }
}

// Initialize the popup when the DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
  new ChromLLMPopup();
});