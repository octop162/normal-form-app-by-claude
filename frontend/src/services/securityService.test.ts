import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { securityService } from './securityService';

// Mock fetch
const mockFetch = vi.fn();
global.fetch = mockFetch;

describe('SecurityService', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    // Clear any cached tokens
    securityService.clearCSRFToken();
    localStorage.clear();
    sessionStorage.clear();
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  describe('CSRF Token Management', () => {
    it('should fetch and cache CSRF token', async () => {
      const mockToken = 'test-csrf-token-123';
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          success: true,
          data: { token: mockToken },
        }),
      });

      const token = await securityService.getCSRFToken();

      expect(token).toBe(mockToken);
      expect(mockFetch).toHaveBeenCalledWith('/api/v1/csrf-token', {
        method: 'GET',
        credentials: 'same-origin',
      });
    });

    it('should return cached token on subsequent calls', async () => {
      const mockToken = 'cached-token-456';
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          success: true,
          data: { token: mockToken },
        }),
      });

      // First call
      const token1 = await securityService.getCSRFToken();
      // Second call
      const token2 = await securityService.getCSRFToken();

      expect(token1).toBe(mockToken);
      expect(token2).toBe(mockToken);
      expect(mockFetch).toHaveBeenCalledTimes(1);
    });

    it('should handle CSRF token fetch failure', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 500,
      });

      await expect(securityService.getCSRFToken()).rejects.toThrow('Failed to fetch CSRF token');
    });

    it('should clear cached token', async () => {
      const mockToken = 'token-to-clear';
      mockFetch.mockResolvedValue({
        ok: true,
        json: async () => ({
          success: true,
          data: { token: mockToken },
        }),
      });

      // Get token (should cache)
      await securityService.getCSRFToken();
      
      // Clear token
      securityService.clearCSRFToken();
      
      // Get token again (should fetch again)
      await securityService.getCSRFToken();

      expect(mockFetch).toHaveBeenCalledTimes(2);
    });
  });

  describe('Input Sanitization', () => {
    it('should detect script tags', () => {
      const maliciousInput = '<script>alert("xss")</script>';
      const result = securityService.containsSuspiciousPatterns(maliciousInput);
      expect(result).toBe(true);
    });

    it('should detect javascript: protocol', () => {
      const maliciousInput = 'javascript:alert("xss")';
      const result = securityService.containsSuspiciousPatterns(maliciousInput);
      expect(result).toBe(true);
    });

    it('should detect SQL injection patterns', () => {
      const maliciousInputs = [
        "' OR '1'='1",
        '; DROP TABLE users;--',
        '1\' UNION SELECT * FROM users--',
      ];

      maliciousInputs.forEach(input => {
        const result = securityService.containsSuspiciousPatterns(input);
        expect(result).toBe(true);
      });
    });

    it('should allow clean input', () => {
      const cleanInputs = [
        'John Doe',
        'john@example.com',
        '東京都千代田区',
        '123-4567',
        'Normal text input',
      ];

      cleanInputs.forEach(input => {
        const result = securityService.containsSuspiciousPatterns(input);
        expect(result).toBe(false);
      });
    });

    it('should sanitize input by removing suspicious patterns', () => {
      const maliciousInput = 'Hello <script>alert("xss")</script> World';
      const sanitized = securityService.sanitizeInput(maliciousInput);
      expect(sanitized).toBe('Hello  World');
    });
  });

  describe('Validation Functions', () => {
    describe('Phone Number Validation', () => {
      it('should validate correct mobile phone numbers', () => {
        const validMobile = securityService.isValidPhoneNumber('090', '1234', '5678');
        expect(validMobile).toBe(true);

        const validMobile2 = securityService.isValidPhoneNumber('080', '9876', '5432');
        expect(validMobile2).toBe(true);

        const validMobile3 = securityService.isValidPhoneNumber('070', '1111', '2222');
        expect(validMobile3).toBe(true);
      });

      it('should validate correct landline phone numbers', () => {
        const validLandline = securityService.isValidPhoneNumber('03', '1234', '5678');
        expect(validLandline).toBe(true);

        const validLandline2 = securityService.isValidPhoneNumber('06', '6789', '0123');
        expect(validLandline2).toBe(true);
      });

      it('should reject toll-free numbers', () => {
        const tollFree1 = securityService.isValidPhoneNumber('0120', '123', '456');
        expect(tollFree1).toBe(false);

        const tollFree2 = securityService.isValidPhoneNumber('0800', '123', '456');
        expect(tollFree2).toBe(false);
      });

      it('should reject invalid mobile patterns', () => {
        const invalid1 = securityService.isValidPhoneNumber('050', '1234', '5678');
        expect(invalid1).toBe(false);

        const invalid2 = securityService.isValidPhoneNumber('090', '123', '5678');
        expect(invalid2).toBe(false);
      });
    });

    describe('Postal Code Validation', () => {
      it('should validate correct postal codes', () => {
        const valid1 = securityService.isValidPostalCode('100', '0001');
        expect(valid1).toBe(true);

        const valid2 = securityService.isValidPostalCode('530', '0001');
        expect(valid2).toBe(true);
      });

      it('should reject invalid postal codes', () => {
        const invalid1 = securityService.isValidPostalCode('12', '0001');
        expect(invalid1).toBe(false);

        const invalid2 = securityService.isValidPostalCode('100', '001');
        expect(invalid2).toBe(false);

        const invalid3 = securityService.isValidPostalCode('abc', '0001');
        expect(invalid3).toBe(false);
      });
    });

    describe('Email Validation', () => {
      it('should validate correct email addresses', () => {
        const validEmails = [
          'test@example.com',
          'user.name@domain.co.jp',
          'user+tag@example.org',
          'user123@test-domain.com',
        ];

        validEmails.forEach(email => {
          const result = securityService.isValidEmail(email);
          expect(result).toBe(true);
        });
      });

      it('should reject invalid email addresses', () => {
        const invalidEmails = [
          'invalid-email',
          '@domain.com',
          'user@',
          'user..name@domain.com',
          'user@domain',
          'user@.com',
        ];

        invalidEmails.forEach(email => {
          const result = securityService.isValidEmail(email);
          expect(result).toBe(false);
        });
      });
    });
  });

  describe('Session Data Management', () => {
    it('should encrypt and store session data', () => {
      const testData = { name: 'John', email: 'john@example.com' };
      const sessionId = 'test-session-123';

      securityService.storeSecureSessionData(sessionId, testData);

      expect(sessionStorage.setItem).toHaveBeenCalledWith(
        expect.stringContaining(sessionId),
        expect.any(String)
      );
    });

    it('should retrieve and decrypt session data', () => {
      const testData = { name: 'John', email: 'john@example.com' };
      const sessionId = 'test-session-456';

      // Mock sessionStorage to return encrypted data
      const encryptedData = JSON.stringify({
        data: btoa(JSON.stringify(testData)),
        timestamp: Date.now(),
        expiry: Date.now() + 4 * 60 * 60 * 1000, // 4 hours
      });

      vi.mocked(sessionStorage.getItem).mockReturnValue(encryptedData);

      const retrievedData = securityService.getSecureSessionData(sessionId);

      expect(retrievedData).toEqual(testData);
    });

    it('should return null for expired session data', () => {
      const testData = { name: 'John', email: 'john@example.com' };
      const sessionId = 'expired-session';

      // Mock sessionStorage to return expired data
      const expiredData = JSON.stringify({
        data: btoa(JSON.stringify(testData)),
        timestamp: Date.now() - 5 * 60 * 60 * 1000, // 5 hours ago
        expiry: Date.now() - 1 * 60 * 60 * 1000, // 1 hour ago (expired)
      });

      vi.mocked(sessionStorage.getItem).mockReturnValue(expiredData);

      const retrievedData = securityService.getSecureSessionData(sessionId);

      expect(retrievedData).toBeNull();
    });

    it('should clear session data', () => {
      const sessionId = 'session-to-clear';

      securityService.clearSecureSessionData(sessionId);

      expect(sessionStorage.removeItem).toHaveBeenCalledWith(
        expect.stringContaining(sessionId)
      );
    });
  });

  describe('Content Security', () => {
    it('should detect potentially malicious content', () => {
      const maliciousContents = [
        '<iframe src="javascript:alert(1)">',
        '<img onerror="alert(1)" src="x">',
        'javascript:void(0)',
        'data:text/html,<script>alert(1)</script>',
        'vbscript:msgbox(1)',
      ];

      maliciousContents.forEach(content => {
        const result = securityService.containsSuspiciousPatterns(content);
        expect(result).toBe(true);
      });
    });

    it('should allow safe content', () => {
      const safeContents = [
        'Hello World',
        'user@example.com',
        '東京都千代田区丸の内1-1-1',
        '03-1234-5678',
        'https://example.com/page',
      ];

      safeContents.forEach(content => {
        const result = securityService.containsSuspiciousPatterns(content);
        expect(result).toBe(false);
      });
    });
  });

  describe('Error Handling', () => {
    it('should handle network errors gracefully', async () => {
      mockFetch.mockRejectedValueOnce(new Error('Network error'));

      await expect(securityService.getCSRFToken()).rejects.toThrow('Failed to fetch CSRF token');
    });

    it('should handle malformed JSON responses', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => {
          throw new Error('Invalid JSON');
        },
      });

      await expect(securityService.getCSRFToken()).rejects.toThrow('Failed to fetch CSRF token');
    });

    it('should handle missing token in response', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          success: true,
          data: {}, // Missing token
        }),
      });

      await expect(securityService.getCSRFToken()).rejects.toThrow('Invalid CSRF token response');
    });
  });
});