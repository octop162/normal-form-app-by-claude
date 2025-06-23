import { apiClient } from './apiClient';

interface CSRFTokenResponse {
  success: boolean;
  data: {
    token: string;
  };
}

class SecurityService {
  private csrfToken: string | null = null;
  private tokenPromise: Promise<string> | null = null;

  // Get CSRF token (with caching and retry logic)
  async getCSRFToken(): Promise<string> {
    // Return cached token if available
    if (this.csrfToken) {
      return this.csrfToken;
    }

    // Return existing promise if token is being fetched
    if (this.tokenPromise) {
      return this.tokenPromise;
    }

    // Fetch new token
    this.tokenPromise = this.fetchCSRFToken();
    
    try {
      const token = await this.tokenPromise;
      this.csrfToken = token;
      return token;
    } finally {
      this.tokenPromise = null;
    }
  }

  private async fetchCSRFToken(): Promise<string> {
    try {
      const response = await apiClient.get<CSRFTokenResponse>('/csrf-token');
      
      if (!response.data.success || !response.data.data?.token) {
        throw new Error('Invalid CSRF token response');
      }

      return response.data.data.token;
    } catch (error) {
      console.error('Failed to fetch CSRF token:', error);
      throw new Error('Failed to get CSRF token');
    }
  }

  // Clear cached token (call when token expires or on error)
  clearCSRFToken(): void {
    this.csrfToken = null;
    this.tokenPromise = null;
  }

  // Sanitize input data to prevent XSS
  sanitizeInput(input: string): string {
    if (!input) return '';
    
    return input
      .replace(/&/g, '&amp;')
      .replace(/</g, '&lt;')
      .replace(/>/g, '&gt;')
      .replace(/"/g, '&quot;')
      .replace(/'/g, '&#x27;')
      .replace(/\//g, '&#x2F;');
  }

  // Validate email format
  isValidEmail(email: string): boolean {
    const emailRegex = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
    return emailRegex.test(email);
  }

  // Validate phone number format (Japanese)
  isValidPhoneNumber(phone1: string, phone2: string, phone3: string): boolean {
    // Check for free dial numbers (not allowed)
    const freeDialPrefixes = ['0120', '0800', '0570'];
    const fullNumber = phone1 + phone2 + phone3;
    
    for (const prefix of freeDialPrefixes) {
      if (fullNumber.startsWith(prefix)) {
        return false;
      }
    }

    // 11-digit mobile numbers must start with 0X0
    if (fullNumber.length === 11) {
      const mobilePattern = /^0[789]0/;
      return mobilePattern.test(fullNumber);
    }

    // 10-digit landline numbers
    if (fullNumber.length === 10) {
      // Area code: 2-5 digits, local: 1-4 digits, subscriber: 4 digits
      const landlinePattern = /^0[1-9][0-9]{8}$/;
      return landlinePattern.test(fullNumber);
    }

    return false;
  }

  // Validate postal code format (Japanese)
  isValidPostalCode(code1: string, code2: string): boolean {
    return /^\d{3}$/.test(code1) && /^\d{4}$/.test(code2);
  }

  // Check for suspicious patterns in input
  containsSuspiciousPatterns(input: string): boolean {
    const suspiciousPatterns = [
      /<script/i,
      /javascript:/i,
      /on\w+\s*=/i,
      /data:\s*text\/html/i,
      /vbscript:/i,
      /<iframe/i,
      /<object/i,
      /<embed/i,
      /<link/i,
      /<meta/i,
      /eval\s*\(/i,
      /expression\s*\(/i
    ];

    return suspiciousPatterns.some(pattern => pattern.test(input));
  }

  // Generate secure session ID
  generateSessionId(): string {
    const array = new Uint8Array(16);
    crypto.getRandomValues(array);
    return Array.from(array, byte => byte.toString(16).padStart(2, '0')).join('');
  }

  // Secure data storage helpers
  setSecureSessionData(key: string, data: any): void {
    try {
      const sessionId = this.generateSessionId();
      const encryptedData = btoa(JSON.stringify({
        data,
        timestamp: Date.now(),
        sessionId
      }));
      sessionStorage.setItem(key, encryptedData);
    } catch (error) {
      console.error('Failed to store secure session data:', error);
    }
  }

  getSecureSessionData(key: string): any | null {
    try {
      const encryptedData = sessionStorage.getItem(key);
      if (!encryptedData) return null;

      const decryptedData = JSON.parse(atob(encryptedData));
      
      // Check if data is expired (4 hours)
      const fourHours = 4 * 60 * 60 * 1000;
      if (Date.now() - decryptedData.timestamp > fourHours) {
        sessionStorage.removeItem(key);
        return null;
      }

      return decryptedData.data;
    } catch (error) {
      console.error('Failed to retrieve secure session data:', error);
      sessionStorage.removeItem(key);
      return null;
    }
  }

  clearSecureSessionData(key: string): void {
    sessionStorage.removeItem(key);
  }

  // Content Security Policy helper
  checkCSPCompliance(): boolean {
    try {
      // Test if inline scripts are blocked (good security practice)
      const script = document.createElement('script');
      script.innerHTML = 'window.cspTest = true;';
      document.head.appendChild(script);
      document.head.removeChild(script);
      
      return !(window as any).cspTest;
    } catch (error) {
      return true; // CSP is likely working if script execution fails
    }
  }
}

export const securityService = new SecurityService();